package coach

import (
	"fmt"
	"github.com/lugobots/lugo4go/v2/field"
	"github.com/lugobots/lugo4go/v2/proto"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewPositioner(t *testing.T) {
	p, err := NewPositioner(MinCols, MinRows, lugo.Team_HOME)
	assert.Nil(t, err)

	myStruct, ok := p.(*positioner)
	assert.True(t, ok)
	assert.Equal(t, field.FieldWidth/int(MinCols), int(myStruct.regionWidth))
	assert.Equal(t, field.FieldHeight/int(MinRows), int(myStruct.regionHeight))
}

func TestNewPositioner_InvalidArgs(t *testing.T) {
	p, err := NewPositioner(MinCols-1, MinRows, lugo.Team_HOME)
	assert.Nil(t, p)
	assert.Equal(t, ErrMinCols, err)

	p, err = NewPositioner(MaxCols+1, MinRows, lugo.Team_HOME)
	assert.Nil(t, p)
	assert.Equal(t, ErrMaxCols, err)

	p, err = NewPositioner(MinCols, MinRows-1, lugo.Team_HOME)
	assert.Nil(t, p)
	assert.Equal(t, ErrMinRows, err)

	p, err = NewPositioner(MinCols, MaxCols+1, lugo.Team_HOME)
	assert.Nil(t, p)
	assert.Equal(t, ErrMaxRows, err)
}

func TestRegion_Center_HomeTeam(t *testing.T) {
	type testCase struct {
		cols             uint8
		rows             uint8
		regionHalfWidth  int32
		regionHalfHeight int32
	}

	testCases := map[string]testCase{
		"minimals":  {cols: MinCols, rows: MinRows, regionHalfWidth: int32(2500), regionHalfHeight: int32(2500)},
		"maximums":  {cols: MaxCols, rows: MaxRows, regionHalfWidth: int32(500), regionHalfHeight: int32(500)},
		"custom-1":  {cols: 10, rows: 10, regionHalfWidth: int32(1000), regionHalfHeight: int32(500)},
		"inexact-2": {cols: 12, rows: 6, regionHalfWidth: int32(833), regionHalfHeight: int32(833)},
	}

	team := lugo.Team_HOME

	for testName, testSettings := range testCases {

		p, err := NewPositioner(testSettings.cols, testSettings.rows, team)
		assert.Nil(t, err)
		expectedPointDefenseRight := lugo.Point{X: testSettings.regionHalfWidth, Y: testSettings.regionHalfHeight}
		expectedPointDefenseLeft := lugo.Point{X: +testSettings.regionHalfWidth, Y: field.FieldHeight - testSettings.regionHalfHeight}
		expectedPointAttackLeft := lugo.Point{X: field.FieldWidth - testSettings.regionHalfWidth, Y: field.FieldHeight - testSettings.regionHalfHeight}
		expectedPointAttackRight := lugo.Point{X: field.FieldWidth - testSettings.regionHalfWidth, Y: testSettings.regionHalfHeight}

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
	}
}

func TestRegion_Center_Away(t *testing.T) {
	type testCase struct {
		cols             uint8
		rows             uint8
		regionHalfWidth  int32
		regionHalfHeight int32
	}

	testCases := map[string]testCase{
		"minimals":  {cols: MinCols, rows: MinRows, regionHalfWidth: int32(2500), regionHalfHeight: int32(2500)},
		"maximums":  {cols: MaxCols, rows: MaxRows, regionHalfWidth: int32(500), regionHalfHeight: int32(500)},
		"custom-1":  {cols: 10, rows: 10, regionHalfWidth: int32(1000), regionHalfHeight: int32(500)},
		"inexact-2": {cols: 12, rows: 6, regionHalfWidth: int32(833), regionHalfHeight: int32(833)},
	}

	team := lugo.Team_AWAY

	for testName, testSettings := range testCases {

		p, err := NewPositioner(testSettings.cols, testSettings.rows, team)
		assert.Nil(t, err)
		expectedPointDefenseRight := lugo.Point{X: field.FieldWidth - testSettings.regionHalfWidth, Y: field.FieldHeight - testSettings.regionHalfHeight}
		expectedPointDefenseLeft := lugo.Point{X: field.FieldWidth - testSettings.regionHalfWidth, Y: testSettings.regionHalfHeight}
		expectedPointAttackLeft := lugo.Point{X: testSettings.regionHalfWidth, Y: testSettings.regionHalfHeight}
		expectedPointAttackRight := lugo.Point{X: testSettings.regionHalfWidth, Y: field.FieldHeight - testSettings.regionHalfHeight}

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
	}
}

func TestPositioner_GetRegion_InvalidArgs(t *testing.T) {
	p, err := NewPositioner(10, 10, lugo.Team_AWAY)
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
		cols             uint8
		rows             uint8
		regionHalfWidth  int32
		regionHalfHeight int32
	}

	testCases := map[string]testCase{
		"minimals":  {cols: MinCols, rows: MinRows, regionHalfWidth: int32(2500), regionHalfHeight: int32(2500)},
		"maximums":  {cols: MaxCols, rows: MaxRows, regionHalfWidth: int32(500), regionHalfHeight: int32(500)},
		"custom-1":  {cols: 10, rows: 10, regionHalfWidth: int32(1000), regionHalfHeight: int32(500)},
		"inexact-2": {cols: 12, rows: 6, regionHalfWidth: int32(833), regionHalfHeight: int32(833)},
	}

	team := lugo.Team_HOME

	for testName, testSettings := range testCases {

		p, err := NewPositioner(testSettings.cols, testSettings.rows, team)
		assert.Nil(t, err)
		pointDefenseRight := lugo.Point{X: testSettings.regionHalfWidth, Y: testSettings.regionHalfHeight}
		pointDefenseLeft := lugo.Point{X: +testSettings.regionHalfWidth, Y: field.FieldHeight - testSettings.regionHalfHeight}
		pointAttackLeft := lugo.Point{X: field.FieldWidth - testSettings.regionHalfWidth, Y: field.FieldHeight - testSettings.regionHalfHeight}
		pointAttackRight := lugo.Point{X: field.FieldWidth - testSettings.regionHalfWidth, Y: testSettings.regionHalfHeight}

		r, err := p.GetPointRegion(pointDefenseRight)
		assert.Nil(t, err)
		assert.Equal(t, uint8(0), r.Col(), testName)
		assert.Equal(t, uint8(0), r.Row(), testName)

		r, err = p.GetPointRegion(pointDefenseLeft)
		assert.Nil(t, err)
		assert.Equal(t, uint8(0), r.Col(), testName)
		assert.Equal(t, testSettings.rows-1, r.Row(), testName)

		r, err = p.GetPointRegion(pointAttackLeft)
		assert.Nil(t, err)
		assert.Equal(t, testSettings.cols-1, r.Col(), testName)
		assert.Equal(t, testSettings.rows-1, r.Row(), testName)

		r, err = p.GetPointRegion(pointAttackRight)
		assert.Nil(t, err)
		assert.Equal(t, testSettings.cols-1, r.Col(), testName)
		assert.Equal(t, uint8(0), r.Row(), testName)
	}
}

