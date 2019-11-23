package lugo

func GetTeam(s *GameSnapshot, side Team_Side) *Team {
	if s == nil {
		return nil
	}
	if side == Team_HOME {
		return s.HomeTeam
	}
	return s.AwayTeam
}

func IsBallHolder(s *GameSnapshot, player *Player) bool {
	if s == nil {
		return false
	}
	return s.Ball != nil &&
		s.Ball.Holder != nil &&
		s.Ball.Holder.TeamSide == player.TeamSide &&
		s.Ball.Holder.Number == player.Number
}

func GetOpponentSide(side Team_Side) Team_Side {
	if side == Team_HOME {
		return Team_AWAY
	}
	return Team_HOME
}

func GetPlayer(s *GameSnapshot, side Team_Side, number uint32) *Player {
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

func MakeOrder_MoveMaxSpeed(origin, target Point) (Order_Move, error) {
	return MakeOrder_Move(origin, target, PlayerMaxSpeed)
}

func MakeOrder_Move(origin, target Point, speed float64) (Order_Move, error) {
	vec, err := NewVector(origin, target)
	if err != nil {
		return Order_Move{}, err
	}
	vel := NewZeroedVelocity(*vec.Normalize())
	vel.Speed = speed
	return Order_Move{Move: &Move{Velocity: &vel}}, nil
}

func MakeOrder_Kick(ball Ball, target Point, speed float64) (Order_Kick, error) {
	ballExpectedDirection, err := NewVector(*ball.Position, target)
	if err != nil {
		return Order_Kick{}, err
	}
	diffVector, err := ballExpectedDirection.Sub(ball.Velocity.Direction)
	if err != nil {
		return Order_Kick{}, err
	}
	vel := NewZeroedVelocity(*diffVector)
	vel.Speed = speed

	return Order_Kick{Kick: &Kick{Velocity: &vel}}, nil
}

func MakeOrder_Jump(origin, target Point, speed float64) (Order_Jump, error) {
	vec, err := NewVector(origin, target)
	if err != nil {
		return Order_Jump{}, err
	}
	vel := NewZeroedVelocity(*vec.Normalize())
	vel.Speed = speed
	return Order_Jump{Jump: &Jump{Velocity: &vel}}, nil
}

func MakeOrder_Catch() (Order_Catch, error) {
	return Order_Catch{Catch: &Catch{}}, nil
}
