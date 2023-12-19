package bot

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/lugobots/lugo4go/v2/pkg/field"

	"github.com/lugobots/lugo4go/v2"
	"github.com/lugobots/lugo4go/v2/mapper"
	"github.com/lugobots/lugo4go/v2/proto"
	"github.com/lugobots/lugo4go/v2/specs"
)

type Bot struct {
	Side   proto.Team_Side
	Number uint32
	Logger lugo4go.Logger
	arr    mapper.Mapper
}

func NewBot(logger lugo4go.Logger, side proto.Team_Side, number uint32) *Bot {
	arr, _ := mapper.NewMapper(mapper.MaxCols, mapper.MaxRows, side)
	rand.Seed(time.Now().UnixNano() * int64(number))

	return &Bot{
		Logger: logger,
		Number: number,
		Side:   side,
		arr:    arr,
	}
}

func (b *Bot) OnDisputing(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *proto.GameSnapshot) error {
	return b.myDecider(ctx, sender, snapshot, lugo4go.DisputingTheBall)
}

func (b *Bot) OnDefending(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *proto.GameSnapshot) error {
	return b.myDecider(ctx, sender, snapshot, lugo4go.Defending)
}

func (b *Bot) OnHolding(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *proto.GameSnapshot) error {
	return b.myDecider(ctx, sender, snapshot, lugo4go.HoldingTheBall)
}

func (b *Bot) OnSupporting(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *proto.GameSnapshot) error {
	return b.myDecider(ctx, sender, snapshot, lugo4go.Supporting)
}

func (b *Bot) AsGoalkeeper(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *proto.GameSnapshot, state lugo4go.PlayerState) error {
	return b.myDecider(ctx, sender, snapshot, state)
}

func (b *Bot) myDecider(ctx context.Context, sender lugo4go.TurnOrdersSender, snapshot *proto.GameSnapshot, state lugo4go.PlayerState) error {
	var orders []proto.PlayerOrder
	// we are going to kick the ball as soon as we catch it
	me := field.GetPlayer(snapshot, b.Side, b.Number)
	if me == nil {
		return errorHandler(b.Logger, errors.New("bot not found in the game snapshot"))
	}
	if state == lugo4go.HoldingTheBall {
		orderToKick, err := field.MakeOrderKick(*snapshot.Ball, field.GetOpponentGoal(me.TeamSide).Center, specs.BallMaxSpeed)
		if err != nil {
			return errorHandler(b.Logger, fmt.Errorf("could not create kick order during turn %d: %s", snapshot.Turn, err))
		}
		orders = []proto.PlayerOrder{orderToKick}
	} else if me.Number == 10 {
		// otherwise, let's run towards the ball like kids
		orderToMove, err := field.MakeOrderMoveMaxSpeed(*me.Position, *snapshot.Ball.Position)
		if err != nil {
			return errorHandler(b.Logger, fmt.Errorf("could not create move order during turn %d: %s", snapshot.Turn, err))
		}
		orders = []proto.PlayerOrder{orderToMove, field.MakeOrderCatch()}
	} else {
		orders = []proto.PlayerOrder{field.MakeOrderCatch()}
		orders = []proto.PlayerOrder{field.MakeOrderCatch()}
		switch rand.Intn(30) {
		case 0:
			orders = append(orders, field.GoRight(b.Side))
		case 1:
			orders = append(orders, field.GoLeft(b.Side))
		case 2:
			orders = append(orders, field.GoForward(b.Side))
		case 3:
			orders = append(orders, field.GoBackward(b.Side))
		}
	}

	resp, err := sender.Send(ctx, orders, "")
	if err != nil {
		return errorHandler(b.Logger, fmt.Errorf("could not send kick order during turn %d: %s", snapshot.Turn, err))
	} else if resp.Code != proto.OrderResponse_SUCCESS {
		return errorHandler(b.Logger, fmt.Errorf("order sent not  order during turn %d: %s", snapshot.Turn, err))
	}
	return nil
}

func errorHandler(logger lugo4go.Logger, err error) error {
	logger.Errorf("bot error: %s", err)
	return err
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
