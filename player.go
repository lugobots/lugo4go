package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/orders"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/arena/talk"
	"github.com/makeitplay/arena/units"
	"log"
	"os"
	"os/signal"
)

// Player acts as a brainless player in the game. This struct implements many methods that does not affect the player
// intelligence/behaviour/decisions. So, it is meant to reduce the developer concerns about communication, protocols,
// attributes, etc, and focusing in the player intelligence.
type Player struct {
	physics.Element
	Id             string             `json:"id"`
	Number         arena.PlayerNumber `json:"number"`
	TeamPlace      arena.TeamPlace    `json:"team_place"`
	OnMessage      func(msg GameMessage)
	OnAnnouncement func(msg GameMessage)
	config         *Configuration
	Talker         talk.Talker
	LastMsg        GameMessage
}

// playerCtx is used to keep the process running while the player is playing
var playerCtx GamerCtx
var stopPlayer context.CancelFunc

// Play make the player start to play
func (p *Player) Play(initialPosition physics.Point, configuration *Configuration) {
	playerCtx, stopPlayer = NewGamerContext(context.Background(), configuration)

	p.config = configuration
	p.TeamPlace = configuration.TeamPlace
	p.Number = configuration.PlayerNumber
	talkerCtx, talker, err := TalkerSetup(playerCtx, configuration, initialPosition)
	if err != nil {
		log.Fatal(err)
	}
	// we have to set the call back function that will process the player behaviour when the game state has been changed
	defer talker.Close()
	p.Talker = talker

	go listenServerMessages(p)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	exitCode := 0
	select {
	case <-signalChan:
		playerCtx.Logger().Print("*********** INTERRUPTION SIGNAL ****************")
		talker.Close()
		stopPlayer()
		exitCode = 1
	case <-talkerCtx.Done():
		playerCtx.Logger().Printf("was connection lost: %s", talkerCtx.Err())
		stopPlayer()
		exitCode = 2
	case <-playerCtx.Done():
		playerCtx.Logger().Printf("player stopped: %s", playerCtx.Err())

	}
	os.Exit(exitCode)
}

// stopToPlay stop the player to play
func (p *Player) stopToPlay(interrupted bool) {
	stopPlayer()
}

// ID returns the player ID, that is the team place and it concatenated.
func (p *Player) ID() string {
	if p.Id == "" {
		p.Id = fmt.Sprintf("%s-%s", p.TeamPlace, p.Number)
	}
	return p.Id
}

// LastServerMessage returns the last message got from the server
func (p *Player) LastServerMessage() GameMessage {
	return p.LastMsg
}

// SendOrders sends a list of orders to the game server, and includes a message to them (only displayed in the game server log)
func (p *Player) SendOrders(message string, ordersList ...orders.Order) {
	msg := PlayerMessage{
		orders.ORDER,
		ordersList,
		message,
	}
	stringed, err := json.Marshal(msg)
	if err != nil {
		playerCtx.Logger().Errorf("Fail generating JSON: %s", err.Error())
		return
	}

	err = p.Talker.Send(stringed)
	if err != nil {
		playerCtx.Logger().Errorf("Fail on sending message: %s", err.Error())
		return
	}
}

// UpdatePosition update the player status after the last game server message
func (p *Player) UpdatePosition(gameInfo GameInfo) {
	status := p.GetMyStatus(gameInfo)
	if status == nil {
		// sometimes the player gets a message before his welcome message be processed, then he is not officially in the game,
		// so, this status is not available yet.
		return
	}

	p.Velocity = status.Velocity
	p.Coords = status.Coords
}

// GetMyStatus retrieve the player status from the game server message
func (p *Player) GetMyStatus(gameInfo GameInfo) *Player {
	myteamInfo := p.GetMyTeamStatus(gameInfo)
	for _, playerInfo := range myteamInfo.Players {
		if playerInfo.ID() == p.ID() {
			return playerInfo
		}
	}
	return nil
}

