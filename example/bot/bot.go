package bot

import (
	"context"
	"errors"
	"fmt"
	"github.com/lugobots/coach"
	"github.com/lugobots/lugo4go/v2/lugo"
	"github.com/lugobots/lugo4go/v2/pkg/field"
	"github.com/lugobots/lugo4go/v2/pkg/util"
	"github.com/lugobots/lugo4go/v2/team"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"io/ioutil"
	"log"
	"path/filepath"
)

var session *tf.Session
var graph *tf.Graph

func init() {
	path := "/home/rubens/projects/reading-modeal-experiments/model/"
	//region op1
	//Load a frozen graph to use for queries
	modelpath := filepath.Join(path, "valendo-2020-07-18-15-59.pb")
	//modelpath := filepath.Join(path, "saved_model_mapped.pb")
	model, err := ioutil.ReadFile(modelpath)
	if err != nil {
		log.Fatal(err)
	}

	// Construct an in-memory graph from the serialized form.
	graph = tf.NewGraph()
	if err := graph.Import(model, ""); err != nil {
		log.Fatalf("graph: %s", err)
	}

	// Create a session for inference over graph.
	session, err = tf.NewSession(graph, nil)
	if err != nil {
		log.Fatal(err)
	}
	//defer session.Close()
}

type Bot struct {
	Side   lugo.Team_Side
	Number uint32
	Logger util.Logger
	arr    team.Positioner
}

func NewBot(logger util.Logger, side lugo.Team_Side, number uint32) *Bot {
	arr, _ := team.NewArrangement(team.MaxCols, team.MaxRows, side)
	return &Bot{
		Logger: logger,
		Number: number,
		Side:   side,
		arr:    arr,
	}
}

func (b *Bot) OnDisputing(ctx context.Context, sender team.OrderSender, snapshot *lugo.GameSnapshot) error {
	return b.myDecider(ctx, sender, snapshot, team.DisputingTheBall)
}

func (b *Bot) OnDefending(ctx context.Context, sender team.OrderSender, snapshot *lugo.GameSnapshot) error {
	return b.myDecider(ctx, sender, snapshot, team.Defending)
}

func (b *Bot) OnHolding(ctx context.Context, sender team.OrderSender, snapshot *lugo.GameSnapshot) error {
	return b.myDecider(ctx, sender, snapshot, team.HoldingTheBall)
}

func (b *Bot) OnSupporting(ctx context.Context, sender team.OrderSender, snapshot *lugo.GameSnapshot) error {
	return b.myDecider(ctx, sender, snapshot, team.Supporting)
}

func (b *Bot) AsGoalkeeper(ctx context.Context, sender team.OrderSender, snapshot *lugo.GameSnapshot, state team.PlayerState) error {
	return b.myDecider(ctx, sender, snapshot, state)
}

