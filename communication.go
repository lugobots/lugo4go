package client

import (
	"encoding/json"
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/orders"
	"runtime/debug"
)

func listenServerMessages(player *Player) {
	for {
		select {
		case bytes := <-player.Talker.Listen():
			var msg GameMessage
			err := json.Unmarshal(bytes, &msg)
			if err != nil {
				playerCtx.Logger().Errorf("Fail on convert wb message: %s (%s)", err.Error(), bytes)
			} else {
				player.onMessage(msg)
			}
		case connError := <-player.Talker.ListenInterruption():
			if player.LastMsg.State == arena.Over {
				playerCtx.Logger().Info("game over")
			} else {
				playerCtx.Logger().Errorf("ws connection lost: %s", connError.Error())
			}
			player.stopToPlay(true)
			return
		}
	}
}

// onMessage is the callback function called when the game server sends a new message
func (p *Player) onMessage(msg GameMessage) {
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
func (p *Player) defaultOnMessage(msg GameMessage) {
	switch msg.Type {
	case orders.WELCOME:
		playerCtx.Logger().Info("Accepted by the game server")
		myStatus := p.GetMyStatus(msg.GameInfo)
		p.Number = myStatus.Number
	case orders.ANNOUNCEMENT:
		if p.OnAnnouncement == nil {
			playerCtx.Logger().Fatal("the player must implement the `OnAnnouncement` method")
		} else {
			p.OnAnnouncement(msg)
		}
	case orders.RIP:
		playerCtx.Logger().Warn("the server died")
		p.stopToPlay(true)
	}
}
