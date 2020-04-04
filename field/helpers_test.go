package field

import (
	"github.com/lugobots/lugo4go/v2/lugo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTeam_GetsTheRightTeam(t *testing.T) {
	expected := lugo.Team{Name: "MAY TEAM"}
	snapshot := lugo.GameSnapshot{HomeTeam: &expected}

	assert.Equal(t, &expected, GetTeam(&snapshot, lugo.Team_HOME))

	expected = lugo.Team{Name: "Another Team"}
	snapshot = lugo.GameSnapshot{AwayTeam: &expected}

	assert.Equal(t, &expected, GetTeam(&snapshot, lugo.Team_AWAY))
}

func TestGetTeam_DoNotPanicInvalidSnapshot(t *testing.T) {
	assert.Nil(t, GetTeam(&lugo.GameSnapshot{}, lugo.Team_HOME))
	assert.Nil(t, GetTeam(nil, lugo.Team_HOME))
}

func TestIsBallHolder_ShouldBeTrue(t *testing.T) {
	expectedHolder := lugo.Player{Number: 3, TeamSide: lugo.Team_AWAY}
	ball := lugo.Ball{Holder: &expectedHolder}
	snapshot := lugo.GameSnapshot{Ball: &ball}

	assert.True(t, IsBallHolder(&snapshot, &expectedHolder))
}

func TestIsBallHolder_ShouldBeFalse_NoHolder(t *testing.T) {
	expectedHolder := lugo.Player{Number: 3, TeamSide: lugo.Team_AWAY}
	ball := lugo.Ball{}
	snapshot := lugo.GameSnapshot{Ball: &ball}

	assert.False(t, IsBallHolder(&snapshot, &expectedHolder))
}

func TestIsBallHolder_ShouldBeFalse_OtherPlayerHolds(t *testing.T) {
	expectedHolder := lugo.Player{Number: 3, TeamSide: lugo.Team_AWAY}
	ball := lugo.Ball{Holder: &lugo.Player{Number: 2, TeamSide: lugo.Team_HOME}}
	snapshot := lugo.GameSnapshot{Ball: &ball}

	assert.False(t, IsBallHolder(&snapshot, &expectedHolder))
}

func TestIsBallHolder_ShouldBeFalse_InvalidInputs(t *testing.T) {
	expectedHolder := lugo.Player{Number: 3, TeamSide: lugo.Team_AWAY}

	assert.False(t, IsBallHolder(&lugo.GameSnapshot{}, &expectedHolder))
	assert.False(t, IsBallHolder(nil, &expectedHolder))
	assert.False(t, IsBallHolder(&lugo.GameSnapshot{Ball: &lugo.Ball{Holder: &expectedHolder}}, nil))
}

func TestGetPlayer(t *testing.T) {
	expectedPlayer := lugo.Player{TeamSide: lugo.Team_HOME, Number: 11}
	snapshot := lugo.GameSnapshot{
		HomeTeam: &lugo.Team{Players: []*lugo.Player{
			&expectedPlayer,
		}},
	}
	assert.Equal(t, &expectedPlayer, GetPlayer(&snapshot, lugo.Team_HOME, 11))
}

func TestGetPlayer_PlayerNotPresent(t *testing.T) {
	snapshot := lugo.GameSnapshot{
		HomeTeam: &lugo.Team{Players: []*lugo.Player{
			{TeamSide: lugo.Team_HOME, Number: 10},
		}},
	}

	assert.Nil(t, GetPlayer(&snapshot, lugo.Team_HOME, 11))
}

func TestGetPlayer_TeamNotPresent(t *testing.T) {
	snapshot := lugo.GameSnapshot{
		HomeTeam: &lugo.Team{Players: []*lugo.Player{
			{TeamSide: lugo.Team_HOME, Number: 10},
		}},
	}

	assert.Nil(t, GetPlayer(&snapshot, lugo.Team_AWAY, 10))
}

func TestGetPlayer_InvalidSnapshot(t *testing.T) {
	assert.Nil(t, GetPlayer(nil, lugo.Team_AWAY, 10))
}

func TestMakeOrder_Move(t *testing.T) {
	expectedOrderA := &lugo.Order_Move{Move: &lugo.Move{
		Velocity: &lugo.Velocity{
			Speed:     100,
			Direction: lugo.North().Copy().Normalize(),
		},
	}}

	got, err := MakeOrderMove(lugo.Point{}, lugo.Point{Y: 100}, 100)
	assert.Nil(t, err)
	assert.Equal(t, expectedOrderA, got)

	expectedOrderB := &lugo.Order_Move{Move: &lugo.Move{
		Velocity: &lugo.Velocity{
			Speed:     40,
			Direction: lugo.South().Copy().Normalize(),
		},
	}}
	got, err = MakeOrderMove(lugo.Point{Y: 100}, lugo.Point{}, 40)
	assert.Nil(t, err)
	assert.Equal(t, expectedOrderB, got)

	expectedOrderC := &lugo.Order_Move{Move: &lugo.Move{
		Velocity: &lugo.Velocity{
			Speed:     40,
			Direction: lugo.SouthEast().Copy().Normalize(),
		},
	}}
	got, err = MakeOrderMove(lugo.Point{Y: 100}, lugo.Point{X: 100}, 40)
	assert.Nil(t, err)
	assert.Equal(t, expectedOrderC, got)
}

