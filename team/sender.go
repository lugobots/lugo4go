package team

import (
	"context"
	"github.com/lugobots/lugo4go/v2/lugo"
)

func NewSender(grpcClient lugo.GameClient) *Sender {
	return &Sender{
		GRPCClient: grpcClient,
	}
}

type Sender struct {
	GRPCClient lugo.GameClient
}

func (s *Sender) Send(ctx context.Context, turn uint32, orders []lugo.PlayerOrder, debugMsg string) (*lugo.OrderResponse, error) {
	orderSet := &lugo.OrderSet{
		Turn:         turn,
		DebugMessage: debugMsg,
		Orders:       []*lugo.Order{},
	}
	for _, order := range orders {
		orderSet.Orders = append(orderSet.Orders, &lugo.Order{Action: order})
	}
	return s.GRPCClient.SendOrders(ctx, orderSet)
}
