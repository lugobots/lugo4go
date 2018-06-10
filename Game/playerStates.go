package Game

import (
	"math"
	"reflect"

	"github.com/makeitplay/commons/BasicTypes"
	"github.com/makeitplay/commons/Units"
	"github.com/makeitplay/commons/Physics"
)

type PlayerState BasicTypes.State

const (
	// attacking, holding the ball, home field
	AtckHoldHse PlayerState = "atk-hld-hs"
	// attacking, holding the ball, foreign field
	AtckHoldFrg PlayerState = "atk-hld-fr"
	// attacking, helping the team, home field
	AtckHelpHse PlayerState = "atk-hlp-hs"
	// attacking, helping the team, foreign field
	AtckHelpFrg PlayerState = "atk-hlp-fr"

	// defading, on my region, home field
	DefdMyrgHse PlayerState = "dfd-mrg-hs"
	// defading, on my region, foreign field
	DefdMyrgFrg PlayerState = "dfd-mrg-fr"
	// defading, on other region, home field
	DefdOtrgHse PlayerState = "dfd-org-hs"
	// defading, on other region, foreign field
	DefdOtrgFrg PlayerState = "dfd-org-fr"

	// disputing, near to the ball, home field
	DsptNfblHse PlayerState = "dsp-nbl-hs"
	// disputing, near to the ball, foreign field
	DsptNfblFrg PlayerState = "dsp-nbl-fr"
	// disputing, far to the ball, home field
	DsptFrblHse PlayerState = "dsp-fbl-hs"
	// disputing, far to the ball, foreign field
	DsptFrblFrg PlayerState = "dsp-fbl-fr"
)

type DistanceScale string

const (
	DISTANCE_SCALE_NEAR DistanceScale = "near"
	DISTANCE_SCALE_FAR  DistanceScale = "far"
	DISTANCE_SCALE_GOOD DistanceScale = "good"
)

//region Attack states
func (p *Player) orderForAtckHoldHse() (msg string, orders []BasicTypes.Order) {
	obstacles := p.watchOpponentOnMyRoute(p.offenseGoalCoods(), ERROR_MARGIN_RUNNING)

	if len(obstacles) == 0 {
		return "I am free yet", []BasicTypes.Order{p.orderAdvance()}
	} else if len(obstacles) == 1 {
		num := int(reflect.ValueOf(obstacles).MapKeys()[0].Int())
		if p.calcDistanceScale(p.findOpponentTeam(p.lastMsg.GameInfo).Players[num].Coords) != DISTANCE_SCALE_FAR {
			return "Dribble this guys (not yet)", []BasicTypes.Order{p.orderPassTheBall()}
		} else {
			return "Advance watching", []BasicTypes.Order{p.orderAdvance()}
		}
	} else {
		nearstObstacle := float64(Units.CourtWidth) //just initializing with a high value
		//num := int(reflect.ValueOf(obstacles).MapKeys()[0].Int())
		for opponentId := range obstacles {
			oppCoord := p.findOpponentTeam(p.lastMsg.GameInfo).Players[opponentId].Coords
			oppDist := p.Coords.DistanceTo(oppCoord)
			if oppDist < nearstObstacle {
				nearstObstacle = oppDist
			}
		}
		if nearstObstacle < Units.PlayerMaxSpeed*2 {
			return "I need help", []BasicTypes.Order{p.orderPassTheBall(), p.orderAdvance()}
		} else {
			return "Advance watching", []BasicTypes.Order{p.orderAdvance()}
		}
	}
}

func (p *Player) orderForAtckHoldFrg() (msg string, orders []BasicTypes.Order) {
	goalCoords := p.offenseGoalCoods()
	goalDistance := p.Coords.DistanceTo(goalCoords)
	if int(math.Abs(goalDistance)) < BallMaxDistance() {
		return "Shoot!", []BasicTypes.Order{p.createKickOrder(goalCoords)}
	} else {
		obstacles := p.watchOpponentOnMyRoute(p.offenseGoalCoods(), ERROR_MARGIN_RUNNING)

		if len(obstacles) == 0 {
			return "I am still free", []BasicTypes.Order{p.orderAdvance()}
		} else if len(obstacles) == 1 {
			num := int(reflect.ValueOf(obstacles).MapKeys()[0].Int())
			if p.calcDistanceScale(p.findOpponentTeam(p.lastMsg.GameInfo).Players[num].Coords) != DISTANCE_SCALE_FAR {
				return "Dribble this guys (not yet)", []BasicTypes.Order{p.orderPassTheBall()}
			} else {
				return "Advace watching", []BasicTypes.Order{p.orderAdvance()}
			}
		} else {
			return "I need help", []BasicTypes.Order{p.orderPassTheBall(), p.orderAdvance()}
		}

	}
}

