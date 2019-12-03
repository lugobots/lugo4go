package bot

import (
	"context"
	"fmt"
	"github.com/lugobots/lugo4go/v2"
	"github.com/lugobots/lugo4go/v2/coach"
	"github.com/lugobots/lugo4go/v2/field"
	"github.com/lugobots/lugo4go/v2/proto"
)

var matchCtx context.Context
var log lugo4go.Logger

type Bot struct {
}

func NewBot(ctx context.Context, logger lugo4go.Logger) coach.Decider {
	log = logger
	matchCtx = ctx
	return Bot{}
}

func (b Bot) OnDisputing(data coach.TurnData) error {
	return myDecider(data)
}

func (b Bot) OnDefending(data coach.TurnData) error {
	return myDecider(data)
}

func (b Bot) OnHolding(data coach.TurnData) error {
	return myDecider(data)
}

func (b Bot) OnSupporting(data coach.TurnData) error {
	return myDecider(data)
}

func (b Bot) AsGoalkeeper(data coach.TurnData) error {
	return myDecider(data)
}

func myDecider(data coach.TurnData) error {
	var orders []proto.PlayerOrder
	// we are going to kick the ball as soon as we catch it
	if field.IsBallHolder(data.Snapshot, data.Me) {
		orderToKick, err := field.MakeOrderKick(*data.Snapshot.Ball, field.GetOpponentGoal(data.Me.TeamSide).Center, field.BallMaxSpeed)
		if err != nil {
			return fmt.Errorf("could not create kick order during turn %d: %s", data.Snapshot.Turn, err)
		}
		orders = []proto.PlayerOrder{orderToKick}
	} else if data.Me.Number == 10 {
		// otherwise, let's run towards the ball like kids
		orderToMove, err := field.MakeOrderMoveMaxSpeed(*data.Me.Position, *data.Snapshot.Ball.Position)
		if err != nil {
			return fmt.Errorf("could not create move order during turn %d: %s", data.Snapshot.Turn, err)
		}
		orders = []proto.PlayerOrder{orderToMove, field.MakeOrderCatch()}
	} else {
		orders = []proto.PlayerOrder{field.MakeOrderCatch()}
	}

	resp, err := data.Sender.Send(matchCtx, orders, "")
	if err != nil {
		return fmt.Errorf("could not send kick order during turn %d: %s", data.Snapshot.Turn, err)
	} else if resp.Code != proto.OrderResponse_SUCCESS {
		return fmt.Errorf("order sent not  order during turn %d: %s", data.Snapshot.Turn, err)
	}
	return nil
}
