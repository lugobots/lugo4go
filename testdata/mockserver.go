package testdata

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/makeitplay/client-player-go/lugo"
	"google.golang.org/grpc"
	"log"
	"net"
)

func NewMockServer(ctx context.Context, ctr *gomock.Controller, port int16) (*MockFootballServer, error) {
	mock := NewMockFootballServer(ctr)
	gRPCServer := grpc.NewServer()
	lugo.RegisterFootballServer(gRPCServer, mock)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	go func() {
		<-ctx.Done()
		gRPCServer.Stop()
	}()
	go func() {
		if err := gRPCServer.Serve(lis); err != nil {
			log.Fatalf("test server has stopped: %s", err)
		}
	}()
	return mock, nil
}
