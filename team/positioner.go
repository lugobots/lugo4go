package team

import (
	"fmt"
	"github.com/lugobots/lugo4go/v2/lugo"
	"github.com/lugobots/lugo4go/v2/pkg/field"
	"math"
)

// Important note: since our bot needs to have the best performance possible. We may ensure that some errors will never
// happen based on our configuration. Once the errors are in a controlled and limited list of methods, we are able
// to ignore the errors during the game, and only test them in our unit tests.
//
// For doing that, in our bot, we may create a limited list of FieldArea coordinates that will be translated to Regions,
// and than test all of them. The other direction is the conversion from a Point to a FieldArea. In this case, we
// assume that the server will only send us valid points (within the field limits).

const (
	// Define the min number of cols allowed on the field division by the Arrangement
	MinCols uint8 = 4
	// Define the min number of rows allowed on the field division by the Arrangement
	MinRows uint8 = 2
	// Define the max number of cols allowed on the field division by the Arrangement
	MaxCols uint8 = 200
	// Define the max number of rows allowed on the field division by the Arrangement
	MaxRows uint8 = 100
)

// NewArrangement creates a new Positioner that will map the field to provide Regions
func NewArrangement(cols, rows uint8, sideRef lugo.Team_Side) (*Arrangement, error) {
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

	return &Arrangement{
		TeamSide:     sideRef,
		cols:         cols,
		rows:         rows,
		regionWidth:  field.FieldWidth / float64(cols),
		regionHeight: field.FieldHeight / float64(rows),
	}, nil
}

type Arrangement struct {
	TeamSide     lugo.Team_Side
	cols         uint8
	rows         uint8
	regionWidth  float64
	regionHeight float64
}

func (p *Arrangement) GetRegion(col, row uint8) (FieldNav, error) {
	if col >= p.cols {
		return nil, ErrMaxCols
	}
	if row >= p.rows {
		return nil, ErrMaxRows
	}

	center := lugo.Point{
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

func (p *Arrangement) GetPointRegion(point lugo.Point) (FieldNav, error) {
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
	center     lugo.Point
	positioner *Arrangement
}

func (r FieldArea) Col() uint8 {
	return r.col
}

func (r FieldArea) Row() uint8 {
	return r.row
}

func (r FieldArea) Center() lugo.Point {
	return r.center
}

func (r FieldArea) String() string {
	return fmt.Sprintf("{%d,%d-%s}", r.col, r.row, r.sideRef)
}

func (r FieldArea) Front() FieldNav {
	if n, err := r.positioner.GetRegion(r.col+1, r.row); err == nil {
		return n
	}
	return r
}

func (r FieldArea) Back() FieldNav {
	if n, err := r.positioner.GetRegion(r.col-1, r.row); err == nil {
		return n
	}
	return r
}

func (r FieldArea) Left() FieldNav {
	if n, err := r.positioner.GetRegion(r.col, r.row+1); err == nil {
		return n
	}
	return r
}

func (r FieldArea) Right() FieldNav {
	if n, err := r.positioner.GetRegion(r.col, r.row-1); err == nil {
		return n
	}
	return r
}

// Invert the coords X and Y as in a mirror to found out the same position seen from the away team field
// Keep in mind that all coords in the field is based on the bottom left corner!
func mirrorCoordsToAway(coords lugo.Point) lugo.Point {
	return lugo.Point{
		X: field.FieldWidth - coords.X,
		Y: field.FieldHeight - coords.Y,
	}
}
