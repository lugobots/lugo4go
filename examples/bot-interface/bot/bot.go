package bot

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.uber.org/zap"

	"github.com/lugobots/lugo4go/v3"
	"github.com/lugobots/lugo4go/v3/mapper"
	"github.com/lugobots/lugo4go/v3/proto"
	"github.com/lugobots/lugo4go/v3/specs"
)

func NewBot(FieldMapper mapper.Mapper, Config lugo4go.Config, Logger *zap.SugaredLogger) *Bot {
	return &Bot{
		FieldMapper: FieldMapper,
		Config:      Config,
		Logger:      Logger,
		random:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

type Bot struct {
	FieldMapper mapper.Mapper
	Config      lugo4go.Config
	Logger      *zap.SugaredLogger

	random *rand.Rand
}

func (b *Bot) OnDisputing(ctx context.Context, snapshot lugo4go.SnapshotInspector) ([]proto.PlayerOrder, string, error) {
	return b.myDecider(ctx, snapshot, lugo4go.DisputingTheBall)
}

func (b *Bot) OnDefending(ctx context.Context, snapshot lugo4go.SnapshotInspector) ([]proto.PlayerOrder, string, error) {
	return b.myDecider(ctx, snapshot, lugo4go.Defending)
}

func (b *Bot) OnHolding(ctx context.Context, snapshot lugo4go.SnapshotInspector) ([]proto.PlayerOrder, string, error) {
	return b.myDecider(ctx, snapshot, lugo4go.HoldingTheBall)
}

func (b *Bot) OnSupporting(ctx context.Context, snapshot lugo4go.SnapshotInspector) ([]proto.PlayerOrder, string, error) {
	return b.myDecider(ctx, snapshot, lugo4go.Supporting)
}

func (b *Bot) AsGoalkeeper(ctx context.Context, snapshot lugo4go.SnapshotInspector, state lugo4go.PlayerState) ([]proto.PlayerOrder, string, error) {
	return b.myDecider(ctx, snapshot, state)
}

func (b *Bot) OnGetReady(ctx context.Context, snapshot lugo4go.SnapshotInspector) {
	b.Logger.Debug("the game is ready to start or the score has changed")
}

func (b *Bot) myDecider(ctx context.Context, inspector lugo4go.SnapshotInspector, state lugo4go.PlayerState) ([]proto.PlayerOrder, string, error) {
	var orders []proto.PlayerOrder
	// we are going to kick the ball as soon as we catch it
	me := inspector.GetMe()

	if state == lugo4go.HoldingTheBall {
		orderToKick, err := inspector.MakeOrderKick(b.FieldMapper.GetOpponentGoal().Center, specs.BallMaxSpeed)
		if err != nil {
			return nil, "", fmt.Errorf("could not create kick order during turn %d: %s", inspector.GetSnapshot().Turn, err)
		}
		return []proto.PlayerOrder{orderToKick}, "kicking the ball", nil
	}
	if me.Number == 10 {
		// otherwise, let's run towards the ball like kids
		orderToMove, err := inspector.MakeOrderMoveMaxSpeed(*inspector.GetBall().Position)
		if err != nil {
			return nil, "", fmt.Errorf("could not create move order during turn %d: %s", inspector.GetSnapshot().Turn, err)
		}
		return []proto.PlayerOrder{orderToMove, inspector.MakeOrderCatch()}, "I am the player 10, running", nil
	}

	orders = []proto.PlayerOrder{inspector.MakeOrderCatch()}
	debugMsg := "keeping direction"
	switch b.random.Intn(30) {
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
