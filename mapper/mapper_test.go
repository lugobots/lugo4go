package mapper

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lugobots/lugo4go/v3/proto"
	"github.com/lugobots/lugo4go/v3/specs"
)

func TestNewPositioner(t *testing.T) {
	p, err := NewMapper(MinCols, MinRows, proto.Team_HOME)
	assert.Nil(t, err)

	assert.Equal(t, specs.FieldWidth/int(MinCols), int(p.regionWidth))
	assert.Equal(t, specs.FieldHeight/int(MinRows), int(p.regionHeight))
}

func TestNewPositioner_InvalidArgs(t *testing.T) {
	p, err := NewMapper(MinCols-1, MinRows, proto.Team_HOME)
	assert.Nil(t, p)
	assert.Equal(t, ErrMinCols, err)

	p, err = NewMapper(MaxCols+1, MinRows, proto.Team_HOME)
	assert.Nil(t, p)
	assert.Equal(t, ErrMaxCols, err)

	p, err = NewMapper(MinCols, MinRows-1, proto.Team_HOME)
	assert.Nil(t, p)
	assert.Equal(t, ErrMinRows, err)

	p, err = NewMapper(MinCols, MaxCols+1, proto.Team_HOME)
	assert.Nil(t, p)
	assert.Equal(t, ErrMaxRows, err)
}

func TestRegion_Center_HomeTeam(t *testing.T) {
	type testCase struct {
		cols             int
		rows             int
		regionHalfWidth  int
		regionHalfHeight int
	}

	testCases := map[string]testCase{
		"minimals":  {cols: MinCols, rows: MinRows, regionHalfWidth: 2500, regionHalfHeight: 2500},
		"maximums":  {cols: MaxCols, rows: MaxRows, regionHalfWidth: 50, regionHalfHeight: 50},
		"custom-1":  {cols: 10, rows: 10, regionHalfWidth: 1000, regionHalfHeight: 500},
		"inexact-2": {cols: 12, rows: 6, regionHalfWidth: 833, regionHalfHeight: 833},
	}

	team := proto.Team_HOME

	for testName, testSettings := range testCases {
		t.Run(testName, func(t *testing.T) {

			p, err := NewMapper(testSettings.cols, testSettings.rows, team)
			assert.Nil(t, err)
			expectedPointDefenseRight := &proto.Point{X: int32(testSettings.regionHalfWidth), Y: int32(testSettings.regionHalfHeight)}
			expectedPointDefenseLeft := &proto.Point{X: int32(+testSettings.regionHalfWidth), Y: int32(specs.MaxYCoordinate - testSettings.regionHalfHeight)}
			expectedPointAttackLeft := &proto.Point{X: int32(specs.MaxXCoordinate - testSettings.regionHalfWidth), Y: int32(specs.MaxYCoordinate - testSettings.regionHalfHeight)}
			expectedPointAttackRight := &proto.Point{X: int32(specs.MaxXCoordinate - testSettings.regionHalfWidth), Y: int32(testSettings.regionHalfHeight)}

			r, err := p.GetRegion(0, 0)
			assert.Nil(t, err)
			assert.Equal(t, expectedPointDefenseRight, r.Center(), testName)

			r, err = p.GetRegion(0, testSettings.rows-1)
			assert.Nil(t, err)
			assert.Equal(t, expectedPointDefenseLeft, r.Center(), testName)

			r, err = p.GetRegion(testSettings.cols-1, testSettings.rows-1)
			assert.Nil(t, err)
			assert.Equal(t, expectedPointAttackLeft, r.Center(), testName)

			r, err = p.GetRegion(testSettings.cols-1, 0)
			assert.Nil(t, err)
			assert.Equal(t, expectedPointAttackRight, r.Center(), testName)
		})
	}
}

