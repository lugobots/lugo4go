package main

import (
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/arena/units"
	"github.com/makeitplay/client-player-go"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
)

var gamer *client.Gamer

func main() {

	serverConfig := new(client.Configuration)
	serverConfig.ParseFromFlags()

	initialPosition := physics.Point{
		PosX: units.FieldWidth / 4,
		PosY: units.FieldHeight / 2,
	}
	if serverConfig.TeamPlace == arena.AwayTeam {
		initialPosition.PosX *= 2
	}
	gamer = &client.Gamer{}
	gamer.OnAnnouncement = reactToNewState
	if err := gamer.Play(initialPosition, serverConfig); err != nil {
		log.Fatal(err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	select {
	case <-signalChan:
		logrus.Print("*********** INTERRUPTION SIGNAL ****************")
		gamer.StopToPlay(true)
	}

}

func reactToNewState(ctx client.TurnContext) {
	// there is a chance of receiving a msg before the user be add to the game state, so it can be nill at the very beginning
	if ctx.Player() == nil {
		return
	}

	ctx.Logger().Info("I got a message")

	player := ctx.Player()
	// for this example, or smart player only reacts when the game server is listening for orders
	if ctx.GameMsg().State == arena.Listening {

		// we are going to kick the ball as soon as we catch it
		if player.IHoldTheBall(ctx.GameMsg().GameInfo.Ball) {
			orderToKick, _ := player.CreateKickOrder(ctx.GameMsg().GameInfo.Ball, player.OpponentGoal().Center, units.BallMaxSpeed)
			gamer.SendOrders("Shoot!", orderToKick)
			return
		}
		// otherwise, let's run towards the ball like kids
		orderToMove, _ := player.CreateMoveOrderMaxSpeed(ctx.GameMsg().GameInfo.Ball.Coords)
		orderToCatch := player.CreateCatchOrder()
		gamer.SendOrders("Catch the ball!", orderToMove, orderToCatch)
	}
}