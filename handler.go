package lugo4go

import (
	"context"

	"github.com/lugobots/lugo4go/v3/proto"
	"github.com/lugobots/lugo4go/v3/specs"
)

func hewTurnHandler(bot Bot, sender OrderSender, logger Logger, playerNumber uint32, side proto.Team_Side) *handler {
	return &handler{
		Logger:       logger,
		Sender:       sender,
		PlayerNumber: playerNumber,
		Side:         side,
		Bot:          bot,
	}
}

// handler is a Lugo4go client handler that allow you to create an interface to follow a basic strategy based on team
// states.
type handler struct {
	Logger       Logger
	Sender       OrderSender
	PlayerNumber uint32
	Side         proto.Team_Side
	Bot          Bot
}

func (h *handler) Handle(ctx context.Context, snapshot *proto.GameSnapshot) {
	var err error
	var state PlayerState

	if snapshot == nil {
		h.Logger.Errorf("error processing turn: %s", ErrNilSnapshot)
		return
	}

	state, err = defineMyState(snapshot, h.PlayerNumber, h.Side)
	if err != nil {
		h.Logger.Errorf("error processing turn %d: %s", snapshot.Turn, err)
		return
	}
	// TODO bad practice - create a SnapshotToolMaker to allow it to be created externally
	snapshotTools, err := newInspector(h.Side, int(h.PlayerNumber), snapshot)
	if err != nil {
		h.Logger.Errorf("failed to create an inspector for the game snapshot: %s", err)
		return
	}

	var orders []proto.PlayerOrder
	var debugMsg string
	if specs.GoalkeeperNumber == h.PlayerNumber {
		orders, debugMsg, err = h.Bot.AsGoalkeeper(ctx, snapshotTools, state)
	} else {
		switch state {
		case Supporting:
			orders, debugMsg, err = h.Bot.OnSupporting(ctx, snapshotTools)
		case HoldingTheBall:
			orders, debugMsg, err = h.Bot.OnHolding(ctx, snapshotTools)
		case Defending:
			orders, debugMsg, err = h.Bot.OnDefending(ctx, snapshotTools)
		case DisputingTheBall:
			orders, debugMsg, err = h.Bot.OnDisputing(ctx, snapshotTools)
		}
	}
	if err != nil {
		h.Logger.Errorf("error processing turn %d: %s", snapshot.Turn, err)
		return
	}
	resp, errSend := h.Sender.Send(ctx, snapshot.Turn, orders, debugMsg)
	if errSend != nil {
		h.Logger.Errorf("error sending orders to turn %d: %s", snapshot.Turn, errSend)
		return
	} else if resp.Code != proto.OrderResponse_SUCCESS {
		h.Logger.Errorf("order not sent during turn %d: %s", snapshot.Turn, resp.String())
		return
	}
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

func defineMyState(snapshot *proto.GameSnapshot, playerNumber uint32, side proto.Team_Side) (PlayerState, error) {
	if snapshot == nil || snapshot.Ball == nil {
		return "", ErrNoBall
	}

	myTeam := snapshot.HomeTeam
	if side == proto.Team_AWAY {
		myTeam = snapshot.AwayTeam
	}

	var me *proto.Player
	for _, player := range myTeam.GetPlayers() {
		if player.Number == playerNumber {
			me = player
			break
		}
	}

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