func (p *Player) orderForAtckHelpHse() (msg string, orders []BasicTypes.Order) {
	if p.isItInMyRegion(p.Coords) {
		switch p.calcDistanceScale(p.lastMsg.GameInfo.Ball.Coords) {
		case DISTANCE_SCALE_FAR:
			msg = "Let's attack!"
			orders = []BasicTypes.Order{p.createMoveOrder(p.lastMsg.GameInfo.Ball.Coords)}
		case DISTANCE_SCALE_NEAR:
			msg = "Given space"
			opositPoint := Physics.NewVector(p.Coords, p.lastMsg.GameInfo.Ball.Coords).Invert().TargetFrom(p.Coords)
			vectorToOpositPoint := Physics.NewVector(p.Coords, p.offenseGoalCoods())
			vectorToOpositPoint.Add(Physics.NewVector(p.Coords, opositPoint))
			orders = []BasicTypes.Order{p.createMoveOrder(vectorToOpositPoint.TargetFrom(p.Coords))}
		case DISTANCE_SCALE_GOOD:
			msg = "Give me the ball!"
			orders = []BasicTypes.Order{p.createMoveOrder(p.lastMsg.GameInfo.Ball.Coords)}
		}
	} else {
		msg = "I'll be right here"
		myRegionVector := Physics.NewVector(p.Coords, p.myRegionCenter()).Invert().TargetFrom(p.Coords)
		offensivePosition := Physics.NewVector(p.Coords, p.offenseGoalCoods())
		offensivePosition.Add(Physics.NewVector(p.Coords, myRegionVector))
		orders = []BasicTypes.Order{p.createMoveOrder(offensivePosition.TargetFrom(p.Coords))}
	}
	return msg, orders
}

func (p *Player) orderForAtckHelpFrg() (msg string, orders []BasicTypes.Order) {
	if p.isItInMyRegion(p.Coords) {
		switch p.calcDistanceScale(p.lastMsg.GameInfo.Ball.Coords) {
		case DISTANCE_SCALE_FAR:
			msg = "Supporting on attack"
			orders = []BasicTypes.Order{p.createMoveOrder(p.lastMsg.GameInfo.Ball.Coords)}
		case DISTANCE_SCALE_NEAR:
			msg = "Helping on attack"

			offensiveZone := Physics.NewVector(p.Coords, p.myRegionCenter())
			offensiveZone.Add(Physics.NewVector(p.Coords, p.offenseGoalCoods()))
			orders = []BasicTypes.Order{p.createMoveOrder(offensiveZone.TargetFrom(p.Coords))}
		case DISTANCE_SCALE_GOOD:
			msg = "Holding positiong for attack"
			offensiveZone := Physics.NewVector(p.Coords, p.lastMsg.GameInfo.Ball.Coords)
			offensiveZone.Add(Physics.NewVector(p.Coords, p.offenseGoalCoods()))
			orders = []BasicTypes.Order{p.createMoveOrder(offensiveZone.TargetFrom(p.Coords))}
		}
	} else {
		regionCenter := p.myRegionCenter()
		return "Backing to my position", []BasicTypes.Order{p.createMoveOrder(regionCenter)}
	}
	return msg, orders
}

//endregion Attack states

//region Defending states
func (p *Player) orderForDefdMyrgHse() (msg string, orders []BasicTypes.Order) {
	orders = []BasicTypes.Order{p.createMoveOrder(p.lastMsg.GameInfo.Ball.Coords)}
	return "Running towards the ball", orders
}

