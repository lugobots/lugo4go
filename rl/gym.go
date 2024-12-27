package rl

import (
	"context"
	"github.com/lugobots/lugo4go/v3"
	"github.com/lugobots/lugo4go/v3/proto"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io"
	"time"
)

type Gym struct {
	grpcConn  *grpc.ClientConn
	Logger    lugo4go.Logger
	assistant proto.RLAssistantClient
	remote    proto.RemoteClient
}

func NewGym(config Config, logger *zap.SugaredLogger) (*Gym, proto.RemoteClient, error) {

	connHandler := grpc.WithInsecure()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Debug("trying to connect to the server")
	grpcConn, err := grpc.DialContext(ctx, config.GRPCAddress, grpc.WithBlock(), grpc.WithInsecure(), connHandler)
	if err != nil {
		return nil, nil, errors.Wrap(err, "did not connect to the game server")
	}
	assistant := proto.NewRLAssistantClient(grpcConn)
	remote := proto.NewRemoteClient(grpcConn)
	return &Gym{
		grpcConn:  grpcConn,
		Logger:    logger,
		assistant: assistant,
		remote:    remote,
	}, remote, nil
}

func (g *Gym) Start(ctx context.Context, trainner BotTrainer, trainingFunction TrainingFunction) error {
	defer g.grpcConn.Close()
	startingCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// TODO bad practice, this constructors should not be here!
	trainingCrl := NewTrainingCrl(trainner, g.remote, g.assistant)

	g.Logger.Debug("starting training session")
	session, err := g.assistant.StartSession(startingCtx, &proto.RLSessionConfig{
		// nothing for now
	})

	//tem que passar o  g.assistant para o ctrl pra ele poder resetar o wrapped bots

	if err != nil {
		return errors.Wrap(err, "failed to start a training session")
	}
	// Everything may happen too fast no the client (gym) side, so we need to know if the server has actually
	// done all its work before the next steps.
	// However, the way gRPC works in golang, as soon as the stream connection is open, the `startSession` method is returned.
	// To work around that, letś wait the first message from te server to ensure the StartSession method on their side is ready.
	// TODO: confirm if the error may be ignored. I believe that as long as the session context is still ok, the
	// receiver method may fail
	_, err = session.Recv()

	if session.Context().Err() != nil {
		return errors.Wrap(session.Context().Err(), "could not start a training session")
	}

	go func() {
		for {
			// keeping the connection alive
			_, err := session.Recv()
			if err != nil && errors.Cause(err) != io.EOF {
				// At this point, we should be able to stop the trainer function to gracefully end the session.
				//However, this is not currently possible with the existing project structure. That said, this is
				// not an issue, as the stop function would fail regardless.
				g.Logger.Errorf("the RL assistant session ended with error: %v", err)
				return
			}
			g.Logger.Infof("meleca")
		}
	}()

	if err != nil {
		return errors.Wrap(err, "failed to start a training session")
	}
	g.Logger.Infof("session started")

	return trainingFunction(trainingCrl)
}
