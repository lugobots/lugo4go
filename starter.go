package lugo4go

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"go.uber.org/zap"

	"github.com/lugobots/lugo4go/v3/mapper"
)

func NewTurnHandlerConfig() (*Starter, error) {
	playerConfig, defaultMapper, logger, err := DefaultInitBundle()
	if err != nil {
		return nil, fmt.Errorf("could not init default config or logger: %w", err)
	}

	return &Starter{
		Config:      playerConfig,
		FieldMapper: defaultMapper,
		Logger:      logger,
	}, nil
}

type Starter struct {
	Config      Config
	FieldMapper mapper.Mapper
	Logger      *zap.SugaredLogger
}

func (s *Starter) Run(handler RawBot) error {
	player, err := NewClient(s.Config, s.Logger)
	if err != nil {
		return fmt.Errorf("could not init the client: %w", err)
	}

	ctx, stop := context.WithCancel(context.Background())
	go func() {
		defer stop()
		if err := player.Play(handler); err != nil {
			s.Logger.Errorf("bot stopped with an error: %s", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	select {
	case <-ctx.Done():
	case <-signalChan:
		if err := player.Stop(); err != nil {
			s.Logger.Errorf("error stopping bot: %s", err)
		}
	}
	s.Logger.Debug("bot stopped")
	return nil
}
