package field

import (
	"github.com/lugobots/lugo4go/v2/proto"
)

func GetTeam(s *proto.GameSnapshot, side proto.Team_Side) *proto.Team {
	if s == nil {
		return nil
	}
	if side == proto.Team_HOME {
		return s.HomeTeam
	}
	return s.AwayTeam
}

func IsBallHolder(s *proto.GameSnapshot, player *proto.Player) bool {
	if s == nil {
		return false
	}
	return s.Ball != nil && player != nil &&
		s.Ball.Holder != nil &&
		s.Ball.Holder.TeamSide == player.TeamSide &&
		s.Ball.Holder.Number == player.Number
}

func GetOpponentSide(side proto.Team_Side) proto.Team_Side {
	if side == proto.Team_HOME {
		return proto.Team_AWAY
	}
	return proto.Team_HOME
}

func GetOpponentGoal(mySide proto.Team_Side) Goal {
	return GetTeamsGoal(GetOpponentSide(mySide))
}

func GetPlayer(s *proto.GameSnapshot, side proto.Team_Side, number uint32) *proto.Player {
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

func MakeOrderMoveMaxSpeed(origin, target proto.Point) (*proto.Order_Move, error) {
	return MakeOrderMove(origin, target, PlayerMaxSpeed)
}

func MakeOrderMove(origin, target proto.Point, speed float64) (*proto.Order_Move, error) {
	vec, err := proto.NewVector(origin, target)
	if err != nil {
		return nil, err
	}
	vel := proto.NewZeroedVelocity(*vec.Normalize())
	vel.Speed = speed
	return &proto.Order_Move{Move: &proto.Move{Velocity: &vel}}, nil
}

func MakeOrderJump(origin, target proto.Point, speed float64) (*proto.Order_Jump, error) {
	vec, err := proto.NewVector(origin, target)
	if err != nil {
		return nil, err
	}
	vel := proto.NewZeroedVelocity(*vec.Normalize())
	vel.Speed = speed
	return &proto.Order_Jump{Jump: &proto.Jump{Velocity: &vel}}, nil
}

func MakeOrderKick(ball proto.Ball, target proto.Point, speed float64) (*proto.Order_Kick, error) {
	ballExpectedDirection, err := proto.NewVector(*ball.Position, target)
	if err != nil {
		return nil, err
	}

	diffVector, err := ballExpectedDirection.Sub(ball.Velocity.Direction)
	if err != nil {
		return nil, err
	}
	vel := proto.NewZeroedVelocity(*diffVector)
	vel.Direction.Normalize()
	vel.Speed = speed

	return &proto.Order_Kick{Kick: &proto.Kick{Velocity: &vel}}, nil
}

func MakeOrderKickMaxSpeed(ball proto.Ball, target proto.Point) (*proto.Order_Kick, error) {
	return MakeOrderKick(ball, target, BallMaxSpeed)

}

func MakeOrderCatch() *proto.Order_Catch {
	return &proto.Order_Catch{Catch: &proto.Catch{}}
}

func GoForward(side proto.Team_Side) *proto.Order_Move {
	forward := proto.East()
	if side == proto.Team_AWAY {
		forward = proto.West()
	}
	vel := proto.NewZeroedVelocity(forward)
	vel.Speed = PlayerMaxSpeed
	return &proto.Order_Move{Move: &proto.Move{Velocity: &vel}}
}

func GoBackward(side proto.Team_Side) *proto.Order_Move {
	backward := proto.West()
	if side == proto.Team_AWAY {
		backward = proto.East()
	}
	vel := proto.NewZeroedVelocity(backward)
	vel.Speed = PlayerMaxSpeed
	return &proto.Order_Move{Move: &proto.Move{Velocity: &vel}}
}

func GoRight(side proto.Team_Side) *proto.Order_Move {
	right := proto.South()
	if side == proto.Team_AWAY {
		right = proto.North()
	}
	vel := proto.NewZeroedVelocity(right)
	vel.Speed = PlayerMaxSpeed
	return &proto.Order_Move{Move: &proto.Move{Velocity: &vel}}
}

func GoLeft(side proto.Team_Side) *proto.Order_Move {
	left := proto.North()
	if side == proto.Team_AWAY {
		left = proto.South()
	}
	vel := proto.NewZeroedVelocity(left)
	vel.Speed = PlayerMaxSpeed
	return &proto.Order_Move{Move: &proto.Move{Velocity: &vel}}
}
