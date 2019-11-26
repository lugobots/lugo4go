package lugo

import "github.com/lugobots/client-player-go/v2/proto"

// Goal is a set of value about a goal from a team
type Goal struct {
	// Center the is coordinate of the center of the goal
	Center proto.Point
	// Place identifies the team of this goal (the team that should defend this goal)
	Place proto.Team_Side
	// TopPole is the coordinates of the pole with a higher Y coordinate
	TopPole proto.Point
	// BottomPole is the coordinates of the pole  with a lower Y coordinate
	BottomPole proto.Point
}

// HomeTeamGoal works as a constant value to help to retrieve a Goal struct with the values of the Home team goal
func HomeTeamGoal() Goal {
	return Goal{
		Place:      proto.Team_HOME,
		Center:     proto.Point{X: 0, Y: FieldHeight / 2},
		TopPole:    proto.Point{X: 0, Y: GoalMaxY},
		BottomPole: proto.Point{X: 0, Y: GoalMinY},
	}
}

// AwayTeamGoal works as a constant value to help to retrieve a Goal struct with the values of the Away team goal
func AwayTeamGoal() Goal {
	return Goal{
		Place:      proto.Team_AWAY,
		Center:     proto.Point{X: FieldWidth, Y: FieldHeight / 2},
		TopPole:    proto.Point{X: FieldWidth, Y: GoalMaxY},
		BottomPole: proto.Point{X: FieldWidth, Y: GoalMinY},
	}
}

// Returns the goal struct to the team side passed as argument
func GetTeamsGoal(side proto.Team_Side) Goal {
	if side == proto.Team_HOME {
		return HomeTeamGoal()
	}
	return AwayTeamGoal()
}

// FieldCenter works as a constant value to help to retrieve a Point struct with the values of the center of the court
func FieldCenter() proto.Point {
	return proto.Point{X: FieldWidth / 2, Y: FieldHeight / 2}
}
