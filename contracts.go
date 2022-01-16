package lugo4go

import (
	"context"
	"github.com/lugobots/lugo4go/v2/proto"
)

type TurnData struct {
	Me       *proto.Player
	Snapshot *proto.GameSnapshot
}

type TurnHandler interface {
	Handle(ctx context.Context, snapshot *proto.GameSnapshot)
}

type OrderSender interface {
	Send(ctx context.Context, turn uint32, orders []proto.PlayerOrder, debugMsg string) (*proto.OrderResponse, error)
}

type TurnOrdersSender interface {
	Send(ctx context.Context, orders []proto.PlayerOrder, debugMsg string) (*proto.OrderResponse, error)
}

type Bot interface {
	OnDisputing(ctx context.Context, sender TurnOrdersSender, snapshot *proto.GameSnapshot) error
	OnDefending(ctx context.Context, sender TurnOrdersSender, snapshot *proto.GameSnapshot) error
	OnHolding(ctx context.Context, sender TurnOrdersSender, snapshot *proto.GameSnapshot) error
	OnSupporting(ctx context.Context, sender TurnOrdersSender, snapshot *proto.GameSnapshot) error
	AsGoalkeeper(ctx context.Context, sender TurnOrdersSender, snapshot *proto.GameSnapshot, state PlayerState) error
}

type Logger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}
