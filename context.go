package client

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

type GamerCtx interface {
	context.Context
	Logger() *logrus.Entry
}

func NewGamerContext(ctx context.Context, config *Configuration) (GamerCtx, context.CancelFunc) {
	ctx, cancelFunc := context.WithCancel(ctx)
	return &ourCtx{
		mainCtx: ctx,
		log:     logrus.WithField("player", fmt.Sprintf("%s-%s", config.TeamPlace, config.PlayerNumber)),
	}, cancelFunc
}

type ourCtx struct {
	mainCtx context.Context
	log     *logrus.Entry
}

func (o *ourCtx) Deadline() (deadline time.Time, ok bool) {
	return o.mainCtx.Deadline()
}

func (o *ourCtx) Done() <-chan struct{} {
	return o.mainCtx.Done()
}

func (o *ourCtx) Err() error {
	return o.mainCtx.Err()
}

func (o *ourCtx) Value(key interface{}) interface{} {
	return o.mainCtx.Value(key)
}

func (o *ourCtx) Logger() *logrus.Entry {
	return o.log
}
