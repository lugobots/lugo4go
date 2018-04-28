package Game

import (
	"os"
	"fmt"
	"math"
	"github.com/maketplay/commons"
	"github.com/maketplay/commons/Physics"
	"github.com/maketplay/commons/BasicTypes"
	"net/url"
	"github.com/maketplay/commons/Units"
	"encoding/json"
	"github.com/maketplay/commons/GameState"
	"github.com/gorilla/websocket"
	"strconv"
)

type Player struct {
	Physics.Element
	Id        int                `json:"id"`
	Number    Units.PlayerNumber `json:"number"`
	TeamPlace Units.TeamPlace    `json:"team_place"`
	state     PlayerState
	config    *Configuration
	GameConn  *websocket.Conn
	//OutputCom *commons.Comunicator
	//InputCom  *commons.Comunicator
	//channel   string
	lastMsg   GameMessage
	readingWs *commons.Task
}

var keepListenning = make(chan bool)

func (p *Player) Start(configuration *Configuration) {
	p.config = configuration
	p.Id = os.Getpid() //easy way to create a unique ID since the player UUID is not public
	p.TeamPlace = configuration.TeamPlace
	commons.Log("Try to join to the team %s ", p.TeamPlace)
	p.initializeCommunicator()
	commons.NickName = fmt.Sprintf("%s-%d", p.TeamPlace, p.Id)
	//p.askToPlay()
	p.keepPlaying()
}

func (p *Player) initializeCommunicator() {
	uri := new(url.URL)
	uri.Scheme = "ws"
	uri.Host = "localhost:8080"
	uri.Path = fmt.Sprintf("/announcements/%s/%s", p.config.Uuid, p.TeamPlace)

	var err error
	p.GameConn, _, err = websocket.DefaultDialer.Dial(uri.String(), nil)
	if err != nil {
		commons.Log("Fail on dial: %s", err.Error())
	}
	commons.RegisterCleaner("Websocket connection", func(interrupted bool) {
		err := p.GameConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			commons.LogError("Fail on closing ws connection: %s", err.Error())
		}
		p.readingWs.Stop()
		p.GameConn.Close()
	})

	p.GameConn.SetCloseHandler(func(code int, text string) error {
		p.readingWs.Stop()
		if code == websocket.CloseNormalClosure {
			commons.Log("Game has closed the websocket connection")
		} else {
			commons.LogError("Lost webscket connection with the game: msgType %d: %s", code, text)
		}
		commons.Cleanup(true)
		os.Exit(0)
		return nil
	})

	go p.websocketListenner()
}

func (p *Player) ResetPosition() {
	if p.TeamPlace == Units.HomeTeam {
		p.Coords = Units.InitialPostionHomeTeam[p.Number]
	} else {
		p.Coords = Units.InitialPostionAwayTeam[p.Number]
	}
}

func (p *Player) onMessage(msg GameMessage) {
	p.lastMsg = msg
	switch msg.Type {
	case BasicTypes.WELCOME:
		commons.LogInfo("Accepted by the game server")
		if myId, ok := msg.Data["id"]; ok {
			i, err := strconv.Atoi(myId)
			if err != nil {
				commons.LogError("Invalid player id: %v", err.Error())
				panic("Invalid player id")
			}
			p.Id = i
		} else {
			commons.LogError("Player id missing in the welcome message")
			panic("Player id missing in the welcome message")
		}
		p.updatePostion(p.lastMsg.GameInfo)
		p.Number = p.findMyStatus(msg.GameInfo).Number
	case BasicTypes.ANNOUNCEMENT:
		commons.LogAnn("ANN %s", string(msg.State))
		switch GameState.State(msg.State) {
		case GameState.GETREADY:
			p.updatePostion(p.lastMsg.GameInfo)
			p.Number = p.findMyStatus(msg.GameInfo).Number
		case GameState.LISTENING:
			p.updatePostion(p.lastMsg.GameInfo)
			p.state = p.determineMyState()
			commons.LogDebug("State: %s", p.state)
			p.madeAMove()
		}
	case BasicTypes.RIP:
		commons.LogError("The server has stopped :/")
		commons.Cleanup(true)
		os.Exit(0)
	}
}

func (p *Player) sendOrders(message string, orders ...BasicTypes.Order) {
	msg := PlayerMessage{
		BasicTypes.ORDER,
		orders,
		message,
	}
	jsonsified, _ := json.Marshal(msg)

	err := p.GameConn.WriteMessage(websocket.TextMessage, jsonsified)
	if err != nil {
		commons.LogError("Fail on sending message: %s", err.Error())
		return
	}
}

func (p *Player) askToPlay() {
	data := BasicTypes.Order{
		Type: BasicTypes.ENTER,
		Data: nil,
	}

	p.sendOrders("Let me play", data)
}

