package lugo4go

import (
	"context"
	"fmt"
	"github.com/lugobots/lugo4go/v2/coach"
	"github.com/lugobots/lugo4go/v2/lugo"
	"github.com/lugobots/lugo4go/v2/pkg/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
	"io"
	"sync"
)

const ProtocolVersion = "2.0"

func NewClient(config util.Config) (*Client, error) {
	var err error
	c := &Client{
		config: config,
	}

	// A bot may eventually do not listen to server Stream (ignoring OnNewTurn). In this case, the Client must stop
	// when the gRPC connection is closed.
	connHandler := grpc.WithStatsHandler(c)
	if config.Insecure {
		c.grpcConn, err = grpc.Dial(config.GRPCAddress, grpc.WithInsecure(), connHandler)
	} else {
		c.grpcConn, err = grpc.Dial(config.GRPCAddress, connHandler)
	}
	if err != nil {
		return nil, err
	}

	c.GRPCClient = lugo.NewGameClient(c.grpcConn)

	if c.Stream, err = c.GRPCClient.JoinATeam(context.Background(), &lugo.JoinRequest{
		Token:           config.Token,
		Number:          config.Number,
		InitPosition:    &config.InitialPosition,
		TeamSide:        config.TeamSide,
		ProtocolVersion: ProtocolVersion,
	}); err != nil {
		return nil, err
	}
	return c, nil
}

type Client struct {
	Stream     lugo.Game_JoinATeamClient
	GRPCClient lugo.GameClient
	grpcConn   *grpc.ClientConn
	Handler    TurnHandler
	Logger     util.Logger
	config     util.Config
}

func (c *Client) PlayWithBot(bot coach.Bot, logger util.Logger) error {
	sender := coach.NewSender(c.GRPCClient)
	handler := coach.NewHandler(bot, sender, logger, c.config.Number, c.config.TeamSide)
	return c.Play(handler)
}

func (c *Client) Play(handler TurnHandler) error {
	var turnCrx context.Context
	var stop context.CancelFunc = func() {}
	m := sync.Mutex{}
	for {
		snapshot, err := c.Stream.Recv()
		stop()
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("%w: %s", ErrGRPCConnectionLost, err)
			}
			return ErrGRPCConnectionClosed
		}
		turnCrx, stop = context.WithCancel(context.Background())
		mustHasStarted := make(chan bool)
		go func() {
			m.Lock()
			close(mustHasStarted)
			defer m.Unlock()
			handler.Handle(turnCrx, snapshot)
		}()
		<-mustHasStarted
	}
}

func (c *Client) Stop() error {
	return c.grpcConn.Close()
}

func (c *Client) TagRPC(ctx context.Context, _ *stats.RPCTagInfo) context.Context {
	return ctx
}

func (c *Client) HandleRPC(context.Context, stats.RPCStats) {

}

func (c *Client) TagConn(ctx context.Context, _ *stats.ConnTagInfo) context.Context {
	return ctx
}

func (c *Client) HandleConn(_ context.Context, sts stats.ConnStats) {
	switch sts.(type) {
	case *stats.ConnEnd:
		_ = c.Stop()
		break
	}
}
