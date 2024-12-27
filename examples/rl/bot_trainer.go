package main

import (
	"context"
	"github.com/lugobots/lugo4go/v3"
	"github.com/lugobots/lugo4go/v3/field"
	"github.com/lugobots/lugo4go/v3/proto"
	"github.com/lugobots/lugo4go/v3/rl"
	"github.com/lugobots/lugo4go/v3/specs"
	"math/rand"
)

const botNumber = 5

type BotTrainer struct {
	remote proto.RemoteClient
	//myTrainingSide proto.Team_Side
}

func NewBotTrainer(remote proto.RemoteClient) *BotTrainer {
	return &BotTrainer{
		remote: remote,
	}
}

func (b *BotTrainer) CreateNewInitialState(data interface{}) (proto.GameSnapshot, error) {

	limitX := specs.FieldWidth * 0.9
	limitY := specs.FieldHeight * 0.9
	for i := uint32(1); i < 11; i++ {
		_, err := b.remote.SetPlayerProperties(context.Background(), &proto.PlayerProperties{
			Side:   proto.Team_HOME,
			Number: i,
			Position: &proto.Point{
				X: rand.Int31n(int32(limitX)),
				Y: rand.Int31n(int32(limitY)),
			},
			Velocity: &proto.Velocity{},
		})
		if err != nil {
			return proto.GameSnapshot{}, err
		}
		_, err = b.remote.SetPlayerProperties(context.Background(), &proto.PlayerProperties{
			Side:   proto.Team_AWAY,
			Number: i,
			Position: &proto.Point{
				X: rand.Int31n(int32(limitX)),
				Y: rand.Int31n(int32(limitY)),
			},
			Velocity: &proto.Velocity{},
		})
		if err != nil {
			return proto.GameSnapshot{}, err
		}

	}
	snapshot, err := b.remote.GetGameSnapshot(context.Background(), &proto.GameSnapshotRequest{})
	if err != nil {
		return proto.GameSnapshot{}, err
	}
	b.remote.SetGameProperties(context.Background(), &proto.GameProperties{
		Turn: 1,
	})

	return *snapshot.GameSnapshot, nil
}

func (b *BotTrainer) GetTrainingState(snapshot proto.GameSnapshot) interface{} {
	return []int{}
}

func (b *BotTrainer) Play(gameSnapshot proto.GameSnapshot, action interface{}) (proto.PlayersOrders, error) {

	playersOrdersBuilder := rl.NewPlayersOrdersBuilder(rl.BotBehaviourStatues)
	inspector, err := lugo4go.NewGameSnapshotInspector(proto.Team_HOME, 5, &gameSnapshot)
	if err != nil {
		return proto.PlayersOrders{}, err
	}

	order := inspector.MakeOrderMoveByDirection(field.Forward, specs.BallMaxSpeed)
	playersOrdersBuilder.AddOrder(5, proto.Team_HOME, []*proto.Order{
		{Action: order},
	})

	return playersOrdersBuilder.Build(), nil
}

func (b *BotTrainer) Evaluate(previousGameSnapshot, newGameSnapshot proto.GameSnapshot, turnOutcome proto.TurnOutcome) (float64, bool, error) {

	return rand.Float64(), newGameSnapshot.Turn > 5, nil
}
