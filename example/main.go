package main

import (
	"context"
	clientGo "github.com/lugobots/lugo4go/v2"
	"github.com/lugobots/lugo4go/v2/example/bot"
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

	fieldGridCols := uint8(8)
	fieldGridRows := uint8(8)

	fieldMapper, _ := field.NewMapper(fieldGridCols, fieldGridRows, playerConfig.TeamSide)

	region, _ := fieldMapper.GetRegion(bot.FieldMap[playerConfig.Number].Col, bot.FieldMap[playerConfig.Number].Row)

	// just creating a position for example purposes
	playerConfig.InitialPosition = region.Center()

	player, err := clientGo.NewClient(playerConfig)
	if err != nil {
		log.Fatalf("could not init the client: %s", err)
	}
	logger.Info("connected to the game server")

	// Creating a bot to play
	myBot := bot.NewBot(logger, playerConfig.TeamSide, playerConfig.Number)

	ctx, stop := context.WithCancel(context.Background())
	go func() {
		defer stop()
		if err := player.PlayWithBot(myBot, logger.Named("bot")); err != nil {
			log.Printf("bot stopped with an error: %s", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	select {
	case <-ctx.Done():
	case <-signalChan:
		logger.Warnf("got interruption signal")
		if err := player.Stop(); err != nil {
			log.Printf("error stopping bot: %s", err)
		}
	}
	logger.Infof("process finished")
}
