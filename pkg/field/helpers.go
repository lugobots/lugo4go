package field

import (
	"github.com/lugobots/lugo4go/v2/lugo"
)

func GetTeam(s *lugo.GameSnapshot, side lugo.Team_Side) *lugo.Team {
	if s == nil {
		return nil
	}
	if side == lugo.Team_HOME {
		return s.HomeTeam
	}
	return s.AwayTeam
}

func IsBallHolder(s *lugo.GameSnapshot, player *lugo.Player) bool {
	if s == nil {
		return false
	}
	return s.Ball != nil && player != nil &&
		s.Ball.Holder != nil &&
		s.Ball.Holder.TeamSide == player.TeamSide &&
		s.Ball.Holder.Number == player.Number
}

func GetOpponentSide(side lugo.Team_Side) lugo.Team_Side {
	if side == lugo.Team_HOME {
		return lugo.Team_AWAY
	}
	return lugo.Team_HOME
}

func GetOpponentGoal(mySide lugo.Team_Side) Goal {
	return GetTeamsGoal(GetOpponentSide(mySide))
}

func GetPlayer(s *lugo.GameSnapshot, side lugo.Team_Side, number uint32) *lugo.Player {
	team := GetTeam(s, side)
	if team == nil {
		return nil
	}
	for _, player := range team.Players {
		if player.Number == number {
			return player
		}
	}
	return nil
}

func MakeOrderMoveMaxSpeed(origin, target lugo.Point) (*lugo.Order_Move, error) {
	return MakeOrderMove(origin, target, PlayerMaxSpeed)
}

func MakeOrderMove(origin, target lugo.Point, speed float64) (*lugo.Order_Move, error) {
	vec, err := lugo.NewVector(origin, target)
	if err != nil {
		return nil, err
	}
	vel := lugo.NewZeroedVelocity(*vec.Normalize())
	vel.Speed = speed
	return &lugo.Order_Move{Move: &lugo.Move{Velocity: &vel}}, nil
}

func MakeOrderJump(origin, target lugo.Point, speed float64) (*lugo.Order_Jump, error) {
	vec, err := lugo.NewVector(origin, target)
	if err != nil {
		return nil, err
	}
	vel := lugo.NewZeroedVelocity(*vec.Normalize())
	vel.Speed = speed
	return &lugo.Order_Jump{Jump: &lugo.Jump{Velocity: &vel}}, nil
}

func MakeOrderKick(ball lugo.Ball, target lugo.Point, speed float64) (*lugo.Order_Kick, error) {
	ballExpectedDirection, err := lugo.NewVector(*ball.Position, target)
	if err != nil {
		return nil, err
	}

	diffVector, err := ballExpectedDirection.Sub(ball.Velocity.Direction)
	if err != nil {
		return nil, err
	}
	vel := lugo.NewZeroedVelocity(*diffVector)
	vel.Direction.Normalize()
	vel.Speed = speed

	return &lugo.Order_Kick{Kick: &lugo.Kick{Velocity: &vel}}, nil
}

func MakeOrderCatch() *lugo.Order_Catch {
	return &lugo.Order_Catch{Catch: &lugo.Catch{}}
}
