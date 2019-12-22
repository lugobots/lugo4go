#!/bin/bash

mockgen -package=lugo4go \
        -source=./interfaces.go \
        -destination=./mocks.go \
        -self_package=github.com/lugobots/lugo4go/v2

mockgen -package=proto \
        -destination=proto/mocks.go \
        -self_package=github.com/lugobots/lugo4go/v2/proto \
        github.com/lugobots/lugo4go/v2/proto PlayerOrder,GameServer,GameClient,Game_JoinATeamClient,Game_JoinATeamServer

mockgen -package=coach \
        -destination=coach/mocks.go \
        -self_package=github.com/lugobots/lugo4go/v2/coach \
        github.com/lugobots/lugo4go/v2/coach Decider