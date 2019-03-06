package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/arena/talk"
	"github.com/makeitplay/arena/units"
	"github.com/sirupsen/logrus"
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
	talker         talk.Talker
	talkerCtx      context.Context
	LastMsg        GameMessage
	logger         *logrus.Entry
}

// playerCtx is used to keep the process running while the player is playing
var playerCtx context.Context
var stopPlayer context.CancelFunc

// ID returns the player ID, that is the team place and it concatenated.
func (p *Player) ID() string {
	if p.Id == "" {
		p.Id = fmt.Sprintf("%s-%s", p.TeamPlace, p.Number)
	}
	return p.Id
}

// Start make the player start to play
func (p *Player) Start(logger *logrus.Logger, configuration *Configuration) {
	p.config = configuration
	playerCtx, stopPlayer = context.WithCancel(context.Background())
	p.Size = units.PlayerSize
	if p.OnAnnouncement == nil {
		log.Fatal("your player must implement the `OnAnnouncement` action")
	}

	p.logger = logger.WithField("player", fmt.Sprintf("%s-%s", p.TeamPlace, p.Number))
	p.logger.Infof("Try to join to the team %s ", p.TeamPlace)
	if !p.initializeCommunicator(p.logger) {
		return
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	exitCode := 0
	select {
	case <-signalChan:
		p.logger.Print("*********** INTERRUPTION SIGNAL ****************")
		p.talker.Close()
	case <-playerCtx.Done():
		p.logger.Info("player stopped")
		p.talker.Close()
	case <-p.talkerCtx.Done():
		p.logger.Warn("communication interrupted")
		exitCode = 1
	}
	os.Exit(exitCode)
}

// LastServerMessage returns the last message got from the server
func (p *Player) LastServerMessage() GameMessage {
	return p.LastMsg
}

// SendOrders sends a list of orders to the game server, and includes a message to them (only displayed in the game server log)
func (p *Player) SendOrders(message string, orders ...arena.Order) {
	msg := PlayerMessage{
		arena.ORDER,
		orders,
		message,
	}
	stringed, err := json.Marshal(msg)
	if err != nil {
		p.logger.Errorf("Fail generating JSON: %s", err.Error())
		return
	}

	err = p.talker.Send(stringed)
	if err != nil {
		p.logger.Errorf("Fail on sending message: %s", err.Error())
		return
	}
}

// stopToPlay stop the player to play
func (p *Player) stopToPlay(interrupted bool) {
	stopPlayer()
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
func (p *Player) CreateMoveOrder(target physics.Point, speed float64) arena.Order {
	vec := physics.NewZeroedVelocity(*physics.NewVector(p.Coords, target).Normalize())
	vec.Speed = speed
	return arena.Order{
		Type: arena.MOVE,
		Data: arena.MoveOrderData{Velocity: vec},
	}
}

// CreateJumpOrder creates a jump order (only allowed to goal keeper
func (p *Player) CreateJumpOrder(target physics.Point, speed float64) arena.Order {
	vec := physics.NewZeroedVelocity(*physics.NewVector(p.Coords, target).Normalize())
	vec.Speed = speed
	return arena.Order{
		Type: arena.MOVE,
		Data: arena.MoveOrderData{Velocity: vec},
	}
}

// CreateMoveOrderMaxSpeed creates a move order with max speed allowed
func (p *Player) CreateMoveOrderMaxSpeed(target physics.Point) arena.Order {
	return p.CreateMoveOrder(target, units.PlayerMaxSpeed)
}

// CreateStopOrder creates a move order with speed zero
func (p *Player) CreateStopOrder(direction physics.Vector) arena.Order {
	vec := p.Velocity.Copy()
	vec.Speed = 0
	vec.Direction = &direction
	return arena.Order{
		Type: arena.MOVE,
		Data: arena.MoveOrderData{Velocity: vec},
	}
}

// CreateKickOrder creates a kick order and try to find the best vector to reach the target
func (p *Player) CreateKickOrder(target physics.Point, speed float64) arena.Order {
	ballExpectedDirection := physics.NewVector(p.LastMsg.GameInfo.Ball.Coords, target)
	diffVector := *ballExpectedDirection.Sub(p.LastMsg.GameInfo.Ball.Velocity.Direction)
	vec := physics.NewZeroedVelocity(diffVector)
	vec.Speed = speed
	return arena.Order{
		Type: arena.KICK,
		Data: arena.KickOrderData{Velocity: vec},
	}
}

// CreateCatchOrder creates the catch order
func (p *Player) CreateCatchOrder() arena.Order {
	return arena.Order{
		Type: arena.CATCH,
		Data: map[string]interface{}{},
	}
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
