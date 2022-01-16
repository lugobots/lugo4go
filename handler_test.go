package lugo4go_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/lugobots/lugo4go/v2"
	"github.com/lugobots/lugo4go/v2/lugo"
	"github.com/lugobots/lugo4go/v2/pkg/field"
	"github.com/lugobots/lugo4go/v2/pkg/util"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCoachDefineMyState_AllStates(t *testing.T) {
	var state lugo4go.PlayerState
	var err error
	home3 := &lugo.Player{Number: 3, TeamSide: lugo.Team_HOME}
	home5 := &lugo.Player{Number: 5, TeamSide: lugo.Team_HOME}
	away5 := &lugo.Player{Number: 5, TeamSide: lugo.Team_AWAY}
	ball := &lugo.Ball{}

	snapshot := &lugo.GameSnapshot{
		Ball:     ball,
		HomeTeam: &lugo.Team{Players: []*lugo.Player{home3, home5}},
		AwayTeam: &lugo.Team{Players: []*lugo.Player{away5}},
	}

	// everyone is disputing the ball
	state, err = lugo4go.DefineMyState(snapshot, 3, lugo.Team_HOME)
	assert.Nil(t, err)
	assert.Equal(t, lugo4go.DisputingTheBall, state)

	state, err = lugo4go.DefineMyState(snapshot, 5, lugo.Team_HOME)
	assert.Nil(t, err)
	assert.Equal(t, lugo4go.DisputingTheBall, state)

	state, err = lugo4go.DefineMyState(snapshot, 5, lugo.Team_AWAY)
	assert.Nil(t, err)
	assert.Equal(t, lugo4go.DisputingTheBall, state)

	ball.Holder = home3

	// Holding
	state, err = lugo4go.DefineMyState(snapshot, 3, lugo.Team_HOME)
	assert.Nil(t, err)
	assert.Equal(t, lugo4go.HoldingTheBall, state)

	// supporting
	state, err = lugo4go.DefineMyState(snapshot, 5, lugo.Team_HOME)
	assert.Nil(t, err)
	assert.Equal(t, lugo4go.Supporting, state)

	//
	state, err = lugo4go.DefineMyState(snapshot, 5, lugo.Team_AWAY)
	assert.Nil(t, err)
	assert.Equal(t, lugo4go.Defending, state)
}

func TestCoachDefineMyState_ErrorInvalidSnapshot(t *testing.T) {
	var err error

	_, err = lugo4go.DefineMyState(nil, 3, lugo.Team_HOME)
	assert.Equal(t, err, lugo4go.ErrNoBall)

	_, err = lugo4go.DefineMyState(&lugo.GameSnapshot{}, 3, lugo.Team_HOME)
	assert.Equal(t, err, lugo4go.ErrNoBall)
}

func TestCoachDefineMyState_ErrorNoPlayer(t *testing.T) {
	var err error

	_, err = lugo4go.DefineMyState(&lugo.GameSnapshot{Ball: &lugo.Ball{}}, 3, lugo.Team_HOME)
	assert.Equal(t, err, lugo4go.ErrPlayerNotFound)
}

