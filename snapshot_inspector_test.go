package lugo4go

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lugobots/lugo4go/v3/field"
	"github.com/lugobots/lugo4go/v3/proto"
	"github.com/lugobots/lugo4go/v3/specs"
)

func TestGetTeam_GetsTheRightTeam(t *testing.T) {
	expected := proto.Team{Name: "MAY TEAM"}

	insp := inspector{
		mySide:   proto.Team_HOME,
		snapshot: &proto.GameSnapshot{HomeTeam: &expected},
	}
	assert.Equal(t, &expected, insp.GetTeam(proto.Team_HOME))

	expected = proto.Team{Name: "Another Team"}
	insp.snapshot = &proto.GameSnapshot{AwayTeam: &expected}

	assert.Equal(t, &expected, insp.GetTeam(proto.Team_AWAY))
}

func TestIsBallHolder_ShouldBeTrue(t *testing.T) {
	expectedHolder := proto.Player{Number: 3, TeamSide: proto.Team_AWAY}
	ball := proto.Ball{Holder: &expectedHolder}
	snapshot := proto.GameSnapshot{Ball: &ball}

	insp := inspector{
		mySide:   proto.Team_HOME,
		snapshot: &snapshot,
	}
	assert.True(t, insp.IsBallHolder(&expectedHolder))
}

func TestIsBallHolder_ShouldBeFalse_NoHolder(t *testing.T) {
	expectedHolder := proto.Player{Number: 3, TeamSide: proto.Team_AWAY}
	ball := proto.Ball{}
	snapshot := proto.GameSnapshot{Ball: &ball}

	insp := inspector{
		mySide:   proto.Team_HOME,
		snapshot: &snapshot,
	}
	assert.False(t, insp.IsBallHolder(&expectedHolder))
}

func TestIsBallHolder_ShouldBeFalse_OtherPlayerHolds(t *testing.T) {
	expectedHolder := proto.Player{Number: 3, TeamSide: proto.Team_AWAY}
	ball := proto.Ball{Holder: &proto.Player{Number: 2, TeamSide: proto.Team_HOME}}
	snapshot := proto.GameSnapshot{Ball: &ball}

	insp := inspector{
		mySide:   proto.Team_HOME,
		snapshot: &snapshot,
	}
	assert.False(t, insp.IsBallHolder(&expectedHolder))
}

func TestIsBallHolder_ShouldBeFalse_InvalidInputs(t *testing.T) {
	expectedHolder := proto.Player{Number: 3, TeamSide: proto.Team_AWAY}
	insp := inspector{
		snapshot: nil,
	}

	assert.False(t, insp.IsBallHolder(&expectedHolder))

	insp = inspector{
		snapshot: nil,
	}

	assert.False(t, insp.IsBallHolder(&expectedHolder))

	insp = inspector{
		snapshot: &proto.GameSnapshot{Ball: &proto.Ball{Holder: &expectedHolder}},
	}
}

func TestGetPlayer(t *testing.T) {
	expectedPlayer := proto.Player{TeamSide: proto.Team_HOME, Number: 11}
	snapshot := proto.GameSnapshot{
		HomeTeam: &proto.Team{Players: []*proto.Player{
			&expectedPlayer,
		}},
	}

	insp := inspector{
		snapshot: &snapshot,
	}
	assert.Equal(t, &expectedPlayer, insp.GetPlayer(proto.Team_HOME, 11))
}

func TestGetPlayer_PlayerNotPresent(t *testing.T) {
	snapshot := proto.GameSnapshot{
		HomeTeam: &proto.Team{Players: []*proto.Player{
			{TeamSide: proto.Team_HOME, Number: 10},
		}},
	}
	insp := inspector{
		snapshot: &snapshot,
	}

	assert.Nil(t, insp.GetPlayer(proto.Team_HOME, 11))
}

func TestGetPlayer_TeamNotPresent(t *testing.T) {
	snapshot := proto.GameSnapshot{
		HomeTeam: &proto.Team{Players: []*proto.Player{
			{TeamSide: proto.Team_HOME, Number: 10},
		}},
	}
	insp := inspector{
		snapshot: &snapshot,
	}

	assert.Nil(t, insp.GetPlayer(proto.Team_AWAY, 10))
}

func TestMakeOrder_Move(t *testing.T) {
	expectedOrderA := &proto.Order_Move{Move: &proto.Move{
		Velocity: &proto.Velocity{
			Speed:     100,
			Direction: proto.North().Copy().Normalize(),
		},
	}}
	insp := inspector{}

	got, err := insp.MakeOrderMoveFromPoint(proto.Point{}, proto.Point{Y: 100}, 100)
	assert.Nil(t, err)
	assert.Equal(t, expectedOrderA, got)

	expectedOrderB := &proto.Order_Move{Move: &proto.Move{
		Velocity: &proto.Velocity{
			Speed:     40,
			Direction: proto.South().Copy().Normalize(),
		},
	}}
	got, err = insp.MakeOrderMoveFromPoint(proto.Point{Y: 100}, proto.Point{}, 40)
	assert.Nil(t, err)
	assert.Equal(t, expectedOrderB, got)

	expectedOrderC := &proto.Order_Move{Move: &proto.Move{
		Velocity: &proto.Velocity{
			Speed:     40,
			Direction: proto.SouthEast().Copy().Normalize(),
		},
	}}
	got, err = insp.MakeOrderMoveFromPoint(proto.Point{Y: 100}, proto.Point{X: 100}, 40)
	assert.Nil(t, err)
	assert.Equal(t, expectedOrderC, got)
}

