package main

import (
	clientGo "github.com/lugobots/lugo4go/v2"
	"github.com/lugobots/lugo4go/v2/coach"
	"github.com/lugobots/lugo4go/v2/example/bot"
	"github.com/lugobots/lugo4go/v2/field"
	"github.com/lugobots/lugo4go/v2/lugo"
	"log"
	"os"
	"os/signal"
)

func main() {
	// DefaultBundle is a shortcut for stuff that usually we define in init functions
	playerConfig, logger, err := clientGo.DefaultBundle()
	if err != nil {
		log.Fatalf("could not init default config or logger: %s", err)
	}

	// just creating a position for example purposes
	playerConfig.InitialPosition = lugo.Point{
		X: field.FieldWidth / 4,
		Y: int32(playerConfig.Number) * field.PlayerSize * 2,
	}

	if playerConfig.TeamSide == lugo.Team_AWAY {
		playerConfig.InitialPosition.X = field.FieldWidth - playerConfig.InitialPosition.X
	}

	// open the connection to the server
	playerCtx, playerClient, err := clientGo.NewClient(playerConfig)
	if err != nil {
		logger.Fatalf("did not connected to the gRPC server at '%s': %s", playerConfig.GRPCAddress, err)
	}

	// Creating a bot to play
	myBot := bot.NewBot(logger)

	// defining the bot as the "decider" interface to be used by the Turn Handler
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
