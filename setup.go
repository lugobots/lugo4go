package client

import (
	"context"
	"fmt"
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/arena/talk"
	"math/rand"
	"net/url"
	"time"
)

func TalkerSetup(mainCtx GamerCtx, config *Configuration, initialPos physics.Point) (context.Context, talk.Talker, error) {
	rand.Seed(time.Now().UnixNano())
	// First we have to get the command line arguments to identify this bot in its game
	uri := new(url.URL)
	uri.Scheme = "ws"
	uri.Host = fmt.Sprintf("%s:%s", config.WSHost, config.WSPort)
	uri.Path = fmt.Sprintf("/announcements/%s/%s", config.UUID, config.TeamPlace)

	playerSpec := arena.PlayerSpecifications{
		Number:          config.PlayerNumber,
		InitialCoords:   initialPos,
		Token:           config.Token,
		ProtocolVersion: "1.0",
	}

	talker := talk.NewTalker(mainCtx.Logger())
	talkerCtx, err := talker.Connect(mainCtx, *uri, playerSpec)
	if err != nil {
		return nil, nil, fmt.Errorf("fail on opening the websocket connection: %s", err)
	}
	return talkerCtx, talker, nil
}
