package lugo4go

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"github.com/lugobots/lugo4go/v3/proto"
)

func TestSender_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGRPCClient := NewMockGameClient(ctrl)
	sender := sender{GRPCClient: mockGRPCClient}

	velocitySample := proto.NewZeroedVelocity(proto.North())
	moveOrder := &proto.Order_Move{
		Move: &proto.Move{
			Velocity: &velocitySample,
		},
	}
	catchOrder := &proto.Order_Catch{}
	kickOrder := &proto.Order_Kick{
		Kick: &proto.Kick{
			Velocity: &velocitySample,
		},
	}
	jumpOder := &proto.Order_Jump{
		Jump: &proto.Jump{
			Velocity: &velocitySample,
		},
	}
	testCases := []struct {
		name string
		//inputs
		turn     uint32
		orders   []proto.PlayerOrder
		debugMsg string

		//expected args
		orderSet *proto.OrderSet

		//expected outputs
		clientResp *proto.OrderResponse
		clientErr  error
	}{
		{
			name:     "no orders",
			turn:     150,
			orders:   []proto.PlayerOrder{},
			debugMsg: "hi",
			orderSet: &proto.OrderSet{
				Turn:         150,
				Orders:       []*proto.Order{},
				DebugMessage: "hi",
			},
			clientResp: &proto.OrderResponse{
				Code:    0,
				Details: "",
			},
			clientErr: nil,
		},
		{
			name:     "invalid input",
			turn:     150,
			orders:   nil,
			debugMsg: "hi2",
			orderSet: &proto.OrderSet{
				Turn:         150,
				Orders:       []*proto.Order{},
				DebugMessage: "hi2",
			},
			clientResp: nil,
			clientErr:  nil,
		},
		{
			name:     "one move order",
			turn:     150,
			orders:   []proto.PlayerOrder{moveOrder},
			debugMsg: "hi2",
			orderSet: &proto.OrderSet{
				Turn: 150,
				Orders: []*proto.Order{
					{Action: moveOrder},
				},
				DebugMessage: "hi2",
			},
			clientResp: nil,
			clientErr:  nil,
		},
		{
			name:     "one catch order",
			turn:     150,
			orders:   []proto.PlayerOrder{catchOrder},
			debugMsg: "hi3",
			orderSet: &proto.OrderSet{
				Turn: 150,
				Orders: []*proto.Order{
					{Action: catchOrder},
				},
				DebugMessage: "hi3",
			},
			clientResp: nil,
			clientErr:  nil,
		},
		{
			name:     "one Kick order",
			turn:     150,
			orders:   []proto.PlayerOrder{kickOrder},
			debugMsg: "hi3",
			orderSet: &proto.OrderSet{
				Turn: 150,
				Orders: []*proto.Order{
					{Action: kickOrder},
				},
				DebugMessage: "hi3",
			},
			clientResp: nil,
			clientErr:  nil,
		},
		{
			name:     "one jump order",
			turn:     150,
			orders:   []proto.PlayerOrder{jumpOder},
			debugMsg: "hi3",
			orderSet: &proto.OrderSet{
				Turn: 150,
				Orders: []*proto.Order{
					{Action: jumpOder},
				},
				DebugMessage: "hi3",
			},
			clientResp: nil,
			clientErr:  nil,
		},
		{
			name:     "two order move and kick",
			turn:     150,
			orders:   []proto.PlayerOrder{moveOrder, kickOrder},
			debugMsg: "hi3",
			orderSet: &proto.OrderSet{
				Turn: 150,
				Orders: []*proto.Order{
					{Action: moveOrder},
					{Action: kickOrder},
				},
				DebugMessage: "hi3",
			},
			clientResp: nil,
			clientErr:  nil,
		},
		{
			name:     "two order kick and move ",
			turn:     150,
			orders:   []proto.PlayerOrder{kickOrder, moveOrder},
			debugMsg: "hi3",
			orderSet: &proto.OrderSet{
				Turn: 150,
				Orders: []*proto.Order{
					{Action: kickOrder},
					{Action: moveOrder},
				},
				DebugMessage: "hi3",
			},
			clientResp: nil,
			clientErr:  nil,
		},
	}

	for _, testCase := range testCases {
		ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)

		mockGRPCClient.EXPECT().SendOrders(ctx, testCase.orderSet).Return(testCase.clientResp, testCase.clientErr)
		resp, err := sender.Send(ctx, testCase.turn, testCase.orders, testCase.debugMsg)
		assert.Equal(t, testCase.clientResp, resp)
		assert.Equal(t, testCase.clientErr, err)
	}
}
