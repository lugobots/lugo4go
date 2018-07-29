package main

import (
	"./Game"
	"github.com/makeitplay/commons"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

var (
	player Game.Player
)

func main() {
	rand.Seed(time.Now().Unix())
	watchInterruptions()
	defer commons.Cleanup(false)
	serverConfig := new(Game.Configuration)
	serverConfig.LoadCmdArg()
	/**********************************************/

	player = Game.Player{}
	player.Start(serverConfig)
}

func watchInterruptions() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			commons.Log("*********** INTERRUPTION SIGNAL ****************")
			commons.Cleanup(true)
			os.Exit(0)
		}
	}()
}
