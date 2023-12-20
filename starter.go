package lugo4go

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/lugobots/lugo4go/v3/mapper"
)

func NewTurnHandlerConfig() (*Starter, error) {
	playerConfig, _, err := DefaultInitBundle()
	if err != nil {
		return nil, fmt.Errorf("could not init default config or logger: %w", err)
	}
	// defining default initial player position
	tempMapper, _ := mapper.NewMapper(11, 7, playerConfig.TeamSide)

	region, _ := tempMapper.GetRegion(FieldMap[playerConfig.Number].Col, FieldMap[playerConfig.Number].Row)

	// just creating a position for example purposes
	playerConfig.InitialPosition = region.Center()

	return &Starter{
		Config:      playerConfig,
		FieldMapper: tempMapper,
	}, nil
}

type Starter struct {
	Config      Config
	FieldMapper mapper.Mapper
}

func (s *Starter) Run(handler TurnHandler) error {
	player, err := NewClient(s.Config)
	if err != nil {
		return fmt.Errorf("could not init the client: %w", err)
	}

	ctx, stop := context.WithCancel(context.Background())
	go func() {
		defer stop()
		if err := player.Play(handler); err != nil {
			log.Printf("bot stopped with an error: %s", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	select {
	case <-ctx.Done():
	case <-signalChan:
		if err := player.Stop(); err != nil {
			log.Printf("error stopping bot: %s", err)
		}
	}
	return nil
}

var FieldMap = map[uint32]struct {
	Col int
	Row int
}{
	2:  {Col: 1, Row: 1},
	3:  {Col: 1, Row: 3},
	4:  {Col: 1, Row: 4},
	5:  {Col: 1, Row: 6},
	6:  {Col: 2, Row: 2},
	7:  {Col: 2, Row: 3},
	8:  {Col: 2, Row: 4},
	9:  {Col: 2, Row: 5},
	10: {Col: 3, Row: 3},
	11: {Col: 3, Row: 4},
}
