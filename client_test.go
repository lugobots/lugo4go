package lugo4go

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/lugobots/lugo4go/v2/lugo"
	"github.com/lugobots/lugo4go/v2/testdata"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"testing"
	"time"
)

const testServerPort = 2222

func NewMockServer(ctx context.Context, ctr *gomock.Controller, port int16) (*lugo.MockGameServer, error) {
	mock := lugo.NewMockGameServer(ctr)
	gRPCServer := grpc.NewServer()
	lugo.RegisterGameServer(gRPCServer, mock)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	go func() {
		<-ctx.Done()
		gRPCServer.Stop()
	}()
	go func() {
		if err := gRPCServer.Serve(lis); err != nil {
			log.Fatalf("test server has stopped: %s", err)
		}
	}()
	return mock, nil
}

func TestNewClient(t *testing.T) {
	// initiates Mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // checks all expected things for mocks

	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	// creates a fake server to test our Client
	srv, err := NewMockServer(ctx, ctrl, testServerPort)
	if err != nil {
		t.Fatalf("did not start mock server: %s", err)
	}

	config := Config{
		GRPCAddress:     fmt.Sprintf(":%d", testServerPort),
		Insecure:        true,
		TeamSide:        lugo.Team_HOME,
		Number:          3,
		InitialPosition: lugo.Point{X: 4000, Y: 4000},
	}

	// it is an async test, we have to wait some stuff be done before finishing the game, but we do not want to freeze
	waiting, done := context.WithTimeout(context.Background(), 500*time.Millisecond)

	// the Client will try to join to a team, so our server need to expect it happens
	srv.EXPECT().JoinATeam(testdata.NewMatcher(func(arg interface{}) bool {
		expectedRequest := &lugo.JoinRequest{
			Number:          config.Number,
			InitPosition:    &config.InitialPosition,
			TeamSide:        config.TeamSide,
			ProtocolVersion: ProtocolVersion,
		}
		defer done()
		return fmt.Sprintf("%s", arg) == fmt.Sprintf("%s", expectedRequest)
	}), gomock.Any()).Return(nil)

	// Now we may create the Client to connect to our fake server
	_, playerClient, err := NewClient_deprecated(config)

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
	expectedResponse := &lugo.OrderResponse{
		Code: lugo.OrderResponse_SUCCESS,
	}
	receivedSnapshot := false

	// defining mocks and expected method calls
	mockLogger := NewMockLogger(ctrl)
	mockStream := lugo.NewMockGame_JoinATeamClient(ctrl)
	mockSender := NewMockOrderSender(ctrl)
	mockLogger.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
	mockStream.EXPECT().Recv().Return(expectedSnapshot, nil)
	mockStream.EXPECT().Recv().Return(nil, io.EOF)
	mockSender.EXPECT().Send(gomock.Any(), []lugo.PlayerOrder{expectedOrder}, expectedDebugMsg).Return(expectedResponse, nil)

	c := &Client{
		stream: mockStream,
		senderBuilder: func(snapshot *lugo.GameSnapshot, logger Logger) OrderSender {
			return mockSender
		},
		ctx: context.Background(),
		stopCtx: func() {

		},
	}

	// it is an async test, we have to wait some stuff be done before finishing the game, but we do not want to freeze
	waiting, done := context.WithTimeout(context.Background(), 500*time.Millisecond)

	c.OnNewTurn(func(ctx context.Context, snapshot *lugo.GameSnapshot, sender OrderSender) {
		if snapshot != expectedSnapshot {
			t.Errorf("Unexpected snapshot - Expected %v, Got %v", expectedSnapshot, snapshot)
			return
		}
		response, err := sender.Send(ctx, []lugo.PlayerOrder{expectedOrder}, expectedDebugMsg)
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
	mockLogger := NewMockLogger(ctrl)
	mockStream := lugo.NewMockGame_JoinATeamClient(ctrl)

	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
	mockStream.EXPECT().Recv().Return(nil, io.EOF)

	wasClosed := false

	c := &Client{
		stream: mockStream,
		stopCtx: func() {
			wasClosed = true
		},
		ctx: context.Background(),
	}

	// it is an async test, we have to wait some stuff be done before finishing the game, but we do not want to freeze
	waiting, done := context.WithTimeout(context.Background(), 200*time.Millisecond)

	c.OnNewTurn(func(ctx context.Context, snapshot *lugo.GameSnapshot, sender OrderSender) {
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

	mockGameConn := lugo.NewMockGameClient(ctrl)
	mockLogger := NewMockLogger(ctrl)

	expectedSnapshot := &lugo.GameSnapshot{Turn: 321}
	expectedContext := context.Background()

	// this is the list that the developer should be concerned about while the bot is been built.
	// the rest of the job will be abstracted by the sender.
	expectedOrderSlice := []lugo.PlayerOrder{
		&lugo.Order_Catch{}, &lugo.Order_Catch{},
	}
	expectedDebugMsg := "it's a nice debug message"

	// The whole work done by the sender is converting a list of PlayerOrders into a more complex format expected
	// by the server, that's the `OrderSet`. Since OrderSet also has a debug msg and the turn number, it also receive
	// the snapshot.
	expectedOrderSet := &lugo.OrderSet{
		Turn:         expectedSnapshot.Turn,
		DebugMessage: expectedDebugMsg,
		Orders: []*lugo.Order{
			{Action: &lugo.Order_Catch{}}, {Action: &lugo.Order_Catch{}},
		},
	}
	expectedServerResponse := &lugo.OrderResponse{
		Code:    lugo.OrderResponse_SUCCESS,
		Details: "nothing else to say",
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

func TestClient_StopsIfGRPCConnectionIsInterrupted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // checks all expected things for mocks

	ctx, stopServer := context.WithCancel(context.Background())
	defer stopServer()

	// creates a fake server to test our Client
	if _, err := NewMockServer(ctx, ctrl, testServerPort); err != nil {
		t.Fatalf("did not start mock server: %s", err)
	}

	config := Config{
		GRPCAddress:     fmt.Sprintf(":%d", testServerPort),
		Insecure:        true,
		TeamSide:        lugo.Team_HOME,
		Number:          3,
		InitialPosition: lugo.Point{X: 4000, Y: 4000},
	}

	// it is an async test, we have to wait some stuff be done before finishing the game, but we do not want to freeze
	//waiting, done := context.WithTimeout(context.Background(), 500*time.Millisecond)

	// Now we may create the Client to connect to our fake server
	clientCtx, _, err := NewClient_deprecated(config)
	if err != nil {
		t.Errorf("Unexpected erro - Expected nil, Got %v", err)
	}
	// let's give some time to the Client stop after the server be stopped
	maxWait, _ := context.WithTimeout(clientCtx, 500*time.Millisecond)

	stopServer()
	<-maxWait.Done()
	assert.Equal(t, context.Canceled, clientCtx.Err())
}
