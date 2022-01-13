package lugo4go

import (
	"context"
	"github.com/lugobots/lugo4go/v2/lugo"
)

type TurnData struct {
	Me       *lugo.Player
	Snapshot *lugo.GameSnapshot
}

type TurnHandler interface {
	Handle(ctx context.Context, snapshot *lugo.GameSnapshot)
}

type OrderSender interface {
	Send(ctx context.Context, turn uint32, orders []lugo.PlayerOrder, debugMsg string) (*lugo.OrderResponse, error)
}

type TurnOrdersSender interface {
	Send(ctx context.Context, orders []lugo.PlayerOrder, debugMsg string) (*lugo.OrderResponse, error)
}

type Bot interface {
	OnDisputing(ctx context.Context, sender TurnOrdersSender, snapshot *lugo.GameSnapshot) error
	OnDefending(ctx context.Context, sender TurnOrdersSender, snapshot *lugo.GameSnapshot) error
	OnHolding(ctx context.Context, sender TurnOrdersSender, snapshot *lugo.GameSnapshot) error
	OnSupporting(ctx context.Context, sender TurnOrdersSender, snapshot *lugo.GameSnapshot) error
	AsGoalkeeper(ctx context.Context, sender TurnOrdersSender, snapshot *lugo.GameSnapshot, state PlayerState) error
}
