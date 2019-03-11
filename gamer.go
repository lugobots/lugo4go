package client

import (
	"os"
	"os/signal"
)

type Gamer interface {
	Start(ctx GamerCtx, configuration *Configuration)
	StopToPlay(interrupted bool)
}

type gamer struct {
	config *Configuration
}

func (g *gamer) Start(ctx GamerCtx, configuration *Configuration) {
	g.config = configuration

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	exitCode := 0
	select {
	case <-signalChan:
		ctx.Logger().Print("*********** INTERRUPTION SIGNAL ****************")
	case <-ctx.Done():
		ctx.Logger().Info("the main context has stopped")
	}
	os.Exit(exitCode)
}

func (*gamer) StopToPlay(interrupted bool) {
	panic("implement me")
}
