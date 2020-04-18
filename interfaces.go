package lugo4go

import (
	"context"
	"github.com/lugobots/lugo4go/v2/lugo"
)

type TurnData struct {
	Me       *lugo.Player
	Snapshot *lugo.GameSnapshot
}

type TurnHandler interface {
	Handle(ctx context.Context, snapshot *lugo.GameSnapshot) // (orders []lugo.PlayerOrder, debugMsg string, err error)

}
