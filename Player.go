package client

import (
	"encoding/json"
	"fmt"
	"github.com/makeitplay/commons"
	"github.com/makeitplay/commons/BasicTypes"
	"github.com/makeitplay/commons/Physics"
	"github.com/makeitplay/commons/Units"
	"github.com/makeitplay/commons/talk"
	"log"
	"math"
	"os"
)

// Player acts as a brainless player in the game. This struct implements many methods that does not affect the player
// intelligence/behaviour/decisions. So, it is meant to reduce the developer concerns about communication, protocols,
// attributes, etc, and focusing in the player intelligence.
type Player struct {
	Physics.Element
	Id             string                  `json:"id"`
	Number         BasicTypes.PlayerNumber `json:"number"`
	TeamPlace      Units.TeamPlace         `json:"team_place"`
	OnMessage      func(msg GameMessage)
	OnAnnouncement func(msg GameMessage)
	config         *Configuration
	talker         *talk.Channel
	LastMsg        GameMessage
	readingWs      *commons.Task
}

// keepListening is used to keep the process running while the player is playing
var keepListening = make(chan bool)

// ID returns the player ID, that is the team place and it concatenated.
func (p *Player) ID() string {
	if p.Id == "" {
		p.Id = fmt.Sprintf("%s-%s", p.TeamPlace, p.Number)
	}
	return p.Id
}

// Start make the player start to play
func (p *Player) Start(configuration *Configuration) {
	p.config = configuration
	p.Size = Units.PlayerSize
	if p.OnAnnouncement == nil {
		log.Fatal("your player must implement the `OnAnnouncement` action")
	}
	commons.NickName = fmt.Sprintf("%s-%s", p.TeamPlace, p.Number)
	commons.Log("Try to join to the team %s ", p.TeamPlace)
	if p.initializeCommunicator() {
		p.keepPlaying()
	}
}

// LastServerMessage returns the last message got from the server
func (p *Player) LastServerMessage() GameMessage {
	return p.LastMsg
}

// SendOrders sends a list of orders to the game server, and includes a message to them (only displayed in the game server log)
func (p *Player) SendOrders(message string, orders ...BasicTypes.Order) {
	msg := PlayerMessage{
		BasicTypes.ORDER,
		orders,
		message,
	}
	stringed, err := json.Marshal(msg)
	if err != nil {
		commons.LogError("Fail generating JSON: %s", err.Error())
		return
	}

	err = p.talker.Send(stringed)
	if err != nil {
		commons.LogError("Fail on sending message: %s", err.Error())
		return
	}
}

// keepPlaying keeps the process running while the player plays
func (p *Player) keepPlaying() {
	commons.RegisterCleaner("Stopping to play", p.stopToPlay)
	for stillUp := range keepListening {
		if !stillUp {
			os.Exit(0)
		}
	}
}

// stopToPlay stop the player to play
func (p *Player) stopToPlay(interrupted bool) {
	keepListening <- false
}

// UpdatePosition update the player status after the last game server message
func (p *Player) UpdatePosition(gameInfo GameInfo) {
	status := p.FindMyStatus(gameInfo)
	if status == nil {
		// sometimes the player gets a message before his welcome message be processed, then he is not officially in the game,
		// so, this status is not available yet.
		return
	}

	p.Velocity = status.Velocity
	p.Coords = status.Coords
}

// FindMyStatus retrieve the player status from the game server message
func (p *Player) FindMyStatus(gameInfo GameInfo) *Player {
	myteamInfo := p.FindMyTeamStatus(gameInfo)
	for _, playerInfo := range myteamInfo.Players {
		if playerInfo.ID() == p.ID() {
			return playerInfo
		}
	}
	return nil
}

// FindMyTeamStatus retrieve the player team status from the game server message
func (p *Player) FindMyTeamStatus(gameInfo GameInfo) Team {
	if p.TeamPlace == Units.HomeTeam {
		return gameInfo.HomeTeam
	}
	return gameInfo.AwayTeam
}

// GetOpponentTeam retrieve the opponent team status from the game server message
func (p *Player) GetOpponentTeam(status GameInfo) Team {
	if p.TeamPlace == Units.HomeTeam {
		return status.AwayTeam
	}
	return status.HomeTeam
}

