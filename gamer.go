package client

import (
	"context"
	"encoding/json"
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/orders"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/arena/talk"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
)

type Responser interface {
	SendOrders(message string, ordersList ...orders.Order)
}

type Gamer struct {
	OnMessage      func(msg GameMessage)
	OnAnnouncement func(turnTx TurnContext)
	config         *Configuration
	Talker         talk.Talker
	LastMsg        GameMessage
}

func listenServerMessages(gamer *Gamer) {
	for {
		select {
		case bytes := <-gamer.Talker.Listen():
			var msg GameMessage
			err := json.Unmarshal(bytes, &msg)
			if err != nil {
				playerCtx.Logger().Errorf("Fail on convert wb message: %s (%s)", err.Error(), bytes)
			} else {
				gamer.onMessage(msg)
			}
		case connError := <-gamer.Talker.ListenInterruption():
			if gamer.LastMsg.State == arena.Over {
				playerCtx.Logger().Info("game over")
				gamer.stopToPlay(false)
			} else {
				playerCtx.Logger().Errorf("ws connection lost: %s", connError.Error())
				gamer.stopToPlay(true)
			}
			return
		}
	}
}

// playerCtx is used to keep the process running while the player is playing
var playerCtx GamerCtx
var stopPlayer context.CancelFunc

// Play make the player start to play
func (p *Gamer) Play(initialPosition physics.Point, configuration *Configuration) {
	playerCtx, stopPlayer = NewGamerContext(context.Background(), configuration)

	p.config = configuration
	talkerCtx, talker, err := TalkerSetup(playerCtx, configuration, initialPosition)
	if err != nil {
		log.Fatal(err)
	}
	// we have to set the call back function that will process the player behaviour when the game state has been changed
	defer talker.Close()
	p.Talker = talker
	go listenServerMessages(p)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	exitCode := 0
	select {
	case <-signalChan:
		playerCtx.Logger().Print("*********** INTERRUPTION SIGNAL ****************")
		talker.Close()
		stopPlayer()
		exitCode = 1
	case <-talkerCtx.Done():
		playerCtx.Logger().Printf("was connection lost: %s", talkerCtx.Err())
		stopPlayer()
		exitCode = 2
	case <-playerCtx.Done():
		playerCtx.Logger().Printf("player stopped: %s", playerCtx.Err())

	}
	os.Exit(exitCode)
}

// stopToPlay stop the player to play
func (p *Gamer) stopToPlay(interrupted bool) {
	stopPlayer()
}

// onMessage is the callback function called when the game server sends a new message
func (p *Gamer) onMessage(msg GameMessage) {
	defer func() {
		if err := recover(); err != nil {
			playerCtx.Logger().Errorf("Panic processing new game message: %s", err)
			debug.PrintStack()
		}
	}()
	p.LastMsg = msg
	if p.OnMessage == nil {
		p.defaultOnMessage(msg)
	} else {
		p.OnMessage(msg)
	}

}

// defaultOnMessage is the default callback to process the new messages got from the game server
func (p *Gamer) defaultOnMessage(msg GameMessage) {
	switch msg.Type {
	case orders.WELCOME:
		playerCtx.Logger().Info("Accepted by the game server")
	case orders.ANNOUNCEMENT:
		if p.OnAnnouncement == nil {
			playerCtx.Logger().Fatal("the player must implement the `OnAnnouncement` method")
		} else {
			turnTx := playerCtx.CreateTurnContext(msg)
			p.OnAnnouncement(turnTx)
		}
	case orders.RIP:
		playerCtx.Logger().Warn("the server died")
		p.stopToPlay(true)
	}
}

// SendOrders sends a list of orders to the game server, and includes a message to them (only displayed in the game server log)
func (p *Gamer) SendOrders(message string, ordersList ...orders.Order) {
	msg := PlayerMessage{
		orders.ORDER,
		ordersList,
		message,
	}
	stringed, err := json.Marshal(msg)
	if err != nil {
		playerCtx.Logger().Errorf("Fail generating JSON: %s", err.Error())
		return
	}

	err = p.Talker.Send(stringed)
	if err != nil {
		playerCtx.Logger().Errorf("Fail on sending message: %s", err.Error())
		return
	}
}