func (b *Bot) myDecider(ctx context.Context, sender team.OrderSender, snapshot *lugo.GameSnapshot, state team.PlayerState) error {
	var orders []lugo.PlayerOrder
	// we are going to kick the ball as soon as we catch it
	me := field.GetPlayer(snapshot, b.Side, b.Number)
	if me == nil {
		return errorHandler(b.Logger, errors.New("bot not found in the game snapshot"))
	}
	if state == team.HoldingTheBall {

		orderToKick, err := whereShouldIGo(snapshot, me.TeamSide, me.Number)

		target := field.GetOpponentGoal(me.TeamSide).BottomPole
		target.Y += 50
		//orderToKick, err := field.MakeOrderKick(*snapshot.Ball, target, field.BallMaxSpeed)
		if err != nil {
			return errorHandler(b.Logger, fmt.Errorf("could not create kick order during turn %d: %s", snapshot.Turn, err))
		}
		orders = []lugo.PlayerOrder{orderToKick}
	} else if me.Number == 10 {
		if snapshot.ShotClock != nil && snapshot.ShotClock.TeamSide != me.TeamSide {
			//p, _ := b.arr.GetPointRegion(*me.Position)
			orderToMove, err := whereShouldIGo(snapshot, me.TeamSide, me.Number)
			//
			//orderToMove, err := field.MakeOrderMoveMaxSpeed(*me.Position, field.FieldCenter())
			if err != nil {
				return errorHandler(b.Logger, fmt.Errorf("staying in the centger %d: %s", snapshot.Turn, err))
			}
			b.Logger.Infof("dir: %v", orderToMove.Move)
			orders = []lugo.PlayerOrder{orderToMove, field.MakeOrderCatch()}
		} else {
			orderToMove, err := field.MakeOrderMoveMaxSpeed(*me.Position, *snapshot.Ball.Position)
			if err != nil {
				return errorHandler(b.Logger, fmt.Errorf("could not create move order during turn %d: %s", snapshot.Turn, err))
			}
			orders = []lugo.PlayerOrder{orderToMove, field.MakeOrderCatch()}
		}

		// otherwise, let's run towards the ball like kids
	} else {
		orders = []lugo.PlayerOrder{field.MakeOrderCatch()}
	}

	resp, err := sender.Send(ctx, snapshot.Turn, orders, "")
	if err != nil {
		return errorHandler(b.Logger, fmt.Errorf("could not send kick order during turn %d: %s", snapshot.Turn, err))
	} else if resp.Code != lugo.OrderResponse_SUCCESS {
		return errorHandler(b.Logger, fmt.Errorf("order sent not  order during turn %d: %s", snapshot.Turn, err))
	}
	return nil
}

type pointCoord []float32

type req struct {
	Instances []pointCoord `json:"instances"`
}

type Result struct {
	Predictions [][]float32 `json:"predictions"`
}

const (
	UNKNOWN  = -1
	LEFT     = 0
	RIGHT    = 1
	FORWARD  = 2
	BACKWARD = 3
)

func (r Result) GetBest() int {
	best := UNKNOWN
	highestV := float32(0.0)
	for ind, value := range r.Predictions[0] {
		if value > highestV {
			best = ind
			highestV = value
		}
	}
	return best
}

var class_names = []string{"left", "right", "forward", "backward"}

func whereShouldIGo(snap *lugo.GameSnapshot, side lugo.Team_Side, referenceNumber uint32) (*lugo.Order_Move, error) {
	//url := "http://localhost:8501/v1/models/saved_model:predict"
	//fmt.Printf("URL: %s", url)

	ref, err := coach.PlayerReference(snap, side, referenceNumber)
	inputs := ref.ExportGroups(
		coach.GroupBall,
		coach.GroupOpponentGoal,
		coach.GroupMyTeam,
		coach.GroupOpponentTeam)

	//inputs, _ := coach.NewMapped(snap, side, referenceNumber)

	tensor, terr := tf.NewTensor([][][]float32{inputs})
	if terr != nil {
		fmt.Printf("Error creating input tensor: %s\n", terr.Error())
		panic(terr)
	}

	output, err := session.Run(
		map[tf.Output]*tf.Tensor{
			graph.Operation("x").Output(0): tensor,
		},
		[]tf.Output{
			//graph.Operation("sequential/dense/MatMul/ReadVariableOp").Output(0),
			//graph.Operation("sequential/dense_1/BiasAdd/ReadVariableOp").Output(0),
			graph.Operation("Identity").Output(0),
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	t := output[0].Value().([][]float32)

	best, score := higherindex(t[0])
	fmt.Printf("Got %d (%f) (%v)\n", best, score, t[0])

	switch best {
	case LEFT:
		return field.GoLeft(side), nil
	case RIGHT:
		return field.GoRight(side), nil
	case FORWARD:
		return field.GoForward(side), nil
	case BACKWARD:
		return field.GoBackward(side), nil
	}
	panic(best)
	return nil, fmt.Errorf("oops")
}

func higherindex(values []float32) (int32, float32) {
	score := float32(-999)
	chosen := int32(-1)
	for i, s := range values {
		if s > score {
			score = s
			chosen = int32(i)
		}
	}
	return chosen, values[chosen]
}

func errorHandler(logger util.Logger, err error) error {
	logger.Errorf("bot error: %s", err)
	return err
}
