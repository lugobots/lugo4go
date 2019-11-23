package lugo

import "testing"

func TestGetTeam_GetsTheRightTeam(t *testing.T) {
	expected := Team{Name: "MAY TEAM"}
	snapshot := GameSnapshot{HomeTeam: &expected}

	got := GetTeam(&snapshot, Team_HOME)
	if got != &expected {
		t.Errorf("Unexpedted Team - Expected %v, got %v", expected, got)
	}
	expected = Team{Name: "Another Team"}
	snapshot = GameSnapshot{AwayTeam: &expected}

	got = GetTeam(&snapshot, Team_AWAY)
	if got != &expected {
		t.Errorf("Unexpedted Team - Expected %v, got %v", expected, got)
	}
}

func TestGetTeam_DoNotPanicInvalidSnapshot(t *testing.T) {
	snapshot := GameSnapshot{}

	got := GetTeam(&snapshot, Team_HOME)
	if got != nil {
		t.Errorf("Unexpedted value - Expected nil, got %v", got)
	}

	got = GetTeam(nil, Team_HOME)
	if got != nil {
		t.Errorf("Unexpedted value - Expected nil, got %v", got)
	}
}

func TestIsBallHolder_ShouldBeTrue(t *testing.T) {
	expectedHolder := Player{Number: 3, TeamSide: Team_AWAY}
	ball := Ball{Holder: &expectedHolder}
	snapshot := GameSnapshot{Ball: &ball}

	if !IsBallHolder(&snapshot, &expectedHolder) {
		t.Errorf("Unexpedted value - Expected true, got false")
	}
}

func TestIsBallHolder_ShouldBeFalse_NoHolder(t *testing.T) {
	expectedHolder := Player{Number: 3, TeamSide: Team_AWAY}
	ball := Ball{}
	snapshot := GameSnapshot{Ball: &ball}

	if IsBallHolder(&snapshot, &expectedHolder) {
		t.Errorf("Unexpedted value - Expected false, got true")
	}
}

func TestIsBallHolder_ShouldBeFalse_OtherPlayerHolds(t *testing.T) {
	expectedHolder := Player{Number: 3, TeamSide: Team_AWAY}
	ball := Ball{Holder: &Player{Number: 2, TeamSide: Team_HOME}}
	snapshot := GameSnapshot{Ball: &ball}

	if IsBallHolder(&snapshot, &expectedHolder) {
		t.Errorf("Unexpedted value - Expected false, got true")
	}
}

func TestIsBallHolder_ShouldBeFalse_InvalidInputs(t *testing.T) {
	expectedHolder := Player{Number: 3, TeamSide: Team_AWAY}
	snapshot := GameSnapshot{}

	if IsBallHolder(&snapshot, &expectedHolder) {
		t.Errorf("Unexpedted value - Expected false, got true")
	}

	if IsBallHolder(nil, &expectedHolder) {
		t.Errorf("Unexpedted value - Expected false, got true")
	}
	if IsBallHolder(&GameSnapshot{Ball: &Ball{Holder: &expectedHolder}}, nil) {
		t.Errorf("Unexpedted value - Expected false, got true")
	}
}
