#!/bin/bash

mockgen -package=testdata -destination=testdata/mock_grpc.go github.com/makeitplay/client-player-go/ops Client,Logger
mockgen -package=testdata -destination=testdata/mock_lugo.go github.com/makeitplay/client-player-go/lugo GameServer,GameClient,Game_JoinATeamClient

