package main

import (
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/arena/units"
	"github.com/makeitplay/client-player-go"
)

var player *client.Player

func main() {

	serverConfig := new(client.Configuration)
	serverConfig.ParseFromFlags()

	initialPosition := physics.Point{
		PosX: units.FieldWidth / 4,
		PosY: units.FieldHeight / 2,
	}

	player := &client.Player{}
	player.OnAnnouncement = reactToNewState
	player.Play(initialPosition, serverConfig)

}

func reactToNewState(msg client.GameMessage) {
	// as soo we get the new game state, we have to update or position in the field
	player.UpdatePosition(msg.GameInfo)

	// for this example, or smart player only reacts when the game server is listening for orders
	if msg.State == arena.Listening {

		// we are going to kick the ball as soon as we catch it
		if player.IHoldTheBall() {
			orderToKick := player.CreateKickOrder(player.OpponentGoal().Center, units.BallMaxSpeed)
			player.SendOrders("Shoot!", orderToKick)
			return
		}
		// otherwise, let's run towards the ball like kids
		orderToMove := player.CreateMoveOrderMaxSpeed(player.LastServerMessage().GameInfo.Ball.Coords)
		orderToCatch := player.CreateCatchOrder()
		player.SendOrders("Catch the ball!", orderToMove, orderToCatch)
	}
}
