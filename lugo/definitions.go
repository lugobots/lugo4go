package lugo

import "time"

// BaseUnit is used to increase the integer units scale and improve the precision when the integer numbers
// come from float  calculations. Some units have to be integer to avoid infinite intervals (e.g. a point in the field)
const BaseUnit = 100

// PlayerSize is the size of each player
const PlayerSize = 4 * BaseUnit

// PlayerReconnectionWaitTime is a penalty time imposed to the player that needs to reconnect during the match.
// this interval ensure players won't drop connection in purpose to be reallocated to their initial position.
const PlayerReconnectionWaitTime = 20 * time.Second

// PlayerMaxSpeed is the max speed that a play may move  by frame
const PlayerMaxSpeed = 100.0

// FieldWidth is the width of the field (horizontal view)
const FieldWidth = 200 * BaseUnit

// FieldHeight is the height of the field (horizontal view)
const FieldHeight = 100 * BaseUnit

// FieldNeutralCenter is the radius of the neutral circle on the center of the field
const FieldNeutralCenter = 100

// BallSize size of the element ball
const BallSize = 2 * BaseUnit

// BallDeceleration is the deceleration rate of the ball speed  by frame
const BallDeceleration = 10.0

// BallMaxSpeed is the max speed of the ball by frame
const BallMaxSpeed = 4.0 * BaseUnit

// BallMinSpeed is the minimal speed of the ball  by frame. When the ball was at this speed or slower, it will be considered stopped.
const BallMinSpeed = 2

// BallTimeInGoalZone is the max number of turns that the ball may be in a goal zone. After that, the ball will be auto kicked
// towards the center of the field.
const BallTimeInGoalZone = 60 // 60 / 20 fps = 3 seconds

// GoalWidth is the goal width
const GoalWidth = 30 * BaseUnit

// GoalMinY is the coordinate Y of the lower pole of the goals
const GoalMinY = (FieldHeight - GoalWidth) / 2

// GoalMaxY is the coordinate Y of the upper pole of the goals
const GoalMaxY = GoalMinY + GoalWidth

// GoalZoneRange is the minimal distance that a player can stay from the opponent goal
const GoalZoneRange = 14 * BaseUnit

// GoalKeeperJumpDuration is the number of turns that the jump takes. A jump cannot be interrupted after has been requested
const GoalKeeperJumpDuration = 3

// GoalKeeperJumpSpeed is the max speed of the goalkeeper during the jump
const GoalKeeperJumpSpeed = 2 * PlayerMaxSpeed
