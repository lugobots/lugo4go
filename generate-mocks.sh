#!/bin/bash

mockgen -package=lugo4go_test \
        -source=./interfaces.go \
        -destination=./mocks_test.go \
        -self_package=github.com/lugobots/lugo4go/v2

mockgen -package=lugo4go_test \
        -destination=./mocks_lugo_test.go \
        github.com/lugobots/lugo4go/v2/lugo PlayerOrder,GameServer,GameClient,Game_JoinATeamClient,\
Game_JoinATeamServer,BroadcastClient,Broadcast_OnEventClient,BroadcastServer,Broadcast_OnEventServer

mockgen -package=team_test \
        -source=team/interfaces.go \
        -destination=./team/mocks_test.go

mockgen -package=team_test \
        -destination=team/mocks_lugo_test.go \
        github.com/lugobots/lugo4go/v2/lugo GameClient

mockgen -package=team_test \
        -destination=team/mocks_log_test.go \
        github.com/lugobots/lugo4go/v2/pkg/util Logger