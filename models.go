package client

import (
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/orders"
	"github.com/makeitplay/arena/physics"
)

// PlayerMessage is the message sent from a player to the game server
type PlayerMessage struct {
	Type   arena.MsgType  `json:"type"`
	Orders []orders.Order `json:"orders"`
	// Debug is a message the will be only visible in the game server log (used for debugging purposes)
	Debug string `json:"message"`
}

// GameMessage is the message sent from the game server to the player
type GameMessage struct {
	Type     arena.MsgType          `json:"type"`
	GameInfo GameInfo               `json:"info"`
	State    arena.GameState        `json:"state"`
	Data     map[string]interface{} `json:"data"`
	// Message is quite useless, but could help the developers to debug the game server messages
	Message string `json:"message"`
}

// GameInfo is the set of values that defines the current game state
type GameInfo struct {
	State arena.GameState `json:"state"`
	// Turn is the sequential number of turns. Read the game documentation to understand what a turn is
	Turn     int  `json:"turn"`
	Ball     Ball `json:"ball"`
	HomeTeam Team `json:"home"`
	AwayTeam Team `json:"away"`
}

type Player struct {
	physics.Element
	Id        string             `json:"id"`
	Number    arena.PlayerNumber `json:"number"`
	TeamPlace arena.TeamPlace    `json:"team_place"`
}
