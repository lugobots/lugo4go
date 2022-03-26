package proto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindIntersection_SameLine(t *testing.T) {
	A1 := Point{}
	A2 := Point{X: 10}

	B1 := Point{}
	B2 := Point{X: 10}

	_, _, err := FindIntersection(A1, A2, B1, B2)
	assert.NotNil(t, err)
}

func TestFindIntersection_SameLineNotOverlapping(t *testing.T) {
	A1 := Point{}
	A2 := Point{X: 10}

	B1 := Point{X: 12}
	B2 := Point{X: 20}

	_, _, err := FindIntersection(A1, A2, B1, B2)
	assert.NotNil(t, err)
}

func TestFindIntersection_ParallelLines(t *testing.T) {
	A1 := Point{}
	A2 := Point{X: 10}

	B1 := Point{Y: 10, X: 12}
	B2 := Point{Y: 10, X: 20}

	_, _, err := FindIntersection(A1, A2, B1, B2)
	assert.NotNil(t, err)
}

func TestFindIntersection_CrossAndTouch(t *testing.T) {
	A1 := Point{}
	A2 := Point{X: 10}

	B1 := Point{X: 5, Y: -5}
	B2 := Point{X: 5, Y: 5}

	p, touch, err := FindIntersection(A1, A2, B1, B2)
	assert.Nil(t, err)
	assert.Equal(t, Point{X: 5, Y: 0}, p)
	assert.True(t, touch)
}

func TestFindIntersection_CrossAndDoNotTouch(t *testing.T) {
	A1 := Point{}
	A2 := Point{X: 10}

	B1 := Point{X: 5, Y: 15}
	B2 := Point{X: 5, Y: 5}

	p, touch, err := FindIntersection(A1, A2, B1, B2)
	assert.Nil(t, err)
	assert.Equal(t, Point{X: 5, Y: 0}, p)
	assert.False(t, touch)
}

func TestFindIntersection_DiagonalTouching(t *testing.T) {
	A1 := Point{}
	A2 := Point{X: 10, Y: 10}

	B1 := Point{X: 0, Y: 10}
	B2 := Point{X: 10, Y: 0}

	p, touch, err := FindIntersection(A1, A2, B1, B2)
	assert.Nil(t, err)
	assert.Equal(t, Point{X: 5, Y: 5}, p)
	assert.True(t, touch)
}

func TestFindIntersection_DiagonalNotTouching(t *testing.T) {
	A1 := Point{}
	A2 := Point{X: 10, Y: 10}

	B1 := Point{X: 0, Y: 30}
	B2 := Point{X: 10, Y: 20}

	p, touch, err := FindIntersection(A1, A2, B1, B2)
	assert.Nil(t, err)
	assert.Equal(t, Point{X: 15, Y: 15}, p)
	assert.False(t, touch)
}
