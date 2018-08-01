package main

import (
	"github.com/makeitplay/client-player-go"
	"github.com/makeitplay/commons/GameState"
	"github.com/makeitplay/commons/Units"
)

var player *client.Player

func main() {
	// First we have to get the command line arguments to identify this bot in its game
	serverConfig := new(client.Configuration)
	serverConfig.LoadCmdArg()

	// then we create a client that will handle the communication for us
	player = new(client.Player)
	player.TeamPlace = serverConfig.TeamPlace
	player.Number = serverConfig.PlayerNumber

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
			orderToKick := player.CreateKickOrder(player.OpponentGoal().Center, Units.BallMaxSpeed)
			player.SendOrders("Shoot!", orderToKick)
			return
		}
		// otherwise, let's run towards the ball like kids
		orderToMove := player.CreateMoveOrderMaxSpeed(player.LastServerMessage().GameInfo.Ball.Coords)
		orderToCatch := player.CreateCatchOrder()
		player.SendOrders("Catch the ball!", orderToMove, orderToCatch)
	}
}
