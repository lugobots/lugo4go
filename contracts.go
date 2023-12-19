package lugo4go

import (
	"context"

	"github.com/lugobots/lugo4go/v3/mapper"
	"github.com/lugobots/lugo4go/v3/proto"
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

type orderSender interface {
	Send(ctx context.Context, turn uint32, orders []proto.PlayerOrder, debugMsg string) (*proto.OrderResponse, error)
}

type TurnOrdersSender interface {
	Send(ctx context.Context, orders []proto.PlayerOrder, debugMsg string) (*proto.OrderResponse, error)
}

type SnapshotInspector interface {
	GetSnapshot() *proto.GameSnapshot
	GetMe() *proto.Player

	GetBall() *proto.Ball
	GetBallHolder() (*proto.Player, bool)
	IsBallHolder(player *proto.Player) bool

	GetTeam(side proto.Team_Side) *proto.Team
	GetMyTeam() *proto.Team
	GetOpponentMyTeam() *proto.Team

	GetMyTeamSide() proto.Team_Side
	GetOpponentSide() proto.Team_Side

	GetPlayer(side proto.Team_Side, number int) *proto.Player
	GetMyTeamPlayers() []*proto.Player
	GetOpponentPlayers() []*proto.Player

	MakeOrderMoveMaxSpeed(target proto.Point) (*proto.Order_Move, error)
	MakeOrderMoveFromPoint(origin, target proto.Point, speed float64) (*proto.Order_Move, error)
	MakeOrderMoveFromVector(vector proto.Vector, speed float64) (*proto.Order_Move, error)
	MakeOrderMoveByDirection(direction mapper.Direction, speed float64) (*proto.Order_Move, error)

	MakeOrderJump(target proto.Point, speed float64) (*proto.Order_Jump, error)

	MakeOrderKick(target proto.Point, speed float64) (*proto.Order_Kick, error)

	MakeOrderKickMaxSpeed(target proto.Point) (*proto.Order_Kick, error)

	MakeOrderCatch() *proto.Order_Catch
}

type Bot interface {
	// OnDisputing is called when no one has the ball possession
	OnDisputing(ctx context.Context, snapshot SnapshotInspector) ([]proto.PlayerOrder, string, error)
	// OnDefending is called when an opponent player has the ball possession
	OnDefending(ctx context.Context, snapshot SnapshotInspector) ([]proto.PlayerOrder, string, error)
	// OnHolding is called when this bot has the ball possession
	OnHolding(ctx context.Context, snapshot SnapshotInspector) ([]proto.PlayerOrder, string, error)
	// OnSupporting is called when a teammate player has the ball possession
	OnSupporting(ctx context.Context, snapshot SnapshotInspector) ([]proto.PlayerOrder, string, error)
	// AsGoalkeeper is only called when this bot is the goalkeeper (number 1). This method is called on every turn,
	// and the player state is passed at the last parameter.
	AsGoalkeeper(ctx context.Context, snapshot SnapshotInspector, state PlayerState) ([]proto.PlayerOrder, string, error)
}

type Logger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}
