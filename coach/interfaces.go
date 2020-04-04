package coach

import (
	"fmt"
	"github.com/lugobots/lugo4go/v2/proto"
)

// Positioner Helps the bots to see the fields from their team perspective instead of using the cartesian plan provided
// by the game server. Instead of base your logic on the axes X and Y, the positioner create a region map based
// on the team side.
// The region coordinates uses the defensive field's right corner as its origin.
// This mechanism if specially useful to define players regions based on their roles, since you do not have to mirror
// the coordinate, neither do extra logic to define regions on the field where the player should be.
type Positioner interface {
	// GetRegion Returns a region based on the coordinates and on the current field division
	GetRegion(col, row uint8) (Region, error)
	// GetPointRegion returns the region where that point is in
	GetPointRegion(point lugo.Point) (Region, error)
}

// Region represent a quadrant on the field. It is not always squared form because you may define how many cols/rows
// the field will be divided in. So, based on that division (e.g. 4 rows, 6 cols) there will be a fixed number of regions
// and their coordinates will be zero-index (e.g. from 0 to 3 rows when divided in 4 rows).
type Region interface {
	fmt.Stringer
	// The col coordinate based on the field division
	Col() uint8
	// The row coordinate based on the field division
	Row() uint8
	// Return the point at the center of the quadrant represented by this Region. It is not always precise.
	Center() lugo.Point

	// The region immediatelly in front of this one from the player perspective
	// Important: The same region is returned if the requested region is not valid
	Front() Region
	// The region immediatelly behind this one from the player perspective
	// Important: The same region is returned if the requested region is not valid
	Back() Region
	// The region immediatelly on left of this one from the player perspective
	// Important: The same region is returned if the requested region is not valid
	Left() Region
	// The region immediatelly on right of this one from the player perspective
	// Important: The same region is returned if the requested region is not valid
	Right() Region
}
