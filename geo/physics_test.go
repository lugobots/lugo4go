package geo

import (
	"github.com/lugobots/lugo4go/v2/proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAngleWithRoute(t *testing.T) {
	type testCase struct {
		direction     lugo.Vector
		from          lugo.Point
		obstacle      lugo.Point
		expectedAngle float64
	}

	testCases := map[string]testCase{
		"in front":   {direction: lugo.North(), from: lugo.Point{}, obstacle: lugo.Point{Y: 100}, expectedAngle: 0},
		"behind":     {direction: lugo.North(), from: lugo.Point{Y: 100}, obstacle: lugo.Point{}, expectedAngle: 180},
		"Right side": {direction: lugo.North(), from: lugo.Point{}, obstacle: lugo.Point{X: 1}, expectedAngle: -90},
		"Left side":  {direction: lugo.North(), from: lugo.Point{X: 1}, obstacle: lugo.Point{}, expectedAngle: 90},
	}

	for caseName, def := range testCases {
		assert.Equal(t, def.expectedAngle, AngleWithRoute(def.direction, def.from, def.obstacle), caseName)
	}
}
