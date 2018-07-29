package client

import (
	"github.com/makeitplay/commons/BasicTypes"
)

// PlayerMessage is the message sent from a player to the game server
type PlayerMessage struct {
	Type   BasicTypes.MsgType `json:"type"`
	Orders []BasicTypes.Order `json:"orders"`
	// Debug is a message the will be only visible in the game server log (used for debugging purposes)
	Debug string `json:"message"`
}

// GameMessage is the message sent from the game server to the player
type GameMessage struct {
	Type     BasicTypes.MsgType `json:"type"`
	GameInfo GameInfo           `json:"info"`
	State    BasicTypes.State   `json:"state"`
	Data     map[string]string  `json:"data"`
	// Message is quite useless, but could help the developers to debug the game server messages
	Message string `json:"message"`
}

// GameInfo is the set of values that defines the current game state
type GameInfo struct {
	State BasicTypes.State `json:"state"`
	// Turn is the sequential number of turns. Read the game documentation to understand what a turn is
	Turn     int  `json:"turn"`
	Ball     Ball `json:"ball"`
	HomeTeam Team `json:"home"`
	AwayTeam Team `json:"away"`
	// RemainingSeconds is a estimation of how long the game will take. However it is not precise.
	RemainingSeconds int `json:"time"`
}