func (p *Player) keepPlaying() {
	commons.RegisterCleaner("Stopping to play", p.stopsPlayer)
	for stillUp := range keepListenning {
		if !stillUp {
			os.Exit(0)
		}
	}
}

func (p *Player) stopsPlayer(interrupted bool) {
	keepListenning <- false
}

func (p *Player) madeAMove() {
	var orders []BasicTypes.Order
	var msg string

	switch p.state {
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
	if p.TeamPlace == Units.HomeTeam {
		p.Coords = p.findMyStatus(status).Coords
	} else {
		p.Coords = p.findMyStatus(status).Coords
	}
}

func (p *Player) findMyStatus(gameInfo GameInfo) *Player {
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

func (p *Player) playerEDecision() (string, []BasicTypes.Order) {
	var orders []BasicTypes.Order
	if p.IHoldTheBall() {
		goalDistance := p.Coords.DistanceTo(p.offenseGoalCoods())
		if int(math.Abs(goalDistance)) < BallMaxDistance() {
			orders = []BasicTypes.Order{p.createKickOrder(p.offenseGoalCoods())}
		} else {
			orders = []BasicTypes.Order{p.orderAdvance()}
		}
	} else {
		ballDistance := p.Coords.DistanceTo(p.lastMsg.GameInfo.Ball.Coords)
		if ballDistance < Units.DistanceCatchBall {
			orders = make([]BasicTypes.Order, 2)
			orders[0] = BasicTypes.Order{
				Type: BasicTypes.CATCH,
				Data: map[string]interface{}{
				},
			}
			orders[1] = p.orderAdvance()
		} else {
			orders = make([]BasicTypes.Order, 1)
			orders[0] = p.createMoveOrder(p.lastMsg.GameInfo.Ball.Coords)
		}
	}
	return "Gerenic order", orders
}

func (p *Player) createMoveOrder(target Physics.Point) BasicTypes.Order {
	return BasicTypes.Order{
		Type: BasicTypes.MOVE,
		Data: map[string]interface{}{
			"x": target.PosX,
			"y": target.PosY,
		},
	}
}

func (p *Player) createKickOrder(target Physics.Point) BasicTypes.Order {
	return BasicTypes.Order{
		Type: BasicTypes.KICK,
		Data: map[string]interface{}{
			"x": target.PosX,
			"y": target.PosY,
		},
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

func (p *Player) determineMyState() PlayerState {
	var isOnMyField bool
	var subState string
	var ballPossess string

	if p.lastMsg.GameInfo.Ball.Holder == nil {
		ballPossess = "dsp" //disputing
		subState = "fbl"    //far
		if int(math.Abs(p.Coords.DistanceTo(p.lastMsg.GameInfo.Ball.Coords))) <= DistanceNearBall {
			subState = "nbl" //near
		}
	} else if p.lastMsg.GameInfo.Ball.Holder.TeamPlace == p.TeamPlace {
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

	if p.TeamPlace == Units.HomeTeam {
		isOnMyField = p.lastMsg.GameInfo.Ball.Coords.PosX <= Units.CourtWidth/2
	} else {
		isOnMyField = p.lastMsg.GameInfo.Ball.Coords.PosX >= Units.CourtWidth/2
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
	if p.TeamPlace == Units.AwayTeam {
		myRagion = MirrorRegion(myRagion)
	}
	return myRagion
}
func MirrorRegion(region PlayerRegion) PlayerRegion {
	return PlayerRegion{
		Units.MirrorCoordToAway(region.CornerB), // have to switch the corner because the convention for Regions
		Units.MirrorCoordToAway(region.CornerA),
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

func (p *Player) offenseGoalCoods() Physics.Point {
	if p.TeamPlace == Units.HomeTeam {
		return Units.AwayTeamGoalcenter
	} else {
		return Units.HomeTeamGoalcenter
	}

}

func (p *Player) deffenseGoalCoods() Physics.Point {
	if p.TeamPlace == Units.HomeTeam {
		return Units.HomeTeamGoalcenter
	} else {
		return Units.AwayTeamGoalcenter
	}
}
func (p *Player) websocketListenner() {
	p.readingWs = commons.NewTask(func(task *commons.Task) {
		msgType, message, err := p.GameConn.ReadMessage()
		if msgType == -1 {
			return
		} else if err != nil {
			commons.LogError("Fail reading websocket message (%d): %s", msgType, err)
		} else {
			var msg GameMessage
			err = json.Unmarshal(message, &msg)
			if err != nil {
				commons.LogError("Fail on convert wb message: %s", err.Error())
			} else {
				p.onMessage(msg)
			}
		}
	})
	p.readingWs.Start()
}
