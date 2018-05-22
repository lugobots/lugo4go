package Game

import (
	"github.com/makeitplay/commons/Units"
	"github.com/makeitplay/commons/Physics"
)

// distance considered "near" for a player to the ball
const DistanceNearBall = Units.CourtHeight / 2 // units float
const ERROR_MARGIN_RUNNING = 20.0
const ERROR_MARGIN_PASSING = 20.0

type PlayerRegion struct {
	CornerA Physics.Point
	CornerB Physics.Point
}

var HomePlayersRegions = map[Units.PlayerNumber]PlayerRegion{
	Units.PositionA: {
		Physics.Point{0, int(Units.CourtHeight * 0.6)},
		Physics.Point{int(Units.CourtWidth * 0.6), Units.CourtHeight},
	},
	Units.PositionB: {
		Physics.Point{0, int(Units.CourtHeight * 0.3)},
		Physics.Point{int(Units.CourtWidth * 0.6), int(Units.CourtHeight * 0.7)},
	},
	Units.PositionC: {
		Physics.Point{0, 0},
		Physics.Point{int(Units.CourtWidth * 0.6), int(Units.CourtHeight * 0.4)},
	},
	Units.PositionD: {
		Physics.Point{int(Units.CourtWidth * 0.4), int(Units.CourtHeight * 0.5)},
		Physics.Point{Units.CourtWidth, Units.CourtHeight},
	},
	Units.PositionE: {
		Physics.Point{int(Units.CourtWidth * 0.4), 0},
		Physics.Point{Units.CourtWidth, int(Units.CourtHeight * 0.5)},
	},
}