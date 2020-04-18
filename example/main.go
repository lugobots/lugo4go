package main

import (
	clientGo "github.com/lugobots/lugo4go/v2"
	"github.com/lugobots/lugo4go/v2/example/bot"
	"github.com/lugobots/lugo4go/v2/lugo"
	"github.com/lugobots/lugo4go/v2/pkg/field"
	"github.com/lugobots/lugo4go/v2/pkg/util"
	"log"
	"os"
	"os/signal"
)

func main() {
	// DefaultInitBundle is a shortcut for stuff that usually we define in init functions
	playerConfig, logger, err := util.DefaultInitBundle()
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

	player, err := clientGo.NewClient(playerConfig)
	if err != nil {
		log.Fatalf("could not init the client: %s", err)
	}
	logger.Info("connected to the game server")

	// Creating a bot to play
	myBot, errs := bot.NewBot(logger, playerConfig.TeamSide, playerConfig.Number)

	watch := make(chan error)
	go func() {
		logger.Info("starting playing")
		if err := player.PlayWithBot(myBot, logger.Named("bot")); err != nil {
			log.Fatalf("could not start to play: %s", err)
		}
		close(watch)
	}()

	go func() {
		for {
			select {
			case err := <-errs:
				//that's a nice place to implement some logic to understand your bot errors
				log.Printf("bot error: %s", err)
			case <-watch:
				return
			}
		}
	}()

	// keep the process alive
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	select {
	case <-signalChan:
		logger.Warnf("got interruption signal")
		player.Stop()
	case <-watch:
		logger.Infof("player client stopped")
	}
	logger.Infof("process finished")
}
