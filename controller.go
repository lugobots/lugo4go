package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/orders"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/arena/talk"
	"github.com/makeitplay/arena/units"
	"github.com/makeitplay/commons/Units"
	"github.com/sirupsen/logrus"
	"net/url"
	"time"
)

type Controller interface {
	SendOrders(place arena.TeamPlace, number arena.PlayerNumber, orderList []orders.Order)
	ToggleDebugMode()
	NextTurn() (newState GameMessage, err error)
	LoadArrangement(name string)
	ResetGame() //score and time
}

type commandArgs map[string]interface{}

type debugCmd struct {
	Cmd  string
	Args map[string]commandArgs
}

type ctrlTeam map[arena.PlayerNumber]Gamer

type controller struct {
	browser         talk.Talker
	teams           map[arena.TeamPlace]ctrlTeam
	lastState       *GameMessage
	listenerStopper context.CancelFunc
}

func NewTestController(ctx context.Context, confg Configuration) (context.Context, Controller, error) {
	subCtx, stop := context.WithCancel(ctx)

	uri := new(url.URL)
	uri.Scheme = "ws"
	uri.Host = fmt.Sprintf("%s:%s", confg.WSHost, confg.WSPort)
	uri.Path = fmt.Sprintf("/ws/%s", confg.UUID)

	talker := talk.NewTalker(logrus.NewEntry(logrus.New()))
	_, err := talker.Connect(subCtx, *uri, arena.PlayerSpecifications{})
	if err != nil {
		return nil, nil, fmt.Errorf("fail on opening the websocket connection: %s", err)
	}

	ctrl := &controller{
		browser: talker,
		teams: map[arena.TeamPlace]ctrlTeam{
			arena.HomeTeam: {},
			arena.AwayTeam: {},
		},
		listenerStopper: func() {

		},
	}
	ctrl.ToggleDebugMode()
	go ctrl.ctrlServerListenner()

	for i := 1; i <= 11; i++ {

		if err := ctrl.addPlayer(subCtx, confg, arena.HomeTeam, i); err != nil {
			stop()
		}
		if err := ctrl.addPlayer(subCtx, confg, arena.AwayTeam, i); err != nil {
			stop()
		}
	}
	msg, err := ctrl.waitListeningState()
	ctrl.lastState = &msg
	return subCtx, ctrl, err
}

func (c *controller) ToggleDebugMode() {
	nextTurnDebugMsg := debugCmd{
		Cmd: "clean-breakpoint",
	}
	jsonMsg, _ := json.Marshal(nextTurnDebugMsg)
	c.lastState = nil
	c.browser.Send(jsonMsg)
}

func (c *controller) NextTurn() (newState GameMessage, err error) {
	nextTurnDebugMsg := debugCmd{
		Cmd: "next-turn",
	}
	jsonMsg, _ := json.Marshal(nextTurnDebugMsg)
	c.lastState = nil
	c.browser.Send(jsonMsg)
	return c.waitListeningState()
}

func (c *controller) SendOrders(place arena.TeamPlace, number arena.PlayerNumber, orderList []orders.Order) {
	gamer := c.teams[place][number]
	gamer.SendOrders("debug", orderList...)
}

func (c *controller) LoadArrangement(name string) {
	panic("implement me")
}

func (c *controller) ResetGame() {
	panic("implement me")
}

func (c *controller) ctrlServerListenner() {
	for {
		select {
		case bytes := <-c.browser.Listen():
			logrus.WithTime(time.Now()).Infof("MSGSGS")
			var msg GameMessage
			err := json.Unmarshal(bytes, &msg)
			if err != nil {
				logrus.Errorf("Fail on convert wb message: %s (%s)", err.Error(), bytes)
			} else {
				logrus.WithTime(time.Now()).Infof("MSG %s", msg.State)
				if msg.State == arena.Listening {
					c.lastState = &msg
					c.listenerStopper()
				}
			}
		case connError := <-c.browser.ListenInterruption():
			logrus.Infof("ws connection lost: %s", connError)
			return
		}
	}
}
func (c *controller) addPlayer(subCtx context.Context, confg Configuration, teamPlace arena.TeamPlace, playerNumber int) error {
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)

	gamer := Gamer{}
	gamer.OnAnnouncement = c.msgReceiver
	initialPosition := physics.Point{
		PosX: arena.FieldCenter.PosX - (playerNumber * Units.PlayerSize),
		PosY: 1,
	}
	if teamPlace == arena.AwayTeam {
		initialPosition.PosX = units.FieldWidth - initialPosition.PosX
	}

	confg.PlayerNumber = arena.PlayerNumber(fmt.Sprintf("%d", playerNumber))
	confg.TeamPlace = teamPlace

	if err := gamer.Play(initialPosition, &confg); err != nil {
		return err
	}

	if _, ok := c.teams[teamPlace]; !ok {
		c.teams[teamPlace] = ctrlTeam{}
	}
	c.teams[teamPlace][confg.PlayerNumber] = gamer
	return nil
}

func (c *controller) msgReceiver(turnTx TurnContext) {
	//logrus.Info(turnTx.GameMsg().Turn())
	//logrus.Info(turnTx.GameMsg().State)
	//logrus.WithTime(time.Now()).Infof("I got it: %s", turnTx.GameMsg().State)
}
func (c *controller) waitListeningState() (GameMessage, error) {
	var listennerCtx context.Context
	listennerCtx, c.listenerStopper = context.WithTimeout(context.Background(), 10*time.Second)
	logrus.WithTime(time.Now()).Info("WAIT")
	<-listennerCtx.Done()
	if listennerCtx.Err() == context.DeadlineExceeded {
		return GameMessage{}, listennerCtx.Err()
	}
	return *c.lastState, nil
}
