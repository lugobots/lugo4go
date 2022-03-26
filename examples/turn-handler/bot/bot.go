package bot

import (
	"context"
	"github.com/lugobots/lugo4go/v2"
	"github.com/lugobots/lugo4go/v2/pkg/field"
	"github.com/lugobots/lugo4go/v2/proto"
)

type Bot struct {
	OrderSender lugo4go.OrderSender
	Side        proto.Team_Side
	Number      uint32
	Logger      lugo4go.Logger
	arr         field.Mapper
}

func NewBot(orderSender lugo4go.OrderSender, logger lugo4go.Logger, side proto.Team_Side, number uint32) *Bot {
	arr, _ := field.NewMapper(field.MaxCols, field.MaxRows, side)
	return &Bot{
		OrderSender: orderSender,
		Logger:      logger,
		Number:      number,
		Side:        side,
		arr:         arr,
	}
}

func (b *Bot) Handle(ctx context.Context, snapshot *proto.GameSnapshot) {
	var orders []proto.PlayerOrder
	// we are going to kick the ball as soon as we catch it
	me := field.GetPlayer(snapshot, b.Side, b.Number)
	if me == nil {
		b.Logger.Errorf("could not find myself in the team")
		return
	}

	if field.IsBallHolder(snapshot, me) {
		orderToKick, err := field.MakeOrderKick(*snapshot.Ball, field.GetOpponentGoal(me.TeamSide).Center, field.BallMaxSpeed)
		if err != nil {
			b.Logger.Errorf("could not create kick order during turn %d: %s", snapshot.Turn, err)
			return
		}
		orders = []proto.PlayerOrder{orderToKick}
	} else if me.Number == 10 {
		// otherwise, let's run towards the ball like kids
		orderToMove, err := field.MakeOrderMoveMaxSpeed(*me.Position, *snapshot.Ball.Position)
		if err != nil {
			b.Logger.Errorf("could not create move order during turn %d: %s", snapshot.Turn, err)
			return
		}
		orders = []proto.PlayerOrder{orderToMove, field.MakeOrderCatch()}
	} else {
		orders = []proto.PlayerOrder{field.MakeOrderCatch()}
	}

	resp, err := b.OrderSender.Send(ctx, snapshot.Turn, orders, "")
	if err != nil {
		b.Logger.Errorf("could not send kick order during turn %d: %s", snapshot.Turn, err)
	} else if resp.Code != proto.OrderResponse_SUCCESS {
		b.Logger.Errorf("order sent not  order during turn %d: %s", snapshot.Turn, err)
	}
	return
}

var FieldMap = map[uint32]struct {
	Col uint8
	Row uint8
}{
	2:  {Col: 1, Row: 1},
	3:  {Col: 1, Row: 3},
	4:  {Col: 1, Row: 4},
	5:  {Col: 1, Row: 6},
	6:  {Col: 2, Row: 2},
	7:  {Col: 2, Row: 3},
	8:  {Col: 2, Row: 4},
	9:  {Col: 2, Row: 5},
	10: {Col: 3, Row: 3},
	11: {Col: 3, Row: 4},
}
