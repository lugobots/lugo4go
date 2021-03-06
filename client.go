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

// ProtocolVersion defines the current game protocol
const ProtocolVersion = "2.0"

// NewClient creates a Lugo4Go client that will hide common logic and let you focus on your bot.
func NewClient(config util.Config) (*Client, error) {
	var err error
	c := &Client{
		config: config,
	}

	// A bot may eventually do not listen to server Stream (ignoring OnNewTurn). In this case, the Client must stop
	// when the gRPC connection is closed.
	connHandler := grpc.WithStatsHandler(c)
	// @todo there are some gRPC options that we should take a look tro improve this part.
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

// Client handle the gRPC stuff and provide you a easy way to handle the game messages
type Client struct {
	Stream     lugo.Game_JoinATeamClient
	GRPCClient lugo.GameClient
	grpcConn   *grpc.ClientConn
	Handler    TurnHandler
	Logger     util.Logger
	config     util.Config
}

// PlayWithBot is a sugared Play mode that uses an TurnHandler from coach package.
// Coach TurnHandler creates basic player states to help the development of new bots.
func (c *Client) PlayWithBot(bot coach.Bot, logger util.Logger) error {
	sender := coach.NewSender(c.GRPCClient)
	handler := coach.NewHandler(bot, sender, logger, c.config.Number, c.config.TeamSide)
	return c.Play(handler)
}

// Play starts the player communication with the server. The TurnHandler will receive the raw snapshot from the
// game server. The context passed to the handler will be canceled as soon a new turn starts.
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
		// to avoid race conditions we need to ensure that the loop can only start after the Go routine has started.
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

// Stop drops the communication with the gRPC server.
func (c *Client) Stop() error {
	return c.grpcConn.Close()
}

// TagRPC implements the interface required by gRPC handler
func (c *Client) TagRPC(ctx context.Context, _ *stats.RPCTagInfo) context.Context {
	return ctx
}

// HandleRPC implements the interface required by gRPC handler
func (c *Client) HandleRPC(context.Context, stats.RPCStats) {

}

// TagConn implements the interface required by gRPC handler
func (c *Client) TagConn(ctx context.Context, _ *stats.ConnTagInfo) context.Context {
	return ctx
}

// HandleConn implements the interface required by gRPC handler
func (c *Client) HandleConn(_ context.Context, sts stats.ConnStats) {
	switch sts.(type) {
	case *stats.ConnEnd:
		_ = c.Stop()
		break
	}
}
