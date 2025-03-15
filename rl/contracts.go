package rl

import (
	"github.com/lugobots/lugo4go/v3/proto"
)

// TrainingController interface defines the methods for controlling the training process
type TrainingController interface {
	// SetEnvironment resets the game to an initial state, passing any necessary data.
	SetEnvironment(data interface{}) error

	// GetState retrieves the inputs used by your model, e.g., tensors for neural networks.
	GetState() interface{}

	// Update passes the action picked by your model and returns the reward and done values.
	Update(action interface{}) (reward float64, done bool, err error)
}

// BotTrainer interface defines the methods for controlling the game as a bot.
type BotTrainer interface {
	// CreateNewInitialState sets up the initial state for each game.
	CreateNewInitialState(data interface{}) (proto.GameSnapshot, error)

	// GetTrainingState returns the input values (e.g., sensor data) based on the current game state.
	GetTrainingState(snapshot proto.GameSnapshot) interface{}

	// Play translates the action chosen by the model to game orders.
	Play(gameSnapshot proto.GameSnapshot, action interface{}) (proto.PlayersOrders, error)

	// Evaluate compares the previous and new game states to determine the reward and whether the game is done.
	Evaluate(previousGameSnapshot, newGameSnapshot proto.GameSnapshot, turnOutcome proto.TurnOutcome) (float64, bool, error)
}

type PlayersOrdersBuilder interface {
	AddOrder(playerNumber int, teamSide proto.Team_Side, orders []*proto.Order) PlayersOrdersBuilder
	SetPlayerBehaviour(playerNumber int, teamSide proto.Team_Side, behaviour string) PlayersOrdersBuilder
	Build() proto.PlayersOrders
}

type TrainingFunction func(TrainingController) error

type Config struct {
	// Full url to the gRPC server
	GRPCAddress string `json:"grpc_address"`
}

const (
	BotBehaviourStatues  = "statues"
	BotBehaviourKids     = "kids"
	BotBehaviourDefenses = "defenses"
	//BotBehaviourCampers  = "campers"
)
