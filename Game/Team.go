package Game

import (
	"errors"
	"strconv"
)

type Team struct {
	Name    TeamName        `json:"name"`
	Score   int             `json:"score"`
	Players map[int]*Player `json:"players"`
}

func (t *Team) AddPlayer(player *Player) (numPlayers int, err error) {
	if t.Players == nil {
		t.Players = map[int]*Player{}
	}

	if len(t.Players) >= MAX_PLAYERS {
		return len(t.Players), errors.New("Cannot accept more players")
	}
	countPlayers := len(t.Players) + 1
	nextNumber := strconv.Itoa(countPlayers)

	player.TeamName = t.Name
	player.Number = PlayerNumber(nextNumber)
	t.Players[player.Id] = player

	return countPlayers, nil
}
