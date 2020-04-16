package coach

import (
	"context"
	"github.com/lugobots/lugo4go/v2/field"
	"github.com/lugobots/lugo4go/v2/lugo"
)

func NewHandler(bot Bot, sender OrderSender, logger lugo.Logger, playerNumber uint32, side lugo.Team_Side) *Handler {
	return &Handler{
		Logger:       logger,
		Sender:       sender,
		PlayerNumber: playerNumber,
		Side:         side,
		Bot:          bot,
	}
}

type Handler struct {
	Logger       lugo.Logger
	Sender       OrderSender
	PlayerNumber uint32
	Side         lugo.Team_Side
	Bot          Bot
}

func (h *Handler) Handle(ctx context.Context, snapshot *lugo.GameSnapshot) {
	var err error
	var state PlayerState

	if field.GoalkeeperNumber == h.PlayerNumber {
		err = h.Bot.AsGoalkeeper(ctx, h.Sender, snapshot)
	} else {
		state, err = DefineMyState(snapshot, h.PlayerNumber, h.Side)
		switch state {
		case Supporting:
			err = h.Bot.OnSupporting(ctx, h.Sender, snapshot)
		case HoldingTheBall:
			err = h.Bot.OnHolding(ctx, h.Sender, snapshot)
		case Defending:
			err = h.Bot.OnDefending(ctx, h.Sender, snapshot)
		case DisputingTheBall:
			err = h.Bot.OnDisputing(ctx, h.Sender, snapshot)
		}
	}
	if err != nil {
		h.Logger.Errorf("error processing turn %d: %s", snapshot.Turn, err)
	}
}
