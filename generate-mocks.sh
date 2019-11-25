#!/bin/bash

mockgen -package=testdata -destination=testdata/mock_grpc.go github.com/makeitplay/client-player-go/ops Client,Logger,OrderSender
mockgen -package=testdata -destination=testdata/mock_proto.go github.com/makeitplay/client-player-go/proto GameServer,GameClient,Game_JoinATeamClient

