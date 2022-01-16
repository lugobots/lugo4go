package field

import (
	"fmt"
	"github.com/lugobots/lugo4go/v2/proto"
)

// Mapper Helps the bots to see the fields from their team perspective instead of using the cartesian plan provided
// by the game server. Instead of base your logic on the axes X and Y, the Map create a FieldArea map based
// on the team side.
// The FieldArea coordinates uses the defensive field's right corner as its origin.
// This mechanism if specially useful to define players regions based on their roles, since you do not have to mirror
// the coordinate, neither do extra logic to define regions on the field where the player should be.
type Mapper interface {
	// GetRegion Returns a FieldArea based on the coordinates and on the current field division
	GetRegion(col, row uint8) (Region, error)
	// GetPointRegion returns the FieldArea where that point is in
	GetPointRegion(point *proto.Point) (Region, error)
}

// Region represent a quadrant on the field. It is not always squared form because you may define how many cols/rows
// the field will be divided in. So, based on that division (e.g. 4 rows, 6 cols) there will be a fixed number of regions
// and their coordinates will be zero-index (e.g. from 0 to 3 rows when divided in 4 rows).
type Region interface {
	fmt.Stringer
	// Col The col coordinate based on the field division
	Col() uint8
	// Row The row coordinate based on the field division
	Row() uint8
	// Center Return the point at the center of the quadrant represented by this Region. It is not always precise.
	Center() *proto.Point

	// Front is the FieldArea immediately in front of this one from the player perspective
	// Important: The same FieldArea is returned if the requested FieldArea is not valid
	Front() Region
	// Back is the FieldArea immediately behind this one from the player perspective
	// Important: The same FieldArea is returned if the requested FieldArea is not valid
	Back() Region
	// Left is the FieldArea immediately on left of this one from the player perspective
	// Important: The same FieldArea is returned if the requested FieldArea is not valid
	Left() Region
	// Right is the FieldArea immediately on right of this one from the player perspective
	// Important: The same FieldArea is returned if the requested FieldArea is not valid
	Right() Region

	// Eq does not check if the passed region is on a map of same size!
	Eq(region Region) bool
}
