package lugo4go

import (
	"context"
	"github.com/lugobots/lugo4go/v2/lugo"
	"google.golang.org/grpc"
)

type TurnData struct {
	Me       *lugo.Player
	Snapshot *lugo.GameSnapshot
}

type Bot interface {
	OnDisputing(ctx context.Context, sender OrderSender, data TurnData) error
	OnDefending(ctx context.Context, sender OrderSender, data TurnData) error
	OnHolding(ctx context.Context, sender OrderSender, data TurnData) error
	OnSupporting(ctx context.Context, sender OrderSender, data TurnData) error
	AsGoalkeeper(ctx context.Context, sender OrderSender, data TurnData) error
}

type OrderSender interface {
	Send(ctx context.Context, orders []lugo.PlayerOrder, debugMsg string) (*lugo.OrderResponse, error)
}

type DecisionMaker func(ctx context.Context, snapshot *lugo.GameSnapshot, sender OrderSender)

type Logger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

type GRPCHandler interface {
	OnNewTurn(DecisionMaker, Logger)
	Stop() error
	GetGRPCConn() *grpc.ClientConn
	GetServiceConn() lugo.GameClient
	// The sender will not need the entire snapshot struct. However there are plans to allow the sender
	// to do mre complex jobs (e.g. having middleware to save status for machine learning). Then, we are
	// passing the snapshot since now, so the new versions will be compatible.
	//	SenderBuilder(builder func(snapshot *lugo.GameSnapshot, logger Logger) OrderSender)
}
