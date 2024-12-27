package main

import (
	"context"
	"github.com/lugobots/lugo4go/v3"
	mapper "github.com/lugobots/lugo4go/v3/field"
	"github.com/lugobots/lugo4go/v3/rl"
	"github.com/pkg/errors"
	"math/rand"
	"strconv"
	"time"
)

const trainingIterations = 100

func main() {

	logger := lugo4go.DefaultLogger(lugo4go.Config{})

	gym, gameRemoteCtrl, err := rl.NewGym(rl.Config{
		GRPCAddress: "localhost:5000",
	}, logger)

	if err != nil {
		logger.Error(err)
		return
	}

	botTrainer := NewBotTrainer(gameRemoteCtrl)

	err = gym.Start(context.Background(), botTrainer, func(trainingController rl.TrainingController) error {
		possibleActions := map[int]mapper.Direction{
			0: mapper.Forward,
			1: mapper.Right,
			2: mapper.Left,
			3: mapper.Backward,
		}

		for i := 0; i < trainingIterations; i++ {
			//time.Sleep(2 * time.Second)
			err := trainingController.SetEnvironment(struct {
				Whatever string
			}{
				Whatever: strconv.Itoa(i),
			})
			if err != nil {
				return errors.Wrap(err, "failed to reset environment")
			}

			for {
				sensors := trainingController.GetState()

				_ = sensors
				action := possibleActions[rand.Int()%4]

				reward, done, err := trainingController.Update(action)
				if err != nil {
					return errors.Wrap(err, "failed to send training action")
				}
				if done {
					break
				}
				//time.Sleep(1 * time.Second)
				logger.With("reward", reward).Info("step reward")
				time.Sleep(200 * time.Millisecond)
			}
		}
		return nil
	})
	if err != nil {
		logger.With("err", err).Info("training has failed")
		return
	}
	logger.Info("done")
}
