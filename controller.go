package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/orders"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/arena/talk"
	"github.com/makeitplay/arena/units"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Controller interface {
	SendOrders(place arena.TeamPlace, number arena.PlayerNumber, orderList []orders.Order)
	//ToggleDebugMode()
	NextTurn() (newState GameMessage, err error)
	LoadArrangement(name string) (newState GameMessage, err error)
	SetBallProperties(v physics.Velocity, position physics.Point) (newState GameMessage, err error)
	SetPlayerPos(place arena.TeamPlace, number arena.PlayerNumber, position physics.Point) (newState GameMessage, err error)
	SetGameTurn(turn int) (newState GameMessage, err error)
	ResetScore() (newState GameMessage, err error)
	SetFrameInterval(time time.Duration)
	GetGamerCtx(place arena.TeamPlace, number arena.PlayerNumber) (ctx TurnContext, err error)
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
	config          Configuration
	intervalTime    time.Duration
}

func (c *controller) GetGamerCtx(place arena.TeamPlace, number arena.PlayerNumber) (ctx TurnContext, err error) {
	if c.lastState == nil {
		return nil, fmt.Errorf("no game state available")
	}

	team, ok := c.teams[place]
	if !ok {
		return nil, fmt.Errorf("unknown team")
	}
	playerGamer, ok := team[number]
	if !ok {
		return nil, fmt.Errorf("unknown player")
	}
	return playerGamer.ctx.CreateTurnContext(*c.lastState), nil
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
		intervalTime: 0,
		config:       confg,
	}
	go ctrl.ctrlServerListenner()

	msg, err := ctrl.doAndWaitListening(func() {
		for i := 1; i <= 11; i++ {

			if err := ctrl.addPlayer(subCtx, confg, arena.HomeTeam, i); err != nil {
				stop()
			}
			if err := ctrl.addPlayer(subCtx, confg, arena.AwayTeam, i); err != nil {
				stop()
			}
		}
	})
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

func (c *controller) SetFrameInterval(time time.Duration) {
	c.intervalTime = time
}
func (c *controller) NextTurn() (newState GameMessage, err error) {
	nextTurnDebugMsg := debugCmd{
		Cmd: "play",
	}
	return c.doAndWaitListening(func() {
		c.sendDebugMsg(nextTurnDebugMsg)
	})
}

func (c *controller) sendDebugMsg(msg debugCmd) (newState GameMessage, err error) {
	uri := new(url.URL)
	uri.Scheme = "http"
	uri.Host = fmt.Sprintf("%s:%s", c.config.WSHost, c.config.WSPort)
	uri.Path = fmt.Sprintf("/%s/debug", c.config.UUID)
	jsonValue, _ := json.Marshal(msg)

	resp, err := http.Post(uri.String(), "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return GameMessage{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		var merda GameMessage
		err := json.Unmarshal(body, &merda)
		if err != nil {
			return GameMessage{}, err
		}
		newState = merda
	}
	return
}

func (c *controller) SetGameTurn(turn int) (newState GameMessage, err error) {
	setProps := debugCmd{
		Cmd: "set-game-properties",
	}
	setProps.Args = map[string]commandArgs{
		"props": map[string]interface{}{
			"turn": turn,
		},
	}
	return c.sendDebugMsg(setProps)
}
func (c *controller) SetPlayerPos(place arena.TeamPlace, number arena.PlayerNumber, position physics.Point) (newState GameMessage, err error) {
	setProps := debugCmd{
		Cmd: "rearrange",
	}
	setProps.Args = map[string]commandArgs{
		fmt.Sprintf("%s-%s", place, number): map[string]interface{}{
			"x": position.PosX,
			"y": position.PosY,
		},
	}
	return c.sendDebugMsg(setProps)
}
func (c *controller) SetBallProperties(v physics.Velocity, position physics.Point) (newState GameMessage, err error) {
	setPos := debugCmd{
		Cmd: "set-ball",
	}
	setPos.Args = map[string]commandArgs{
		"coords": map[string]interface{}{
			"x": position.PosX,
			"y": position.PosY,
		},
	}

	setVel := debugCmd{
		Cmd: "set-ball",
	}
	setVel.Args = map[string]commandArgs{
		"velocity": map[string]interface{}{
			"v": v.Speed,
			"x": v.Direction.GetX(),
			"y": v.Direction.GetY(),
		},
	}

	if newState, err := c.sendDebugMsg(setPos); err != nil {
		return newState, err
	}

	return c.sendDebugMsg(setVel)
}

func (c *controller) SendOrders(place arena.TeamPlace, number arena.PlayerNumber, orderList []orders.Order) {
	gamer := c.teams[place][number]
	gamer.SendOrders("debug", orderList...)
	time.Sleep(50 * time.Millisecond) //there is a little chance of the next instruction be executed before the message be sent
}

func (c *controller) LoadArrangement(name string) (newState GameMessage, err error) {
	setProps := debugCmd{
		Cmd: "load-positions",
	}
	setProps.Args = map[string]commandArgs{
		"file": map[string]interface{}{
			"name": name,
		},
	}
	return c.sendDebugMsg(setProps)
}

func (c *controller) ResetScore() (newState GameMessage, err error) {
	setProps := debugCmd{
		Cmd: "set-game-properties",
	}
	setProps.Args = map[string]commandArgs{
		"props": map[string]interface{}{
			"score": map[string]int{"home": 0, "away": 0},
		},
	}
	return c.sendDebugMsg(setProps)
}

func (c *controller) ctrlServerListenner() {
	for {
		select {
		case msgInBytes := <-c.browser.Listen():
			var msg GameMessage
			err := json.Unmarshal(msgInBytes, &msg)
			if err != nil {
				logrus.Errorf("Fail on convert wb message: %s (%s)", err.Error(), msgInBytes)
			} else {
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
	gamer := Gamer{}
	gamer.OnAnnouncement = c.msgReceiver
	//	gamer.LogLevel = logrus.PanicLevel

	initialPosition := physics.Point{
		PosX: arena.FieldCenter.PosX - (playerNumber * units.PlayerSize),
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

// WARNING: if this method be called passing a cb that does not change the game server state, it can hang forever
func (c *controller) doAndWaitListening(cb func()) (GameMessage, error) {
	var listennerCtx context.Context
	listennerCtx, c.listenerStopper = context.WithTimeout(context.Background(), 10*time.Second)

	c.lastState = nil
	cb()

	<-listennerCtx.Done()
	if listennerCtx.Err() == context.DeadlineExceeded {
		return GameMessage{}, listennerCtx.Err()
	}
	if c.intervalTime > 0 {
		time.Sleep(c.intervalTime)
	}
	return *c.lastState, nil
}
