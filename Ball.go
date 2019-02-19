package client

import (
	"github.com/makeitplay/arena/Physics"
)

// Ball is the ball :-)
type Ball struct {
	Physics.Element
	// Holder identifies the player who is holding the ball
	Holder *Player
}
