package lugo4go

import (
	"context"
	"fmt"
	"github.com/lugobots/lugo4go/v2/coach"
	"github.com/lugobots/lugo4go/v2/lugo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
	"io"
)

const ProtocolVersion = "2.0"

func NewClient(config lugo.Config) (*Client, error) {
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
	Logger     lugo.Logger
	config     lugo.Config
}

func (c Client) PlayWithBot(bot coach.Bot, logger lugo.Logger) error {
	sender := coach.NewSender(c.GRPCClient)
	handler := coach.NewHandler(bot, sender, logger, c.config.Number, c.config.TeamSide)
	return c.Play(handler)
}
func (c Client) Play(handler TurnHandler) error {
	var turnCrx context.Context
	var stop context.CancelFunc = func() {}
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
			close(mustHasStarted)
			handler.Handle(turnCrx, snapshot)
		}()
		<-mustHasStarted
	}
}

func (c Client) Stop() error {
	return c.grpcConn.Close()
}

//
//func (c Client) GetGRPCConn() *grpc.ClientConn {
//	return c.grpcConn
//}

//func (c Client) GetServiceConn() lugo.GameClient {
//	return c.GRPCClient
//}

//func (c Client) SenderBuilder(builder func(snapshot *lugo.GameSnapshot, logger Logger) OrderSender) {
//	c.Sender = builder
//}

//type sender struct {
//	snapshot *lugo.GameSnapshot
//	logger   Logger
//	gameConn lugo.GameClient
//}
//
//func (s sender) Send(ctx context.Context, orders []lugo.PlayerOrder, debugMsg string) (*lugo.OrderResponse, error) {
//	orderSet := &lugo.OrderSet{
//		Turn:         s.snapshot.Turn,
//		DebugMessage: debugMsg,
//		Orders:       []*lugo.Order{},
//	}
//	for _, order := range orders {
//		orderSet.Orders = append(orderSet.Orders, &lugo.Order{Action: order})
//	}
//	s.logger.Debugf("sending orders for turn %d", s.snapshot.Turn)
//	return s.gameConn.SendOrders(ctx, orderSet)
//}

func (c *Client) TagRPC(ctx context.Context, t *stats.RPCTagInfo) context.Context {
	return ctx
}

func (c *Client) HandleRPC(context.Context, stats.RPCStats) {

}

func (c *Client) TagConn(ctx context.Context, t *stats.ConnTagInfo) context.Context {
	return ctx
}

func (c *Client) HandleConn(ctx context.Context, sts stats.ConnStats) {
	switch sts.(type) {
	case *stats.ConnEnd:
		_ = c.Stop()
		break
	}
}