// GetMyTeamStatus retrieve the player team status from the game server message
func (p *Player) GetMyTeamStatus(gameInfo GameInfo) Team {
	if p.TeamPlace == arena.HomeTeam {
		return gameInfo.HomeTeam
	}
	return gameInfo.AwayTeam
}

// GetOpponentTeam retrieve the opponent team status from the game server message
func (p *Player) GetOpponentTeam(status GameInfo) Team {
	if p.TeamPlace == arena.HomeTeam {
		return status.AwayTeam
	}
	return status.HomeTeam
}

// FindOpponentPlayer retrieve a specific opponent player status from the game server message
func (p *Player) FindOpponentPlayer(status GameInfo, playerNumber arena.PlayerNumber) *Player {
	teamInfo := p.GetOpponentTeam(status)
	for _, playerInfo := range teamInfo.Players {
		if playerInfo.Number == playerNumber {
			return playerInfo
		}
	}
	return nil
}

// CreateMoveOrder creates a move order
func (p *Player) CreateMoveOrder(target physics.Point, speed float64) orders.Order {
	vec := physics.NewZeroedVelocity(*physics.NewVector(p.Coords, target).Normalize())
	vec.Speed = speed
	return orders.Order{
		Type: orders.MOVE,
		Data: orders.MoveOrderData{Velocity: vec},
	}
}

// CreateJumpOrder creates a jump order (only allowed to goal keeper
func (p *Player) CreateJumpOrder(target physics.Point, speed float64) orders.Order {
	vec := physics.NewZeroedVelocity(*physics.NewVector(p.Coords, target).Normalize())
	vec.Speed = speed
	return orders.Order{
		Type: orders.MOVE,
		Data: orders.MoveOrderData{Velocity: vec},
	}
}

// CreateMoveOrderMaxSpeed creates a move order with max speed allowed
func (p *Player) CreateMoveOrderMaxSpeed(target physics.Point) orders.Order {
	return p.CreateMoveOrder(target, units.PlayerMaxSpeed)
}

// CreateStopOrder creates a move order with speed zero
func (p *Player) CreateStopOrder(direction physics.Vector) orders.Order {
	vec := p.Velocity.Copy()
	vec.Speed = 0
	vec.Direction = &direction
	return orders.NewMoveOrder(vec)
}

// CreateKickOrder creates a kick order and try to find the best vector to reach the target
func (p *Player) CreateKickOrder(target physics.Point, speed float64) orders.Order {
	ballExpectedDirection := physics.NewVector(p.LastMsg.GameInfo.Ball.Coords, target)
	diffVector := *ballExpectedDirection.Sub(p.LastMsg.GameInfo.Ball.Velocity.Direction)
	vec := physics.NewZeroedVelocity(diffVector)
	vec.Speed = speed

	return orders.NewKickOrder(vec)
}

// CreateCatchOrder creates the catch order
func (p *Player) CreateCatchOrder() orders.Order {
	return orders.NewCatchOrder()
}

// IHoldTheBall returns true when the player is holding the ball
func (p *Player) IHoldTheBall() bool {
	return p.LastMsg.GameInfo.Ball.Holder != nil && p.LastMsg.GameInfo.Ball.Holder.ID() == p.ID()
}

// OpponentGoal returns the Goal os the opponent
func (p *Player) OpponentGoal() arena.Goal {
	if p.TeamPlace == arena.HomeTeam {
		return arena.AwayTeamGoal
	}
	return arena.HomeTeamGoal
}

// DefenseGoal returns the player team goal
func (p *Player) DefenseGoal() arena.Goal {
	if p.TeamPlace == arena.HomeTeam {
		return arena.HomeTeamGoal
	}
	return arena.AwayTeamGoal
}

// IsGoalkeeper returns true if the player is the goalkeeper
func (p *Player) IsGoalkeeper() bool {
	return p.Number == arena.GoalkeeperNumber
}
