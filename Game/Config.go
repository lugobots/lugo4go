package Game

import "flag"

type Configuration struct {
	TeamName TeamName
	Uuid string //this value will be automatically given to your binary by the server :) You may ignore it locally

	QueueUser string
	QueuePassword string
	QueueHost string
	QueueVHost string
	QueuePort string

	OutputExchange string
	OutputQueue string
	InputExchange string
	InputQueue string
}

func (c *Configuration) LoadCmdArg() {
	//mandatory
	var name string
	flag.StringVar(&name, 		"team", 		c.QueueUser, "Team (home or away). (Auto-provided in production)")
	if name != string(HomeTeam) && name == string(AwayTeam) {
		panic("Invalid team option. Must be either home or away")
	}

	flag.StringVar(&c.Uuid,	 			"uui", 		c.QueueUser, "Uuid for this player instance. (Auto-provided in production)")
	if c.Uuid == "" {
		panic("You must provide a player uuid")
	}
	//optional
	flag.StringVar(&c.QueueUser, 		"QueueUser", 		c.QueueUser, "AMQP username")
	flag.StringVar(&c.QueuePassword, 	"QueuePassword", 	c.QueuePassword, "AMQP Password")
	flag.StringVar(&c.QueueHost, 		"QueueHost", 		c.QueueHost, "AMQP server host")
	flag.StringVar(&c.QueueVHost, 		"QueueVHost", 		c.QueueVHost, "AMQP ")
	flag.StringVar(&c.QueuePort, 		"QueuePort", 		c.QueuePort, "The match Uuid (useless locally)")
	flag.StringVar(&c.OutputExchange, 	"OutputExchange", 	c.OutputExchange, "The match Uuid (useless locally)")
	flag.StringVar(&c.OutputQueue, 		"OutputQueue", 	c.OutputQueue, "The match Uuid (useless locally)")
	flag.StringVar(&c.InputExchange, 	"InputExchange", 	c.InputExchange, "The match Uuid (useless locally)")
	flag.StringVar(&c.InputQueue, 		"InputQueue", 		c.InputExchange, "The match Uuid (useless locally)")
	flag.Parse()
}
