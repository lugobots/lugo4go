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
	"github.com/makeitplay/commons/GameState"
)

func (p *Player) initializeCommunicator() {
	uri := new(url.URL)
	uri.Scheme = "ws"
	uri.Host = "localhost:8080"
	uri.Path = fmt.Sprintf("/announcements/%s/%s", p.config.Uuid, p.TeamPlace)
	region := p.myRegion()
	p.talker = talk.NewTalkChannel(*uri, BasicTypes.PlayerSpecifications{
		Number:        p.Number,
		InitialCoords: region.InitialPosition(),
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
	p.lastMsg = msg
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
		} else {
			commons.LogError("Player id missing in the welcome message")
			panic("Player id missing in the welcome message")
		}
		p.updatePosition(p.lastMsg.GameInfo)
		p.Number = p.FindMyStatus(msg.GameInfo).Number
	case BasicTypes.ANNOUNCEMENT:
		commons.LogBroadcast("ANN %s", string(msg.State))
		switch GameState.State(msg.State) {
		case GameState.GETREADY:
			//p.updatePosition(p.lastMsg.GameInfo)
			//p.Number = p.FindMyStatus(msg.GameInfo).Number
		case GameState.LISTENING:
			p.updatePosition(p.lastMsg.GameInfo)
			p.state = p.determineMyState()
			commons.LogDebug("State: %s", p.state)
			p.TakeAnAction()
		}
	case BasicTypes.RIP:
		commons.LogError("The server has stopped :/")
		commons.Cleanup(true)
		os.Exit(0)
	}
}
