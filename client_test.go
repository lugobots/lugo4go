package client

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/lugobots/client-player-go/v2/lugo"
	"github.com/lugobots/client-player-go/v2/proto"
	"github.com/lugobots/client-player-go/v2/testdata"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
	"time"
)

const testServerPort = 2222

func TestNewClient(t *testing.T) {
	// initiates Mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // checks all expected things for mocks

	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	// creates a fake server to test our client
	srv, err := testdata.NewMockServer(ctx, ctrl, testServerPort)
	if err != nil {
		t.Fatalf("did not start mock server: %s", err)
	}

	config := Config{
		GRPCAddress:     fmt.Sprintf(":%d", testServerPort),
		Insecure:        true,
		TeamSide:        proto.Team_HOME,
		Number:          3,
		InitialPosition: proto.Point{X: 4000, Y: 4000},
	}

	// it is an async test, we have to wait some stuff be done before finishing the game, but we do not want to freeze
	waiting, done := context.WithTimeout(context.Background(), 500*time.Millisecond)

	// the client will try to join to a team, so our server need to expect it happens
	srv.EXPECT().JoinATeam(testdata.NewMatcher(func(arg interface{}) bool {
		expectedRequest := &proto.JoinRequest{
			Number:          config.Number,
			InitPosition:    &config.InitialPosition,
			TeamSide:        config.TeamSide,
			ProtocolVersion: ProtocolVersion,
		}
		defer done()
		return fmt.Sprintf("%s", arg) == fmt.Sprintf("%s", expectedRequest)
	}), gomock.Any()).Return(nil)

	// Now we may create the client to connect to our fake server
	_, playerClient, err := NewClient(config)

	// This last lines may run really quickly, and the server may not have ran the expected methods yet
	// Let's give some time to the server run it before finish the test function
	<-waiting.Done()
	if err != nil {
		t.Fatalf("did not connect to the server: %s", err)
	}
	if err := playerClient.Stop(); err != nil {
		t.Fatalf("did not connect to the server: %s", err)
	}
}

func TestClient_OnNewTurn(t *testing.T) {
	// initiates Mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // checks all expected things for mocks

	// defining expectations

	expectedSnapshot := &proto.GameSnapshot{Turn: 200}
	expectedOrder := &proto.Order_Catch{}
	expectedDebugMsg := "a-important-msg"
	expectedResponse := &proto.OrderResponse{
		Code: proto.OrderResponse_SUCCESS,
	}
	receivedSnapshot := false

	// defining mocks and expected method calls
	mockLogger := testdata.NewMockLogger(ctrl)
	mockStream := testdata.NewMockGame_JoinATeamClient(ctrl)
	mockSender := testdata.NewMockOrderSender(ctrl)
	mockLogger.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
	mockStream.EXPECT().Recv().Return(expectedSnapshot, nil)
	mockStream.EXPECT().Recv().Return(nil, io.EOF)
	mockSender.EXPECT().Send(gomock.Any(), []proto.PlayerOrder{expectedOrder}, expectedDebugMsg).Return(expectedResponse, nil)

	c := &client{
		stream: mockStream,
		senderBuilder: func(snapshot *proto.GameSnapshot, logger lugo.Logger) lugo.OrderSender {
			return mockSender
		},
		ctx: context.Background(),
		stopCtx: func() {

		},
	}

	// it is an async test, we have to wait some stuff be done before finishing the game, but we do not want to freeze
	waiting, done := context.WithTimeout(context.Background(), 500*time.Millisecond)

	c.OnNewTurn(func(snapshot *proto.GameSnapshot, sender lugo.OrderSender) {
		if snapshot != expectedSnapshot {
			t.Errorf("Unexpected snapshot - Expected %v, Got %v", expectedSnapshot, snapshot)
			return
		}
		response, err := sender.Send(waiting, []proto.PlayerOrder{expectedOrder}, expectedDebugMsg)
		if err != nil {
			t.Errorf("Unexpected erro - Expected nil, Got %v", err)
		}

		if response != expectedResponse {
			t.Errorf("Unexpected response - Expected %v, Got %v", expectedResponse, response)
		}
		receivedSnapshot = true
		done()
	}, mockLogger)

	// This last lines may run really quickly, and the mock may not have ran the expected methods yet
	// Let's give some time to the server run it before finish the test function
	<-waiting.Done()
	if !receivedSnapshot {
		t.Error("Expected has received msg")
	}
	if waiting.Err() != context.Canceled {
		t.Errorf("Unexpected waiting - Expected %v, Got %v", context.Canceled, waiting.Err())
	}
}

func TestClient_ShouldStopItsContext(t *testing.T) {
	// initiates Mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // checks all expected things for mocks

	// defining mocks and expected method calls
	mockLogger := testdata.NewMockLogger(ctrl)
	mockStream := testdata.NewMockGame_JoinATeamClient(ctrl)

	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
	mockStream.EXPECT().Recv().Return(nil, io.EOF)

	wasClosed := false

	c := &client{
		stream: mockStream,
		stopCtx: func() {
			wasClosed = true
		},
	}

	// it is an async test, we have to wait some stuff be done before finishing the game, but we do not want to freeze
	waiting, done := context.WithTimeout(context.Background(), 200*time.Millisecond)

	c.OnNewTurn(func(snapshot *proto.GameSnapshot, sender lugo.OrderSender) {
		t.Error("The DecisionMaker should not be called")
		done()
	}, mockLogger)

	// This last lines may run really quickly, and the mock may not have ran the expected methods yet
	// Let's give some time to the server run it before finish the test function
	<-waiting.Done()
	if waiting.Err() != context.DeadlineExceeded {
		t.Errorf("Unexpected waiting - Expected %v, Got %v", context.DeadlineExceeded, waiting.Err())
	}

	if !wasClosed {
		t.Error("Unexpected context to be closed")
	}
}

func TestSender_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // checks all expected things for mocks

	mockGameConn := testdata.NewMockGameClient(ctrl)
	mockLogger := testdata.NewMockLogger(ctrl)

	expectedSnapshot := &proto.GameSnapshot{Turn: 321}
	expectedContext := context.Background()

	// this is the list that the developer should be concerned about while the bot is been built.
	// the rest of the job will be abstracted by the sender.
	expectedOrderSlice := []proto.PlayerOrder{
		&proto.Order_Catch{}, &proto.Order_Catch{},
	}
	expectedDebugMsg := "it's a nice debug message"

	// The whole work done by the sender is converting a list of PlayerOrders into a more complex format expected
	// by the server, that's the `OrderSet`. Since OrderSet also has a debug msg and the turn number, it also receive
	// the snapshot.
	expectedOrderSet := &proto.OrderSet{
		Turn:         expectedSnapshot.Turn,
		DebugMessage: expectedDebugMsg,
		Orders: []*proto.Order{
			{Action: &proto.Order_Catch{}}, {Action: &proto.Order_Catch{}},
		},
	}
	expectedServerResponse := &proto.OrderResponse{
		Code:    proto.OrderResponse_SUCCESS,
		Details: "nonthing else to say",
	}

	mockGameConn.EXPECT().SendOrders(expectedContext, expectedOrderSet).Return(expectedServerResponse, nil)
	mockLogger.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes()

	orderSender := &sender{
		gameConn: mockGameConn,
		snapshot: expectedSnapshot,
		logger:   mockLogger,
	}

	serverResponse, err := orderSender.Send(expectedContext, expectedOrderSlice, expectedDebugMsg)

	assert.Nil(t, err)
	assert.Equal(t, expectedServerResponse, serverResponse)
}
