package Game

import (
	"github.com/makeitplay/commons/Physics"
)

type Ball struct {
	Physics.Element
	Holder *Player
}
