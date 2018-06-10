package Game

import (
	"strconv"
	"github.com/makeitplay/commons/Units"
	"github.com/makeitplay/commons/BasicTypes"
)

type Team struct {
	Name    Units.TeamPlace        `json:"name"`
	Score   int             `json:"score"`
	Players map[int]*Player `json:"players"`
}

func (t *Team) AddPlayer(player *Player) (numPlayers int, err error) {
	if t.Players == nil {
		t.Players = map[int]*Player{}
	}

	countPlayers := len(t.Players) + 1
	nextNumber := strconv.Itoa(countPlayers)

	player.TeamPlace = t.Name
	player.Number = BasicTypes.PlayerNumber(nextNumber)
	t.Players[player.Id] = player

	return countPlayers, nil
}
