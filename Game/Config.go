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

	QueueUser     string
	QueuePassword string
	QueueHost     string
	QueueVHost    string
	QueuePort     string

	OutputExchange string
	OutputQueue    string
	InputExchange  string
	InputQueue     string
}

func (c *Configuration) LoadCmdArg() {
	//mandatory
	var name string
	var number int

	flag.StringVar(&name, "team", "home", "Team (home or away). (Auto-provided in production)")
	flag.IntVar(&number, "number", 0, "Player's number")

	flag.StringVar(&c.Uuid, "uui", c.Uuid, "Uuid for this player instance. (Auto-provided in production)")
	flag.StringVar(&c.QueueUser, "QueueUser", c.QueueUser, "AMQP username")
	flag.StringVar(&c.QueuePassword, "QueuePassword", c.QueuePassword, "AMQP Password")
	flag.StringVar(&c.QueueHost, "QueueHost", c.QueueHost, "AMQP server host")
	flag.StringVar(&c.QueueVHost, "QueueVHost", c.QueueVHost, "AMQP ")
	flag.StringVar(&c.QueuePort, "QueuePort", c.QueuePort, "The match Uuid (useless locally)")
	flag.StringVar(&c.OutputExchange, "OutputExchange", c.OutputExchange, "The match Uuid (useless locally)")
	flag.StringVar(&c.OutputQueue, "OutputQueue", c.OutputQueue, "The match Uuid (useless locally)")
	flag.StringVar(&c.InputExchange, "InputExchange", c.InputExchange, "The match Uuid (useless locally)")
	flag.StringVar(&c.InputQueue, "InputQueue", c.InputQueue, "The match Uuid (useless locally)")
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
