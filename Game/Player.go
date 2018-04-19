package Game

import (
	"github.com/rubens21/test1/App"
	"os"
	"fmt"
	"math"
	"github.com/rubens21/test1/Lib/Physics"
)

type Player struct {
	Physics.Element
	Id              int          `json:"id"`
	Number          PlayerNumber `json:"number"`
	TeamName        TeamName     `json:"team_name"`
	State           PlayerState  `json:"state"`
	//gameComunicator *Comunicator
	channel         string
	//lastMsg         Msg
}

var keepListenning = make(chan bool)



Parou quando tava limpando para ficar so o que interessa para o player



func (p *Player) Start(teamName TeamName) {
	p.Id = os.Getpid()
	p.TeamName = teamName
	//p.gameComunicator = CreateComunicator(MAINTOPIC, p.onMessage)
	App.NickName = fmt.Sprintf("%s-%d", p.TeamName, p.Id)
	p.askToPlay()
	p.keepPlaying()
}

func (p *Player) ResetPosition() {
	if p.TeamName == HomeTeam {
		p.Coords = InitialPostionHomeTeam[p.Number]
	} else {
		p.Coords = InitialPostionAwayTeam[p.Number]
	}
}

func (p *Player) onMessage(msg Msg) {
	p.lastMsg = msg

	switch msg.Type {
	case ANNOUNCEMENT:
		switch GameState(msg.State) {
		case GETREADY:
			p.updatePostion(p.lastMsg.GameInfo)
			p.Number = p.findMyStatus(msg.GameInfo).Number
		case LISTENING:
			p.updatePostion(p.lastMsg.GameInfo)
			p.State = p.determineMyState()
			p.madeAMove()
		}
	case RIP:
		App.Log("Sorry, guys! I'm out")
		App.Cleanup()
		os.Exit(0)
	}
}

func (p *Player) sendOrders(message string, orders ...Order) {
	p.gameComunicator.Send(Msg{
		ORDER,
		GameInfo{},
		State(p.State),
		p.Id,
		orders,
		message,
	})
}

func (p *Player) askToPlay() {
	data := Order{ENTER, map[string]interface{}{
		"teamName": p.TeamName,
		"id":       p.Id,
	}}
	App.Log("Try to join to the team %s with the id %d", p.TeamName, p.Id)
	p.sendOrders("Let me play", data)
}

func (p *Player) keepPlaying() {
	App.RegisterCleaner(p.stopsPlayer)
	for stillUp := range keepListenning {
		if !stillUp {
			os.Exit(0)
		}
	}
}

func (p *Player) stopsPlayer() {
	App.Log("Stopping player")
	keepListenning <- false
}

func (p *Player) madeAMove() {
	var orders []Order
	var msg string

	switch p.State {
	case AtckHoldHse:
		msg, orders = p.orderForAtckHoldHse()
	case AtckHoldFrg:
		msg, orders = p.orderForAtckHoldFrg()
	case AtckHelpHse:
		msg, orders = p.orderForAtckHelpHse()
	case AtckHelpFrg:
		msg, orders = p.orderForAtckHelpFrg()
	case DefdMyrgHse:
		msg, orders = p.orderForDefdMyrgHse()
		orders = append(orders, p.createCatchOrder())
	case DefdMyrgFrg:
		msg, orders = p.orderForDefdMyrgFrg()
		orders = append(orders, p.createCatchOrder())
	case DefdOtrgHse:
		msg, orders = p.orderForDefdOtrgHse()
		orders = append(orders, p.createCatchOrder())
	case DefdOtrgFrg:
		msg, orders = p.orderForDefdOtrgFrg()
		orders = append(orders, p.createCatchOrder())
	case DsptNfblHse:
		msg, orders = p.orderForDsptNfblHse()
		orders = append(orders, p.createCatchOrder())
	case DsptNfblFrg:
		msg, orders = p.orderForDsptNfblFrg()
		orders = append(orders, p.createCatchOrder())
	case DsptFrblHse:
		msg, orders = p.orderForDsptFrblHse()
		orders = append(orders, p.createCatchOrder())
	case DsptFrblFrg:
		msg, orders = p.orderForDsptFrblFrg()
		orders = append(orders, p.createCatchOrder())
	}
	p.sendOrders(msg, orders...)

}

func (p *Player) updatePostion(status GameInfo) {
	if p.TeamName == HomeTeam {
		p.Coords = p.findMyStatus(status).Coords
	} else {
		p.Coords = p.findMyStatus(status).Coords
	}
}

