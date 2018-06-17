package Game

import (
	"github.com/makeitplay/commons/BasicTypes"
	"github.com/makeitplay/commons"
	"math"
	"github.com/makeitplay/commons/Units"
)
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

func (p *Player) TakeAnAction() {
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
	commons.LogDebug("Sending order")
	p.sendOrders(msg, orders...)

}
