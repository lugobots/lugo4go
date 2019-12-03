package coach

import (
	"github.com/lugobots/lugo4go/v2/field"
	"github.com/lugobots/lugo4go/v2/proto"
	"math"
)

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
func NewPositioner(cols, rows uint8, sideRef proto.Team_Side) (Positioner, error) {
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
	sideRef      proto.Team_Side
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

	center := proto.Point{
		X: int32(math.Round(float64(col)*p.regionWidth + p.regionWidth/2)),
		Y: int32(math.Round(float64(row)*p.regionHeight + p.regionHeight/2)),
	}
	if p.sideRef == proto.Team_AWAY {
		center = mirrorCoordsToAway(center)
	}

	return region{
		col:     col,
		row:     row,
		sideRef: p.sideRef,
		center:  center,
	}, nil
}

func (p *positioner) GetPointRegion(point proto.Point) (Region, error) {
	if p.sideRef == proto.Team_AWAY {
		point = mirrorCoordsToAway(point)
	}
	cx := float64(point.X) / p.regionWidth
	cy := float64(point.Y) / p.regionHeight
	col := uint8(math.Min(cx, float64(p.cols-1)))
	row := uint8(math.Min(cy, float64(p.rows-1)))
	return p.GetRegion(col, row)
}

type region struct {
	col     uint8
	row     uint8
	sideRef proto.Team_Side
	center  proto.Point
}

func (r region) Col() uint8 {
	return r.col
}

func (r region) Row() uint8 {
	return r.row
}

func (r region) Center() proto.Point {
	return r.center
}

// Invert the coords X and Y as in a mirror to found out the same position seen from the away team field
// Keep in mind that all coords in the field is based on the bottom left corner!
func mirrorCoordsToAway(coords proto.Point) proto.Point {
	return proto.Point{
		X: field.FieldWidth - coords.X,
		Y: field.FieldHeight - coords.Y,
	}
}