func TestMakeOrder_Move_ShouldNotMakeInvalidMovement(t *testing.T) {
	order, err := MakeOrderMove(lugo.Point{X: 40, Y: 50}, lugo.Point{X: 40, Y: 50}, 100)

	assert.NotNil(t, err)
	assert.Nil(t, order)
}

func TestMakeOrder_Jump(t *testing.T) {

	expectedOrderA := &lugo.Order_Jump{Jump: &lugo.Jump{
		Velocity: &lugo.Velocity{
			Speed:     100,
			Direction: lugo.North().Copy().Normalize(),
		},
	}}

	got, err := MakeOrderJump(lugo.Point{}, lugo.Point{Y: 100}, 100)
	assert.Nil(t, err)
	assert.Equal(t, expectedOrderA, got)

	expectedOrderB := &lugo.Order_Jump{Jump: &lugo.Jump{
		Velocity: &lugo.Velocity{
			Speed:     40,
			Direction: lugo.South().Copy().Normalize(),
		},
	}}
	got, err = MakeOrderJump(lugo.Point{Y: 100}, lugo.Point{}, 40)
	assert.Nil(t, err)
	assert.Equal(t, expectedOrderB, got)

	expectedOrderC := &lugo.Order_Jump{Jump: &lugo.Jump{
		Velocity: &lugo.Velocity{
			Speed:     40,
			Direction: lugo.SouthEast().Copy().Normalize(),
		},
	}}
	got, err = MakeOrderJump(lugo.Point{Y: 100}, lugo.Point{X: 100}, 40)
	assert.Nil(t, err)
	assert.Equal(t, expectedOrderC, got)
}

func TestMakeOrder_Jump_ShouldNotMakeInvalidMovement(t *testing.T) {
	order, err := MakeOrderJump(lugo.Point{X: 40, Y: 50}, lugo.Point{X: 40, Y: 50}, 100)

	assert.NotNil(t, err)
	assert.Nil(t, order)
}

func TestMakeOrder_Kick_SameDirection(t *testing.T) {
	expectedOrderA := &lugo.Order_Kick{Kick: &lugo.Kick{
		Velocity: &lugo.Velocity{
			Speed:     BallMaxSpeed,
			Direction: lugo.North().Copy().Normalize(),
		},
	}}

	origin := FieldCenter()
	ball := lugo.Ball{Position: &origin, Velocity: lugo.NewZeroedVelocity(*lugo.North().Copy()).Copy()}

	got, err := MakeOrderKick(ball, lugo.Point{X: origin.X, Y: origin.Y + 100}, BallMaxSpeed)
	assert.Nil(t, err)
	assert.Equal(t, expectedOrderA, got)
}

func TestMakeOrder_Kick_DiffDirection(t *testing.T) {
	originPoint := lugo.Point{X: 1, Y: 1}
	originVelocity := lugo.Velocity{
		Speed: 100,
		Direction: &lugo.Vector{ // Going east
			X: 1,
			Y: 0,
		},
	}

	targetPoint := lugo.Point{X: 2, Y: 2} // this point is on northeast from the original potin

	// this is the final direction we desire the ball goes in
	desiredDirection, err := lugo.NewVector(originPoint, targetPoint)
	if err != nil {
		t.Fatalf("bad test: %s", err)
	}

	// However, remember that the velocity will be summed! So, we should send the complement
	complementVector, err := desiredDirection.Sub(originVelocity.Direction)
	if err != nil {
		t.Fatalf("bad test: %s", err)
	}

	expectedOrderA := &lugo.Order_Kick{Kick: &lugo.Kick{
		Velocity: &lugo.Velocity{
			Speed: BallMaxSpeed,
			// we expect that the function created a complementary Vector,
			// so we do not have to think about it during the development
			Direction: complementVector.Normalize(),
		},
	}}

	ball := lugo.Ball{Position: &originPoint, Velocity: &originVelocity}

	got, err := MakeOrderKick(ball, targetPoint, BallMaxSpeed)
	assert.Nil(t, err)
	assert.Equal(t, expectedOrderA, got)
}
