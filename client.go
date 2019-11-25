package client

import (
	"context"
	"github.com/makeitplay/client-player-go/lugo"
	"github.com/makeitplay/client-player-go/ops"
	"google.golang.org/grpc"
	"io"
)

const ProtocolVersion = "2.0"

func NewClient(config Config) (context.Context, ops.Client, error) {
	var err error
	c := &client{}

	if config.Insecure {
		c.grpcConn, err = grpc.Dial(config.GRPCAddress, grpc.WithInsecure())
	} else {
		c.grpcConn, err = grpc.Dial(config.GRPCAddress)
	}
	if err != nil {
		return nil, nil, err
	}

	c.gameConn = lugo.NewGameClient(c.grpcConn)

	c.senderBuilder = func(snapshot *lugo.GameSnapshot, logger ops.Logger) ops.OrderSender {
		return &sender{
			gameConn: c.gameConn,
			snapshot: snapshot,
			logger:   logger,
		}
	}

	c.ctx, c.stopCtx = context.WithCancel(context.Background())
	if c.stream, err = c.gameConn.JoinATeam(c.ctx, &lugo.JoinRequest{
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
	stream        lugo.Game_JoinATeamClient
	gameConn      lugo.GameClient
	grpcConn      *grpc.ClientConn
	ctx           context.Context
	stopCtx       context.CancelFunc
	senderBuilder func(snapshot *lugo.GameSnapshot, logger ops.Logger) ops.OrderSender
	sender        ops.OrderSender
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
				c.stopCtx()
				return
			}
			log.Debugf("calling DecisionMaker for turn %d", snapshot.Turn)
			decider(snapshot, c.senderBuilder(snapshot, log))
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

func (c client) GetServiceConn() lugo.GameClient {
	return c.gameConn
}

func (c client) SenderBuilder(builder func(snapshot *lugo.GameSnapshot, logger ops.Logger) ops.OrderSender) {
	c.senderBuilder = builder
}

type sender struct {
	snapshot *lugo.GameSnapshot
	logger   ops.Logger
	gameConn lugo.GameClient
}

func (s sender) Send(ctx context.Context, orders []lugo.PlayerOrder, debugMsg string) (*lugo.OrderResponse, error) {
	orderSet := &lugo.OrderSet{
		Turn:         s.snapshot.Turn,
		DebugMessage: debugMsg,
		Orders:       []*lugo.Order{},
	}
	for _, order := range orders {
		orderSet.Orders = append(orderSet.Orders, &lugo.Order{Action: order})
	}
	s.logger.Debugf("sending orders for turn %d", s.snapshot.Turn)
	return s.gameConn.SendOrders(ctx, orderSet)
}
