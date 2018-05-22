package Game

import (
	"github.com/makeitplay/commons/Physics"
	"github.com/makeitplay/commons/Units"
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
	for power >= Units.BallMinSpeed {
		distance += int(Units.BallMaxSpeed * power)
		power *= Units.BallDeceleration
	}
	return distance
}
