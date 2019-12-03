package main

import (
	"context"
	clientGo "github.com/lugobots/lugo4go/v2"
	"github.com/lugobots/lugo4go/v2/coach"
	"github.com/lugobots/lugo4go/v2/example/bot"
	"github.com/lugobots/lugo4go/v2/field"
	"github.com/lugobots/lugo4go/v2/proto"
	"log"
	"os"
	"os/signal"
)

var logger clientGo.Logger
var playerClient clientGo.Client
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
		X: field.FieldWidth / 4,
		Y: int32(playerConfig.Number) * field.PlayerSize * 2,
	}

	if playerConfig.TeamSide == proto.Team_AWAY {
		playerConfig.InitialPosition.X = field.FieldWidth - playerConfig.InitialPosition.X
	}

	playerCtx, playerClient, err = clientGo.NewClient(playerConfig)
	if err != nil {
		logger.Fatalf("did not connected to the gRPC server at '%s': %s", playerConfig.GRPCAddress, err)
	}
	myBot := bot.NewBot(playerCtx, logger)
	playerClient.OnNewTurn(coach.DefaultTurnHandler(myBot, playerConfig, logger), logger)

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
