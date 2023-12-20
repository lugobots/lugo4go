package lugo4go

import (
	"context"
	"fmt"

	"github.com/lugobots/lugo4go/v3/proto"
	"github.com/lugobots/lugo4go/v3/specs"
)

func hewRawBotWrapper(bot Bot, logger Logger, playerNumber uint32, side proto.Team_Side) *rawBotWrapper {
	return &rawBotWrapper{
		Logger:       logger,
		PlayerNumber: playerNumber,
		Side:         side,
		Bot:          bot,
	}
}

// rawBotWrapper is a Lugo4go client that allow you to create an interface to follow a basic strategy based on team
// states.
type rawBotWrapper struct {
	Logger       Logger
	PlayerNumber uint32
	Side         proto.Team_Side
	Bot          Bot
}

func (h *rawBotWrapper) GetReadyHandler(ctx context.Context, inspector SnapshotInspector) {
	h.Bot.OnGetReady(ctx, inspector)
}

func (h *rawBotWrapper) TurnHandler(ctx context.Context, inspector SnapshotInspector) ([]proto.PlayerOrder, string, error) {
	if inspector == nil {
		return nil, "", fmt.Errorf("error processing turn: %s", ErrNilSnapshot)

	}

	state, err := defineMyState(inspector.GetSnapshot(), h.PlayerNumber, h.Side)
	if err != nil {
		return nil, "", fmt.Errorf("error processing turn %d: %s", inspector.GetSnapshot().Turn, err)

	}

	if specs.GoalkeeperNumber == h.PlayerNumber {
		return h.Bot.AsGoalkeeper(ctx, inspector, state)
	} else {
		switch state {
		case Supporting:
			return h.Bot.OnSupporting(ctx, inspector)
		case HoldingTheBall:
			return h.Bot.OnHolding(ctx, inspector)
		case Defending:
			return h.Bot.OnDefending(ctx, inspector)
		case DisputingTheBall:
			return h.Bot.OnDisputing(ctx, inspector)
		}
	}
	return nil, "", fmt.Errorf("unknown player state '%s'", state)
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
