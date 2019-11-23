package client

import (
	"context"
	"github.com/makeitplay/client-player-go/lugo"
	"github.com/makeitplay/client-player-go/ops"
	"google.golang.org/grpc"
	"io"
)

const ProtocolVersion = "2.0"

type Config struct {
	// Full url to the gRPC server
	GRPCHost        string
	Insecure        bool
	TeamSide        lugo.Team_Side
	Number          uint32
	InitialPosition *lugo.Point
}

func NewClient(config Config) (context.Context, ops.Client, error) {
	var err error
	c := &client{}

	if config.Insecure {
		c.grpcConn, err = grpc.Dial(config.GRPCHost, grpc.WithInsecure())
	} else {
		c.grpcConn, err = grpc.Dial(config.GRPCHost)
	}
	if err != nil {
		return nil, nil, err
	}

	c.gameConn = lugo.NewGameClient(c.grpcConn)
	c.ctx, c.stopCtx = context.WithCancel(context.Background())
	if c.stream, err = c.gameConn.JoinATeam(c.ctx, &lugo.JoinRequest{
		Number:          config.Number,
		InitPosition:    config.InitialPosition,
		TeamSide:        config.TeamSide,
		ProtocolVersion: ProtocolVersion,
	}); err != nil {
		return nil, nil, err
	}
	return c.ctx, c, nil
}

type client struct {
	stream   lugo.Game_JoinATeamClient
	gameConn lugo.GameClient
	grpcConn *grpc.ClientConn
	ctx      context.Context
	stopCtx  context.CancelFunc
}

func (c client) OnNewTurn(decider ops.DecisionMaker, log ops.Logger) {
	go func() {
		for {
			snapshot, err := c.stream.Recv()
			if err != nil {
				if err == io.EOF {
					log.Infof("gRPC connection closed")
				} else {
					log.Errorf("gRPC stream error: %s", err)
				}
				return
			}
			log.Debugf("calling DecisionMaker for turn %d", snapshot.Turn)
			decider(snapshot, func(orderSet *lugo.OrderSet) (response *lugo.OrderResponse, e error) {
				log.Debugf("sending orders for turn %d", snapshot.Turn)
				return c.gameConn.SendOrders(c.ctx, orderSet)
			})
		}
	}()
}

func (c client) Stop() error {
	c.stopCtx()
	return c.grpcConn.Close()
}
