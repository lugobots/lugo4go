package lugo4go_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/lugobots/lugo4go/v2"
	util2 "github.com/lugobots/lugo4go/v2/pkg/util"
	"github.com/lugobots/lugo4go/v2/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"testing"
	"time"
)

const testServerPort = 2222

func NewMockServer(ctx context.Context, ctr *gomock.Controller, port int16) (*MockGameServer, error) {
	mock := NewMockGameServer(ctr)
	gRPCServer := grpc.NewServer()
	proto.RegisterGameServer(gRPCServer, mock)

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

func TestNewRawClient(t *testing.T) {
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

	config := util2.Config{
		GRPCAddress:     fmt.Sprintf(":%d", testServerPort),
		Insecure:        true,
		TeamSide:        proto.Team_HOME,
		Number:          3,
		InitialPosition: &proto.Point{X: 4000, Y: 4000},
	}

	// it is an async test, we have to wait some stuff be done before finishing the game, but we do not want to freeze
	waiting, done := context.WithTimeout(context.Background(), 500*time.Millisecond)

	// the Client will try to join to a team, so our server need to expect it happens
	srv.EXPECT().JoinATeam(NewMatcher(func(arg interface{}) bool {
		expectedRequest := &proto.JoinRequest{
			Number:          config.Number,
			InitPosition:    config.InitialPosition,
			TeamSide:        config.TeamSide,
			ProtocolVersion: lugo4go.ProtocolVersion,
		}
		defer done()
		return fmt.Sprintf("%s", arg) == fmt.Sprintf("%s", expectedRequest)
	}), gomock.Any()).Return(nil)

	// Now we may create the Client to connect to our fake server
	playerClient, err := lugo4go.NewClient(config)

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

func TestClient_PlayCallsHandlerForEachMessage(t *testing.T) {
	// initiates Mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // checks all expected things for mocks

	// defining mocks and expected method calls
	mockStream := NewMockGame_JoinATeamClient(ctrl)
	mockGRPCClient := NewMockGameClient(ctrl)
	mockHandler := NewMockTurnHandler(ctrl)

	c := &lugo4go.Client{
		Stream:     mockStream,
		GRPCClient: mockGRPCClient,
		Handler:    mockHandler,
	}

	// it is an async test, we have to wait some stuff be done before finishing the game, but we do not want to freeze
	waiting, done := context.WithTimeout(context.Background(), 500*time.Millisecond)

	// defining expectations
	expectedSnapshot := &proto.GameSnapshot{Turn: 200}
	mockStream.EXPECT().Recv().Return(expectedSnapshot, nil)
	mockStream.EXPECT().Recv().DoAndReturn(func() {
		//let's pretend some interval between messages
		time.Sleep(50 * time.Millisecond)
	}).Return(nil, io.EOF)

	mockHandler.EXPECT().Handle(gomock.Any(), expectedSnapshot)

	err := c.Play(mockHandler)
	done()
	assert.Equal(t, lugo4go.ErrGRPCConnectionClosed, err)
	if waiting.Err() != context.Canceled {
		t.Errorf("Unexpected waiting - Expected %v, Got %v", context.Canceled, waiting.Err())
	}
}

func TestClient_PlayReturnsTheRightError(t *testing.T) {
	// initiates Mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // checks all expected things for mocks

	// defining mocks and expected method calls
	mockStream := NewMockGame_JoinATeamClient(ctrl)
	mockGRPCClient := NewMockGameClient(ctrl)
	mockHandler := NewMockTurnHandler(ctrl)

	c := &lugo4go.Client{
		Stream:     mockStream,
		GRPCClient: mockGRPCClient,
		Handler:    mockHandler,
	}

	// it is an async test, we have to wait some stuff be done before finishing the game, but we do not want to freeze
	waiting, done := context.WithTimeout(context.Background(), 500*time.Millisecond)

	// defining expectations
	expectedError := errors.New("some-error")
	mockStream.EXPECT().Recv().DoAndReturn(func() {
		//let's pretend some interval between messages
		time.Sleep(50 * time.Millisecond)
	}).Return(nil, expectedError)

	err := c.Play(mockHandler)
	done()
	assert.True(t, errors.Is(err, lugo4go.ErrGRPCConnectionLost))
	if waiting.Err() != context.Canceled {
		t.Errorf("Unexpected waiting - Expected %v, Got %v", context.Canceled, waiting.Err())
	}
}

func TestClient_PlayShouldStopContextWhenANewTurnStarts(t *testing.T) {
	// initiates Mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // checks all expected things for mocks

	// defining mocks and expected method calls
	mockStream := NewMockGame_JoinATeamClient(ctrl)
	mockGRPCClient := NewMockGameClient(ctrl)
	mockHandler := NewMockTurnHandler(ctrl)

	c := &lugo4go.Client{
		Stream:     mockStream,
		GRPCClient: mockGRPCClient,
		Handler:    mockHandler,
	}

	// it is an async test, we have to wait some stuff be done before finishing the game, but we do not want to freeze
	waiting, done := context.WithTimeout(context.Background(), 500*time.Millisecond)

	// defining expectations
	expectedSnapshotA := &proto.GameSnapshot{Turn: 200}
	expectedSnapshotB := &proto.GameSnapshot{Turn: 201}
	mockStream.EXPECT().Recv().Return(expectedSnapshotA, nil)
	mockStream.EXPECT().Recv().Return(expectedSnapshotB, nil)
	mockStream.EXPECT().Recv().Return(nil, io.EOF)

	firstHandlerIsExpired := false
	holder := make(chan bool, 1)
	// this is the first time the handler will be called
	// we expect that it finishes immediately after the stream gets a new turn msg
	// if it does not, the next handler will unblock it, but it will be considered an error
	mockHandler.EXPECT().
		Handle(gomock.Any(), expectedSnapshotA).
		DoAndReturn(func(ctx context.Context, snapshot *proto.GameSnapshot) {
			select {
			case <-ctx.Done():
				firstHandlerIsExpired = true
			case <-holder:
				firstHandlerIsExpired = false
			}
		})

	// the second call to handler will close our channel just to ensure anything will be left behind in your test
	mockHandler.EXPECT().
		Handle(gomock.Any(), expectedSnapshotB).
		DoAndReturn(func(ctx context.Context, snapshot *proto.GameSnapshot) {
			close(holder)
		})

	err := c.Play(mockHandler)
	done()
	assert.Equal(t, lugo4go.ErrGRPCConnectionClosed, err)
	assert.True(t, firstHandlerIsExpired)
	if waiting.Err() != context.Canceled {
		t.Errorf("Unexpected waiting - Expected %v, Got %v", context.Canceled, waiting.Err())
	}
}
