package Game

import (
	"flag"
	"github.com/makeitplay/commons/Units"
	"fmt"
	"github.com/makeitplay/commons/BasicTypes"
	"log"
	"strconv"
)

type Configuration struct {
	TeamPlace Units.TeamPlace
	PlayerNumber BasicTypes.PlayerNumber
	Uuid      string //this value will be automatically given to your binary by the server :) You may ignore it locally

}

func (c *Configuration) LoadCmdArg() {
	//mandatory
	var name string
	var number int

	flag.StringVar(&name, "team", "home", "Team (home or away). (Auto-provided in production)")
	flag.IntVar(&number, "number", 0, "Player's number")

	flag.StringVar(&c.Uuid, "uui", "local", "Uuid for this player instance. (Auto-provided in production)")
	flag.Parse()

	if name != string(Units.HomeTeam) && name != string(Units.AwayTeam) {
		log.Fatal("Invalid team option {" + name + "}. Must be either home or away")
	}
	if number < 1 || number > 11 {
		log.Fatal( fmt.Errorf("invalid player number {%d}. Must be 1 to 11", number))
	}
	c.PlayerNumber = BasicTypes.PlayerNumber(strconv.Itoa(number))
	c.TeamPlace = Units.TeamPlace(name)
}
