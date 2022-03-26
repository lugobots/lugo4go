package proto

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVelocity_GetNextPoint(t *testing.T) {

	testCases := []struct {
		name     string
		velocity Velocity
		from     Point
		expected Point
	}{
		{"to north", Velocity{Direction: North().Copy(), Speed: 10}, Point{}, Point{X: 0, Y: 10}},
		{"to south", Velocity{Direction: South().Copy(), Speed: 300.1}, Point{X: 100, Y: 100}, Point{X: 100, Y: -200}},
		{"to west", Velocity{Direction: West().Copy(), Speed: 43.6}, Point{X: 20, Y: 233}, Point{X: -24, Y: 233}},
		{"to east", Velocity{Direction: East().Copy(), Speed: 0.4}, Point{}, Point{X: 0, Y: 0}},
		{"to northeast", Velocity{Direction: NorthEast().Copy(), Speed: 5}, Point{X: 10, Y: 40}, Point{X: 14, Y: 44}},
		{"to southwest", Velocity{Direction: SouthWest().Copy(), Speed: 5}, Point{}, Point{X: -4, Y: -4}},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.expected, testCase.velocity.GetNextPoint(testCase.from), fmt.Sprintf("Case %s has failed", testCase.name))
	}

}
