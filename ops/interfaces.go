package ops

import (
	"context"
	"github.com/makeitplay/client-player-go/lugo"
	"google.golang.org/grpc"
)

type DecisionMaker func(snapshot *lugo.GameSnapshot, sender OrderSender)

type Logger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

type Client interface {
	OnNewTurn(DecisionMaker, Logger)
	Stop() error
	GetGRPCConn() *grpc.ClientConn
	GetServiceConn() lugo.GameClient
	SenderBuilder(builder func(snapshot *lugo.GameSnapshot, logger Logger) OrderSender)
}

type OrderSender interface {
	Send(ctx context.Context, orders []lugo.PlayerOrder, debugMsg string) (*lugo.OrderResponse, error)
}
