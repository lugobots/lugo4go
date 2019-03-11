package client

import (
	"fmt"
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/orders"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/arena/units"
)

// Player acts as a brainless player in the game. This struct implements many methods that does not affect the player
// intelligence/behaviour/decisions. So, it is meant to reduce the developer concerns about communication, protocols,
// attributes, etc, and focusing in the player intelligence.

// ID returns the player ID, that is the team place and it concatenated.
func (p *Player) ID() string {
	if p.Id == "" {
		p.Id = fmt.Sprintf("%s-%s", p.TeamPlace, p.Number)
	}
	return p.Id
}

// GetMyTeamStatus retrieve the player team status from the game server message
func (p *Player) GetMyTeamStatus(gameInfo GameInfo) Team {
	if p.TeamPlace == arena.HomeTeam {
		return gameInfo.HomeTeam
	}
	return gameInfo.AwayTeam
}

// GetOpponentTeam retrieve the opponent team status from the game server message
func (p *Player) GetOpponentTeam(status GameInfo) Team {
	if p.TeamPlace == arena.HomeTeam {
		return status.AwayTeam
	}
	return status.HomeTeam
}

// FindOpponentPlayer retrieve a specific opponent player status from the game server message
func (p *Player) FindOpponentPlayer(status GameInfo, playerNumber arena.PlayerNumber) *Player {
	teamInfo := p.GetOpponentTeam(status)
	for _, playerInfo := range teamInfo.Players {
		if playerInfo.Number == playerNumber {
			return playerInfo
		}
	}
	return nil
}

// CreateMoveOrder creates a move order
func (p *Player) CreateMoveOrder(target physics.Point, speed float64) (orders.Order, error) {
	vec, err := physics.NewVector(p.Coords, target)
	if err != nil {
		return orders.Order{}, err
	}
	vel := physics.NewZeroedVelocity(*vec.Normalize())
	vel.Speed = speed
	return orders.NewMoveOrder(vel), nil
}

// CreateJumpOrder creates a jump order (only allowed to goal keeper
func (p *Player) CreateJumpOrder(target physics.Point, speed float64) (orders.Order, error) {
	vec, err := physics.NewVector(p.Coords, target)
	if err != nil {
		return orders.Order{}, err
	}
	vel := physics.NewZeroedVelocity(*vec.Normalize())
	vel.Speed = speed
	return orders.NewMoveOrder(vel), nil
}

// CreateMoveOrderMaxSpeed creates a move order with max speed allowed
func (p *Player) CreateMoveOrderMaxSpeed(target physics.Point) (orders.Order, error) {
	return p.CreateMoveOrder(target, units.PlayerMaxSpeed)
}

// CreateStopOrder creates a move order with speed zero
func (p *Player) CreateStopOrder(direction physics.Vector) orders.Order {
	vec := p.Velocity.Copy()
	vec.Speed = 0
	vec.Direction = &direction
	return orders.NewMoveOrder(vec)
}

// CreateKickOrder creates a kick order and try to find the best vector to reach the target
func (p *Player) CreateKickOrder(ball Ball, target physics.Point, speed float64) (orders.Order, error) {
	ballExpectedDirection, err := physics.NewVector(ball.Coords, target)
	if err != nil {
		return orders.Order{}, err
	}
	diffVector, err := ballExpectedDirection.Sub(ball.Velocity.Direction)
	if err != nil {
		return orders.Order{}, err
	}
	vec := physics.NewZeroedVelocity(*diffVector)
	vec.Speed = speed

	return orders.NewKickOrder(vec), nil
}

// CreateCatchOrder creates the catch order
func (p *Player) CreateCatchOrder() orders.Order {
	return orders.NewCatchOrder()
}

// IHoldTheBall returns true when the player is holding the ball
func (p *Player) IHoldTheBall(ball Ball) bool {
	return ball.Holder != nil && ball.Holder.ID() == p.ID()
}

// OpponentGoal returns the Goal os the opponent
func (p *Player) OpponentGoal() arena.Goal {
	if p.TeamPlace == arena.HomeTeam {
		return arena.AwayTeamGoal
	}
	return arena.HomeTeamGoal
}

// DefenseGoal returns the player team goal
func (p *Player) DefenseGoal() arena.Goal {
	if p.TeamPlace == arena.HomeTeam {
		return arena.HomeTeamGoal
	}
	return arena.AwayTeamGoal
}

// IsGoalkeeper returns true if the player is the goalkeeper
func (p *Player) IsGoalkeeper() bool {
	return p.Number == arena.GoalkeeperNumber
}
