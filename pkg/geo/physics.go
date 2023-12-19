package geo

import "github.com/lugobots/lugo4go/v3/proto"

// AngleWithRoute this function returns the angle to a particular target to a given a direction and an origin point,
// The angle adopts the direction as the base axis, so a positive angle indicates the obstacle is on the left side,
// while a negative angle indicates that the obstacle if on the right side.
//
// This function is specially useful when a player have opponent player at some point between him and the goal.
// The angle between the route to the goal and the opponent may be used to decide to change its route.
func AngleWithRoute(direction proto.Vector, from, obstacle proto.Point) float64 {
	angleToObstacle, err := proto.NewVector(from, obstacle)
	if err != nil {
		return 0
	}
	return direction.AngleWith(angleToObstacle)
}
