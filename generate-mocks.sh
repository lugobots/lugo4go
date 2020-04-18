#!/bin/bash

mockgen -package=lugo4go_test \
        -source=./interfaces.go \
        -destination=./mocks_test.go \
        -self_package=github.com/lugobots/lugo4go/v2

mockgen -package=lugo4go_test \
        -destination=./mocks_lugo_test.go \
        github.com/lugobots/lugo4go/v2/lugo PlayerOrder,GameServer,GameClient,Game_JoinATeamClient,\
Game_JoinATeamServer,BroadcastClient,Broadcast_OnEventClient,BroadcastServer,Broadcast_OnEventServer

mockgen -package=coach_test \
        -source=coach/interfaces.go \
        -destination=./coach/mocks_test.go

mockgen -package=coach_test \
        -destination=coach/mocks_lugo_test.go \
        github.com/lugobots/lugo4go/v2/lugo Logger,GameClient