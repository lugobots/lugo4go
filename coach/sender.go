package coach

import (
	"context"
	"github.com/lugobots/lugo4go/v2/lugo"
)

func NewSender(grpcClient lugo.GameClient) *Sender {
	return &Sender{
		grpcClient: grpcClient,
	}
}

type Sender struct {
	grpcClient lugo.GameClient
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
	return s.grpcClient.SendOrders(ctx, orderSet)
}