func (p *Player) orderForDefdMyrgFrg() (msg string, orders []BasicTypes.Order) {
	switch p.calcDistanceScale(p.lastMsg.GameInfo.Ball.Coords) {
	case DISTANCE_SCALE_NEAR:
		// too close
		msg = "Pressing the player"
		orders = []BasicTypes.Order{p.createMoveOrder(p.lastMsg.GameInfo.Ball.Coords)}
	case DISTANCE_SCALE_FAR:
		//get closer
		msg = "Back to my position!"
		var backOffPos Physics.Point
		backOffPos = *p.myRegion().CentralDefense()
		orders = []BasicTypes.Order{p.createMoveOrder(backOffPos)}
	case DISTANCE_SCALE_GOOD:
		msg = "Holding positiong"
	}
	//nothing more smart than that so far. stay stopped
	return msg, orders
}

func (p *Player) orderForDefdOtrgHse() (msg string, orders []BasicTypes.Order) {

	if p.calcDistanceScale(p.lastMsg.GameInfo.Ball.Coords) == DISTANCE_SCALE_NEAR {
		msg = "Defensing while back off"
		backOffDir := Physics.NewVector(p.Coords, p.deffenseGoalCoods())
		backOffDir.Add(Physics.NewVector(p.Coords, p.lastMsg.GameInfo.Ball.Coords))
		orders = []BasicTypes.Order{p.createMoveOrder(backOffDir.TargetFrom(p.Coords))}
	} else {
		msg = "Back off!"
		backOffDir := Physics.NewVector(p.Coords, p.deffenseGoalCoods())
		backOffDir.Add(Physics.NewVector(p.Coords, p.myRegionCenter()))
		orders = []BasicTypes.Order{p.createMoveOrder(backOffDir.TargetFrom(p.Coords))}
	}
	//nothing more smart than that so far. stay stopped
	return msg, orders
}

func (p *Player) orderForDefdOtrgFrg() (msg string, orders []BasicTypes.Order) {
	if p.calcDistanceScale(p.lastMsg.GameInfo.Ball.Coords) == DISTANCE_SCALE_NEAR {
		msg = "Defensing while back off"
		backOffDir := Physics.NewVector(p.Coords, p.deffenseGoalCoods())
		backOffDir.Add(Physics.NewVector(p.Coords, p.lastMsg.GameInfo.Ball.Coords))
		orders = []BasicTypes.Order{p.createMoveOrder(backOffDir.TargetFrom(p.Coords))}
	} else {
		msg = "Back off!"
		backOffDir := Physics.NewVector(p.Coords, p.deffenseGoalCoods())
		backOffDir.Add(Physics.NewVector(p.Coords, p.myRegionCenter()))
		orders = []BasicTypes.Order{p.createMoveOrder(backOffDir.TargetFrom(p.Coords))}
	}
	//nothing more smart than that so far. stay stopped
	return msg, orders
}
//endregion Defending states

//region Disputing states
func (p *Player) orderForDsptNfblHse() (msg string, orders []BasicTypes.Order) {
	myDistance := p.Coords.DistanceTo(p.lastMsg.GameInfo.Ball.Coords)
	playerCloser := 0
	for _, teamMate := range p.findMyTeam(p.lastMsg.GameInfo).Players {
		if teamMate.Coords.DistanceTo(p.lastMsg.GameInfo.Ball.Coords) < myDistance {
			playerCloser++
			if playerCloser > 2 {
				return "Holding position for suport", []BasicTypes.Order{p.createMoveOrder(p.myRegionCenter())}
			}
		}
	}
	msg = "Disputing for the ball"
	orders = []BasicTypes.Order{p.createMoveOrder(p.lastMsg.GameInfo.Ball.Coords)}
	return msg, orders
}

func (p *Player) orderForDsptNfblFrg() (msg string, orders []BasicTypes.Order) {
	return p.orderForDsptNfblHse()
}

