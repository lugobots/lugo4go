package client

import (
	"github.com/makeitplay/arena/physics"
)

// Ball is the ball :-)
type Ball struct {
	physics.Element
	// Holder identifies the player who is holding the ball
	Holder *Player
}