func TestRegion_Center_Away(t *testing.T) {
	type testCase struct {
		cols             int
		rows             int
		regionHalfWidth  int
		regionHalfHeight int
	}

	testCases := map[string]testCase{
		"minimals":  {cols: MinCols, rows: MinRows, regionHalfWidth: 2500, regionHalfHeight: 2500},
		"maximums":  {cols: MaxCols, rows: MaxRows, regionHalfWidth: 50, regionHalfHeight: 50},
		"custom-1":  {cols: 10, rows: 10, regionHalfWidth: 1000, regionHalfHeight: 500},
		"inexact-2": {cols: 12, rows: 6, regionHalfWidth: 833, regionHalfHeight: 833},
	}

	team := proto.Team_AWAY

	for testName, testSettings := range testCases {
		t.Run(testName, func(t *testing.T) {
			p, err := NewMapper(testSettings.cols, testSettings.rows, team)
			assert.Nil(t, err)
			expectedPointDefenseRight := proto.Point{X: int32(specs.MaxXCoordinate - testSettings.regionHalfWidth), Y: int32(specs.MaxYCoordinate - testSettings.regionHalfHeight)}
			expectedPointDefenseLeft := proto.Point{X: int32(specs.MaxXCoordinate - testSettings.regionHalfWidth), Y: int32(testSettings.regionHalfHeight)}
			expectedPointAttackLeft := proto.Point{X: int32(testSettings.regionHalfWidth), Y: int32(testSettings.regionHalfHeight)}
			expectedPointAttackRight := proto.Point{X: int32(testSettings.regionHalfWidth), Y: int32(specs.MaxYCoordinate - testSettings.regionHalfHeight)}

			r, err := p.GetRegion(0, 0)
			assert.Nil(t, err)
			assert.Equal(t, &expectedPointDefenseRight, r.Center(), testName)

			r, err = p.GetRegion(0, testSettings.rows-1)
			assert.Nil(t, err)
			assert.Equal(t, &expectedPointDefenseLeft, r.Center(), testName)

			r, err = p.GetRegion(testSettings.cols-1, testSettings.rows-1)
			assert.Nil(t, err)
			assert.Equal(t, &expectedPointAttackLeft, r.Center(), testName)

			r, err = p.GetRegion(testSettings.cols-1, 0)
			assert.Nil(t, err)
			assert.Equal(t, &expectedPointAttackRight, r.Center(), testName)
		})
	}
}

func TestPositioner_GetRegion_InvalidArgs(t *testing.T) {
	p, err := NewMapper(10, 10, proto.Team_AWAY)
	assert.Nil(t, err)

	r, err := p.GetRegion(11, 5)
	assert.Nil(t, r)
	assert.Equal(t, ErrMaxCols, err)

	r, err = p.GetRegion(10, 5)
	assert.Nil(t, r)
	assert.Equal(t, ErrMaxCols, err)

	r, err = p.GetRegion(9, 11)
	assert.Nil(t, r)
	assert.Equal(t, ErrMaxRows, err)

	r, err = p.GetRegion(9, 10)
	assert.Nil(t, r)
	assert.Equal(t, ErrMaxRows, err)

}

func TestPositioner_GetPointRegion_HomeTeam(t *testing.T) {
	type testCase struct {
		cols             int
		rows             int
		regionHalfWidth  int
		regionHalfHeight int
	}

	testCases := map[string]testCase{
		"minimals":  {cols: MinCols, rows: MinRows, regionHalfWidth: 2500, regionHalfHeight: 2500},
		"maximums":  {cols: MaxCols, rows: MaxRows, regionHalfWidth: 50, regionHalfHeight: 50},
		"custom-1":  {cols: 10, rows: 10, regionHalfWidth: 1000, regionHalfHeight: 500},
		"inexact-2": {cols: 12, rows: 6, regionHalfWidth: 833, regionHalfHeight: 833},
	}

	team := proto.Team_HOME

	for testName, testSettings := range testCases {

		p, err := NewMapper(testSettings.cols, testSettings.rows, team)
		assert.Nil(t, err)
		pointDefenseRight := &proto.Point{X: int32(testSettings.regionHalfWidth), Y: int32(testSettings.regionHalfHeight)}
		pointDefenseLeft := &proto.Point{X: int32(+testSettings.regionHalfWidth), Y: int32(specs.FieldHeight - testSettings.regionHalfHeight)}
		pointAttackLeft := &proto.Point{X: int32(specs.FieldWidth - testSettings.regionHalfWidth), Y: int32(specs.FieldHeight - testSettings.regionHalfHeight)}
		pointAttackRight := &proto.Point{X: int32(specs.FieldWidth - testSettings.regionHalfWidth), Y: int32(testSettings.regionHalfHeight)}

		r, err := p.GetPointRegion(pointDefenseRight)
		assert.Nil(t, err)
		assert.Equal(t, int(0), r.Col(), testName)
		assert.Equal(t, int(0), r.Row(), testName)

		r, err = p.GetPointRegion(pointDefenseLeft)
		assert.Nil(t, err)
		assert.Equal(t, int(0), r.Col(), testName)
		assert.Equal(t, testSettings.rows-1, r.Row(), testName)

		r, err = p.GetPointRegion(pointAttackLeft)
		assert.Nil(t, err)
		assert.Equal(t, testSettings.cols-1, r.Col(), testName)
		assert.Equal(t, testSettings.rows-1, r.Row(), testName)

		r, err = p.GetPointRegion(pointAttackRight)
		assert.Nil(t, err)
		assert.Equal(t, testSettings.cols-1, r.Col(), testName)
		assert.Equal(t, int(0), r.Row(), testName)
	}
}

