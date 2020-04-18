package coach

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
// For doing that, in our bot, we may create a limited list of region coordinates that will be translated to Regions,
// and than test all of them. The other direction is the conversion from a Point to a region. In this case, we
// assume that the server will only send us valid points (within the field limits).

const (
	// Define the min number of cols allowed on the field division by the positioner
	MinCols uint8 = 4
	// Define the min number of rows allowed on the field division by the positioner
	MinRows uint8 = 2
	// Define the max number of cols allowed on the field division by the positioner
	MaxCols uint8 = 20
	// Define the max number of rows allowed on the field division by the positioner
	MaxRows uint8 = 10
)

// NewPositioner creates a new Positioner that will map the field to provide Regions
func NewPositioner(cols, rows uint8, sideRef lugo.Team_Side) (Positioner, error) {
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

	return &positioner{
		sideRef:      sideRef,
		cols:         cols,
		rows:         rows,
		regionWidth:  field.FieldWidth / float64(cols),
		regionHeight: field.FieldHeight / float64(rows),
	}, nil
}

type positioner struct {
	sideRef      lugo.Team_Side
	cols         uint8
	rows         uint8
	regionWidth  float64
	regionHeight float64
}

func (p *positioner) GetRegion(col, row uint8) (Region, error) {
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
	if p.sideRef == lugo.Team_AWAY {
		center = mirrorCoordsToAway(center)
	}

	return region{
		col:        col,
		row:        row,
		sideRef:    p.sideRef,
		center:     center,
		positioner: p,
	}, nil
}

func (p *positioner) GetPointRegion(point lugo.Point) (Region, error) {
	if p.sideRef == lugo.Team_AWAY {
		point = mirrorCoordsToAway(point)
	}
	cx := float64(point.X) / p.regionWidth
	cy := float64(point.Y) / p.regionHeight
	col := uint8(math.Min(cx, float64(p.cols-1)))
	row := uint8(math.Min(cy, float64(p.rows-1)))
	return p.GetRegion(col, row)
}

type region struct {
	col        uint8
	row        uint8
	sideRef    lugo.Team_Side
	center     lugo.Point
	positioner *positioner
}

func (r region) Col() uint8 {
	return r.col
}

func (r region) Row() uint8 {
	return r.row
}

func (r region) Center() lugo.Point {
	return r.center
}

func (r region) String() string {
	return fmt.Sprintf("{%d,%d-%s}", r.col, r.row, r.sideRef)
}

func (r region) Front() Region {
	if n, err := r.positioner.GetRegion(r.col+1, r.row); err == nil {
		return n
	}
	return r
}

func (r region) Back() Region {
	if n, err := r.positioner.GetRegion(r.col-1, r.row); err == nil {
		return n
	}
	return r
}

func (r region) Left() Region {
	if n, err := r.positioner.GetRegion(r.col, r.row+1); err == nil {
		return n
	}
	return r
}

func (r region) Right() Region {
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
