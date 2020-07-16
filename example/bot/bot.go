package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lugobots/coach"
	"github.com/lugobots/lugo4go/v2/lugo"
	"github.com/lugobots/lugo4go/v2/pkg/field"
	"github.com/lugobots/lugo4go/v2/pkg/util"
	"github.com/lugobots/lugo4go/v2/team"
	"io/ioutil"
	"net/http"
)

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
		target := field.GetOpponentGoal(me.TeamSide).BottomPole
		target.Y += 50
		orderToKick, err := field.MakeOrderKick(*snapshot.Ball, target, field.BallMaxSpeed)
		if err != nil {
			return errorHandler(b.Logger, fmt.Errorf("could not create kick order during turn %d: %s", snapshot.Turn, err))
		}
		orders = []lugo.PlayerOrder{orderToKick}
	} else if me.Number == 10 {
		if snapshot.ShotClock != nil && snapshot.ShotClock.TeamSide != me.TeamSide {
			p, _ := b.arr.GetPointRegion(*me.Position)
			dir, err := whereShouldIGo(p, me.TeamSide)

			//orderToMove, err := field.MakeOrderMoveMaxSpeed(*me.Position, field.FieldCenter())
			if err != nil {
				return errorHandler(b.Logger, fmt.Errorf("staying in the centger %d: %s", snapshot.Turn, err))
			}
			orders = []lugo.PlayerOrder{dir, field.MakeOrderCatch()}
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

func whereShouldIGo(region team.FieldNav, side lugo.Team_Side) (*lugo.Order_Move, error) {
	url := "http://localhost:8501/v1/models/saved_model:predict"
	fmt.Printf("URL: %s", url)

	r := req{Instances: []pointCoord{
		{float32(region.Col()) / float32(team.MaxCols), float32(region.Row()) / float32(team.MaxRows)},
	}}

	//var jsonStr = []byte(`{"Instances":[[0.2374449339,0.7168710297]]}`)
	jsonStr, err := json.Marshal(&r)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n\nReq:  %s\n", jsonStr)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Printf("response Status:  %s\n", resp.Status)
	fmt.Printf("response Headers: %s\n", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)

	result := Result{}

	if err := json.Unmarshal(body, &result); err != nil {
		panic(err)
	}

	fmt.Printf("response Body:  %s\n", string(body))
	fmt.Printf("Result:  %v (%s)\n", result.GetBest(), class_names[result.GetBest()])

	switch result.GetBest() {
	case LEFT:
		return field.GoLeft(side), nil
	case RIGHT:
		return field.GoRight(side), nil
	case FORWARD:
		return field.GoForward(side), nil
	case BACKWARD:
		return field.GoBackward(side), nil
	}
	return nil, fmt.Errorf("oops")
}

func errorHandler(logger util.Logger, err error) error {
	logger.Errorf("bot error: %s", err)
	return err
}