func TestPositioner_GetPointRegion_AwayTeam(t *testing.T) {
	type testCase struct {
		cols             int
		rows             int
		regionHalfWidth  int
		regionHalfHeight int
	}

	testCases := map[string]testCase{
		"minimals":  {cols: MinCols, rows: MinRows, regionHalfWidth: int(2500), regionHalfHeight: int(2500)},
		"maximums":  {cols: MaxCols, rows: MaxRows, regionHalfWidth: int(50), regionHalfHeight: int(50)},
		"custom-1":  {cols: 10, rows: 10, regionHalfWidth: int(1000), regionHalfHeight: int(500)},
		"inexact-2": {cols: 12, rows: 6, regionHalfWidth: int(833), regionHalfHeight: int(833)},
	}

	team := proto.Team_AWAY

	for testName, testSettings := range testCases {

		p, err := NewMapper(testSettings.cols, testSettings.rows, team)
		assert.Nil(t, err)
		pointDefenseRight := &proto.Point{X: int32(specs.FieldWidth - testSettings.regionHalfWidth), Y: int32(specs.FieldHeight - testSettings.regionHalfHeight)}
		pointDefenseLeft := &proto.Point{X: int32(specs.FieldWidth - testSettings.regionHalfWidth), Y: int32(testSettings.regionHalfHeight)}
		pointAttackLeft := &proto.Point{X: int32(testSettings.regionHalfWidth), Y: int32(testSettings.regionHalfHeight)}
		pointAttackRight := &proto.Point{X: int32(testSettings.regionHalfWidth), Y: int32(specs.FieldHeight - testSettings.regionHalfHeight)}

		r, err := p.GetPointRegion(pointDefenseRight)
		assert.Nil(t, err)
		assert.Equal(t, int(0), r.Col(), testName)
		assert.Equal(t, int(0), r.Row(), testName)

		r, err = p.GetPointRegion(pointDefenseLeft)
		assert.Nil(t, err)
		assert.Equal(t, int(0), r.Col(), testName)
		assert.Equal(t, testSettings.rows-1, r.Row(), testName)

		r, err = p.GetPointRegion(pointAttackLeft)
		assert.Nil(t, err)
		assert.Equal(t, testSettings.cols-1, r.Col(), testName)
		assert.Equal(t, testSettings.rows-1, r.Row(), testName)

		r, err = p.GetPointRegion(pointAttackRight)
		assert.Nil(t, err)
		assert.Equal(t, testSettings.cols-1, r.Col(), testName)
		assert.Equal(t, int(0), r.Row(), testName)
	}
}

func TestRegion_Front(t *testing.T) {
	maxCol := int(9)
	maxRow := int(9)

	type testCase struct {
		col          int
		row          int
		expectedCol  int
		expectedRow  int
		expectedSame bool
	}

	testCases := map[string]testCase{
		"center":             {5, 5, 6, 5, false},
		"back-right-corner":  {0, 0, 1, 0, false},
		"back-left-corner":   {0, 9, 1, 9, false},
		"front-right-corner": {9, 0, 9, 0, true},
		"front-left-corner":  {9, 9, 9, 9, true},
	}
	testTeamRegions := func(teamSide proto.Team_Side) {
		p, err := NewMapper(maxCol+1, maxRow+1, teamSide)
		assert.Nil(t, err)
		for testName, testSettings := range testCases {
			regionTestCase, err := p.GetRegion(testSettings.col, testSettings.row)
			assert.Nil(t, err, fmt.Sprintf("%s: test settings are invalid", testName))
			regionActual := regionTestCase.Front()
			assert.Equal(t, testSettings.expectedCol, regionActual.Col(), testName)
			assert.Equal(t, testSettings.expectedRow, regionActual.Row(), testName)
			assert.Equal(t, testSettings.expectedSame, reflect.DeepEqual(regionActual, regionTestCase))
		}
	}
	testTeamRegions(proto.Team_HOME)
	testTeamRegions(proto.Team_AWAY)
}

