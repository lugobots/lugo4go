package field

import "time"

const (
	// BaseUnit is used to increase the integer units scale and improve the precision when the integer numbers
	// come from float  calculations. Some units have to be integer to avoid infinite intervals (e.g. a point in the field)
	BaseUnit = 100

	// PlayerSize is the size of each player
	PlayerSize = 4 * BaseUnit

	// PlayerReconnectionWaitTime is a penalty time imposed to the player that needs to reconnect during the match.
	// this interval ensure players won't drop connection in purpose to be reallocated to their initial position.
	PlayerReconnectionWaitTime = 20 * time.Second

	// max number of players in a team by mach
	MaxPlayers = 11

	// min number of players in a team by mach, if a team gets to have less to this number, the team loses by W.O.
	MinPlayers = 6

	// PlayerMaxSpeed is the max speed that a play may move  by frame
	PlayerMaxSpeed = 100.0

	// FieldWidth is the width of the field (horizontal view)
	FieldWidth = 200 * BaseUnit

	// FieldHeight is the height of the field (horizontal view)
	FieldHeight = 100 * BaseUnit

	// FieldNeutralCenter is the radius of the neutral circle on the center of the field
	FieldNeutralCenter = 100

	// BallSize size of the element ball
	BallSize = 2 * BaseUnit

	// BallDeceleration is the deceleration rate of the ball speed  by frame
	BallDeceleration = 10.0

	// BallMaxSpeed is the max speed of the ball by frame
	BallMaxSpeed = 4.0 * BaseUnit

	// BallMinSpeed is the minimal speed of the ball  by frame. When the ball was at this speed or slower, it will be considered stopped.
	BallMinSpeed = 2

	// BallTimeInGoalZone is the max number of turns that the ball may be in a goal zone. After that, the ball will be auto kicked
	// towards the center of the field.
	BallTimeInGoalZone = 40 // 40 / 20 fps = 2 seconds

	// GoalWidth is the goal width
	GoalWidth = 30 * BaseUnit

	// GoalMinY is the coordinate Y of the lower pole of the goals
	GoalMinY = (FieldHeight - GoalWidth) / 2

	// GoalMaxY is the coordinate Y of the upper pole of the goals
	GoalMaxY = GoalMinY + GoalWidth

	// GoalZoneRange is the minimal distance that a player can stay from the opponent goal
	GoalZoneRange = 14 * BaseUnit

	// GoalKeeperJumpDuration is the number of turns that the jump takes. A jump cannot be interrupted after has been requested
	GoalKeeperJumpDuration = 3

	// GoalKeeperJumpSpeed is the max speed of the goalkeeper during the jump
	GoalKeeperJumpSpeed = 2 * PlayerMaxSpeed

	// GoalkeeperNumber defines the goalkeeper number
	GoalkeeperNumber = uint32(1)
)
