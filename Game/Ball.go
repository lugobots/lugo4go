package Game

import (
	"github.com/maketplay/commons/Physics"
	"github.com/maketplay/commons/Units"
)

var maxDistance = 0

type Ball struct {
	Physics.Element
	Vector *Physics.Vector `json:"vector"`
	Holder *Player
}

func BallMaxDistance() int {
	if maxDistance == 0 {
		maxDistance = calcMaxBallDistance()
	}
	return maxDistance
}

func calcMaxBallDistance() int {
	power := 1.0
	distance := 0
	for power >= Units.BallMinPower {
		distance += int(Units.BallSpeed * power)
		power *= Units.BallSlowerRatio
	}
	return distance
}