func TestRegion_Back(t *testing.T) {
	maxCol := 9
	maxRow := 9

	type testCase struct {
		col          int
		row          int
		expectedCol  int
		expectedRow  int
		expectedSame bool
	}

	testCases := map[string]testCase{
		"center":             {5, 5, 4, 5, false},
		"back-right-corner":  {0, 0, 0, 0, true},
		"back-left-corner":   {0, 9, 0, 9, true},
		"front-right-corner": {9, 0, 8, 0, false},
		"front-left-corner":  {9, 9, 8, 9, false},
	}
	testTeamRegions := func(teamSide proto.Team_Side) {
		p, err := NewMapper(maxCol+1, maxRow+1, teamSide)
		assert.Nil(t, err)
		for testName, testSettings := range testCases {
			regionTestCase, err := p.GetRegion(testSettings.col, testSettings.row)
			assert.Nil(t, err, fmt.Sprintf("%s: test settings are invalid", testName))
			regionActual := regionTestCase.Back()
			assert.Equal(t, testSettings.expectedCol, regionActual.Col(), testName)
			assert.Equal(t, testSettings.expectedRow, regionActual.Row(), testName)
			assert.Equal(t, testSettings.expectedSame, reflect.DeepEqual(regionActual, regionTestCase))
		}
	}
	testTeamRegions(proto.Team_HOME)
	testTeamRegions(proto.Team_AWAY)
}

func TestRegion_Left(t *testing.T) {
	maxCol := int(9)
	maxRow := int(9)
	type testCase struct {
		col          int
		row          int
		expectedCol  int
		expectedRow  int
		expectedSame bool
	}

	testCases := map[string]testCase{
		"center":             {5, 5, 5, 6, false},
		"back-right-corner":  {0, 0, 0, 1, false},
		"back-left-corner":   {0, 9, 0, 9, true},
		"front-right-corner": {9, 0, 9, 1, false},
		"front-left-corner":  {9, 9, 9, 9, true},
	}
	testTeamRegions := func(teamSide proto.Team_Side) {
		p, err := NewMapper(maxCol+1, maxRow+1, teamSide)
		assert.Nil(t, err)
		for testName, testSettings := range testCases {
			regionTestCase, err := p.GetRegion(testSettings.col, testSettings.row)
			assert.Nil(t, err, fmt.Sprintf("%s: test settings are invalid", testName))
			regionActual := regionTestCase.Left()
			assert.Equal(t, testSettings.expectedCol, regionActual.Col(), testName)
			assert.Equal(t, testSettings.expectedRow, regionActual.Row(), testName)
			assert.Equal(t, testSettings.expectedSame, reflect.DeepEqual(regionActual, regionTestCase))
		}
	}
	testTeamRegions(proto.Team_HOME)
	testTeamRegions(proto.Team_AWAY)
}

func TestRegion_Right(t *testing.T) {
	maxCol := int(9)
	maxRow := int(9)
	type testCase struct {
		col          int
		row          int
		expectedCol  int
		expectedRow  int
		expectedSame bool
	}

	testCases := map[string]testCase{
		"center":             {5, 5, 5, 4, false},
		"back-right-corner":  {0, 0, 0, 0, true},
		"back-left-corner":   {0, 9, 0, 8, false},
		"front-right-corner": {9, 0, 9, 0, true},
		"front-left-corner":  {9, 9, 9, 8, false},
	}
	testTeamRegions := func(teamSide proto.Team_Side) {
		p, err := NewMapper(maxCol+1, maxRow+1, teamSide)
		assert.Nil(t, err)
		for testName, testSettings := range testCases {
			regionTestCase, err := p.GetRegion(testSettings.col, testSettings.row)
			assert.Nil(t, err, fmt.Sprintf("%s: test settings are invalid", testName))
			regionActual := regionTestCase.Right()
			assert.Equal(t, testSettings.expectedCol, regionActual.Col(), testName)
			assert.Equal(t, testSettings.expectedRow, regionActual.Row(), testName)
			assert.Equal(t, testSettings.expectedSame, reflect.DeepEqual(regionActual, regionTestCase))
		}
	}
	testTeamRegions(proto.Team_HOME)
	testTeamRegions(proto.Team_AWAY)
}
