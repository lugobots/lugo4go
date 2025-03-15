package rl

import "github.com/lugobots/lugo4go/v3/proto"

type PlayersOrdersSet struct {
	PlayersOrders proto.PlayersOrders
}

func NewPlayersOrdersBuilder(defaultBotBehaviour string) PlayersOrdersBuilder {
	return &PlayersOrdersSet{
		PlayersOrders: proto.PlayersOrders{
			DefaultBehaviour: defaultBotBehaviour,
		},
	}
}

func (p *PlayersOrdersSet) AddOrder(playerNumber int, teamSide proto.Team_Side, orders []*proto.Order) PlayersOrdersBuilder {
	p.PlayersOrders.PlayersOrders = append(p.PlayersOrders.PlayersOrders, &proto.PlayerOrdersOnRLSession{
		TeamSide: teamSide,
		Number:   uint32(playerNumber),
		//Behaviour: "",
		Orders: orders,
	})
	return p
}

func (p *PlayersOrdersSet) SetPlayerBehaviour(playerNumber int, teamSide proto.Team_Side, behaviour string) PlayersOrdersBuilder {
	p.PlayersOrders.PlayersOrders = append(p.PlayersOrders.PlayersOrders, &proto.PlayerOrdersOnRLSession{
		TeamSide:  teamSide,
		Number:    uint32(playerNumber),
		Behaviour: behaviour,
	})
	return p
}

func (p *PlayersOrdersSet) Build() proto.PlayersOrders {
	return p.PlayersOrders
}
