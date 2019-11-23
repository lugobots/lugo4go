package lugo

const (
	// GoalkeeperNumber defines the goalkeeper number
	GoalkeeperNumber = uint32(1)
)

// Goal is a set of value about a goal from a team
type Goal struct {
	// Center the is coordinate of the center of the goal
	Center Point
	// Place identifies the team of this goal (the team that should defend this goal)
	Place Team_Side
	// TopPole is the coordinates of the pole with a higher Y coordinate
	TopPole Point
	// BottomPole is the coordinates of the pole  with a lower Y coordinate
	BottomPole Point
}

// HomeTeamGoal works as a constant value to help to retrieve a Goal struct with the values of the Home team goal
func HomeTeamGoal() Goal {
	return Goal{
		Place:      Team_HOME,
		Center:     Point{X: 0, Y: FieldHeight / 2},
		TopPole:    Point{X: 0, Y: GoalMaxY},
		BottomPole: Point{X: 0, Y: GoalMinY},
	}
}

// AwayTeamGoal works as a constant value to help to retrieve a Goal struct with the values of the Away team goal
func AwayTeamGoal() Goal {
	return Goal{
		Place:      Team_AWAY,
		Center:     Point{X: FieldWidth, Y: FieldHeight / 2},
		TopPole:    Point{X: FieldWidth, Y: GoalMaxY},
		BottomPole: Point{X: FieldWidth, Y: GoalMinY},
	}
}

// Returns the goal struct to the team side passed as argument
func GetTeamsGoal(side Team_Side) Goal {
	if side == Team_HOME {
		return HomeTeamGoal()
	}
	return AwayTeamGoal()
}

// FieldCenter works as a constant value to help to retrieve a Point struct with the values of the center of the court
func FieldCenter() Point {
	return Point{X: FieldWidth / 2, Y: FieldHeight / 2}
}
