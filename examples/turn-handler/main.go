package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"go.uber.org/zap"

	clientGo "github.com/lugobots/lugo4go/v3"
	"github.com/lugobots/lugo4go/v3/field"
	"github.com/lugobots/lugo4go/v3/proto"
	"github.com/lugobots/lugo4go/v3/specs"
)

func main() {

	connectionStarter, defaultFieldMapper, err := clientGo.NewDefaultStarter()
	if err != nil {
		log.Fatalf("failed to load the bot configuration: %s", err)
	}

	//
	// Optional: define your own field mapper
	// defaultFieldMapper, err = field.NewMapper(NUM_COLS, NUM_ROWS, connectionStarter.Config.TeamSide)
	// if err != nil {
	// 	log.Fatalf("failed to create a field mapper: %s", err)
	// }

	if err := connectionStarter.RunJustTurnHandler(&BasicBot{
		FieldMapper: defaultFieldMapper,
		Config:      connectionStarter.Config,
		Logger:      connectionStarter.Logger,
	}); err != nil {
		log.Fatalf("bot stopped: %s", err)
	}
}

type BasicBot struct {
	FieldMapper field.Mapper
	Config      clientGo.Config
	Logger      *zap.SugaredLogger
}

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (t *BasicBot) GetReadyHandler(_ context.Context, _ clientGo.SnapshotInspector) {
	t.Logger.Debug("the game is ready to start or the score has changed")
}

func (t *BasicBot) TurnHandler(_ context.Context, inspector clientGo.SnapshotInspector) ([]proto.PlayerOrder, string, error) {
	var orders []proto.PlayerOrder
	// we are going to kick the ball as soon as we catch it
	me := inspector.GetMe()

	if inspector.IsBallHolder(me) {
		orderToKick, err := inspector.MakeOrderKick(t.FieldMapper.GetAttackGoal().Center, specs.BallMaxSpeed)
		if err != nil {
			t.Logger.Errorf("could not create kick order during turn %d: %s", inspector.GetSnapshot().Turn, err)
			return nil, "", err
		}
		return []proto.PlayerOrder{orderToKick}, "just kick it", nil
	}

	if me.Number == 10 {
		// otherwise, let's run towards the ball like kids
		orderToMove, err := inspector.MakeOrderMoveMaxSpeed(*inspector.GetBall().Position)
		if err != nil {
			t.Logger.Errorf("could not create move order during turn %d: %s", inspector.GetSnapshot().Turn, err)
			return nil, "", err
		}
		return []proto.PlayerOrder{orderToMove, inspector.MakeOrderCatch()}, "advancing because I am the number 10", nil
	}

	orders = []proto.PlayerOrder{inspector.MakeOrderCatch()}
	debugMsg := "keeping direction"
	switch random.Intn(30) {
	case 0:
		orders = append(orders, inspector.MakeOrderMoveByDirection(field.Forward, specs.BallMaxSpeed))
		debugMsg = "moving Forward"
	case 1:
		orders = append(orders, inspector.MakeOrderMoveByDirection(field.Backward, specs.BallMaxSpeed))
		debugMsg = "moving Backward"
	case 2:
		orders = append(orders, inspector.MakeOrderMoveByDirection(field.Right, specs.BallMaxSpeed))
		debugMsg = "moving to the right"
	case 3:
		orders = append(orders, inspector.MakeOrderMoveByDirection(field.Left, specs.BallMaxSpeed))
		debugMsg = "moving to the left"
	}

	return orders, debugMsg, nil
}
