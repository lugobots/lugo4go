package ops

import "github.com/makeitplay/client-player-go/lugo"

type OrderSender func(in *lugo.OrderSet) (*lugo.OrderResponse, error)
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
}
