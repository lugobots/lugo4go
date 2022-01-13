package lugo4go_test

import (
	"github.com/golang/mock/gomock"
	"github.com/lugobots/lugo4go/v2"
	"github.com/lugobots/lugo4go/v2/lugo"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestSender_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGRPCClient := NewMockGameClient(ctrl)
	sender := lugo4go.Sender{GRPCClient: mockGRPCClient}

	velocitySample := lugo.NewZeroedVelocity(lugo.North())
	moveOrder := &lugo.Order_Move{
		Move: &lugo.Move{
			Velocity: &velocitySample,
		},
	}
	catchOrder := &lugo.Order_Catch{}
	kickOrder := &lugo.Order_Kick{
		Kick: &lugo.Kick{
			Velocity: &velocitySample,
		},
	}
	jumpOder := &lugo.Order_Jump{
		Jump: &lugo.Jump{
			Velocity: &velocitySample,
		},
	}
	testCases := []struct {
		name string
		//inputs
		turn     uint32
		orders   []lugo.PlayerOrder
		debugMsg string

		//expected args
		orderSet *lugo.OrderSet

		//expected outputs
		clientResp *lugo.OrderResponse
		clientErr  error
	}{
		{
			name:     "no orders",
			turn:     150,
			orders:   []lugo.PlayerOrder{},
			debugMsg: "hi",
			orderSet: &lugo.OrderSet{
				Turn:         150,
				Orders:       []*lugo.Order{},
				DebugMessage: "hi",
			},
			clientResp: &lugo.OrderResponse{
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
			orderSet: &lugo.OrderSet{
				Turn:         150,
				Orders:       []*lugo.Order{},
				DebugMessage: "hi2",
			},
			clientResp: nil,
			clientErr:  nil,
		},
		{
			name:     "one move order",
			turn:     150,
			orders:   []lugo.PlayerOrder{moveOrder},
			debugMsg: "hi2",
			orderSet: &lugo.OrderSet{
				Turn: 150,
				Orders: []*lugo.Order{
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
			orders:   []lugo.PlayerOrder{catchOrder},
			debugMsg: "hi3",
			orderSet: &lugo.OrderSet{
				Turn: 150,
				Orders: []*lugo.Order{
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
			orders:   []lugo.PlayerOrder{kickOrder},
			debugMsg: "hi3",
			orderSet: &lugo.OrderSet{
				Turn: 150,
				Orders: []*lugo.Order{
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
			orders:   []lugo.PlayerOrder{jumpOder},
			debugMsg: "hi3",
			orderSet: &lugo.OrderSet{
				Turn: 150,
				Orders: []*lugo.Order{
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
			orders:   []lugo.PlayerOrder{moveOrder, kickOrder},
			debugMsg: "hi3",
			orderSet: &lugo.OrderSet{
				Turn: 150,
				Orders: []*lugo.Order{
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
			orders:   []lugo.PlayerOrder{kickOrder, moveOrder},
			debugMsg: "hi3",
			orderSet: &lugo.OrderSet{
				Turn: 150,
				Orders: []*lugo.Order{
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
