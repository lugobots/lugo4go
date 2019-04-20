package main

import (
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/orders"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/arena/units"
	"github.com/makeitplay/client-player-go"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	"strconv"
)

var gamer *client.Gamer

func main() {

	// Loads the configuration
	serverConfig := new(client.Configuration)
	serverConfig.ParseFromFlags()
	serverConfig.LogLevel = logrus.DebugLevel

	// just creating a position based on the player number
	pos, _ := strconv.Atoi(string(serverConfig.PlayerNumber))
	initialPosition := physics.Point{
		PosX: units.FieldWidth / 4,
		PosY: pos * units.PlayerSize * 2, //(units.FieldHeight / 4) - (pos * units.PlayerSize),
	}

	if serverConfig.TeamPlace == arena.AwayTeam {
		initialPosition.PosX = units.FieldWidth - initialPosition.PosX
	}

	// initialising the player
	gamer = &client.Gamer{}
	gamer.OnAnnouncement = reactToNewState
	gameCtx, err := gamer.Play(initialPosition, serverConfig)
	if err != nil {
		log.Fatal(err)
	}

	// keep the process alive
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	select {
	case <-signalChan:
		logrus.Print("*********** INTERRUPTION SIGNAL ****************")
		gamer.StopToPlay(true)
	case <-gameCtx.Done():
		logrus.Print("*********** Game stopped ****************")
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

		var orderList []orders.Order
		ctx.Logger().Infof("I am player %s", player.Number)
		// otherwise, let's run towards the ball like kids
		if player.Number == arena.PlayerNumber("10") {
			orderToMove, _ := player.CreateMoveOrderMaxSpeed(ctx.GameMsg().GameInfo.Ball.Coords)
			orderList = append(orderList, orderToMove)
		}
		orderToCatch := player.CreateCatchOrder()
		orderList = append(orderList, orderToCatch)
		gamer.SendOrders("Catch the ball!", orderList...)
	}
}
