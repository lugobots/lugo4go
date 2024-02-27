package lugo4go

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/lugobots/lugo4go/v3/proto"
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

	config := Config{
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
			Number:          uint32(config.Number),
			InitPosition:    config.InitialPosition,
			TeamSide:        config.TeamSide,
			ProtocolVersion: ProtocolVersion,
		}
		defer done()
		return fmt.Sprintf("%s", arg) == fmt.Sprintf("%s", expectedRequest)
	}), gomock.Any()).Return(nil)

	// Now we may create the Client to connect to our fake server
	playerClient, err := NewClient(config, DefaultLogger(config))

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

func TestClient_PlayEndsTheConnectionCorrectly(t *testing.T) {
	// initiates Mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // checks all expected things for mocks

	// defining mocks and expected method calls
	mockStream := NewMockGame_JoinATeamClient(ctrl)
	mockGRPCClient := NewMockGameClient(ctrl)
	mockHandler := NewMockRawBot(ctrl)
	mockSender := NewMockOrderSender(ctrl)

	c := &Client{
		Stream:     mockStream,
		GRPCClient: mockGRPCClient,
		config: Config{
			TeamSide: proto.Team_AWAY,
			Number:   5,
		},
		Logger: DefaultLogger(Config{}),
		Sender: mockSender,
	}

	// it is an async test, we have to wait some stuff be done before finishing the game, but we do not want to freeze
	waiting, done := context.WithTimeout(context.Background(), 500*time.Millisecond)

	// defining expectations
	expectedSnapshot := &proto.GameSnapshot{
		State: proto.GameSnapshot_LISTENING,
		Turn:  200,
		AwayTeam: &proto.Team{
			Players: []*proto.Player{
				{Number: 5},
			},
		},
	}
	mockStream.EXPECT().Recv().Return(expectedSnapshot, nil)
	mockStream.EXPECT().Recv().DoAndReturn(func() {
		//let's pretend some interval between messages
		time.Sleep(50 * time.Millisecond)
	}).Return(nil, io.EOF)

	mockHandler.EXPECT().TurnHandler(gomock.Any(), gomock.Any())
	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), nil, "").Return(&proto.OrderResponse{Code: proto.OrderResponse_SUCCESS}, nil)

	err := c.Play(mockHandler)
	done()
	assert.Equal(t, ErrGRPCConnectionClosed, err)
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
	mockHandler := NewMockRawBot(ctrl)

	c := &Client{
		Stream:     mockStream,
		GRPCClient: mockGRPCClient,
		Logger:     DefaultLogger(Config{}),
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
	assert.True(t, errors.Is(err, ErrGRPCConnectionLost))
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
	mockHandler := NewMockRawBot(ctrl)
	mockSender := NewMockOrderSender(ctrl)

	c := &Client{
		Stream:     mockStream,
		GRPCClient: mockGRPCClient,
		config: Config{
			TeamSide: proto.Team_AWAY,
			Number:   5,
		},
		Sender: mockSender,
		Logger: DefaultLogger(Config{}),
	}

	// it is an async test, we have to wait some stuff be done before finishing the game, but we do not want to freeze
	waiting, done := context.WithTimeout(context.Background(), 500*time.Millisecond)

	awayTeam := &proto.Team{
		Players: []*proto.Player{
			{Number: 5},
		},
	}

	// defining expectations
	expectedSnapshotA := &proto.GameSnapshot{State: proto.GameSnapshot_LISTENING, Turn: 200, AwayTeam: awayTeam}
	expectedSnapshotB := &proto.GameSnapshot{State: proto.GameSnapshot_LISTENING, Turn: 201, AwayTeam: awayTeam}
	mockStream.EXPECT().Recv().Return(expectedSnapshotA, nil)
	mockStream.EXPECT().Recv().Return(expectedSnapshotB, nil)
	mockStream.EXPECT().Recv().Return(nil, io.EOF)

	firstHandlerIsExpired := false
	holder := make(chan bool, 1)
	// this is the first time the rawBotWrapper will be called
	// we expect that it finishes immediately after the stream gets a new turn msg
	// if it does not, the next rawBotWrapper will unblock it, but it will be considered an error
	mockHandler.EXPECT().
		TurnHandler(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, inspector SnapshotInspector) {
			select {
			case <-ctx.Done():
				firstHandlerIsExpired = true
			case <-holder:
				firstHandlerIsExpired = false
			}

		}).Return(nil, "", nil)

	// the second call to rawBotWrapper will close our channel just to ensure anything will be left behind in your test
	mockHandler.EXPECT().
		TurnHandler(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, inspector SnapshotInspector) {
			close(holder)
		}).Return(nil, "", nil)

	mockSender.EXPECT().Send(gomock.Any(), gomock.Any(), nil, "").Return(&proto.OrderResponse{Code: proto.OrderResponse_SUCCESS}, nil).AnyTimes()

	err := c.Play(mockHandler)
	done()
	assert.Equal(t, ErrGRPCConnectionClosed, err)
	assert.True(t, firstHandlerIsExpired)
	if waiting.Err() != context.Canceled {
		t.Errorf("Unexpected waiting - Expected %v, Got %v", context.Canceled, waiting.Err())
	}
}
