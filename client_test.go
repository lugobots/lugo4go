package client_test

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/makeitplay/client-player-go"
	"github.com/makeitplay/client-player-go/lugo"
	"github.com/makeitplay/client-player-go/testdata"
	"testing"
	"time"
)

const testServerPort = 2222

func TestNewClient_ShouldJoinToATeam(t *testing.T) {
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

	config := client.Config{
		GRPCHost:        fmt.Sprintf(":%d", testServerPort),
		Insecure:        true,
		TeamSide:        lugo.Team_HOME,
		Number:          3,
		InitialPosition: &lugo.Point{X: 4000, Y: 4000},
	}

	// the client will try to join to a team, so our server need to expect it happens
	srv.EXPECT().JoinATeam(testdata.NewMatcher(func(arg interface{}) bool {
		expectedRequest := &lugo.JoinRequest{
			Number:          config.Number,
			InitPosition:    config.InitialPosition,
			TeamSide:        config.TeamSide,
			ProtocolVersion: client.ProtocolVersion,
		}
		return fmt.Sprintf("%s", arg) == fmt.Sprintf("%s", expectedRequest)
	}), gomock.Any()).Return(nil)

	// Now we may create the client to connect to our fake server
	playerClient, closer, err := client.NewClient(ctx, config)
	defer playerClient.Stop()

	// This last lines may run really quickly, and the server may not have rnu the expected methods yet
	// Let's give some time to the server run it before finish the test function
	time.Sleep(100 * time.Millisecond)
	if err != nil {
		t.Fatalf("did not connect to the server: %s", err)
	}
	if err := closer.Close(); err != nil {
		t.Fatalf("did not connect to the server: %s", err)
	}
}
