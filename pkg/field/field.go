package field

import "github.com/lugobots/lugo4go/v2/lugo"

// Goal is a set of value about a goal from a team
type Goal struct {
	// Center the is coordinate of the center of the goal
	Center lugo.Point
	// Place identifies the team of this goal (the team that should defend this goal)
	Place lugo.Team_Side
	// TopPole is the coordinates of the pole with a higher Y coordinate
	TopPole lugo.Point
	// BottomPole is the coordinates of the pole  with a lower Y coordinate
	BottomPole lugo.Point
}

// HomeTeamGoal works as a constant value to help to retrieve a Goal struct with the values of the Home team goal
func HomeTeamGoal() Goal {
	return Goal{
		Place:      lugo.Team_HOME,
		Center:     lugo.Point{X: 0, Y: FieldHeight / 2},
		TopPole:    lugo.Point{X: 0, Y: GoalMaxY},
		BottomPole: lugo.Point{X: 0, Y: GoalMinY},
	}
}

// AwayTeamGoal works as a constant value to help to retrieve a Goal struct with the values of the Away team goal
func AwayTeamGoal() Goal {
	return Goal{
		Place:      lugo.Team_AWAY,
		Center:     lugo.Point{X: FieldWidth, Y: FieldHeight / 2},
		TopPole:    lugo.Point{X: FieldWidth, Y: GoalMaxY},
		BottomPole: lugo.Point{X: FieldWidth, Y: GoalMinY},
	}
}

// Returns the goal struct to the team side passed as argument
func GetTeamsGoal(side lugo.Team_Side) Goal {
	if side == lugo.Team_HOME {
		return HomeTeamGoal()
	}
	return AwayTeamGoal()
}

// FieldCenter works as a constant value to help to retrieve a Point struct with the values of the center of the court
func FieldCenter() lugo.Point {
	return lugo.Point{X: FieldWidth / 2, Y: FieldHeight / 2}
}
