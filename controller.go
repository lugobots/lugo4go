package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/talk"
	"github.com/sirupsen/logrus"
	"net/url"
)

type Controller interface {
	NextTurn()
}

type commandArgs map[string]interface{}

type debugCmd struct {
	Cmd  string
	Args map[string]commandArgs
}

type controller struct {
	talker talk.Talker
}

func NewTestController(ctx context.Context, host, port, uuid string) (context.Context, Controller, error) {
	uri := new(url.URL)
	uri.Scheme = "ws"
	uri.Host = fmt.Sprintf("%s:%s", host, port)
	uri.Path = fmt.Sprintf("/ws/%s", uuid)

	talker := talk.NewTalker(logrus.NewEntry(logrus.New()))
	talkerCtx, err := talker.Connect(ctx, *uri, arena.PlayerSpecifications{})
	if err != nil {
		return nil, nil, fmt.Errorf("fail on opening the websocket connection: %s", err)
	}
	ctrl := &controller{
		talker: talker,
	}
	go ctrl.ctrlServerListenner()

	return talkerCtx, ctrl, nil
}

func (c *controller) NextTurn() {
	nextTurnDebugMsg := debugCmd{
		Cmd: "next-turn",
		//Cmd: "clean-breakpoint",
	}
	jsonMsg, _ := json.Marshal(nextTurnDebugMsg)
	c.talker.Send(jsonMsg)
}

func (c *controller) ctrlServerListenner() {
	for {
		select {
		case bytes := <-c.talker.Listen():
			var msg GameMessage
			err := json.Unmarshal(bytes, &msg)
			if err != nil {
				logrus.Errorf("Fail on convert wb message: %s (%s)", err.Error(), bytes)
			} else {
				logrus.Info(msg.State)
			}
		case connError := <-c.talker.ListenInterruption():
			logrus.Infof("ws connection lost: %s", connError)
			return
		}
	}
}
