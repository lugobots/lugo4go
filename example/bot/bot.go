package bot

import (
	"context"
	"errors"
	"fmt"
	"github.com/lugobots/lugo4go/v2/coach"
	"github.com/lugobots/lugo4go/v2/lugo"
	"github.com/lugobots/lugo4go/v2/pkg/field"
	"github.com/lugobots/lugo4go/v2/pkg/util"
)

type Bot struct {
	Side         lugo.Team_Side
	Number       uint32
	Logger       util.Logger
	errorsStream chan error
}

func NewBot(logger util.Logger, side lugo.Team_Side, number uint32) (*Bot, <-chan error) {
	err := make(chan error)
	return &Bot{
		Logger:       logger,
		Number:       number,
		Side:         side,
		errorsStream: err,
	}, err
}

func (b *Bot) OnDisputing(ctx context.Context, sender coach.OrderSender, snapshot *lugo.GameSnapshot) error {
	return b.myDecider(ctx, sender, snapshot, coach.DisputingTheBall)
}

func (b *Bot) OnDefending(ctx context.Context, sender coach.OrderSender, snapshot *lugo.GameSnapshot) error {
	return b.myDecider(ctx, sender, snapshot, coach.Defending)
}

func (b *Bot) OnHolding(ctx context.Context, sender coach.OrderSender, snapshot *lugo.GameSnapshot) error {
	return b.myDecider(ctx, sender, snapshot, coach.HoldingTheBall)
}

func (b *Bot) OnSupporting(ctx context.Context, sender coach.OrderSender, snapshot *lugo.GameSnapshot) error {
	return b.myDecider(ctx, sender, snapshot, coach.Supporting)
}

func (b *Bot) AsGoalkeeper(ctx context.Context, sender coach.OrderSender, snapshot *lugo.GameSnapshot, state coach.PlayerState) error {
	return b.myDecider(ctx, sender, snapshot, state)
}

func (b *Bot) myDecider(ctx context.Context, sender coach.OrderSender, snapshot *lugo.GameSnapshot, state coach.PlayerState) error {
	var orders []lugo.PlayerOrder
	// we are going to kick the ball as soon as we catch it
	me := field.GetPlayer(snapshot, b.Side, b.Number)
	if me == nil {
		return errorHandler(b.errorsStream, errors.New("bot not found in the game snapshot"))
	}
	if state == coach.HoldingTheBall {
		orderToKick, err := field.MakeOrderKick(*snapshot.Ball, field.GetOpponentGoal(me.TeamSide).Center, field.BallMaxSpeed)
		if err != nil {
			return errorHandler(b.errorsStream, fmt.Errorf("could not create kick order during turn %d: %s", snapshot.Turn, err))
		}
		orders = []lugo.PlayerOrder{orderToKick}
	} else if me.Number == 10 {
		// otherwise, let's run towards the ball like kids
		orderToMove, err := field.MakeOrderMoveMaxSpeed(*me.Position, *snapshot.Ball.Position)
		if err != nil {
			return errorHandler(b.errorsStream, fmt.Errorf("could not create move order during turn %d: %s", snapshot.Turn, err))
		}
		orders = []lugo.PlayerOrder{orderToMove, field.MakeOrderCatch()}
	} else {
		orders = []lugo.PlayerOrder{field.MakeOrderCatch()}
	}

	resp, err := sender.Send(ctx, snapshot.Turn, orders, "")
	if err != nil {
		return errorHandler(b.errorsStream, fmt.Errorf("could not send kick order during turn %d: %s", snapshot.Turn, err))
	} else if resp.Code != lugo.OrderResponse_SUCCESS {
		return errorHandler(b.errorsStream, fmt.Errorf("order sent not  order during turn %d: %s", snapshot.Turn, err))
	}
	return nil
}

func errorHandler(stream chan error, err error) error {
	select {
	case stream <- err:
	default:
	}
	return err
}
