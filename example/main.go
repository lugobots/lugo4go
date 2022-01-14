package main

import (
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

	fieldMapper, _ := field.NewMapper(8, 4, playerConfig.TeamSide)

	region, _ := fieldMapper.GetRegion(uint8(1+(playerConfig.Number)%2), uint8(playerConfig.Number%4))

	// just creating a position for example purposes
	playerConfig.InitialPosition = region.Center()
	//&lugo.Point{
	//X: field.FieldWidth / 4,
	//Y: int32(playerConfig.Number) * field.PlayerSize * 2,
	//}

	//if playerConfig.TeamSide == lugo.Team_AWAY {
	//	playerConfig.InitialPosition.X = field.FieldWidth - playerConfig.InitialPosition.X
	//}

	player, err := clientGo.NewClient(playerConfig)
	if err != nil {
		log.Fatalf("could not init the client: %s", err)
	}
	logger.Info("connected to the game server")

	// Creating a bot to play
	myBot := bot.NewBot(logger, playerConfig.TeamSide, playerConfig.Number)

	errChan := make(chan error)
	go func() {
		errChan <- player.PlayWithBot(myBot, logger.Named("bot"))
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	select {
	case err := <-errChan:
		//that's a nice place to implement some logic to understand your bot errors
		log.Printf("bot error: %s", err)
	case <-signalChan:
		logger.Warnf("got interruption signal")
		if err := player.Stop(); err != nil {
			log.Printf("error stopping bot: %s", err)
		}
	}
	logger.Infof("process finished")
}