// GetOpponentPlayer retrieve a specific opponent player status from the game server message
func (p *Player) GetOpponentPlayer(status GameInfo, playerNumber BasicTypes.PlayerNumber) *Player {
	teamInfo := p.GetOpponentTeam(status)
	for _, playerInfo := range teamInfo.Players {
		if playerInfo.Number == playerNumber {
			return playerInfo
		}
	}
	return nil
}

// CreateMoveOrder creates a move order
func (p *Player) CreateMoveOrder(target Physics.Point, speed float64) BasicTypes.Order {
	vec := Physics.NewZeroedVelocity(*Physics.NewVector(p.Coords, target).Normalize())
	vec.Speed = speed
	return BasicTypes.Order{
		Type: BasicTypes.MOVE,
		Data: BasicTypes.MoveOrderData{Velocity: vec},
	}
}

// CreateMoveOrderMaxSpeed creates a move order with max speed allowed
func (p *Player) CreateMoveOrderMaxSpeed(target Physics.Point) BasicTypes.Order {
	return p.CreateMoveOrder(target, Units.PlayerMaxSpeed)
}

// CreateStopOrder creates a move order with speed zero
func (p *Player) CreateStopOrder(direction Physics.Vector) BasicTypes.Order {
	vec := p.Velocity.Copy()
	vec.Speed = 0
	vec.Direction = &direction
	return BasicTypes.Order{
		Type: BasicTypes.MOVE,
		Data: BasicTypes.MoveOrderData{Velocity: vec},
	}
}

// CreateKickOrder creates a kick order and try to find the best vector to reach the target
func (p *Player) CreateKickOrder(target Physics.Point, speed float64) BasicTypes.Order {
	ballExpectedDirection := Physics.NewVector(p.LastMsg.GameInfo.Ball.Coords, target)
	diffVector := *ballExpectedDirection.Sub(p.LastMsg.GameInfo.Ball.Velocity.Direction)
	vec := Physics.NewZeroedVelocity(diffVector)
	vec.Speed = speed
	return BasicTypes.Order{
		Type: BasicTypes.KICK,
		Data: BasicTypes.KickOrderData{Velocity: vec},
	}
}

// CreateCatchOrder creates the catch order
func (p *Player) CreateCatchOrder() BasicTypes.Order {
	return BasicTypes.Order{
		Type: BasicTypes.CATCH,
		Data: map[string]interface{}{},
	}
}

// IHoldTheBall returns true when the player is holding the ball
func (p *Player) IHoldTheBall() bool {
	return p.LastMsg.GameInfo.Ball.Holder != nil && p.LastMsg.GameInfo.Ball.Holder.ID() == p.ID()
}

// FindNearestMate returns the nearest team player and the distance to him
func (p *Player) FindNearestMate() (distance float64, player *Player) {
	var nearestPlayer *Player
	//starting from the worst case
	nearestDistance := math.Hypot(float64(Units.CourtHeight), float64(Units.CourtWidth))
	myTeam := p.FindMyTeamStatus(p.LastMsg.GameInfo)

	for _, player := range myTeam.Players {
		distance := math.Abs(p.Coords.DistanceTo(player.Coords))
		if distance <= nearestDistance && player.ID() != p.ID() {
			nearestDistance = distance
			nearestPlayer = player
		}
	}
	return nearestDistance, nearestPlayer
}

// OpponentGoal returns the Goal os the opponent
func (p *Player) OpponentGoal() BasicTypes.Goal {
	if p.TeamPlace == Units.HomeTeam {
		return commons.AwayTeamGoal
	}
	return commons.HomeTeamGoal
}

// DefenseGoal returns the player team goal
func (p *Player) DefenseGoal() BasicTypes.Goal {
	if p.TeamPlace == Units.HomeTeam {
		return commons.HomeTeamGoal
	}
	return commons.AwayTeamGoal
}

// IsGoalkeeper returns true if the player is the goalkeeper
func (p *Player) IsGoalkeeper() bool {
	return p.Number == commons.GoalkeeperNumber
}
