package client

import (
	"context"
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/orders"
	"github.com/makeitplay/arena/physics"
	"testing"
	"time"
)

func TestController_NextTurn(t *testing.T) {
	ctx, stop := context.WithCancel(context.Background())

	defer stop()
	serverConfig := new(Configuration)
	serverConfig.WSPort = "8080"
	serverConfig.WSHost = "localhost"
	serverConfig.UUID = "local"
	_, ctrl, err := NewTestController(ctx, *serverConfig)
	if err != nil {
		t.Fatal(err.Error())
	}

	v := physics.NewZeroedVelocity(physics.North)
	ctrl.SetBallProperties(v, arena.HomeTeamGoal.Center)

	for i := 0; i < 5; i++ {
		ctrl.SendOrders(arena.HomeTeam, arena.GoalkeeperNumber, []orders.Order{})
		time.Sleep(2 * time.Second)
		ctrl.NextTurn()
	}

}