func TestHandler_Handle_ShouldCallRightMethod(t *testing.T) {
	// initiates Mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // checks all expected things for mocks

	mockLog := NewMockLogger(ctrl)
	mockBot := NewMockBot(ctrl)
	mockBotGoalkeeper := NewMockBot(ctrl)
	mockSender := NewMockOrderSender(ctrl)

	config := util.Config{Number: 4, TeamSide: lugo.Team_HOME}

	handler := lugo4go.NewHandler(mockBot, mockSender, mockLog, config.Number, config.TeamSide)
	goalkeeperHandler := lugo4go.NewHandler(mockBotGoalkeeper, mockSender, mockLog, field.GoalkeeperNumber, config.TeamSide)

	ctx, stop := context.WithTimeout(context.Background(), 1*time.Second)
	defer stop()

	me := &lugo.Player{Number: config.Number, TeamSide: config.TeamSide}
	goalKeeper := &lugo.Player{Number: field.GoalkeeperNumber, TeamSide: config.TeamSide}
	myFriend := &lugo.Player{Number: 5, TeamSide: config.TeamSide}
	myOpponent := &lugo.Player{Number: 5, TeamSide: lugo.Team_AWAY}

	ball := &lugo.Ball{}
	snapshot := &lugo.GameSnapshot{
		Ball:     ball,
		HomeTeam: &lugo.Team{Players: []*lugo.Player{me, goalKeeper, myFriend}},
		AwayTeam: &lugo.Team{Players: []*lugo.Player{myOpponent}},
	}

	// disputing
	mockBot.EXPECT().OnDisputing(ctx, gomock.Any(), snapshot)
	handler.Handle(ctx, snapshot)

	mockBotGoalkeeper.EXPECT().AsGoalkeeper(ctx, gomock.Any(), snapshot, lugo4go.DisputingTheBall)
	goalkeeperHandler.Handle(ctx, snapshot)

	// supporting
	ball.Holder = myFriend
	mockBot.EXPECT().OnSupporting(ctx, gomock.Any(), snapshot)
	handler.Handle(ctx, snapshot)

	mockBotGoalkeeper.EXPECT().AsGoalkeeper(ctx, gomock.Any(), snapshot, lugo4go.Supporting)
	goalkeeperHandler.Handle(ctx, snapshot)

	// Defending
	ball.Holder = myOpponent
	mockBot.EXPECT().OnDefending(ctx, gomock.Any(), snapshot)
	handler.Handle(ctx, snapshot)

	mockBotGoalkeeper.EXPECT().AsGoalkeeper(ctx, gomock.Any(), snapshot, lugo4go.Defending)
	goalkeeperHandler.Handle(ctx, snapshot)

	// holding
	ball.Holder = me
	mockBot.EXPECT().OnHolding(ctx, gomock.Any(), snapshot)
	handler.Handle(ctx, snapshot)

	mockBotGoalkeeper.EXPECT().AsGoalkeeper(ctx, gomock.Any(), snapshot, lugo4go.Supporting)
	goalkeeperHandler.Handle(ctx, snapshot)

	// supporting (goalkeeper holding the ball
	ball.Holder = goalKeeper
	mockBot.EXPECT().OnSupporting(ctx, gomock.Any(), snapshot)
	handler.Handle(ctx, snapshot)

	mockBotGoalkeeper.EXPECT().AsGoalkeeper(ctx, gomock.Any(), snapshot, lugo4go.HoldingTheBall)
	goalkeeperHandler.Handle(ctx, snapshot)
}

func TestHandler_Handle_ShouldLogErrors(t *testing.T) {
	// initiates Mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // checks all expected things for mocks

	mockLog := NewMockLogger(ctrl)
	mockBot := NewMockBot(ctrl)
	mockSender := NewMockOrderSender(ctrl)
	//mockSender := NewMockOrderSender(ctrl)

	config := util.Config{Number: 4, TeamSide: lugo.Team_HOME}
	handler := lugo4go.NewHandler(mockBot, mockSender, mockLog, config.Number, config.TeamSide)

	//ball := &lugo.Ball{}
	//snapshot := &lugo.GameSnapshot{
	//	Ball:     ball,
	//}
	ctx, stop := context.WithTimeout(context.Background(), 1*time.Second)
	defer stop()

	t.Run("no snapshot ", func(t *testing.T) {
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any()).Do(func(s string, args ...interface{}) {
			e, ok := args[0].(error)
			assert.True(t, ok)
			assert.Equal(t, lugo4go.ErrNilSnapshot, e)
		})
		handler.Handle(ctx, nil)
	})
	snapshot := &lugo.GameSnapshot{
		//Ball:     ball,
	}

	t.Run("no ball ", func(t *testing.T) {
		//ball := &lugo.Ball{}
		//
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any()).Do(func(s string, args ...interface{}) {
			e, ok := args[1].(error)
			assert.True(t, ok)
			assert.Equal(t, lugo4go.ErrNoBall, e)
		})
		handler.Handle(ctx, snapshot)
	})

	snapshot.Ball = &lugo.Ball{}
	snapshot.HomeTeam = &lugo.Team{Players: []*lugo.Player{{
		TeamSide: config.TeamSide,
		Number:   config.Number},
	}}
	expectedError := errors.New("some-error")
	mockBot.EXPECT().OnDisputing(gomock.Any(), gomock.Any(), snapshot).Return(expectedError)
	t.Run("bot method fails ", func(t *testing.T) {
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any()).Do(func(s string, args ...interface{}) {
			e, ok := args[1].(error)
			assert.True(t, ok)
			assert.Equal(t, expectedError, e)
		})
		handler.Handle(ctx, snapshot)
	})
}