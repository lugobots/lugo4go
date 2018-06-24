package Game

import (
	"net/url"
	"fmt"
	"github.com/makeitplay/commons/talk"
	"github.com/makeitplay/commons/BasicTypes"
	"encoding/json"
	"github.com/makeitplay/commons"
	"strconv"
	"os"
)

func (p *Player) initializeCommunicator() {
	uri := new(url.URL)
	uri.Scheme = "ws"
	uri.Host = "localhost:8080"
	uri.Path = fmt.Sprintf("/announcements/%s/%s", p.config.Uuid, p.TeamPlace)
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
		panic(err)
	} else {
		commons.RegisterCleaner("Websocket connection", func(interrupted bool) {
			p.talker.CloseConnection()
		})
	}
}

func (p *Player) onMessage(msg GameMessage) {
	p.LastMsg = msg
	if p.OnMessage == nil {
		p.defaultOnMessage(msg)
	} else {
		p.OnMessage(msg)
	}

}

func (p *Player) defaultOnMessage(msg GameMessage){
	switch msg.Type {
	case BasicTypes.WELCOME:
		commons.LogInfo("Accepted by the game server")
		if myId, ok := msg.Data["id"]; ok {
			i, err := strconv.Atoi(myId)
			if err != nil {
				commons.LogError("Invalid player id: %v", err.Error())
				panic("Invalid player id")
			}
			p.Id = i
			commons.LogDebug("-----------------------------Setou o ID %d", i)
		} else {
			commons.LogError("Player id missing in the welcome message")
			panic("Player id missing in the welcome message")
		}

		p.Number = p.FindMyStatus(msg.GameInfo).Number
	case BasicTypes.ANNOUNCEMENT:
		if p.OnAnnouncement == nil {
			panic("the player must implement the `OnAnnouncement` method")
		} else {
			commons.LogDebug("-----------------------------CHAMOU")
			p.OnAnnouncement(msg)
		}
	case BasicTypes.RIP:
		commons.LogError("The server has stopped :/")
		commons.Cleanup(true)
		os.Exit(0)
	}
}
