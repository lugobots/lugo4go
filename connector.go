package lugo4go

import (
	"context"
	"encoding/json"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/lugobots/lugo4go/v2/lugo"
	"log"
	"math/rand"
	"time"
)

type GymEnv struct {
}

func (g GymEnv) GetSetup(ctx context.Context, request *lugo.SetupRequest) (*lugo.SetupResponse, error) {
	min := [][]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	max := [][]int{
		{10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10},
	}

	log.Println("got a call")

	minmums, _ := json.Marshal(min)
	maximums, _ := json.Marshal(max)

	return &lugo.SetupResponse{
		RewardThreshold: 0,
		ActionSpace:     3,
		ObservationSpaceMinimums: &any.Any{
			Value: minmums,
		},
		ObservationSpaceMaximums: &any.Any{
			Value: maximums,
		},
		MinReward: 0,
		MaxReward: 10,
	}, nil
}

func (g GymEnv) Step(ctx context.Context, request *lugo.StepRequest) (*lugo.StepResponse, error) {
	log.Println("got a SUPER call")
	rand.Seed(time.Now().UTC().UnixNano())

	ob := [][]int{
		{0, 2, 1, 0, 4, 0, 2, 4, 0, 0, 5, 6},
	}
	observation, _ := json.Marshal(ob)

	resp := &lugo.StepResponse{
		Reward: rand.Float32(),
		Ob: &any.Any{
			Value: observation,
		},
	}
	return resp, nil
}

func (g GymEnv) Reset(ctx context.Context, request *lugo.ResetRequest) (*lugo.ResetResponse, error) {
	panic("implement me")
}

func (g GymEnv) Render(ctx context.Context, request *lugo.RenderRequest) (*lugo.RenderResponse, error) {
	panic("implement me")
}

func (g GymEnv) Close(ctx context.Context, request *lugo.CloseRequest) (*lugo.CloseResponse, error) {
	panic("implement me")
}
