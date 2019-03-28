package client

import (
	"context"
	"fmt"
	"github.com/makeitplay/arena"
	"github.com/sirupsen/logrus"
	"time"
)

type GamerCtx interface {
	context.Context
	Logger() *logrus.Entry
	CreateTurnContext(msg GameMessage) TurnContext
}

type TurnContext interface {
	context.Context
	Logger() *logrus.Entry
	Player() *Player
	GameMsg() *GameMessage
}

func NewGamerContext(ctx context.Context, config *Configuration) (GamerCtx, context.CancelFunc) {
	ctx, cancelFunc := context.WithCancel(ctx)
	logger := logrus.New()
	logger.SetLevel(config.LogLevel)

	return &gameCtx{
		config:  config,
		mainCtx: ctx,
		log:     logger.WithField("player", fmt.Sprintf("%s-%s", config.TeamPlace, config.PlayerNumber)),
	}, cancelFunc
}

type gameCtx struct {
	config  *Configuration
	mainCtx context.Context
	log     *logrus.Entry
}

func (o *gameCtx) Deadline() (deadline time.Time, ok bool) {
	return o.mainCtx.Deadline()
}

func (o *gameCtx) Done() <-chan struct{} {
	return o.mainCtx.Done()
}

func (o *gameCtx) Err() error {
	return o.mainCtx.Err()
}

func (o *gameCtx) Value(key interface{}) interface{} {
	return o.mainCtx.Value(key)
}

func (o *gameCtx) Logger() *logrus.Entry {
	return o.log
}

func (o *gameCtx) CreateTurnContext(msg GameMessage) TurnContext {
	teamState := msg.GameInfo.HomeTeam
	if o.config.TeamPlace == arena.AwayTeam {
		teamState = msg.GameInfo.AwayTeam
	}
	var player *Player
	for _, playerInfo := range teamState.Players {
		if playerInfo.Number == o.config.PlayerNumber {
			player = playerInfo
		}
	}

	return &turnCtx{
		gameCtx: o,
		log:     o.log.WithField("turn", msg.GameInfo.Turn),
		msg:     &msg,
		player:  player, //remember that this value can be nil at the very first msgs
	}
}

type turnCtx struct {
	*gameCtx
	log    *logrus.Entry
	msg    *GameMessage
	player *Player
}

func (t *turnCtx) Logger() *logrus.Entry {
	return t.log
}

func (t *turnCtx) Player() *Player {
	return t.player
}

func (t *turnCtx) GameMsg() *GameMessage {
	return t.msg
}
