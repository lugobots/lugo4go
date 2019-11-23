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
