package client

import (
	"encoding/json"
	"fmt"
	"github.com/makeitplay/commons"
	"github.com/makeitplay/commons/BasicTypes"
	"github.com/makeitplay/commons/talk"
	"net/url"
	"os"
	"runtime/debug"
)

// initializeCommunicator initialize a communication with the game server
func (p *Player) initializeCommunicator() bool {
	uri := new(url.URL)
	uri.Scheme = "ws"
	uri.Host = fmt.Sprintf("%s:%s", p.config.WSHost, p.config.WSPort)
	uri.Path = fmt.Sprintf("/announcements/%s/%s", p.config.UUID, p.TeamPlace)
	p.talker = talk.NewTalkChannel(*uri, BasicTypes.PlayerSpecifications{
		Number:        p.Number,
		InitialCoords: p.Coords,
	})

	err := p.talker.OpenConnection(func(bytes []byte) {
		var msg GameMessage
		err := json.Unmarshal(bytes, &msg)
		if err != nil {
			commons.LogError("Fail on convert wb message: %s (%s)", err.Error(), bytes)
		} else {
			p.onMessage(msg)
		}
	})

	if err != nil {
		commons.LogError("Fail on opening the websocket connection: %s", err)
		return false
	}
	commons.RegisterCleaner("Websocket connection", func(interrupted bool) {
		p.talker.CloseConnection()
	})
	return true
}

// onMessage is the callback function called when the game server sends a new message
func (p *Player) onMessage(msg GameMessage) {
	defer func() {
		if err := recover(); err != nil {
			commons.LogError("Panic processing new game message: %s", err)
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
		commons.LogInfo("Accepted by the game server")
		myStatus := p.FindMyStatus(msg.GameInfo)
		p.Number = myStatus.Number
	case BasicTypes.ANNOUNCEMENT:
		if p.OnAnnouncement == nil {
			panic("the player must implement the `OnAnnouncement` method")
		} else {
			p.OnAnnouncement(msg)
		}
	case BasicTypes.RIP:
		commons.LogError("The server has stopped :/")
		commons.Cleanup(true)
		os.Exit(0)
	}
}
