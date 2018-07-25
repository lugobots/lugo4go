package Game

import (
	"github.com/makeitplay/commons/Units"
)

type Team struct {
	Name    Units.TeamPlace `json:"name"`
	Score   int             `json:"score"`
	Players []*Player       `json:"players"`
}

func (t *Team) AddPlayer(player *Player) (numPlayers int, err error) {
	player.TeamPlace = t.Name
	t.Players = append(t.Players, player)

	return len(t.Players), nil
}
