package coach

import (
	"github.com/lugobots/lugo4go/v2"
	"github.com/lugobots/lugo4go/v2/proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefineMyState_AllStates(t *testing.T) {
	var state PlayerState
	var err error
	home3 := &proto.Player{Number: 3, TeamSide: proto.Team_HOME}
	home5 := &proto.Player{Number: 5, TeamSide: proto.Team_HOME}
	away5 := &proto.Player{Number: 5, TeamSide: proto.Team_AWAY}
	ball := &proto.Ball{}

	snapshot := &proto.GameSnapshot{
		Ball:     ball,
		HomeTeam: &proto.Team{Players: []*proto.Player{home3, home5}},
		AwayTeam: &proto.Team{Players: []*proto.Player{away5}},
	}

	// everyone is disputing the ball
	state, err = DefineMyState(lugo4go.Config{Number: 3, TeamSide: proto.Team_HOME}, snapshot)
	assert.Nil(t, err)
	assert.Equal(t, DisputingTheBall, state)

	state, err = DefineMyState(lugo4go.Config{Number: 5, TeamSide: proto.Team_HOME}, snapshot)
	assert.Nil(t, err)
	assert.Equal(t, DisputingTheBall, state)

	state, err = DefineMyState(lugo4go.Config{Number: 5, TeamSide: proto.Team_AWAY}, snapshot)
	assert.Nil(t, err)
	assert.Equal(t, DisputingTheBall, state)

	ball.Holder = home3

	// Holding
	state, err = DefineMyState(lugo4go.Config{Number: 3, TeamSide: proto.Team_HOME}, snapshot)
	assert.Nil(t, err)
	assert.Equal(t, HoldingTheBall, state)

	// supporting
	state, err = DefineMyState(lugo4go.Config{Number: 5, TeamSide: proto.Team_HOME}, snapshot)
	assert.Nil(t, err)
	assert.Equal(t, Supporting, state)

	//
	state, err = DefineMyState(lugo4go.Config{Number: 5, TeamSide: proto.Team_AWAY}, snapshot)
	assert.Nil(t, err)
	assert.Equal(t, Defending, state)

}

func TestDefineMyState_ErrorInvalidSnapshot(t *testing.T) {
	var err error

	_, err = DefineMyState(lugo4go.Config{Number: 3, TeamSide: proto.Team_HOME}, nil)
	assert.Equal(t, err, ErrNoBall)

	_, err = DefineMyState(lugo4go.Config{Number: 3, TeamSide: proto.Team_HOME}, &proto.GameSnapshot{})
	assert.Equal(t, err, ErrNoBall)
}

func TestDefineMyState_ErrorNoPlayer(t *testing.T) {
	var err error

	_, err = DefineMyState(lugo4go.Config{Number: 3, TeamSide: proto.Team_HOME}, &proto.GameSnapshot{Ball: &proto.Ball{}})
	assert.Equal(t, err, ErrPlayerNotFound)
}
