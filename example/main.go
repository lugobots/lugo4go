package main

import (
	"context"
	clientGo "github.com/makeitplay/client-player-go"
	"github.com/makeitplay/client-player-go/lugo"
	"github.com/makeitplay/client-player-go/proto"
	"log"
	"os"
	"os/signal"
)

var logger lugo.Logger
var playerClient lugo.Client
var playerCtx context.Context
var playerConfig clientGo.Config

func main() {
	var err error
	// DefaultBundle is a shot cut for stuff that usually we define in init functions
	playerConfig, logger, err = clientGo.DefaultBundle()
	if err != nil {
		log.Fatalf("could not init default config or logger: %s", err)
	}

	// just creating a position based on the player number
	playerConfig.InitialPosition = proto.Point{
		X: lugo.FieldWidth / 4,
		Y: int32(playerConfig.Number) * lugo.PlayerSize * 2,
	}

	if playerConfig.TeamSide == proto.Team_AWAY {
		playerConfig.InitialPosition.X = lugo.FieldWidth - playerConfig.InitialPosition.X
	}

	playerCtx, playerClient, err = clientGo.NewClient(playerConfig)
	if err != nil {
		logger.Fatalf("did not connected to the gRPC server at '%s': %s", playerConfig.GRPCAddress, err)
	}
	playerClient.OnNewTurn(myDecider, logger)

	// keep the process alive
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	select {
	case <-signalChan:
		logger.Warnf("got interruption signal")
		if err := playerClient.Stop(); err != nil {
			logger.Errorf("error stopping the player client: %s", err)
		}
	case <-playerCtx.Done():
		logger.Infof("player client stopped")
	}
	logger.Infof("process finished")
}

func myDecider(snapshot *proto.GameSnapshot, sender lugo.OrderSender) {
	me := lugo.GetPlayer(snapshot, playerConfig.TeamSide, playerConfig.Number)
	if me == nil {
		logger.Fatalf("i did not find my self in the game")
		return
	}
	var orders []proto.PlayerOrder
	// we are going to kick the ball as soon as we catch it
	if lugo.IsBallHolder(snapshot, me) {
		orderToKick, err := lugo.MakeOrderKick(*snapshot.Ball, lugo.GetOpponentGoal(me.TeamSide).Center, lugo.BallMaxSpeed)
		if err != nil {
			logger.Warnf("could not create kick order during turn %d: %s", snapshot.Turn, err)
			return
		}
		orders = []proto.PlayerOrder{orderToKick}
	} else if me.Number == 10 {
		// otherwise, let's run towards the ball like kids
		orderToMove, err := lugo.MakeOrderMoveMaxSpeed(*me.Position, *snapshot.Ball.Position)
		if err != nil {
			logger.Warnf("could not create move order during turn %d: %s", snapshot.Turn, err)
			return
		}
		orders = []proto.PlayerOrder{orderToMove, lugo.MakeOrderCatch()}
	} else {
		orders = []proto.PlayerOrder{lugo.MakeOrderCatch()}
	}

	resp, err := sender.Send(playerCtx, orders, "")
	if err != nil {
		logger.Warnf("could not send kick order during turn %d: %s", snapshot.Turn, err)
	} else if resp.Code != proto.OrderResponse_SUCCESS {
		logger.Warnf("order sent not  order during turn %d: %s", snapshot.Turn, err)
	}
}
