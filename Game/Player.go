package Game

import (
	"os"
	"math"
	"github.com/makeitplay/commons"
	"github.com/makeitplay/commons/Physics"
	"github.com/makeitplay/commons/BasicTypes"
	"github.com/makeitplay/commons/Units"
	"encoding/json"
	"github.com/makeitplay/commons/talk"
	"fmt"
	"log"
)

type Player struct {
	Physics.Element
	Id             int                     `json:"id"`
	Number         BasicTypes.PlayerNumber `json:"number"`
	TeamPlace      Units.TeamPlace         `json:"team_place"`
	OnMessage      func(msg GameMessage)
	OnAnnouncement func(msg GameMessage)
	config         *Configuration
	talker         *talk.Channel
	LastMsg        GameMessage
	readingWs      *commons.Task
}

var keepListening = make(chan bool)

func (p *Player) Start(configuration *Configuration) {
	p.config = configuration
	if p.OnAnnouncement == nil {
		log.Fatal("your player must implement the `OnAnnouncement` action")
	}
	commons.NickName = fmt.Sprintf("%s-%s", p.TeamPlace, p.Number)
	commons.Log("Try to join to the team %s ", p.TeamPlace)
	p.initializeCommunicator()
	p.keepPlaying()
}

func (p *Player) LastServerMessage() GameMessage {
	return p.LastMsg
}

func (p *Player) SendOrders(message string, orders ...BasicTypes.Order) {
	commons.LogDebug(message)
	msg := PlayerMessage{
		BasicTypes.ORDER,
		orders,
		message,
	}
	jsonsified, err := json.Marshal(msg)
	if err != nil {
		commons.LogError("Fail generating JSON: %s", err.Error())
		return
	}

	err = p.talker.Send(jsonsified)
	if err != nil {
		commons.LogError("Fail on sending message: %s", err.Error())
		return
	}
}

func (p *Player) keepPlaying() {
	commons.RegisterCleaner("Stopping to play", p.stopToPlay)
	for stillUp := range keepListening {
		if !stillUp {
			os.Exit(0)
		}
	}
}

func (p *Player) stopToPlay(interrupted bool) {
	keepListening <- false
}

func (p *Player) UpdatePosition(gameInfo GameInfo) {
	status := p.FindMyStatus(gameInfo)
	p.Velocity = status.Velocity
	p.Coords = status.Coords
}

func (p *Player) FindMyStatus(gameInfo GameInfo) *Player {
	return p.GetMyTeam(gameInfo).Players[p.Id]
}

func (p *Player) GetMyTeam(gameInfo GameInfo) Team {
	if p.TeamPlace == Units.HomeTeam {
		return gameInfo.HomeTeam
	} else {
		return gameInfo.AwayTeam
	}
}

func (p *Player) GetOpponentTeam(status GameInfo) Team {
	if p.TeamPlace == Units.HomeTeam {
		return status.AwayTeam
	} else {
		return status.HomeTeam
	}
}

func (p *Player) CreateMoveOrder(target Physics.Point) BasicTypes.Order {
	vec := Physics.NewZeroedVelocity(*Physics.NewVector(p.Coords, target).Normalize())
	vec.Speed = Units.PlayerMaxSpeed
	return BasicTypes.Order{
		Type: BasicTypes.MOVE,
		Data: BasicTypes.MoveOrderData{Velocity: vec},
	}
}

func (p *Player) CreateStopOrder(direction Physics.Vector) BasicTypes.Order {
	vec := p.Velocity.Copy()
	vec.Speed = 0
	vec.Direction = &direction
	return BasicTypes.Order{
		Type: BasicTypes.MOVE,
		Data: BasicTypes.MoveOrderData{Velocity: vec},
	}
}

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

func (p *Player) CreateCatchOrder() BasicTypes.Order {
	return BasicTypes.Order{
		Type: BasicTypes.CATCH,
		Data: map[string]interface{}{
		},
	}
}

func (p *Player) IHoldTheBall() bool {
	return p.LastMsg.GameInfo.Ball.Holder != nil && p.LastMsg.GameInfo.Ball.Holder.Id == p.Id
}

func (p *Player) FindNearestMate() (distance float64, player *Player) {
	var nearestPlayer *Player
	//starting from the worst case
	nearestDistance := math.Hypot(float64(Units.CourtHeight), float64(Units.CourtWidth))
	myTeam := p.GetMyTeam(p.LastMsg.GameInfo)

	for playerId, player := range myTeam.Players {
		distance := math.Abs(p.Coords.DistanceTo(player.Coords))
		if distance <= nearestDistance && playerId != p.Id {
			nearestDistance = distance
			nearestPlayer = player
		}
	}
	return nearestDistance, nearestPlayer
}

func (p *Player) OpponentGoal() BasicTypes.Goal {
	if p.TeamPlace == Units.HomeTeam {
		return commons.AwayTeamGoal
	} else {
		return commons.HomeTeamGoal
	}
}

func (p *Player) DefenseGoal() BasicTypes.Goal {
	if p.TeamPlace == Units.HomeTeam {
		return commons.HomeTeamGoal
	} else {
		return commons.AwayTeamGoal
	}
}
