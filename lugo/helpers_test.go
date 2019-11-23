package lugo

import "testing"

func TestGetTeam_GetsTheRightTeam(t *testing.T) {
	expected := Team{Name: "MAY TEAM"}
	snapshot := GameSnapshot{HomeTeam: &expected}

	if got := GetTeam(&snapshot, Team_HOME); got != &expected {
		t.Errorf("Unexpected Team - Expected %v, got %v", expected, got)
	}

	expected = Team{Name: "Another Team"}
	snapshot = GameSnapshot{AwayTeam: &expected}

	if got := GetTeam(&snapshot, Team_AWAY); got != &expected {
		t.Errorf("Unexpected Team - Expected %v, got %v", expected, got)
	}
}

func TestGetTeam_DoNotPanicInvalidSnapshot(t *testing.T) {
	snapshot := GameSnapshot{}

	if got := GetTeam(&snapshot, Team_HOME); got != nil {
		t.Errorf("Unexpected value - Expected nil, got %v", got)
	}

	if got := GetTeam(nil, Team_HOME); got != nil {
		t.Errorf("Unexpected value - Expected nil, got %v", got)
	}
}

func TestIsBallHolder_ShouldBeTrue(t *testing.T) {
	expectedHolder := Player{Number: 3, TeamSide: Team_AWAY}
	ball := Ball{Holder: &expectedHolder}
	snapshot := GameSnapshot{Ball: &ball}

	if !IsBallHolder(&snapshot, &expectedHolder) {
		t.Errorf("Unexpected value - Expected true, got false")
	}
}

func TestIsBallHolder_ShouldBeFalse_NoHolder(t *testing.T) {
	expectedHolder := Player{Number: 3, TeamSide: Team_AWAY}
	ball := Ball{}
	snapshot := GameSnapshot{Ball: &ball}

	if IsBallHolder(&snapshot, &expectedHolder) {
		t.Errorf("Unexpected value - Expected false, got true")
	}
}

func TestIsBallHolder_ShouldBeFalse_OtherPlayerHolds(t *testing.T) {
	expectedHolder := Player{Number: 3, TeamSide: Team_AWAY}
	ball := Ball{Holder: &Player{Number: 2, TeamSide: Team_HOME}}
	snapshot := GameSnapshot{Ball: &ball}

	if IsBallHolder(&snapshot, &expectedHolder) {
		t.Errorf("Unexpected value - Expected false, got true")
	}
}

func TestIsBallHolder_ShouldBeFalse_InvalidInputs(t *testing.T) {
	expectedHolder := Player{Number: 3, TeamSide: Team_AWAY}
	snapshot := GameSnapshot{}

	if IsBallHolder(&snapshot, &expectedHolder) {
		t.Errorf("Unexpected value - Expected false, got true")
	}

	if IsBallHolder(nil, &expectedHolder) {
		t.Errorf("Unexpected value - Expected false, got true")
	}
	if IsBallHolder(&GameSnapshot{Ball: &Ball{Holder: &expectedHolder}}, nil) {
		t.Errorf("Unexpected value - Expected false, got true")
	}
}

func TestGetPlayer(t *testing.T) {
	expectedPlayer := Player{TeamSide: Team_HOME, Number: 11}
	snapshot := GameSnapshot{
		HomeTeam: &Team{Players: []*Player{
			&expectedPlayer,
		}},
	}

	if got := GetPlayer(&snapshot, Team_HOME, 11); got != &expectedPlayer {
		t.Errorf("Unexpected value - Expected %v, got %v", expectedPlayer, got)
	}
}

func TestGetPlayer_PlayerNotPresent(t *testing.T) {
	snapshot := GameSnapshot{
		HomeTeam: &Team{Players: []*Player{
			{TeamSide: Team_HOME, Number: 10},
		}},
	}

	if got := GetPlayer(&snapshot, Team_HOME, 11); got != nil {
		t.Errorf("Unexpected value - Expected nil, got %v", got)
	}
}

func TestGetPlayer_TeamNotPresent(t *testing.T) {
	snapshot := GameSnapshot{
		HomeTeam: &Team{Players: []*Player{
			{TeamSide: Team_HOME, Number: 10},
		}},
	}

	if got := GetPlayer(&snapshot, Team_AWAY, 10); got != nil {
		t.Errorf("Unexpected value - Expected nil, got %v", got)
	}
}

func TestGetPlayer_InvalidSnapshot(t *testing.T) {
	if got := GetPlayer(nil, Team_AWAY, 11); got != nil {
		t.Errorf("Unexpected value - Expected nil, got %v", got)
	}
}
