package lugo4go

import (
	"context"
	"github.com/lugobots/lugo4go/v2/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
	"io"
)

const ProtocolVersion = "2.0"

func NewClient(config Config) (context.Context, Client, error) {
	var err error
	c := &client{}

	// A bot may eventually do not listen to server stream (ignoring OnNewTurn). In this case, the client must stop
	// when the gRPC connection is closed.
	connHandler := grpc.WithStatsHandler(c)
	if config.Insecure {
		c.grpcConn, err = grpc.Dial(config.GRPCAddress, grpc.WithInsecure(), connHandler)
	} else {
		c.grpcConn, err = grpc.Dial(config.GRPCAddress, connHandler)
	}
	if err != nil {
		return nil, nil, err
	}

	c.gameConn = proto.NewGameClient(c.grpcConn)

	c.senderBuilder = func(snapshot *proto.GameSnapshot, logger Logger) OrderSender {
		return &sender{
			gameConn: c.gameConn,
			snapshot: snapshot,
			logger:   logger,
		}
	}

	c.ctx, c.stopCtx = context.WithCancel(context.Background())
	if c.stream, err = c.gameConn.JoinATeam(c.ctx, &proto.JoinRequest{
		Token:           config.Token,
		Number:          config.Number,
		InitPosition:    &config.InitialPosition,
		TeamSide:        config.TeamSide,
		ProtocolVersion: ProtocolVersion,
	}); err != nil {
		return nil, nil, err
	}
	return c.ctx, c, nil
}

type client struct {
	stream        proto.Game_JoinATeamClient
	gameConn      proto.GameClient
	grpcConn      *grpc.ClientConn
	ctx           context.Context
	stopCtx       context.CancelFunc
	senderBuilder func(snapshot *proto.GameSnapshot, logger Logger) OrderSender
	sender        OrderSender
}

func (c client) OnNewTurn(decider DecisionMaker, log Logger) {
	var turnCrx context.Context
	var stop context.CancelFunc = func() {}
	go func() {
		for {
			snapshot, err := c.stream.Recv()
			stop()
			turnCrx, stop = context.WithCancel(c.ctx)
			if err != nil {
				if err == io.EOF {
					log.Infof("gRPC connection closed")
				} else {
					log.Errorf("gRPC stream error: %s", err)
				}
				c.stopCtx()
				return
			}
			log.Debugf("calling DecisionMaker for turn %d", snapshot.Turn)
			go decider(turnCrx, snapshot, c.senderBuilder(snapshot, log))
		}
	}()
}

func (c client) Stop() error {
	c.stopCtx()
	return c.grpcConn.Close()
}

func (c client) GetGRPCConn() *grpc.ClientConn {
	return c.grpcConn
}

func (c client) GetServiceConn() proto.GameClient {
	return c.gameConn
}

func (c client) SenderBuilder(builder func(snapshot *proto.GameSnapshot, logger Logger) OrderSender) {
	c.senderBuilder = builder
}

type sender struct {
	snapshot *proto.GameSnapshot
	logger   Logger
	gameConn proto.GameClient
}

func (s sender) Send(ctx context.Context, orders []proto.PlayerOrder, debugMsg string) (*proto.OrderResponse, error) {
	orderSet := &proto.OrderSet{
		Turn:         s.snapshot.Turn,
		DebugMessage: debugMsg,
		Orders:       []*proto.Order{},
	}
	for _, order := range orders {
		orderSet.Orders = append(orderSet.Orders, &proto.Order{Action: order})
	}
	s.logger.Debugf("sending orders for turn %d", s.snapshot.Turn)
	return s.gameConn.SendOrders(ctx, orderSet)
}

func (c *client) TagRPC(ctx context.Context, t *stats.RPCTagInfo) context.Context {
	return ctx
}

func (c *client) HandleRPC(context.Context, stats.RPCStats) {

}

func (c *client) TagConn(ctx context.Context, t *stats.ConnTagInfo) context.Context {
	return ctx
}

func (c *client) HandleConn(ctx context.Context, sts stats.ConnStats) {
	switch sts.(type) {
	case *stats.ConnEnd:
		_ = c.Stop()
		break
	}
}
