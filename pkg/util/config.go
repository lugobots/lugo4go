package util

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/lugobots/lugo4go/v2/pkg/field"
	"github.com/lugobots/lugo4go/v2/proto"
	"io/ioutil"
	"strings"
)

// Configuration is the set of values expected as a initial configuration of the player
type Config struct {
	// Full url to the gRPC server
	GRPCAddress     string          `json:"grpc_address"`
	Insecure        bool            `json:"insecure"`
	Token           string          `json:"token"`
	TeamSide        proto.Team_Side `json:"-"`
	Number          uint32          `json:"-"`
	InitialPosition *proto.Point    `json:"-"`
}

type jsonConfig struct {
	Config
	Team string `json:"team"`
	Num  uint   `json:"number"`
}

// REVIEW ParseConfigFlags is a helper that sets flags to make the configuration be overwritten by command line.
// REVIEW Note that it won't be used in production, The config file is the only official way to configure it.

///REVIEW/////

func (c *Config) ParseConfigFlags() (filepath string, err error) {
	intermediate := jsonConfig{}
	flag.StringVar(&filepath, "config-file", "", "Path to the config file")
	flag.StringVar(&intermediate.Team, "team", "home", "Team (home or away)")
	flag.UintVar(&intermediate.Num, "number", 1, "Player's number (1-11)")
	flag.StringVar(&intermediate.GRPCAddress, "grpc_address", "localhost:9090", "Address to connect to the game server")
	flag.StringVar(&intermediate.Token, "token", "", "Token used by the server to identify the right connection")
	flag.BoolVar(&intermediate.Insecure, "insecure", true, "Allow insecure connections (important for development environments)")
	flag.Parse()

	if err := c.transcribe(intermediate); err != nil {
		return "", err
	}
	return filepath, nil
}

func (c *Config) transcribe(intermediate jsonConfig) error {
	side, ok := proto.Team_Side_value[strings.ToUpper(intermediate.Team)]
	if !ok {
		return fmt.Errorf("invalid team option '%s'. Must be either HOME or AWAY", intermediate.Team)
	}
	c.TeamSide = proto.Team_Side(side)

	if intermediate.Num < 1 || intermediate.Num > field.MaxPlayers {
		return fmt.Errorf("invalid player number '%d'. Must be 1 to %d", intermediate.Num, field.MaxPlayers)
	}
	c.Number = uint32(intermediate.Num)

	c.GRPCAddress = intermediate.GRPCAddress
	c.Token = intermediate.Token
	c.Insecure = intermediate.Insecure
	return nil
}

func LoadConfig(filepath string, config *Config) error {
	intermediate := jsonConfig{}
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("error loading the config file at %s: %s", filepath, err)
	} else if err := json.Unmarshal(content, &intermediate); err != nil {
		return fmt.Errorf("error parsing the config file at %s: %s", filepath, err)
	}

	if err := config.transcribe(intermediate); err != nil {
		return err
	}
	return nil
}
