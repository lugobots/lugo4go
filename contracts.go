package lugo4go

import (
	"context"
	"github.com/lugobots/lugo4go/v2/proto"
)

type TurnData struct {
	Me       *proto.Player
	Snapshot *proto.GameSnapshot
}

// TurnHandler is required by the Lugo4Go client to handle each turn snapshot
type TurnHandler interface {

	// Handle is called every turn with the new game state
	Handle(ctx context.Context, snapshot *proto.GameSnapshot)
}

type OrderSender interface {
	Send(ctx context.Context, turn uint32, orders []proto.PlayerOrder, debugMsg string) (*proto.OrderResponse, error)
}

type TurnOrdersSender interface {
	Send(ctx context.Context, orders []proto.PlayerOrder, debugMsg string) (*proto.OrderResponse, error)
}

type Bot interface {
	// OnDisputing is called when no one has the ball possession
	OnDisputing(ctx context.Context, sender TurnOrdersSender, snapshot *proto.GameSnapshot) error
	// OnDefending is called when an opponent player has the ball possession
	OnDefending(ctx context.Context, sender TurnOrdersSender, snapshot *proto.GameSnapshot) error
	// OnHolding is called when this bot has the ball possession
	OnHolding(ctx context.Context, sender TurnOrdersSender, snapshot *proto.GameSnapshot) error
	// OnSupporting is called when a teammate player has the ball possession
	OnSupporting(ctx context.Context, sender TurnOrdersSender, snapshot *proto.GameSnapshot) error
	// AsGoalkeeper is only called when this bot is the goalkeeper (number 1). This method is called on every turn,
	// and the player state is passed at the last parameter.
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
