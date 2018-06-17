package Game

import (
	"github.com/makeitplay/commons/Units"
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


	player.TeamPlace = t.Name
	t.Players[player.Id] = player

	return len(t.Players), nil
}