func TestMakeOrder_Move_ShouldNotMakeInvalidMovement(t *testing.T) {
	insp := inspector{}
	order, err := insp.MakeOrderMoveFromPoint(proto.Point{X: 40, Y: 50}, proto.Point{X: 40, Y: 50}, 100)

	assert.NotNil(t, err)
	assert.Nil(t, order)
}

func TestMakeOrder_Jump(t *testing.T) {
	insp := inspector{
		me: &proto.Player{},
	}

	expectedOrderA := &proto.Order_Jump{Jump: &proto.Jump{
		Velocity: &proto.Velocity{
			Speed:     100,
			Direction: proto.North().Copy().Normalize(),
		},
	}}

	insp.me.Position = &proto.Point{}
	got, err := insp.MakeOrderJump(proto.Point{Y: 100}, 100)
	assert.Nil(t, err)
	assert.Equal(t, expectedOrderA, got)

	expectedOrderB := &proto.Order_Jump{Jump: &proto.Jump{
		Velocity: &proto.Velocity{
			Speed:     40,
			Direction: proto.South().Copy().Normalize(),
		},
	}}

	insp.me.Position = &proto.Point{Y: 100}
	got, err = insp.MakeOrderJump(proto.Point{}, 40)
	assert.Nil(t, err)
	assert.Equal(t, expectedOrderB, got)

	expectedOrderC := &proto.Order_Jump{Jump: &proto.Jump{
		Velocity: &proto.Velocity{
			Speed:     40,
			Direction: proto.SouthEast().Copy().Normalize(),
		},
	}}

	insp.me.Position = &proto.Point{Y: 100}
	got, err = insp.MakeOrderJump(proto.Point{X: 100}, 40)
	assert.Nil(t, err)
	assert.Equal(t, expectedOrderC, got)
}

func TestMakeOrder_Jump_ShouldNotMakeInvalidMovement(t *testing.T) {
	insp := inspector{
		me: &proto.Player{
			Position: &proto.Point{X: 40, Y: 50},
		},
	}

	order, err := insp.MakeOrderJump(proto.Point{X: 40, Y: 50}, 100)

	assert.NotNil(t, err)
	assert.Nil(t, order)
}

func TestMakeOrder_Kick_SameDirection(t *testing.T) {

	expectedOrderA := &proto.Order_Kick{Kick: &proto.Kick{
		Velocity: &proto.Velocity{
			Speed:     specs.BallMaxSpeed,
			Direction: proto.North().Copy().Normalize(),
		},
	}}

	origin := field.FieldCenter()
	ball := proto.Ball{Position: &origin, Velocity: proto.NewZeroedVelocity(*proto.North().Copy()).Copy()}
	insp := inspector{
		snapshot: &proto.GameSnapshot{
			Ball: &ball,
		},
	}

	got, err := insp.MakeOrderKick(proto.Point{X: origin.X, Y: origin.Y + 100}, specs.BallMaxSpeed)
	assert.Nil(t, err)
	assert.Equal(t, expectedOrderA, got)
}

func TestMakeOrder_Kick_DiffDirection(t *testing.T) {
	originPoint := proto.Point{X: 1, Y: 1}
	originVelocity := proto.Velocity{
		Speed: 100,
		Direction: &proto.Vector{ // Going east
			X: 1,
			Y: 0,
		},
	}

	targetPoint := proto.Point{X: 2, Y: 2} // this point is on northeast from the original point

	// this is the final direction we desire the ball goes in
	desiredDirection, err := proto.NewVector(originPoint, targetPoint)
	if err != nil {
		t.Fatalf("bad test: %s", err)
	}

	// However, remember that the velocity will be summed! So, we should send the complement
	complementVector, err := desiredDirection.Sub(originVelocity.Direction)
	if err != nil {
		t.Fatalf("bad test: %s", err)
	}

	expectedOrderA := &proto.Order_Kick{Kick: &proto.Kick{
		Velocity: &proto.Velocity{
			Speed: specs.BallMaxSpeed,
			// we expect that the function created a complementary Vector,
			// so we do not have to think about it during the development
			Direction: complementVector.Normalize(),
		},
	}}

	ball := proto.Ball{Position: &originPoint, Velocity: &originVelocity}

	insp := inspector{
		snapshot: &proto.GameSnapshot{
			Ball: &ball,
		},
	}

	got, err := insp.MakeOrderKick(targetPoint, specs.BallMaxSpeed)
	assert.Nil(t, err)
	assert.Equal(t, expectedOrderA, got)
}
