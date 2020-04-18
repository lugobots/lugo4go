#!/bin/bash

#mockgen -package=lugo4go \
#        -source=./interfaces.go \
#        -destination=./mocks.go \
#        -self_package=github.com/lugobots/lugo4go/v2
#
#mockgen -package=lugo \
#        -destination=lugo/mocks.go \
#        -self_package=github.com/lugobots/lugo4go/v2/lugo \
#        github.com/lugobots/lugo4go/v2/lugo PlayerOrder,GameServer,GameClient,Game_JoinATeamClient,\
#Game_JoinATeamServer,BroadcastClient,Broadcast_OnEventClient,BroadcastServer,Broadcast_OnEventServer,Logger
#
mockgen -package=coach_test \
        -source=coach/interfaces.go \
        -destination=./coach/mocks_test.go

cd internal/mocks
mockgen -package=mocks \
        -source=../../coach/interfaces.go \
        -destination=mocks.go