func TestPositioner_GetPointRegion_AwayTeam(t *testing.T) {
	type testCase struct {
		cols             uint8
		rows             uint8
		regionHalfWidth  int32
		regionHalfHeight int32
	}

	testCases := map[string]testCase{
		"minimals":  {cols: MinCols, rows: MinRows, regionHalfWidth: int32(2500), regionHalfHeight: int32(2500)},
		"maximums":  {cols: MaxCols, rows: MaxRows, regionHalfWidth: int32(500), regionHalfHeight: int32(500)},
		"custom-1":  {cols: 10, rows: 10, regionHalfWidth: int32(1000), regionHalfHeight: int32(500)},
		"inexact-2": {cols: 12, rows: 6, regionHalfWidth: int32(833), regionHalfHeight: int32(833)},
	}

	team := lugo.Team_AWAY

	for testName, testSettings := range testCases {

		p, err := NewPositioner(testSettings.cols, testSettings.rows, team)
		assert.Nil(t, err)
		pointDefenseRight := lugo.Point{X: field.FieldWidth - testSettings.regionHalfWidth, Y: field.FieldHeight - testSettings.regionHalfHeight}
		pointDefenseLeft := lugo.Point{X: field.FieldWidth - testSettings.regionHalfWidth, Y: testSettings.regionHalfHeight}
		pointAttackLeft := lugo.Point{X: testSettings.regionHalfWidth, Y: testSettings.regionHalfHeight}
		pointAttackRight := lugo.Point{X: testSettings.regionHalfWidth, Y: field.FieldHeight - testSettings.regionHalfHeight}

		r, err := p.GetPointRegion(pointDefenseRight)
		assert.Nil(t, err)
		assert.Equal(t, uint8(0), r.Col(), testName)
		assert.Equal(t, uint8(0), r.Row(), testName)

		r, err = p.GetPointRegion(pointDefenseLeft)
		assert.Nil(t, err)
		assert.Equal(t, uint8(0), r.Col(), testName)
		assert.Equal(t, testSettings.rows-1, r.Row(), testName)

		r, err = p.GetPointRegion(pointAttackLeft)
		assert.Nil(t, err)
		assert.Equal(t, testSettings.cols-1, r.Col(), testName)
		assert.Equal(t, testSettings.rows-1, r.Row(), testName)

		r, err = p.GetPointRegion(pointAttackRight)
		assert.Nil(t, err)
		assert.Equal(t, testSettings.cols-1, r.Col(), testName)
		assert.Equal(t, uint8(0), r.Row(), testName)
	}
}

