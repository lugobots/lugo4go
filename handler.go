package lugo4go

import (
	"context"
	"github.com/lugobots/lugo4go/v2/lugo"
	"github.com/lugobots/lugo4go/v2/pkg/field"
)

func NewHandler(bot Bot, sender OrderSender, logger Logger, playerNumber uint32, side lugo.Team_Side) *Handler {
	return &Handler{
		Logger:       logger,
		Sender:       sender,
		PlayerNumber: playerNumber,
		Side:         side,
		Bot:          bot,
	}
}

// Handler is a Lugo4go client handler that allow you to create an interface to follow a basic strategy based on team
// states.
type Handler struct {
	Logger       Logger
	Sender       OrderSender
	PlayerNumber uint32
	Side         lugo.Team_Side
	Bot          Bot
}

func (h *Handler) Handle(ctx context.Context, snapshot *lugo.GameSnapshot) {
	var err error
	var state PlayerState

	if snapshot == nil {
		h.Logger.Errorf("error processing turn: %s", ErrNilSnapshot)
		return
	}

	state, err = DefineMyState(snapshot, h.PlayerNumber, h.Side)
	if err != nil {
		h.Logger.Errorf("error processing turn %d: %s", snapshot.Turn, err)
		return
	}
	if field.GoalkeeperNumber == h.PlayerNumber {
		err = h.Bot.AsGoalkeeper(ctx, wrapSender(h.Sender, snapshot.Turn), snapshot, state)
	} else {
		switch state {
		case Supporting:
			err = h.Bot.OnSupporting(ctx, wrapSender(h.Sender, snapshot.Turn), snapshot)
		case HoldingTheBall:
			err = h.Bot.OnHolding(ctx, wrapSender(h.Sender, snapshot.Turn), snapshot)
		case Defending:
			err = h.Bot.OnDefending(ctx, wrapSender(h.Sender, snapshot.Turn), snapshot)
		case DisputingTheBall:
			err = h.Bot.OnDisputing(ctx, wrapSender(h.Sender, snapshot.Turn), snapshot)
		}
	}
	if err != nil {
		h.Logger.Errorf("error processing turn %d: %s", snapshot.Turn, err)
	}
}

func wrapSender(sender OrderSender, turn uint32) senderWrapper {
	return senderWrapper{
		sender: sender,
		turn:   turn,
	}
}

type senderWrapper struct {
	sender OrderSender
	turn   uint32
}

func (s senderWrapper) Send(ctx context.Context, orders []lugo.PlayerOrder, debugMsg string) (*lugo.OrderResponse, error) {
	return s.sender.Send(ctx, s.turn, orders, debugMsg)
}

// PlayerState defines states specific for players
type PlayerState string

const (
	// Supporting identifies the player supporting the teammate
	Supporting PlayerState = "supporting"
	// HoldingTheBall identifies the player holding	the ball
	HoldingTheBall PlayerState = "holding"
	// Defending identifies the player defending against the opponent team
	Defending PlayerState = "defending"
	// DisputingTheBall identifies the player disputing the ball
	DisputingTheBall PlayerState = "disputing"
)

func DefineMyState(snapshot *lugo.GameSnapshot, playerNumber uint32, side lugo.Team_Side) (PlayerState, error) {
	if snapshot == nil || snapshot.Ball == nil {
		return "", ErrNoBall
	}

	me := field.GetPlayer(snapshot, side, playerNumber)
	if me == nil {
		return "", ErrPlayerNotFound
	}

	ballHolder := snapshot.Ball.Holder

	if ballHolder == nil {
		return DisputingTheBall, nil
	} else if ballHolder.TeamSide == side {
		if ballHolder.Number == playerNumber {
			return HoldingTheBall, nil
		}
		return Supporting, nil
	}
	return Defending, nil
}