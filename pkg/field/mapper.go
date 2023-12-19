package field

import (
	"fmt"
	"math"

	"github.com/lugobots/lugo4go/v2/proto"
)

// Important note: since our bot needs to have the best performance possible. We may ensure that some errors will never
// happen based on our configuration. Once the errors are in a controlled and limited list of methods, we are able
// to ignore the errors during the game, and only test them in our unit tests.
//
// For doing that, in our bot, we may create a limited list of MapArea coordinates that will be translated to Regions,
// and then test all of them. The other direction is the conversion from a Point to a MapArea. In this case, we
// assume that the server will only send us valid points (within the field limits).

const (
	// MinCols Define the min number of cols allowed on the field division by the Map
	MinCols int = 4
	// MinRows Define the min number of rows allowed on the field division by the Map
	MinRows int = 2
	// MaxCols Define the max number of cols allowed on the field division by the Map
	MaxCols int = 200
	// MaxRows Define the max number of rows allowed on the field division by the Map
	MaxRows int = 100
)

// NewMapper creates a new Mapper that will map the field to provide Regions
func NewMapper(cols, rows int, sideRef proto.Team_Side) (*Map, error) {
	if cols < MinCols {
		return nil, ErrMinCols
	}
	if cols > MaxCols {
		return nil, ErrMaxCols
	}
	if rows < MinRows {
		return nil, ErrMinRows
	}
	if rows > MaxRows {
		return nil, ErrMaxRows
	}

	return &Map{
		TeamSide:     sideRef,
		cols:         cols,
		rows:         rows,
		regionWidth:  MaxXCoordinate / float64(cols),
		regionHeight: MaxYCoordinate / float64(rows),
	}, nil
}

type Map struct {
	TeamSide     proto.Team_Side
	cols         int
	rows         int
	regionWidth  float64
	regionHeight float64
}

func (p *Map) GetRegion(col, row int) (Region, error) {
	if col >= p.cols {
		return nil, ErrMaxCols
	}
	if row >= p.rows {
		return nil, ErrMaxRows
	}

	if col < 0 {
		col = 0
	}
	if row < 0 {
		row = 0
	}

	center := &proto.Point{
		X: int32(math.Round(float64(col)*p.regionWidth + p.regionWidth/2)),
		Y: int32(math.Round(float64(row)*p.regionHeight + p.regionHeight/2)),
	}
	if p.TeamSide == proto.Team_AWAY {
		center = mirrorCoordsToAway(center)
	}

	return MapArea{
		col:        col,
		row:        row,
		sideRef:    p.TeamSide,
		center:     center,
		positioner: p,
	}, nil
}

func (p *Map) GetPointRegion(point *proto.Point) (Region, error) {
	if p.TeamSide == proto.Team_AWAY {
		point = mirrorCoordsToAway(point)
	}
	cx := float64(point.X) / p.regionWidth
	cy := float64(point.Y) / p.regionHeight
	col := int(math.Min(cx, float64(p.cols-1)))
	row := int(math.Min(cy, float64(p.rows-1)))
	return p.GetRegion(col, row)
}

type MapArea struct {
	col        int
	row        int
	sideRef    proto.Team_Side
	center     *proto.Point
	positioner *Map
}

func (r MapArea) Eq(region Region) bool {
	return region.Col() == r.Col() && region.Row() == r.Row()
}

func (r MapArea) Col() int {
	return r.col
}

func (r MapArea) Row() int {
	return r.row
}

func (r MapArea) Center() *proto.Point {
	return r.center.Copy()
}

func (r MapArea) String() string {
	return fmt.Sprintf("{%d,%d-%s}", r.col, r.row, r.sideRef)
}

func (r MapArea) Front() Region {
	if n, err := r.positioner.GetRegion(r.col+1, r.row); err == nil {
		return n
	}
	return r
}

func (r MapArea) Back() Region {
	if n, err := r.positioner.GetRegion(r.col-1, r.row); err == nil {
		return n
	}
	return r
}

func (r MapArea) Left() Region {
	if n, err := r.positioner.GetRegion(r.col, r.row+1); err == nil {
		return n
	}
	return r
}

func (r MapArea) Right() Region {
	if n, err := r.positioner.GetRegion(r.col, r.row-1); err == nil {
		return n
	}
	return r
}

// Invert the coords X and Y as in a mirror to found out the same position seen from the away team field
// Keep in mind that all coords in the field are based in the bottom left corner!
func mirrorCoordsToAway(coords *proto.Point) *proto.Point {
	return &proto.Point{
		X: MaxXCoordinate - coords.X,
		Y: MaxYCoordinate - coords.Y,
	}
}
