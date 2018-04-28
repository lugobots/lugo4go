package main

import (
	"os"
	"math/rand"
	"time"
	"os/signal"
	"./Game"
	"github.com/maketplay/commons"
)

var (
	player Game.Player
)

func main() {
	rand.Seed(time.Now().Unix())
	watchInterruptions()
	defer commons.Cleanup(false)
	serverConfig := new(Game.Configuration)
	commons.Load(serverConfig)
	serverConfig.LoadCmdArg()
	/**********************************************/

	player = Game.Player{}
	commons.NickName = "New Player"
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
