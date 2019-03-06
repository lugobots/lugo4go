package client

import (
	"github.com/makeitplay/arena"
)

// Team groups the player team info based on the status sent by the game server
type Team struct {
	Name    arena.TeamPlace `json:"name"`
	Score   int             `json:"score"`
	Players []*Player       `json:"players"`
}

// AddPlayer add to the team struct a player based on the game server messages
func (t *Team) AddPlayer(player *Player) (numPlayers int, err error) {
	player.TeamPlace = t.Name
	t.Players = append(t.Players, player)

	return len(t.Players), nil
}
