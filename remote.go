package lugo4go

import (
	"context"
	"github.com/lugobots/lugo4go/v2/lugo"
	"github.com/lugobots/lugo4go/v2/pkg/util"
	"google.golang.org/grpc"
	"time"
)

func ConnectRemote(config util.Config) (lugo.RemoteClient, *grpc.ClientConn, error) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	var grpcConn *grpc.ClientConn
	var err error
	// @todo there are some gRPC options that we should take a look tro improve this part.
	if config.Insecure {
		grpcConn, err = grpc.DialContext(ctx, config.GRPCAddress, grpc.WithBlock(), grpc.WithInsecure())
	} else {
		grpcConn, err = grpc.DialContext(ctx, config.GRPCAddress, grpc.WithBlock())
	}
	if err != nil {
		return nil, nil, err
	}
	grpcClient := lugo.NewRemoteClient(grpcConn)
	return grpcClient, grpcConn, nil
}

//
//type Remote struct {
//	GRPCClient lugo.RemoteClient
//	grpcConn   *grpc.ClientConn
//	Logger     util.Logger
//	config     util.Config
//}
//
//func (r *Remote) Stop() error {
//	return r.grpcConn.Close()
//}
//
//func (r *Remote) TagRPC(ctx context.Context, _ *stats.RPCTagInfo) context.Context {
//	return ctx
//}
//
//func (r *Remote) HandleRPC(context.Context, stats.RPCStats) {
//
//}
//
//func (r *Remote) TagConn(ctx context.Context, _ *stats.ConnTagInfo) context.Context {
//	return ctx
//}
//
//func (r *Remote) HandleConn(_ context.Context, sts stats.ConnStats) {
//	switch sts.(type) {
//	case *stats.ConnEnd:
//		_ = r.Stop()
//		break
//	}
//}