func (p *Player) orderForDsptFrblHse() (msg string, orders []BasicTypes.Order) {
	msg = "Try to catch the ball"
	if p.isItInMyRegion(p.lastMsg.GameInfo.Ball.Coords) {
		backOffDir := Physics.NewVector(p.Coords, p.deffenseGoalCoods())
		backOffDir.Add(Physics.NewVector(p.Coords, p.lastMsg.GameInfo.Ball.Coords))
		orders = []BasicTypes.Order{p.createMoveOrder(backOffDir.TargetFrom(p.Coords))}
	} else {
		orders = []BasicTypes.Order{p.createMoveOrder(p.myRegionCenter())}
	}
	return msg, orders
}

func (p *Player) orderForDsptFrblFrg() (msg string, orders []BasicTypes.Order) {
	msg = "Watch out the ball"
	backOffDir := Physics.NewVector(p.Coords, p.myRegionCenter())
	backOffDir.Add(Physics.NewVector(p.Coords, p.lastMsg.GameInfo.Ball.Coords))
	orders = []BasicTypes.Order{p.createMoveOrder(backOffDir.TargetFrom(p.Coords))}
	return msg, orders
}
//endregion Disputing states

//region helpers
func (p *Player) orderAdvance() BasicTypes.Order {
	return p.createMoveOrder(p.offenseGoalCoods())
}

func (p *Player) orderPassTheBall() BasicTypes.Order {
	bestCandidate := new(Player)
	bestScore := 0
	for _, playerMate := range p.findMyTeam(p.lastMsg.GameInfo).Players {
		if playerMate.Id == p.Id {
			continue
		}
		obstaclesFromMe := p.watchOpponentOnMyRoute(playerMate.Coords, ERROR_MARGIN_PASSING)
		obstaclesToGoal := playerMate.watchOpponentOnMyRoute(p.offenseGoalCoods(), ERROR_MARGIN_RUNNING)
		distanceFromMe := p.Coords.DistanceTo(playerMate.Coords)
		distanceToGoal := playerMate.Coords.DistanceTo(p.offenseGoalCoods())
		score := 1000
		score -= len(obstaclesFromMe) * 10
		score -= len(obstaclesToGoal) * 5
		score -= int(distanceFromMe * 0.5)
		score -= int(distanceToGoal * 0.5)

		//App.Log("=Player %s | %d obs, %d obs2, %d DfomMe, %d DfomGoal  = Total %d",
		//	playerMate.Number,
		//	len(obstaclesFromMe),
		//	len(obstaclesToGoal),
		//	int(distanceFromMe),
		//	int(distanceToGoal),
		//	score,
		//	)

		if score > bestScore {
			bestScore = score
			bestCandidate = playerMate
		}
	}

	//App.Log("\n=Best candidate %d ", bestCandidate.Number)
	return p.createKickOrder(bestCandidate.Coords)
}

// calc a distance scale where the player could target
func (p *Player) calcDistanceScale(target Physics.Point) DistanceScale {
	distance := math.Abs(p.Coords.DistanceTo(target))
	// try to be closer the player
	fieldDiagonal := math.Hypot(float64(Units.CourtHeight), float64(Units.CourtWidth))
	toFar := fieldDiagonal / 3
	toNear := fieldDiagonal / 5

	if distance >= toFar {
		return DISTANCE_SCALE_FAR
	} else if distance < toNear {
		return DISTANCE_SCALE_NEAR
	} else {
		return DISTANCE_SCALE_GOOD
	}
}

// Opponent id and angle between it and the target
func (p *Player) watchOpponentOnMyRoute(target Physics.Point, errMarginDegree float64) map[int]float64 {
	opponentTeam := p.findOpponentTeam(p.lastMsg.GameInfo)
	opponents := make(map[int]float64)
	for _, opponent := range opponentTeam.Players {
		angle, isObstacle := p.IsObstacle(target, opponent.Coords, errMarginDegree)
		//App.Log("===")
		//App.Log("===")
		//App.Log("From %d,%d to => %d,%d (obstacle %d,%d): %f degrees %v",
		//	p.Coords.PosX,
		//	p.Coords.PosY,
		//	target.PosX,
		//	target.PosY,
		//	opponent.Coords.PosX,
		//	opponent.Coords.PosY,
		//	angle,
		//	isObstacle,
		//	)
		//
		//App.Log("===")
		//App.Log("===")

		if isObstacle {
			opponents[opponent.Id] = angle
		}
	}
	return opponents
}
//endregion
