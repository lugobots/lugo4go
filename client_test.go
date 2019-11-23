package client

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/makeitplay/client-player-go/lugo"
	"github.com/makeitplay/client-player-go/ops"
	"github.com/makeitplay/client-player-go/testdata"
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
		GRPCHost:        fmt.Sprintf(":%d", testServerPort),
		Insecure:        true,
		TeamSide:        lugo.Team_HOME,
		Number:          3,
		InitialPosition: &lugo.Point{X: 4000, Y: 4000},
	}

	// it is an async test, we have to wait some stuff be done before finishing the game, but we do not want to freeze
	waiting, done := context.WithTimeout(context.Background(), 500*time.Millisecond)

	// the client will try to join to a team, so our server need to expect it happens
	srv.EXPECT().JoinATeam(testdata.NewMatcher(func(arg interface{}) bool {
		expectedRequest := &lugo.JoinRequest{
			Number:          config.Number,
			InitPosition:    config.InitialPosition,
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

	expectedSnapshot := &lugo.GameSnapshot{Turn: 200}
	expectedOrder := &lugo.Order_Catch{}
	expectedDebugMsg := "a-important-msg"
	expectedOrderSet := &lugo.OrderSet{
		Turn:         200,
		DebugMessage: expectedDebugMsg,
		Orders:       []*lugo.Order{{Action: expectedOrder}},
	}
	expectedResponse := &lugo.OrderResponse{
		Code: lugo.OrderResponse_SUCCESS,
	}
	receivedSnapshot := false

	// defining mocks and expected method calls
	mockLogger := testdata.NewMockLogger(ctrl)
	mockStream := testdata.NewMockGame_JoinATeamClient(ctrl)
	mockGameClient := testdata.NewMockGameClient(ctrl)
	mockLogger.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
	mockStream.EXPECT().Recv().Return(expectedSnapshot, nil)
	mockStream.EXPECT().Recv().Return(nil, io.EOF)
	mockGameClient.EXPECT().SendOrders(gomock.Any(), expectedOrderSet).Return(expectedResponse, nil)

	c := &client{
		stream:   mockStream,
		gameConn: mockGameClient,
		ctx:      context.Background(),
		stopCtx: func() {

		},
	}

	// it is an async test, we have to wait some stuff be done before finishing the game, but we do not want to freeze
	waiting, done := context.WithTimeout(context.Background(), 500*time.Millisecond)

	c.OnNewTurn(func(snapshot *lugo.GameSnapshot, sender ops.OrderSender) {
		if snapshot != expectedSnapshot {
			t.Errorf("Unexpected snapshot - Expected %v, Got %v", expectedSnapshot, snapshot)
			return
		}
		response, err := sender(expectedDebugMsg, expectedOrder)
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

	c.OnNewTurn(func(snapshot *lugo.GameSnapshot, sender ops.OrderSender) {
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
