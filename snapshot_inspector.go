package lugo4go

import (
	"fmt"

	"github.com/lugobots/lugo4go/v3/mapper"
	"github.com/lugobots/lugo4go/v3/proto"
	"github.com/lugobots/lugo4go/v3/specs"
)

type inspector struct {
	mySide   proto.Team_Side
	myNumber int
	me       *proto.Player
	snapshot *proto.GameSnapshot
}

func newInspector(botSide proto.Team_Side, playerNumber int, gameSnapshot *proto.GameSnapshot) (*inspector, error) {
	tools := &inspector{mySide: botSide, myNumber: playerNumber, snapshot: gameSnapshot}

	me := tools.GetPlayer(botSide, playerNumber)
	if me == nil {
		return nil, fmt.Errorf("could not to find the player %s-%d", botSide, playerNumber)
	}
	tools.me = me
	return tools, nil
}

func (i *inspector) GetSnapshot() *proto.GameSnapshot {
	return i.snapshot
}

func (i *inspector) GetMe() *proto.Player {
	return i.me
}

func (i *inspector) GetBall() *proto.Ball {
	return i.snapshot.GetBall()
}

func (i *inspector) GetBallHolder() (*proto.Player, bool) {
	holder := i.snapshot.GetBall().GetHolder()
	return holder, holder != nil
}

func (i *inspector) IsBallHolder(player *proto.Player) bool {
	holder := i.snapshot.GetBall().GetHolder()
	return holder != nil && holder.TeamSide == player.TeamSide && holder.Number == player.Number
}

func (i *inspector) GetTeam(side proto.Team_Side) *proto.Team {
	if side == proto.Team_HOME {
		return i.snapshot.HomeTeam
	}
	return i.snapshot.AwayTeam
}

func (i *inspector) GetMyTeam() *proto.Team {
	return i.GetTeam(i.mySide)
}

func (i *inspector) GetOpponentMyTeam() *proto.Team {
	return i.GetTeam(i.GetOpponentSide())
}

func (i *inspector) GetMyTeamSide() proto.Team_Side {
	return i.mySide
}

func (i *inspector) GetOpponentSide() proto.Team_Side {
	if i.mySide == proto.Team_HOME {
		return proto.Team_AWAY
	}
	return proto.Team_HOME
}

func (i *inspector) GetPlayer(side proto.Team_Side, number int) *proto.Player {
	team := i.GetTeam(side)
	for _, player := range team.GetPlayers() {
		if int(player.Number) == number {
			return player
		}
	}
	return nil
}

func (i *inspector) GetMyTeamPlayers() []*proto.Player {
	return i.GetTeam(i.mySide).Players
}

func (i *inspector) GetOpponentPlayers() []*proto.Player {
	return i.GetTeam(i.GetOpponentSide()).Players
}

func (i *inspector) GetMyTeamGoalkeeper() *proto.Player {
	return i.GetPlayer(i.GetMyTeamSide(), int(specs.GoalkeeperNumber))
}

func (i *inspector) GetOpponentGoalkeeper() *proto.Player {
	return i.GetPlayer(i.GetOpponentSide(), int(specs.GoalkeeperNumber))
}

func (i *inspector) MakeOrderMoveMaxSpeed(target proto.Point) (*proto.Order_Move, error) {
	return i.MakeOrderMoveFromPoint(*i.me.Position, target, specs.PlayerMaxSpeed)
}

func (i *inspector) MakeOrderMoveFromPoint(origin, target proto.Point, speed float64) (*proto.Order_Move, error) {
	vec, err := proto.NewVector(origin, target)
	if err != nil {
		return nil, err
	}
	vel := proto.NewZeroedVelocity(*vec.Normalize())
	vel.Speed = speed
	return &proto.Order_Move{Move: &proto.Move{Velocity: &vel}}, nil
}

func (i *inspector) MakeOrderMoveFromVector(vector proto.Vector, speed float64) *proto.Order_Move {
	targetPoint := vector.TargetFrom(*i.me.Position)
	// no need to check for errors since a vector will always lead to a valid destination
	order, _ := i.MakeOrderMoveFromPoint(*i.me.Position, targetPoint, speed)
	return order
}

func (i *inspector) MakeOrderMoveByDirection(direction mapper.Direction, speed float64) *proto.Order_Move {
	directionTarget := directionOrientationMap[i.mySide][direction]
	// no need to check for errors since the vector is known and valid
	return i.MakeOrderMoveFromVector(proto.Vector(directionTarget), speed)
}

func (i *inspector) MakeOrderMoveToStop() *proto.Order_Move {
	myDirection := i.GetMe().GetVelocity().GetDirection()
	if myDirection == nil {
		v := proto.Vector(directionOrientationMap[i.mySide][mapper.Forward])
		myDirection = &v
	}
	return i.MakeOrderMoveFromVector(*myDirection, 0)
}

func (i *inspector) MakeOrderJump(target proto.Point, speed float64) (*proto.Order_Jump, error) {
	vec, err := proto.NewVector(*i.me.Position, target)
	if err != nil {
		return nil, err
	}
	vel := proto.NewZeroedVelocity(*vec.Normalize())
	vel.Speed = speed
	return &proto.Order_Jump{Jump: &proto.Jump{Velocity: &vel}}, nil
}

func (i *inspector) MakeOrderKick(target proto.Point, speed float64) (*proto.Order_Kick, error) {
	ballExpectedDirection, err := proto.NewVector(*i.snapshot.GetBall().Position, target)
	if err != nil {
		return nil, err
	}

	diffVector, err := ballExpectedDirection.Sub(i.snapshot.GetBall().Velocity.Direction)
	if err != nil {
		return nil, err
	}
	vel := proto.NewZeroedVelocity(*diffVector)
	vel.Direction.Normalize()
	vel.Speed = speed

	return &proto.Order_Kick{Kick: &proto.Kick{Velocity: &vel}}, nil
}

func (i *inspector) MakeOrderKickMaxSpeed(target proto.Point) (*proto.Order_Kick, error) {
	return i.MakeOrderKick(target, specs.BallMaxSpeed)
}

func (i *inspector) MakeOrderCatch() *proto.Order_Catch {
	return &proto.Order_Catch{Catch: &proto.Catch{}}
}

var directionOrientationMap = map[proto.Team_Side]map[mapper.Direction]mapper.Orientation{
	proto.Team_HOME: {
		mapper.Forward:       mapper.East,
		mapper.Backward:      mapper.West,
		mapper.Left:          mapper.North,
		mapper.Right:         mapper.South,
		mapper.BackwardLeft:  mapper.NorthWest,
		mapper.BackwardRight: mapper.SouthWest,
		mapper.ForwardLeft:   mapper.NorthEast,
		mapper.ForwardRight:  mapper.SouthEast,
	},
	proto.Team_AWAY: {
		mapper.Forward:       mapper.West,
		mapper.Backward:      mapper.East,
		mapper.Left:          mapper.South,
		mapper.Right:         mapper.North,
		mapper.BackwardLeft:  mapper.SouthEast,
		mapper.BackwardRight: mapper.NorthEast,
		mapper.ForwardLeft:   mapper.SouthWest,
		mapper.ForwardRight:  mapper.NorthWest,
	},
}
