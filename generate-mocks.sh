#!/bin/bash

mockgen -package=testdata -destination=testdata/mock_grpc.go github.com/makeitplay/client-player-go Client
mockgen -package=testdata -destination=testdata/mock_lugo.go github.com/makeitplay/client-player-go/lugo FootballServer,Football_JoinATeamClient

