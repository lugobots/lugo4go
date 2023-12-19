package proto

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVector_AngleWith_ZeroDegree(t *testing.T) {
	type tTable struct {
		vecA *Vector
		vecB *Vector
		ang  float64
	}
	testTable := map[string]tTable{}

	caseSample := tTable{}
	caseSample.vecA, _ = NewVector(Point{X: 0, Y: 0}, Point{X: 1, Y: 0})
	caseSample.vecB, _ = NewVector(Point{X: 0, Y: 0}, Point{X: 1, Y: 0})
	caseSample.ang = 0.0
	testTable["Same direction East"] = caseSample

	caseSample = tTable{}
	caseSample.vecA, _ = NewVector(Point{X: 0, Y: 0}, Point{X: 0, Y: 1})
	caseSample.vecB, _ = NewVector(Point{X: 0, Y: 0}, Point{X: 0, Y: 1})
	caseSample.ang = 0.0
	testTable["Same direction North"] = caseSample

	caseSample = tTable{}
	caseSample.vecA, _ = NewVector(Point{X: 0, Y: 0}, Point{X: -5, Y: -10})
	caseSample.vecB, _ = NewVector(Point{X: 0, Y: 0}, Point{X: -5, Y: -10})
	caseSample.ang = 0.0
	testTable["Same direction Southwest"] = caseSample

	caseSample = tTable{}
	caseSample.vecA, _ = NewVector(Point{X: 0, Y: 0}, Point{X: 1, Y: 0})
	caseSample.vecB, _ = NewVector(Point{X: 0, Y: 0}, Point{X: 0, Y: 1})
	caseSample.ang = 90.0
	testTable["90 degree North"] = caseSample

	caseSample = tTable{}
	caseSample.vecA, _ = NewVector(Point{X: 0, Y: 0}, Point{X: 1, Y: 0})
	caseSample.vecB, _ = NewVector(Point{X: 0, Y: 0}, Point{X: 0, Y: -1})
	caseSample.ang = -90.0
	testTable["90 degree South"] = caseSample

	caseSample = tTable{}
	caseSample.vecA, _ = NewVector(Point{X: 0, Y: 0}, Point{X: 1, Y: 0})
	caseSample.vecB, _ = NewVector(Point{X: 0, Y: 0}, Point{X: -1, Y: 0})
	caseSample.ang = 180
	testTable["180 degrees"] = caseSample

	caseSample = tTable{}
	caseSample.vecA, _ = NewVector(Point{X: 0, Y: 0}, Point{X: 1, Y: 0})
	caseSample.vecB, _ = NewVector(Point{X: 0, Y: 0}, Point{X: 1, Y: 1})
	caseSample.ang = 45
	testTable["45 degrees Northeast"] = caseSample

	caseSample = tTable{}
	caseSample.vecA, _ = NewVector(Point{X: 0, Y: 0}, Point{X: 1, Y: 0})
	caseSample.vecB, _ = NewVector(Point{X: 0, Y: 0}, Point{X: 1, Y: -1})
	caseSample.ang = -45
	testTable["45 degrees Southeast"] = caseSample

	caseSample = tTable{}
	caseSample.vecA, _ = NewVector(Point{X: 0, Y: 0}, Point{X: 1, Y: 0})
	caseSample.vecB, _ = NewVector(Point{X: 0, Y: 0}, Point{X: -1, Y: 1})
	caseSample.ang = 135
	testTable["135 degrees Northwest"] = caseSample

	caseSample = tTable{}
	caseSample.vecA, _ = NewVector(Point{X: 0, Y: 0}, Point{X: 1, Y: 0})
	caseSample.vecB, _ = NewVector(Point{X: 0, Y: 0}, Point{X: -1, Y: -1})
	caseSample.ang = -135
	testTable["135 degrees Southwest"] = caseSample

	caseSample = tTable{}
	caseSample.vecA, _ = NewVector(Point{X: 0, Y: 0}, Point{X: 1, Y: 1})
	caseSample.vecB, _ = NewVector(Point{X: 0, Y: 0}, Point{X: -1, Y: 1})
	caseSample.ang = 90
	testTable["90 both not zero"] = caseSample

	caseSample = tTable{}
	caseSample.vecA, _ = NewVector(Point{X: 0, Y: 0}, Point{X: 0, Y: 1})
	caseSample.vecB, _ = NewVector(Point{X: 0, Y: 0}, Point{X: 1, Y: 0})
	caseSample.ang = -90
	testTable["-90 degrees Wast"] = caseSample

	caseSample = tTable{}
	caseSample.vecA, _ = NewVector(Point{X: 0, Y: 0}, Point{X: 0, Y: 1})
	caseSample.vecB, _ = NewVector(Point{X: 0, Y: 0}, Point{X: -1, Y: 0})
	caseSample.ang = 90
	testTable["-90 degrees east"] = caseSample

	for title, conditions := range testTable {
		actualAng := conditions.vecA.AngleWith(conditions.vecB)
		assert.Equal(t, conditions.ang, actualAng, title)
	}

}

func TestVector_AddAngle(t *testing.T) {
	vecA, _ := NewVector(Point{X: 0, Y: 0}, Point{X: 100, Y: 0})

	vecA.AddAngleDegree(90)
	assert.Equal(t, float64(90), math.Round(vecA.AngleDegrees()))
	assert.True(t, vecA.X <= 0.00000001)
	assert.Equal(t, float64(100), vecA.Y)
	assert.Equal(t, float64(100), vecA.Length())

	vecA.AddAngleDegree(90)
	assert.Equal(t, float64(180), math.Round(vecA.AngleDegrees()))
	assert.Equal(t, float64(-100), vecA.X)
	assert.True(t, vecA.Y <= 0.00000001)
	assert.Equal(t, float64(100), vecA.Length())

	vecA.AddAngleDegree(90)
	assert.Equal(t, float64(-90), math.Round(vecA.AngleDegrees()))
	assert.True(t, vecA.X <= 0.00000001)
	assert.Equal(t, float64(-100), vecA.Y)
	assert.Equal(t, float64(100), vecA.Length())

	vecA.AddAngleDegree(90)
	assert.Equal(t, float64(0), math.Round(vecA.AngleDegrees()))
	assert.Equal(t, float64(100), vecA.X)
	assert.True(t, vecA.Y <= 0.00000001)
	assert.Equal(t, float64(100), vecA.Length())

	vecA.AddAngleDegree(45)
	assert.Equal(t, float64(45), math.Round(vecA.AngleDegrees()))
	assert.True(t, math.Abs(vecA.Y-vecA.X) <= 0.00000001)
	assert.Equal(t, float64(100), vecA.Length())

}
