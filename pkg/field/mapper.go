package field

import (
	"fmt"
	"github.com/lugobots/lugo4go/v2/lugo"
	"math"
)

// Important note: since our bot needs to have the best performance possible. We may ensure that some errors will never
// happen based on our configuration. Once the errors are in a controlled and limited list of methods, we are able
// to ignore the errors during the game, and only test them in our unit tests.
//
// For doing that, in our bot, we may create a limited list of FieldArea coordinates that will be translated to Regions,
// and then test all of them. The other direction is the conversion from a Point to a FieldArea. In this case, we
// assume that the server will only send us valid points (within the field limits).

const (
	// MinCols Define the min number of cols allowed on the field division by the Map
	MinCols uint8 = 4
	// MinRows Define the min number of rows allowed on the field division by the Map
	MinRows uint8 = 2
	// MaxCols Define the max number of cols allowed on the field division by the Map
	MaxCols uint8 = 200
	// MaxRows Define the max number of rows allowed on the field division by the Map
	MaxRows uint8 = 100
)

// NewMapper creates a new Mapper that will map the field to provide Regions
func NewMapper(cols, rows uint8, sideRef lugo.Team_Side) (*Map, error) {
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
		regionWidth:  FieldWidth / float64(cols),
		regionHeight: FieldHeight / float64(rows),
	}, nil
}

type Map struct {
	TeamSide     lugo.Team_Side
	cols         uint8
	rows         uint8
	regionWidth  float64
	regionHeight float64
}

func (p *Map) GetRegion(col, row uint8) (Region, error) {
	if col >= p.cols {
		return nil, ErrMaxCols
	}
	if row >= p.rows {
		return nil, ErrMaxRows
	}

	center := &lugo.Point{
		X: int32(math.Round(float64(col)*p.regionWidth + p.regionWidth/2)),
		Y: int32(math.Round(float64(row)*p.regionHeight + p.regionHeight/2)),
	}
	if p.TeamSide == lugo.Team_AWAY {
		center = mirrorCoordsToAway(center)
	}

	return FieldArea{
		col:        col,
		row:        row,
		sideRef:    p.TeamSide,
		center:     center,
		positioner: p,
	}, nil
}

func (p *Map) GetPointRegion(point *lugo.Point) (Region, error) {
	if p.TeamSide == lugo.Team_AWAY {
		point = mirrorCoordsToAway(point)
	}
	cx := float64(point.X) / p.regionWidth
	cy := float64(point.Y) / p.regionHeight
	col := uint8(math.Min(cx, float64(p.cols-1)))
	row := uint8(math.Min(cy, float64(p.rows-1)))
	return p.GetRegion(col, row)
}

type FieldArea struct {
	col        uint8
	row        uint8
	sideRef    lugo.Team_Side
	center     *lugo.Point
	positioner *Map
}

func (r FieldArea) Eq(region Region) bool {
	return region.Col() == r.col && region.Row() == r.Row()
}

func (r FieldArea) Col() uint8 {
	return r.col
}

func (r FieldArea) Row() uint8 {
	return r.row
}

func (r FieldArea) Center() *lugo.Point {
	return r.center.Copy()
}

func (r FieldArea) String() string {
	return fmt.Sprintf("{%d,%d-%s}", r.col, r.row, r.sideRef)
}

func (r FieldArea) Front() Region {
	if n, err := r.positioner.GetRegion(r.col+1, r.row); err == nil {
		return n
	}
	return r
}

func (r FieldArea) Back() Region {
	if n, err := r.positioner.GetRegion(r.col-1, r.row); err == nil {
		return n
	}
	return r
}

func (r FieldArea) Left() Region {
	if n, err := r.positioner.GetRegion(r.col, r.row+1); err == nil {
		return n
	}
	return r
}

func (r FieldArea) Right() Region {
	if n, err := r.positioner.GetRegion(r.col, r.row-1); err == nil {
		return n
	}
	return r
}

// Invert the coords X and Y as in a mirror to found out the same position seen from the away team field
// Keep in mind that all coords in the field are based in the bottom left corner!
func mirrorCoordsToAway(coords *lugo.Point) *lugo.Point {
	return &lugo.Point{
		X: FieldWidth - coords.X,
		Y: FieldHeight - coords.Y,
	}
}
