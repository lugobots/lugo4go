package main

import (
	"github.com/makeitplay/arena/GameState"
	"github.com/makeitplay/arena/Physics"
	"github.com/makeitplay/arena/Units"
	"github.com/makeitplay/client-player-go"
	"math/rand"
	"time"
)

var player *client.Player

func main() {
	rand.Seed(time.Now().UnixNano())
	// First we have to get the command line arguments to identify this bot in its game
	serverConfig := new(client.Configuration)
	serverConfig.ParseFromFlags()

	// then we create a client that will handle the communication for us
	player = new(client.Player)
	player.TeamPlace = serverConfig.TeamPlace
	player.Number = serverConfig.PlayerNumber
	// this will be our bot initial position
	player.Coords = Physics.Point{
		PosX: rand.Int() % units.CourtWidth,
		PosY: rand.Int() % units.CourtHeight,
	}

	// we have to set the call back function that will process the player behaviour when the game state has been changed
	player.OnAnnouncement = reactToNewState
	player.Start(serverConfig)
}

func reactToNewState(msg client.GameMessage) {
	// as soo we get the new game state, we have to update or position in the field
	player.UpdatePosition(msg.GameInfo)

	// for this example, or smart player only reacts when the game server is listening for orders
	if GameState.State(msg.State) == GameState.Listening {

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
