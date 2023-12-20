package main

import (
	"context"
	"log"

	clientGo "github.com/lugobots/lugo4go/v3"
	"github.com/lugobots/lugo4go/v3/mapper"
	"github.com/lugobots/lugo4go/v3/proto"
)

func main() {

	connectionStarter, err := clientGo.NewTurnHandlerConfig()
	if err != nil {
		log.Fatalf("failed to load the bot configuration: %s", err)
	}

	//
	// Optional: define your own field mapper
	//
	//playerMapper, err := mapper.NewMapper(32, 15, connectionStarter.Config.TeamSide)
	//if err != nil {
	//	log.Fatalf("failed to create a field mapper: %s", err)
	//}
	//connectionStarter.FieldMapper = playerMapper
	//

	if err := connectionStarter.Run(&TurnHandler{
		FieldMapper: connectionStarter.FieldMapper,
		config:      connectionStarter.Config,
	}); err != nil {
		log.Fatalf("bot stopped: %s", err)
	}
}

type TurnHandler struct {
	FieldMapper mapper.Mapper
	config      clientGo.Config
}

func (t *TurnHandler) Handle(ctx context.Context, snapshot clientGo.SnapshotInspector) ([]proto.PlayerOrder, string, error) {
	//TODO implement me
	panic("implement me")
}
