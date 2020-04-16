#!/bin/bash

mockgen -package=lugo4go \
        -source=./interfaces.go \
        -destination=./mocks.go \
        -self_package=github.com/lugobots/lugo4go/v2

mockgen -package=lugo \
        -destination=lugo/mocks.go \
        -self_package=github.com/lugobots/lugo4go/v2/lugo \
        github.com/lugobots/lugo4go/v2/lugo PlayerOrder,GameServer,GameClient,Game_JoinATeamClient,\
Game_JoinATeamServer,BroadcastClient,Broadcast_OnEventClient,BroadcastServer,Broadcast_OnEventServer

mockgen -package=coach \
        -source=coach/interfaces.go \
        -destination=./coach/mocks.go \
        -self_package=github.com/lugobots/lugo4go/v2/coach