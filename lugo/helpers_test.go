package lugo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTeam_GetsTheRightTeam(t *testing.T) {
	expected := Team{Name: "MAY TEAM"}
	snapshot := GameSnapshot{HomeTeam: &expected}

	assert.Equal(t, &expected, GetTeam(&snapshot, Team_HOME))

	expected = Team{Name: "Another Team"}
	snapshot = GameSnapshot{AwayTeam: &expected}

	assert.Equal(t, &expected, GetTeam(&snapshot, Team_AWAY))
}

func TestGetTeam_DoNotPanicInvalidSnapshot(t *testing.T) {
	assert.Nil(t, GetTeam(&GameSnapshot{}, Team_HOME))
	assert.Nil(t, GetTeam(nil, Team_HOME))
}

func TestIsBallHolder_ShouldBeTrue(t *testing.T) {
	expectedHolder := Player{Number: 3, TeamSide: Team_AWAY}
	ball := Ball{Holder: &expectedHolder}
	snapshot := GameSnapshot{Ball: &ball}

	assert.True(t, IsBallHolder(&snapshot, &expectedHolder))
}

func TestIsBallHolder_ShouldBeFalse_NoHolder(t *testing.T) {
	expectedHolder := Player{Number: 3, TeamSide: Team_AWAY}
	ball := Ball{}
	snapshot := GameSnapshot{Ball: &ball}

	assert.False(t, IsBallHolder(&snapshot, &expectedHolder))
}

func TestIsBallHolder_ShouldBeFalse_OtherPlayerHolds(t *testing.T) {
	expectedHolder := Player{Number: 3, TeamSide: Team_AWAY}
	ball := Ball{Holder: &Player{Number: 2, TeamSide: Team_HOME}}
	snapshot := GameSnapshot{Ball: &ball}

	assert.False(t, IsBallHolder(&snapshot, &expectedHolder))
}

func TestIsBallHolder_ShouldBeFalse_InvalidInputs(t *testing.T) {
	expectedHolder := Player{Number: 3, TeamSide: Team_AWAY}

	assert.False(t, IsBallHolder(&GameSnapshot{}, &expectedHolder))
	assert.False(t, IsBallHolder(nil, &expectedHolder))
	assert.False(t, IsBallHolder(&GameSnapshot{Ball: &Ball{Holder: &expectedHolder}}, nil))
}

func TestGetPlayer(t *testing.T) {
	expectedPlayer := Player{TeamSide: Team_HOME, Number: 11}
	snapshot := GameSnapshot{
		HomeTeam: &Team{Players: []*Player{
			&expectedPlayer,
		}},
	}
	assert.Equal(t, &expectedPlayer, GetPlayer(&snapshot, Team_HOME, 11))
}

func TestGetPlayer_PlayerNotPresent(t *testing.T) {
	snapshot := GameSnapshot{
		HomeTeam: &Team{Players: []*Player{
			{TeamSide: Team_HOME, Number: 10},
		}},
	}

	assert.Nil(t, GetPlayer(&snapshot, Team_HOME, 11))
}

func TestGetPlayer_TeamNotPresent(t *testing.T) {
	snapshot := GameSnapshot{
		HomeTeam: &Team{Players: []*Player{
			{TeamSide: Team_HOME, Number: 10},
		}},
	}

	assert.Nil(t, GetPlayer(&snapshot, Team_AWAY, 10))
}

func TestGetPlayer_InvalidSnapshot(t *testing.T) {
	assert.Nil(t, GetPlayer(nil, Team_AWAY, 10))
}

func TestMakeOrder_Move(t *testing.T) {
	expectedOrderA := Order_Move{Move: &Move{
		Velocity: &Velocity{
			Speed:     100,
			Direction: North.Normalize(),
		},
	}}

	got, err := MakeOrderMove(Point{}, Point{Y: 100}, 100)
	assert.Nil(t, err)
	assert.Equal(t, expectedOrderA, got)

	expectedOrderB := Order_Move{Move: &Move{
		Velocity: &Velocity{
			Speed:     40,
			Direction: South.Normalize(),
		},
	}}
	got, err = MakeOrderMove(Point{Y: 100}, Point{}, 40)
	assert.Nil(t, err)
	assert.Equal(t, expectedOrderB, got)

	expectedOrderC := Order_Move{Move: &Move{
		Velocity: &Velocity{
			Speed:     40,
			Direction: SouthEast.Normalize(),
		},
	}}
	got, err = MakeOrderMove(Point{Y: 100}, Point{X: 100}, 40)
	assert.Nil(t, err)
	assert.Equal(t, expectedOrderC, got)
}

func TestMakeOrder_Move_ShouldNotMakeInvalidMovement(t *testing.T) {
	_, err := MakeOrderMove(Point{X: 40, Y: 50}, Point{X: 40, Y: 50}, 100)

	assert.NotNil(t, err)
}

func TestMakeOrder_Jump(t *testing.T) {
	expectedOrderA := Order_Jump{Jump: &Jump{
		Velocity: &Velocity{
			Speed:     100,
			Direction: North.Normalize(),
		},
	}}

	got, err := MakeOrderJump(Point{}, Point{Y: 100}, 100)
	assert.Nil(t, err)
	assert.Equal(t, expectedOrderA, got)

	expectedOrderB := Order_Jump{Jump: &Jump{
		Velocity: &Velocity{
			Speed:     40,
			Direction: South.Normalize(),
		},
	}}
	got, err = MakeOrderJump(Point{Y: 100}, Point{}, 40)
	assert.Nil(t, err)
	assert.Equal(t, expectedOrderB, got)

	expectedOrderC := Order_Jump{Jump: &Jump{
		Velocity: &Velocity{
			Speed:     40,
			Direction: SouthEast.Normalize(),
		},
	}}
	got, err = MakeOrderJump(Point{Y: 100}, Point{X: 100}, 40)
	assert.Nil(t, err)
	assert.Equal(t, expectedOrderC, got)
}

func TestMakeOrder_Jump_ShouldNotMakeInvalidMovement(t *testing.T) {
	_, err := MakeOrderJump(Point{X: 40, Y: 50}, Point{X: 40, Y: 50}, 100)

	assert.NotNil(t, err)
}
