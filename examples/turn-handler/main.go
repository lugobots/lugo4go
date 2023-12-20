package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"go.uber.org/zap"

	clientGo "github.com/lugobots/lugo4go/v3"
	"github.com/lugobots/lugo4go/v3/mapper"
	"github.com/lugobots/lugo4go/v3/proto"
	"github.com/lugobots/lugo4go/v3/specs"
)

func main() {

	connectionStarter, err := clientGo.NewTurnHandlerConfig()
	if err != nil {
		log.Fatalf("failed to load the bot configuration: %s", err)
	}

	//
	// Optional: define your own field mapper
	//
	//playerMapper, err := mapper.NewMapper(32, 15, connectionStarter.Config.TeamSide)
	//if err != nil {
	//	log.Fatalf("failed to create a field mapper: %s", err)
	//}
	//connectionStarter.FieldMapper = playerMapper
	//

	if err := connectionStarter.Run(&TurnHandler{
		FieldMapper: connectionStarter.FieldMapper,
		Config:      connectionStarter.Config,
		Logger:      connectionStarter.Logger,
	}); err != nil {
		log.Fatalf("bot stopped: %s", err)
	}
}

type TurnHandler struct {
	FieldMapper mapper.Mapper
	Config      clientGo.Config
	Logger      *zap.SugaredLogger
}

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (t *TurnHandler) GetReadyHandler(ctx context.Context, snapshot clientGo.SnapshotInspector) {
	t.Logger.Debug("the game is ready to start or the score has changed")
}

func (t *TurnHandler) TurnHandler(ctx context.Context, inspector clientGo.SnapshotInspector) ([]proto.PlayerOrder, string, error) {
	var orders []proto.PlayerOrder
	// we are going to kick the ball as soon as we catch it
	me := inspector.GetMe()

	if inspector.IsBallHolder(me) {
		orderToKick, err := inspector.MakeOrderKick(t.FieldMapper.GetOpponentGoal().Center, specs.BallMaxSpeed)
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
		orders = append(orders, inspector.MakeOrderMoveByDirection(mapper.Forward, specs.BallMaxSpeed))
		debugMsg = "moving Forward"
	case 1:
		orders = append(orders, inspector.MakeOrderMoveByDirection(mapper.Backward, specs.BallMaxSpeed))
		debugMsg = "moving Backward"
	case 2:
		orders = append(orders, inspector.MakeOrderMoveByDirection(mapper.Right, specs.BallMaxSpeed))
		debugMsg = "moving to the right"
	case 3:
		orders = append(orders, inspector.MakeOrderMoveByDirection(mapper.Left, specs.BallMaxSpeed))
		debugMsg = "moving to the left"
	}

	return orders, debugMsg, nil
}
