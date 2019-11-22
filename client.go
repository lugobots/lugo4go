package client

import (
	"context"
	"github.com/makeitplay/client-player-go/lugo"
	"google.golang.org/grpc"
	"io"
)

const ProtocolVersion = "2.0"

type DecisionMaker func(snapshot *lugo.GameSnapshot) (orders []lugo.Order, debugMsg string)

type Logger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

type Client interface {
	OnNewTurn(func(snapshot *lugo.GameSnapshot) DecisionMaker)
	Stop()
}

type Config struct {
	// Full url to the gRPC server
	GRPCHost        string
	Insecure        bool
	TeamSide        lugo.Team_Side
	Number          uint32
	InitialPosition *lugo.Point
}

func NewClient(ctx context.Context, config Config) (Client, io.Closer, error) {
	var err error
	var conn *grpc.ClientConn
	if config.Insecure {
		conn, err = grpc.Dial(config.GRPCHost, grpc.WithInsecure())
	} else {
		conn, err = grpc.Dial(config.GRPCHost)
	}
	if err != nil {
		return nil, nil, err
	}

	c := &client{}
	if c.stream, err = lugo.NewFootballClient(conn).JoinATeam(ctx, &lugo.JoinRequest{
		Number:          config.Number,
		InitPosition:    config.InitialPosition,
		TeamSide:        config.TeamSide,
		ProtocolVersion: ProtocolVersion,
	}); err != nil {
		return nil, nil, err
	}
	return c, conn, nil
}

type client struct {
	stream lugo.Football_JoinATeamClient
}

func (c client) OnNewTurn(func(snapshot *lugo.GameSnapshot) DecisionMaker) {
	panic("implement me")
}

func (c client) Stop() {
	return
}
