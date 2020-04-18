package coach_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/lugobots/lugo4go/v2/coach"
	"github.com/lugobots/lugo4go/v2/field"
	"github.com/lugobots/lugo4go/v2/lugo"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCoachDefineMyState_AllStates(t *testing.T) {
	var state coach.PlayerState
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
	state, err = coach.DefineMyState(snapshot, 3, lugo.Team_HOME)
	assert.Nil(t, err)
	assert.Equal(t, coach.DisputingTheBall, state)

	state, err = coach.DefineMyState(snapshot, 5, lugo.Team_HOME)
	assert.Nil(t, err)
	assert.Equal(t, coach.DisputingTheBall, state)

	state, err = coach.DefineMyState(snapshot, 5, lugo.Team_AWAY)
	assert.Nil(t, err)
	assert.Equal(t, coach.DisputingTheBall, state)

	ball.Holder = home3

	// Holding
	state, err = coach.DefineMyState(snapshot, 3, lugo.Team_HOME)
	assert.Nil(t, err)
	assert.Equal(t, coach.HoldingTheBall, state)

	// supporting
	state, err = coach.DefineMyState(snapshot, 5, lugo.Team_HOME)
	assert.Nil(t, err)
	assert.Equal(t, coach.Supporting, state)

	//
	state, err = coach.DefineMyState(snapshot, 5, lugo.Team_AWAY)
	assert.Nil(t, err)
	assert.Equal(t, coach.Defending, state)
}

func TestCoachDefineMyState_ErrorInvalidSnapshot(t *testing.T) {
	var err error

	_, err = coach.DefineMyState(nil, 3, lugo.Team_HOME)
	assert.Equal(t, err, coach.ErrNoBall)

	_, err = coach.DefineMyState(&lugo.GameSnapshot{}, 3, lugo.Team_HOME)
	assert.Equal(t, err, coach.ErrNoBall)
}

func TestCoachDefineMyState_ErrorNoPlayer(t *testing.T) {
	var err error

	_, err = coach.DefineMyState(&lugo.GameSnapshot{Ball: &lugo.Ball{}}, 3, lugo.Team_HOME)
	assert.Equal(t, err, coach.ErrPlayerNotFound)
}

func TestHandler_Handle_ShouldCallRightMethod(t *testing.T) {
	// initiates Mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // checks all expected things for mocks

	mockLog := NewMockLogger(ctrl)
	mockBot := NewMockBot(ctrl)
	mockBotGoalkeeper := NewMockBot(ctrl)
	mockSender := NewMockOrderSender(ctrl)

	config := lugo.Config{Number: 4, TeamSide: lugo.Team_HOME}

	handler := coach.NewHandler(mockBot, mockSender, mockLog, config.Number, config.TeamSide)
	goalkeeperHandler := coach.NewHandler(mockBotGoalkeeper, mockSender, mockLog, field.GoalkeeperNumber, config.TeamSide)

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
	mockBot.EXPECT().OnDisputing(ctx, mockSender, snapshot)
	handler.Handle(ctx, snapshot)

	mockBotGoalkeeper.EXPECT().AsGoalkeeper(ctx, mockSender, snapshot, coach.DisputingTheBall)
	goalkeeperHandler.Handle(ctx, snapshot)

	// supporting
	ball.Holder = myFriend
	mockBot.EXPECT().OnSupporting(ctx, mockSender, snapshot)
	handler.Handle(ctx, snapshot)

	mockBotGoalkeeper.EXPECT().AsGoalkeeper(ctx, mockSender, snapshot, coach.Supporting)
	goalkeeperHandler.Handle(ctx, snapshot)

	// Defending
	ball.Holder = myOpponent
	mockBot.EXPECT().OnDefending(ctx, mockSender, snapshot)
	handler.Handle(ctx, snapshot)

	mockBotGoalkeeper.EXPECT().AsGoalkeeper(ctx, mockSender, snapshot, coach.Defending)
	goalkeeperHandler.Handle(ctx, snapshot)

	// holding
	ball.Holder = me
	mockBot.EXPECT().OnHolding(ctx, mockSender, snapshot)
	handler.Handle(ctx, snapshot)

	mockBotGoalkeeper.EXPECT().AsGoalkeeper(ctx, mockSender, snapshot, coach.Supporting)
	goalkeeperHandler.Handle(ctx, snapshot)

	// supporting (goalkeeper holding the ball
	ball.Holder = goalKeeper
	mockBot.EXPECT().OnSupporting(ctx, mockSender, snapshot)
	handler.Handle(ctx, snapshot)

	mockBotGoalkeeper.EXPECT().AsGoalkeeper(ctx, mockSender, snapshot, coach.HoldingTheBall)
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

	config := lugo.Config{Number: 4, TeamSide: lugo.Team_HOME}
	handler := coach.NewHandler(mockBot, mockSender, mockLog, config.Number, config.TeamSide)

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
			assert.Equal(t, coach.ErrNilSnapshot, e)
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
			assert.Equal(t, coach.ErrNoBall, e)
		})
		handler.Handle(ctx, snapshot)
	})

	snapshot.Ball = &lugo.Ball{}
	snapshot.HomeTeam = &lugo.Team{Players: []*lugo.Player{{
		TeamSide: config.TeamSide,
		Number:   config.Number},
	}}
	expectedError := errors.New("some-error")
	mockBot.EXPECT().OnDisputing(gomock.Any(), mockSender, snapshot).Return(expectedError)
	t.Run("bot method fails ", func(t *testing.T) {
		mockLog.EXPECT().Errorf(gomock.Any(), gomock.Any()).Do(func(s string, args ...interface{}) {
			e, ok := args[1].(error)
			assert.True(t, ok)
			assert.Equal(t, expectedError, e)
		})
		handler.Handle(ctx, snapshot)
	})
}
