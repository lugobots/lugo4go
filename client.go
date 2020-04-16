package lugo4go

import (
	"context"
	"fmt"
	"github.com/lugobots/lugo4go/v2/lugo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
	"io"
)

const ProtocolVersion = "2.0"

func NewRawClient(config Config, handler TurnHandler) (*Client, error) {
	var err error
	c := &Client{
		Handler: handler,
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

	//c.Sender = func(snapshot *lugo.GameSnapshot, logger Logger) OrderSender {
	//	return &sender{
	//		gameConn: c.GRPCClient,
	//		snapshot: snapshot,
	//		logger:   logger,
	//	}
	//}

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
	Logger     Logger
	//Sender     OrderSender
	//ErrorHandler func(response *lugo.OrderResponse, err error)
}

func (c Client) Play() error {
	var turnCrx context.Context
	var stop context.CancelFunc = func() {}
	for {
		snapshot, err := c.Stream.Recv()
		stop()
		if err != nil {
			if err != io.EOF {
				//err = fmt.Errorf("connection error: %w", err)
				return fmt.Errorf("%w: %s", ErrGRPCConnectionLost, err)
			}
			return ErrGRPCConnectionClosed
		}
		turnCrx, stop = context.WithCancel(context.Background())
		mustHasStarted := make(chan bool)
		go func() {
			close(mustHasStarted)
			c.Handler.Handle(turnCrx, snapshot, c.GRPCClient)
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
