package Game

import (
	"os"
	"math"
	"github.com/makeitplay/commons"
	"github.com/makeitplay/commons/Physics"
	"github.com/makeitplay/commons/BasicTypes"
	"github.com/makeitplay/commons/Units"
	"encoding/json"
	"github.com/makeitplay/client-player-go/tatics"
	"github.com/makeitplay/commons/talk"
	"fmt"
)

type Player struct {
	Physics.Element
	Id        int                `json:"id"`
	Number    BasicTypes.PlayerNumber `json:"number"`
	TeamPlace Units.TeamPlace    `json:"team_place"`
	state     PlayerState
	config    *Configuration
	talker *talk.Channel
	lastMsg   GameMessage
	readingWs *commons.Task
}

var keepListening = make(chan bool)

func (p *Player) Start(configuration *Configuration) {
	p.config = configuration
	p.TeamPlace = configuration.TeamPlace
	p.Number = configuration.PlayerNumber
	commons.NickName = fmt.Sprintf("%s-%s", p.TeamPlace, p.Number)
	commons.Log("Try to join to the team %s ", p.TeamPlace)
	p.initializeCommunicator()
	p.keepPlaying()
}


func (p *Player) ResetPosition() {
	region := p.myRegion()
	p.Coords = region.InitialPosition()
}


func (p *Player) sendOrders(message string, orders ...BasicTypes.Order) {
	msg := PlayerMessage{
		BasicTypes.ORDER,
		orders,
		message,
	}
	jsonsified, _ := json.Marshal(msg)

	err := p.talker.Send(jsonsified)
	commons.LogDebug("==== ORDER SENT === ")
	if err != nil {
		commons.LogError("Fail on sending message: %s", err.Error())
		return
	}
}

func (p *Player) keepPlaying() {
	commons.RegisterCleaner("Stopping to play", p.stopsPlayer)
	for stillUp := range keepListening {
		if !stillUp {
			os.Exit(0)
		}
	}
}

func (p *Player) stopsPlayer(interrupted bool) {
	keepListening <- false
}



func (p *Player) updatePosition(status GameInfo) {
	if p.TeamPlace == Units.HomeTeam {
		p.Coords = p.FindMyStatus(status).Coords
	} else {
		p.Coords = p.FindMyStatus(status).Coords
	}
}

func (p *Player) FindMyStatus(gameInfo GameInfo) *Player {
	return p.findMyTeam(gameInfo).Players[p.Id]
}

func (p *Player) findMyTeam(gameInfo GameInfo) Team {
	if p.TeamPlace == Units.HomeTeam {
		return gameInfo.HomeTeam
	} else {
		return gameInfo.AwayTeam
	}
}

func (p *Player) findOpponentTeam(status GameInfo) Team {
	if p.TeamPlace == Units.HomeTeam {
		return status.AwayTeam
	} else {
		return status.HomeTeam
	}
}


func (p *Player) createMoveOrder(target Physics.Point) BasicTypes.Order {
	vec := Physics.NewZeroedVelocity(*Physics.NewVector(p.Coords, target))
	vec.Speed = Units.PlayerMaxSpeed
	return BasicTypes.Order{
		Type: BasicTypes.MOVE,
		Data: BasicTypes.MoveOrderData{Velocity: vec},
	}
}

func (p *Player) createKickOrder(target Physics.Point) BasicTypes.Order {
	vec := Physics.NewZeroedVelocity(*Physics.NewVector(p.Coords, target).Normalize())
	vec.Speed = Units.BallMaxSpeed
	return BasicTypes.Order{
		Type: BasicTypes.KICK,
		Data: BasicTypes.KickOrderData{Velocity: vec},
	}
}

func (p *Player) createCatchOrder() BasicTypes.Order {
	return BasicTypes.Order{
		Type: BasicTypes.CATCH,
		Data: map[string]interface{}{
		},
	}
}

func (p *Player) IHoldTheBall() bool {
	return p.lastMsg.GameInfo.Ball.Holder != nil && p.lastMsg.GameInfo.Ball.Holder.Id == p.Id
}


func (p *Player) isItInMyRegion(coords Physics.Point) bool {
	myRagion := p.myRegion()
	isInX := coords.PosX >= myRagion.CornerA.PosX && coords.PosX <= myRagion.CornerB.PosX
	isInY := coords.PosY >= myRagion.CornerA.PosY && coords.PosY <= myRagion.CornerB.PosY
	return isInX && isInY
}

func (p *Player) myRegionCenter() Physics.Point {
	myRegiao := p.myRegion()
	//regionDiagonal := math.Hypot(float64(myRegiao.CornerA.PosX), float64(myRegiao.CornerB.PosY))
	halfXDistance := (myRegiao.CornerB.PosX - myRegiao.CornerA.PosX) / 2
	halfYDistance := (myRegiao.CornerB.PosY - myRegiao.CornerA.PosY) / 2
	return Physics.Point{
		PosX: int(myRegiao.CornerA.PosX + halfXDistance),
		PosY: int(myRegiao.CornerA.PosY + halfYDistance),
	}
}

func (p *Player) myRegion() tatics.PlayerRegion {
	myRagion := tatics.HomePlayersRegions[p.Number]
	if p.TeamPlace == Units.AwayTeam {
		myRagion = MirrorRegion(myRagion)
	}
	return myRagion
}
func MirrorRegion(region tatics.PlayerRegion) tatics.PlayerRegion {
	return tatics.PlayerRegion{
		CornerA: tatics.MirrorCoordToAway(region.CornerA), // have to switch the corner because the convention for Regions
		CornerB: tatics.MirrorCoordToAway(region.CornerB),
	}
}

func (p *Player) findNearestMate() (distance float64, player *Player) {
	var nearestPlayer *Player
	//starting from the worst case
	nearestDistance := math.Hypot(float64(Units.CourtHeight), float64(Units.CourtWidth))
	myTeam := p.findMyTeam(p.lastMsg.GameInfo)

	for playerId, player := range myTeam.Players {
		distance := math.Abs(p.Coords.DistanceTo(player.Coords))
		if distance <= nearestDistance && playerId != p.Id {
			nearestDistance = distance
			nearestPlayer = player
		}
	}
	return nearestDistance, nearestPlayer
}

func (p *Player) offenseGoalCoords() Physics.Point {
	if p.TeamPlace == Units.HomeTeam {
		return commons.AwayTeamGoal.Center
	} else {
		return commons.HomeTeamGoal.Center
	}

}

func (p *Player) defenseGoalCoords() Physics.Point {
	if p.TeamPlace == Units.HomeTeam {
		return commons.HomeTeamGoal.Center
	} else {
		return  commons.AwayTeamGoal.Center
	}
}
