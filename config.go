package client

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/lugobots/client-player-go/v2/proto"
	"io/ioutil"
	"strings"
)

// Configuration is the set of values expected as a initial configuration of the player
type Config struct {
	// Full url to the gRPC server
	GRPCAddress     string `json:"grpc_address"`
	Insecure        bool   `json:"insecure"`
	Token           string `json:"token"`
	TeamSide        proto.Team_Side
	Number          uint32      `json:"number"`
	InitialPosition proto.Point `json:"-"`
}

type jsonConfig struct {
	Config
	Team string `json:"team"`
}

func LoadConfig(filepath string) (c Config, e error) {
	content, err := ioutil.ReadFile(filepath)

	config := jsonConfig{}
	if err != nil {
		e = fmt.Errorf("error loading the config file at %s: %s", filepath, err)
	} else if err := json.Unmarshal(content, &config); err != nil {
		e = fmt.Errorf("error parsing the config file at %s: %s", filepath, err)
	} else if _, ok := proto.Team_Side_name[int32(c.TeamSide)]; !ok {
		e = fmt.Errorf("invalid team side in config file at %s", filepath)
	} else if config.Number < 1 || config.Number > 11 {
		e = fmt.Errorf("invalid player number in config file at %s: %d", filepath, config.Number)
	}

	side, ok := proto.Team_Side_value[strings.ToUpper(config.Team)]
	if !ok {
		e = fmt.Errorf("invalid team option '%s'. Must be either HOME or AWAY", config.Team)
	}
	c = config.Config
	c.TeamSide = proto.Team_Side(side)
	return
}

// ParseConfigFlags is a helper that sets flags to make the configuration be overwritten by command line.
// Note that it won't be used in production, The config file is the only official way to configure it.
func (c *Config) ParseConfigFlags() error {
	var name string
	var number int

	flag.StringVar(&name, "team", "home", "Team (home or away)")
	flag.IntVar(&number, "number", 0, "Player's number")
	flag.StringVar(&c.GRPCAddress, "grpc_address", "localhost:8080", "Address to connect to the game server")
	flag.StringVar(&c.Token, "token", "", "Token used by the server to identify the right connection")
	flag.BoolVar(&c.Insecure, "insecure", true, "Allow insecure connections (important for development environments)")

	flag.Parse()

	side, ok := proto.Team_Side_value[strings.ToUpper(name)]
	if !ok {
		return fmt.Errorf("invalid team option '%s'. Must be either HOME or AWAY", name)
	}

	if number < 1 || number > 11 {
		return fmt.Errorf("invalid player number '%d'. Must be 1 to 11", number)
	}

	c.TeamSide = proto.Team_Side(side)
	c.Number = uint32(number)
	return nil
}
