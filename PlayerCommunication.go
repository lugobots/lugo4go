package client

import (
	"encoding/json"
	"fmt"
	"github.com/makeitplay/arena/BasicTypes"
	"github.com/makeitplay/arena/GameState"
	"github.com/makeitplay/arena/talk"
	"github.com/sirupsen/logrus"
	"net/url"
	"runtime/debug"
)

// initializeCommunicator initialize a communication with the game server
func (p *Player) initializeCommunicator(logger *logrus.Entry) bool {
	uri := new(url.URL)
	uri.Scheme = "ws"
	uri.Host = fmt.Sprintf("%s:%s", p.config.WSHost, p.config.WSPort)
	uri.Path = fmt.Sprintf("/announcements/%s/%s", p.config.UUID, p.TeamPlace)

	p.talker = talk.NewTalker(logger, func(bytes []byte) {
		var msg GameMessage
		err := json.Unmarshal(bytes, &msg)
		if err != nil {
			p.logger.Errorf("Fail on convert wb message: %s (%s)", err.Error(), bytes)
		} else {
			p.onMessage(msg)
		}
	}, func() {
		if GameState.State(p.LastMsg.State) == GameState.Over {
			logger.Info("game over")
		}
		p.stopToPlay(true)
	})

	playerSpec := BasicTypes.PlayerSpecifications{
		Number:          p.Number,
		InitialCoords:   p.Coords,
		Token:           p.config.Token,
		ProtocolVersion: "1.0",
	}

	var err error
	if p.talkerCtx, err = p.talker.Connect(*uri, playerSpec); err != nil {
		logger.Errorf("Fail on opening the websocket connection: %s", err)
		return false
	}
	return true
}

// onMessage is the callback function called when the game server sends a new message
func (p *Player) onMessage(msg GameMessage) {
	defer func() {
		if err := recover(); err != nil {
			p.logger.Errorf("Panic processing new game message: %s", err)
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
	case BasicTypes.WELCOME:
		p.logger.Info("Accepted by the game server")
		myStatus := p.GetMyStatus(msg.GameInfo)
		p.Number = myStatus.Number
	case BasicTypes.ANNOUNCEMENT:
		if p.OnAnnouncement == nil {
			panic("the player must implement the `OnAnnouncement` method")
		} else {
			p.OnAnnouncement(msg)
		}
	case BasicTypes.RIP:
		p.logger.Warn("the server died")
		p.stopToPlay(true)
	}
}