func (p *Player) findMyStatus(gameInfo GameInfo) *Player {
	return p.findMyTeam(gameInfo).Players[p.Id]
}

func (p *Player) findMyTeam(gameInfo GameInfo) Team {
	if p.TeamName == HomeTeam {
		return gameInfo.HomeTeam
	} else {
		return gameInfo.AwayTeam
	}
}

func (p *Player) findOpponentTeam(status GameInfo) Team {
	if p.TeamName == HomeTeam {
		return status.AwayTeam
	} else {
		return status.HomeTeam
	}
}

func (p *Player) playerEDecision() (string, []Order) {
	var orders []Order
	if p.IHoldTheBall() {
		goalDistance := p.Coords.DistanceTo(p.offenseGoalCoods())
		if int(math.Abs(goalDistance)) < BallMaxDistance() {
			orders = []Order{p.createKickOrder(p.offenseGoalCoods())}
		} else {
			orders = []Order{p.orderAdvance()}
		}
	} else {
		ballDistance := p.Coords.DistanceTo(p.lastMsg.GameInfo.Ball.Coords)
		if ballDistance < DISTANCE_CATCH_BALL {
			orders = make([]Order, 2)
			orders[0] = Order{
				CATCH,
				map[string]interface{}{
				},
			}
			orders[1] = p.orderAdvance()
		} else {
			orders = make([]Order, 1)
			orders[0] = p.createMoveOrder(p.lastMsg.GameInfo.Ball.Coords)
		}
	}
	return "Gerenic order", orders
}

func (p *Player) createMoveOrder(target Physics.Point) Order {
	return Order{
		MOVE,
		map[string]interface{}{
			"x": target.PosX,
			"y": target.PosY,
		},
	}
}

func (p *Player) createKickOrder(target Physics.Point) Order {
	return Order{
		KICK,
		map[string]interface{}{
			"x": target.PosX,
			"y": target.PosY,
		},
	}
}

func (p *Player) createCatchOrder() Order {
	return Order{
		CATCH,
		map[string]interface{}{
		},
	}
}

func (p *Player) IHoldTheBall() bool {
	return p.lastMsg.GameInfo.Ball.Holder != nil && p.lastMsg.GameInfo.Ball.Holder.Id == p.Id
}

func (p *Player) determineMyState() PlayerState {
	var isOnMyField bool
	var subState string
	var ballPossess string

	if p.lastMsg.GameInfo.Ball.Holder == nil {
		ballPossess = "dsp" //disputing
		subState = "fbl"    //far
		if int(math.Abs(p.Coords.DistanceTo(p.lastMsg.GameInfo.Ball.Coords))) <= DISTANCE_NEAR_BALL {
			subState = "nbl" //near
		}
	} else if p.lastMsg.GameInfo.Ball.Holder.TeamName == p.TeamName {
		ballPossess = "atk" //attacking
		subState = "hlp"    //helping
		if p.lastMsg.GameInfo.Ball.Holder.Id == p.Id {
			subState = "hld" //holdin
		}
	} else {
		ballPossess = "dfd"
		subState = "org"
		if p.isItInMyRegion(p.lastMsg.GameInfo.Ball.Coords) {
			subState = "mrg"
		}
	}

	if p.TeamName == HomeTeam {
		isOnMyField = p.lastMsg.GameInfo.Ball.Coords.PosX <= COURT_WIDTH/2
	} else {
		isOnMyField = p.lastMsg.GameInfo.Ball.Coords.PosX >= COURT_WIDTH/2
	}
	fieldState := "fr"
	if isOnMyField {
		fieldState = "hs"
	}
	return PlayerState(ballPossess + "-" + subState + "-" + fieldState)
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

func (p *Player) myRegion() PlayerRegion {
	myRagion := HomePlayersRegions[p.Number]
	if p.TeamName == AwayTeam {
		myRagion = MirrorRegion(myRagion)
	}
	return myRagion
}

func (p *Player) findNearestMate() (distance float64, player *Player) {
	var nearestPlayer *Player
	//starting from the worst case
	nearestDistance := math.Hypot(float64(COURT_HEIGHT), float64(COURT_WIDTH))
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

func (p *Player) offenseGoalCoods() Physics.Point {
	if p.TeamName == HomeTeam {
		return AWAY_TEAM_GOALCENTER
	} else {
		return HOME_TEAM_GOALCENTER
	}

}

func (p *Player) deffenseGoalCoods() Physics.Point {
	if p.TeamName == HomeTeam {
		return HOME_TEAM_GOALCENTER
	} else {
		return AWAY_TEAM_GOALCENTER
	}

}
