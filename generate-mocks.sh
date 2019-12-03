#!/bin/bash

mockgen -package=lugo4go \
        -source=./interfaces.go \
        -destination=./mocks.go \
        -self_package=github.com/lugobots/lugo4go/v2

mockgen -package=proto \
        -destination=proto/mocks.go \
        -self_package=github.com/lugobots/lugo4go/v2/proto \
        github.com/lugobots/lugo4go/v2/proto GameServer,GameClient,Game_JoinATeamClient