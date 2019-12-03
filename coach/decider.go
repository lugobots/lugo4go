package coach

import (
	"github.com/lugobots/lugo4go/v2"
	"github.com/lugobots/lugo4go/v2/field"
	"github.com/lugobots/lugo4go/v2/proto"
)

// PlayerState defines states specific for players
type PlayerState string

const (
	// Supporting identifies the player supporting the team mate
	Supporting PlayerState = "supporting"
	// HoldingTheBall identifies the player holding	the ball
	HoldingTheBall PlayerState = "holding"
	// Defending identifies the player defending against the opponent team
	Defending PlayerState = "defending"
	// DisputingTheBall identifies the player disputing the ball
	DisputingTheBall PlayerState = "disputing"
)

type TurnData struct {
	Me       *proto.Player
	Snapshot *proto.GameSnapshot
	Sender   lugo4go.OrderSender
}

type Decider interface {
	OnDisputing(data TurnData) error
	OnDefending(data TurnData) error
	OnHolding(data TurnData) error
	OnSupporting(data TurnData) error
	AsGoalkeeper(data TurnData) error
}

func DefineMyState(config lugo4go.Config, snapshot *proto.GameSnapshot) (PlayerState, error) {
	if snapshot == nil || snapshot.Ball == nil {
		return "", ErrNoBall
	}

	me := field.GetPlayer(snapshot, config.TeamSide, config.Number)
	if me == nil {
		return "", ErrPlayerNotFound
	}

	ballHolder := snapshot.Ball.Holder

	if ballHolder == nil {
		return DisputingTheBall, nil
	} else if ballHolder.TeamSide == config.TeamSide {
		if ballHolder.Number == config.Number {
			return HoldingTheBall, nil
		}
		return Supporting, nil
	}
	return Defending, nil
}

func DefaultTurnHandler(decider Decider, config lugo4go.Config, logger lugo4go.Logger) lugo4go.DecisionMaker {
	goalkeeper := field.GoalkeeperNumber == config.Number // it is obviously not processed every turn
	return func(snapshot *proto.GameSnapshot, sender lugo4go.OrderSender) {
		var err error
		var state PlayerState
		turnData := TurnData{
			Me:       field.GetPlayer(snapshot, config.TeamSide, config.Number),
			Snapshot: snapshot,
			Sender:   sender,
		}
		if turnData.Me == nil {
			logger.Warnf("i did not find my self in the game")
			return
		}

		if goalkeeper {
			err = decider.AsGoalkeeper(turnData)
		} else {
			state, err = DefineMyState(config, snapshot)
			switch state {
			case Supporting:
				err = decider.OnDefending(turnData)
			case HoldingTheBall:
				err = decider.OnDefending(turnData)
			case Defending:
				err = decider.OnDefending(turnData)
			case DisputingTheBall:
				err = decider.OnDefending(turnData)
			}
		}
		if err != nil {
			logger.Errorf("error processing turn %d: %s", snapshot.Turn, err)
		}

	}
}