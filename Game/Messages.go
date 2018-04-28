package Game

import (
	"github.com/maketplay/commons/BasicTypes"
)

type PlayerMessage struct {
	Type     BasicTypes.MsgType `json:"type"`
	PlayerId int                `json:"player_id"`
	Orders   []BasicTypes.Order `json:"orders"`
	Debug    string             `json:"message"`
}

type GameMessage struct {
	Type     BasicTypes.MsgType     `json:"type"`
	GameInfo GameInfo               `json:"info"`
	State    BasicTypes.State       `json:"state"`
	Data     map[string]string		`json:"data"`
	Message  string                 `json:"message"`
}

type GameInfo struct {
	State            BasicTypes.State `json:"state"`
	Ball             Ball             `json:"ball"`
	HomeTeam         Team             `json:"home"`
	AwayTeam         Team             `json:"away"`
	RemainingSeconds int              `json:"time"`
}
