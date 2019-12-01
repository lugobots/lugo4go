package testdata

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/lugobots/lugo4go/v2/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func NewMockServer(ctx context.Context, ctr *gomock.Controller, port int16) (*MockGameServer, error) {
	mock := NewMockGameServer(ctr)
	gRPCServer := grpc.NewServer()
	proto.RegisterGameServer(gRPCServer, mock)

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
