package lugo4go

import (
	"context"

	"github.com/lugobots/lugo4go/v2/proto"
)

func NewSender(grpcClient proto.GameClient) *Sender {
	return &Sender{
		GRPCClient: grpcClient,
	}
}

type Sender struct {
	GRPCClient proto.GameClient
}

func (s *Sender) Send(ctx context.Context, turn uint32, orders []proto.PlayerOrder, debugMsg string) (*proto.OrderResponse, error) {
	orderSet := &proto.OrderSet{
		Turn:         turn,
		DebugMessage: debugMsg,
		Orders:       []*proto.Order{},
	}
	for _, order := range orders {
		orderSet.Orders = append(orderSet.Orders, &proto.Order{Action: order})
	}
	return s.GRPCClient.SendOrders(ctx, orderSet)
}