func TestRegion_Front(t *testing.T) {
	maxCol := uint8(9)
	maxRow := uint8(9)

	type testCase struct {
		col          uint8
		row          uint8
		expectedCol  uint8
		expectedRow  uint8
		expectedSame bool
	}

	testCases := map[string]testCase{
		"center":             {5, 5, 6, 5, false},
		"back-right-corner":  {0, 0, 1, 0, false},
		"back-left-corner":   {0, 9, 1, 9, false},
		"front-right-corner": {9, 0, 9, 0, true},
		"front-left-corner":  {9, 9, 9, 9, true},
	}
	testTeamRegions := func(teamSide lugo.Team_Side) {
		p, err := NewPositioner(maxCol+1, maxRow+1, teamSide)
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
	testTeamRegions(lugo.Team_HOME)
	testTeamRegions(lugo.Team_AWAY)
}

func TestRegion_Back(t *testing.T) {
	maxCol := uint8(9)
	maxRow := uint8(9)

	type testCase struct {
		col          uint8
		row          uint8
		expectedCol  uint8
		expectedRow  uint8
		expectedSame bool
	}

	testCases := map[string]testCase{
		"center":             {5, 5, 4, 5, false},
		"back-right-corner":  {0, 0, 0, 0, true},
		"back-left-corner":   {0, 9, 0, 9, true},
		"front-right-corner": {9, 0, 8, 0, false},
		"front-left-corner":  {9, 9, 8, 9, false},
	}
	testTeamRegions := func(teamSide lugo.Team_Side) {
		p, err := NewPositioner(maxCol+1, maxRow+1, teamSide)
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
	testTeamRegions(lugo.Team_HOME)
	testTeamRegions(lugo.Team_AWAY)
}

func TestRegion_Left(t *testing.T) {
	maxCol := uint8(9)
	maxRow := uint8(9)
	type testCase struct {
		col          uint8
		row          uint8
		expectedCol  uint8
		expectedRow  uint8
		expectedSame bool
	}

	testCases := map[string]testCase{
		"center":             {5, 5, 5, 6, false},
		"back-right-corner":  {0, 0, 0, 1, false},
		"back-left-corner":   {0, 9, 0, 9, true},
		"front-right-corner": {9, 0, 9, 1, false},
		"front-left-corner":  {9, 9, 9, 9, true},
	}
	testTeamRegions := func(teamSide lugo.Team_Side) {
		p, err := NewPositioner(maxCol+1, maxRow+1, teamSide)
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
	testTeamRegions(lugo.Team_HOME)
	testTeamRegions(lugo.Team_AWAY)
}

func TestRegion_Right(t *testing.T) {
	maxCol := uint8(9)
	maxRow := uint8(9)
	type testCase struct {
		col          uint8
		row          uint8
		expectedCol  uint8
		expectedRow  uint8
		expectedSame bool
	}

	testCases := map[string]testCase{
		"center":             {5, 5, 5, 4, false},
		"back-right-corner":  {0, 0, 0, 0, true},
		"back-left-corner":   {0, 9, 0, 8, false},
		"front-right-corner": {9, 0, 9, 0, true},
		"front-left-corner":  {9, 9, 9, 8, false},
	}
	testTeamRegions := func(teamSide lugo.Team_Side) {
		p, err := NewPositioner(maxCol+1, maxRow+1, teamSide)
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
	testTeamRegions(lugo.Team_HOME)
	testTeamRegions(lugo.Team_AWAY)
}
