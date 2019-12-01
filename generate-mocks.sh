#!/bin/bash

mockgen -package=testdata -destination=testdata/mock_lugo.go github.com/lugobots/lugo4go/v2/lugo Client,Logger,OrderSender
mockgen -package=testdata -destination=testdata/mock_proto.go github.com/lugobots/lugo4go/v2/proto GameServer,GameClient,Game_JoinATeamClient

