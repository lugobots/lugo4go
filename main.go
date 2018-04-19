package the_dummies

import (
	"os"
	"math/rand"
	"time"
	"os/signal"
	"github.com/ant21games/the-dummies/App"
	"github.com/ant21games/the-dummies/Game"
	"flag"
)



var (
	player Game.Player
	Settings Game.StartUpSettings
	teamName Game.TeamName
)

func main() {
	rand.Seed(time.Now().Unix())
	watchInterruptions()
	defer App.Cleanup()
	Settings = LoadSetting()
	player = Game.Player{}
	App.NickName = "New Player"
	player.Start(teamName)
}

func watchInterruptions() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			App.Log("*********** INTERRUPTION SIGNAL ****************")
			App.Cleanup()
			os.Exit(0)
		}
	}()
}

func LoadSetting() Game.StartUpSettings {
	var team, title string
	var port int64
	flag.StringVar(&team, "team", string(Game.HomeTeam), "(home or away)")
	flag.StringVar(&title, "title", string(Game.HomeTeam), "Team's name")
	flag.Int64Var(&port, "port", 8080, "Port server")
	flag.Parse()

	return Game.StartUpSettings{Game.TeamName(team), int(port), title}
}