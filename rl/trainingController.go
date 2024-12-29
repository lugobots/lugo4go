package rl

import (
	"context"
	"fmt"
	"github.com/lugobots/lugo4go/v3/proto"
	"github.com/pkg/errors"
	"time"
)

type TrainingCrl struct {
	remote         proto.RemoteClient
	assistant      proto.RLAssistantClient
	latestSnapShot proto.GameSnapshot
	//snapshotLocker sync.RWMutex
	trainner        BotTrainer
	trainingContext context.Context
}

func NewTrainingCrl(ctx context.Context, trainner BotTrainer, remote proto.RemoteClient, assistant proto.RLAssistantClient) *TrainingCrl {
	return &TrainingCrl{
		remote:          remote,
		assistant:       assistant,
		trainingContext: ctx,
		trainner:        trainner,
	}
}

func (t *TrainingCrl) SetEnvironment(data interface{}) error {
	//t.snapshotLocker.Lock()
	//defer t.snapshotLocker.Unlock()
	var err error
	t.latestSnapShot, err = t.trainner.CreateNewInitialState(data)
	if err != nil {
		return errors.Wrap(err, "failed to set environment")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("**** Environment reset ******")
	if _, err := t.assistant.ResetEnv(ctx, &proto.RLResetConfig{}); err != nil {
		return errors.Wrap(err, "failed to reset RL assistant")
	}
	return nil
}

func (t *TrainingCrl) GetState() interface{} {
	//t.snapshotLocker.
	//defer t.snapshotLocker.Unlock()
	return t.trainner.GetTrainingState(t.latestSnapShot)

}

func (t *TrainingCrl) Update(action interface{}) (reward float64, done bool, err error) {

	playersOrders, err := t.trainner.Play(t.latestSnapShot, action)
	if err != nil {
		return 0, false, errors.Wrap(err, "trainer bot failed to play")
	}
	turnOutcome, err := t.assistant.SendPlayersOrders(t.trainingContext, &playersOrders)

	if err != nil {
		return 0, false, errors.Wrap(err, "RL assistant failed to send the orders")
	}
	previousGameSnapshot := t.latestSnapShot
	t.latestSnapShot = *turnOutcome.GameSnapshot
	return t.trainner.Evaluate(previousGameSnapshot, *turnOutcome.GameSnapshot, *turnOutcome)
}
