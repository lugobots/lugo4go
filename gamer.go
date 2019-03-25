package client

import (
	"context"
	"encoding/json"
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/orders"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/arena/talk"
	"runtime/debug"
)

type Responder interface {
	SendOrders(message string, ordersList ...orders.Order)
}

type Gamer struct {
	OnMessage      func(msg GameMessage)
	OnAnnouncement func(turnTx TurnContext)
	config         *Configuration
	Talker         talk.Talker
	LastMsg        GameMessage
	stop           context.CancelFunc
	ctx            GamerCtx
}

func listenServerMessages(gamer *Gamer) {
	for {
		select {
		case bytes := <-gamer.Talker.Listen():
			var msg GameMessage
			err := json.Unmarshal(bytes, &msg)
			if err != nil {
				gamer.ctx.Logger().Errorf("Fail on convert wb message: %s (%s)", err.Error(), bytes)
			} else {
				gamer.onMessage(msg)
			}
		case connError := <-gamer.Talker.ListenInterruption():
			if gamer.LastMsg.State == arena.Over {
				gamer.ctx.Logger().Info("game over")
				gamer.StopToPlay(false)
			} else {
				gamer.ctx.Logger().Errorf("ws connection lost: %s", connError.Error())
				gamer.StopToPlay(true)
			}
			return
		}
	}
}

// Play make the player start to play
func (p *Gamer) Play(initialPosition physics.Point, configuration *Configuration) error {
	p.ctx, p.stop = NewGamerContext(context.Background(), configuration)

	p.config = configuration
	talkerCtx, talker, err := TalkerSetup(p.ctx, configuration, initialPosition)
	if err != nil {
		return err
	}

	p.Talker = talker
	go listenServerMessages(p)

	go func() {
		select {
		case <-talkerCtx.Done():
			p.ctx.Logger().Printf("was connection lost: %s", talkerCtx.Err())
			p.stop()
		case <-p.ctx.Done():
			talker.Close()
			p.ctx.Logger().Printf("player stopped: %s", p.ctx.Err())

		}
	}()
	return nil
}

// StopToPlay stop the player to play
func (p *Gamer) StopToPlay(interrupted bool) {
	p.stop()
}

// onMessage is the callback function called when the game server sends a new message
func (p *Gamer) onMessage(msg GameMessage) {
	defer func() {
		if err := recover(); err != nil {
			p.ctx.Logger().Errorf("Panic processing new game message: %s", err)
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
		p.ctx.Logger().Info("Accepted by the game server")
	case orders.ANNOUNCEMENT:
		if p.OnAnnouncement == nil {
			p.ctx.Logger().Fatal("the player must implement the `OnAnnouncement` method")
		} else {
			turnTx := p.ctx.CreateTurnContext(msg)
			p.OnAnnouncement(turnTx)
		}
	case orders.RIP:
		p.ctx.Logger().Warn("the server died")
		p.StopToPlay(true)
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
		p.ctx.Logger().Errorf("Fail generating JSON: %s", err.Error())
		return
	}

	err = p.Talker.Send(stringed)
	if err != nil {
		p.ctx.Logger().Errorf("Fail on sending message: %s", err.Error())
		return
	}
}
