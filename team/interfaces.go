package team

import (
	"context"
	"fmt"
	"github.com/lugobots/lugo4go/v2/lugo"
)

type OrderSender interface {
	Send(ctx context.Context, turn uint32, orders []lugo.PlayerOrder, debugMsg string) (*lugo.OrderResponse, error)
}

type TurnOrdersSender interface {
	Send(ctx context.Context, orders []lugo.PlayerOrder, debugMsg string) (*lugo.OrderResponse, error)
}

type Bot interface {
	OnDisputing(ctx context.Context, sender TurnOrdersSender, snapshot *lugo.GameSnapshot) error
	OnDefending(ctx context.Context, sender TurnOrdersSender, snapshot *lugo.GameSnapshot) error
	OnHolding(ctx context.Context, sender TurnOrdersSender, snapshot *lugo.GameSnapshot) error
	OnSupporting(ctx context.Context, sender TurnOrdersSender, snapshot *lugo.GameSnapshot) error
	AsGoalkeeper(ctx context.Context, sender TurnOrdersSender, snapshot *lugo.GameSnapshot, state PlayerState) error
}

// Positioner Helps the bots to see the fields from their team perspective instead of using the cartesian plan provided
// by the game server. Instead of base your logic on the axes X and Y, the Arrangement create a FieldArea map based
// on the team side.
// The FieldArea coordinates uses the defensive field's right corner as its origin.
// This mechanism if specially useful to define players regions based on their roles, since you do not have to mirror
// the coordinate, neither do extra logic to define regions on the field where the player should be.
type Positioner interface {
	// GetRegion Returns a FieldArea based on the coordinates and on the current field division
	GetRegion(col, row uint8) (FieldNav, error)
	// GetPointRegion returns the FieldArea where that point is in
	GetPointRegion(point *lugo.Point) (FieldNav, error)
}

// FieldNav represent a quadrant on the field. It is not always squared form because you may define how many cols/rows
// the field will be divided in. So, based on that division (e.g. 4 rows, 6 cols) there will be a fixed number of regions
// and their coordinates will be zero-index (e.g. from 0 to 3 rows when divided in 4 rows).
type FieldNav interface {
	fmt.Stringer
	// Col The col coordinate based on the field division
	Col() uint8
	// Row The row coordinate based on the field division
	Row() uint8
	// Center Return the point at the center of the quadrant represented by this FieldNav. It is not always precise.
	Center() *lugo.Point

	// Front is the FieldArea immediately in front of this one from the player perspective
	// Important: The same FieldArea is returned if the requested FieldArea is not valid
	Front() FieldNav
	// Back is the FieldArea immediately behind this one from the player perspective
	// Important: The same FieldArea is returned if the requested FieldArea is not valid
	Back() FieldNav
	// Left is the FieldArea immediately on left of this one from the player perspective
	// Important: The same FieldArea is returned if the requested FieldArea is not valid
	Left() FieldNav
	// Right is the FieldArea immediately on right of this one from the player perspective
	// Important: The same FieldArea is returned if the requested FieldArea is not valid
	Right() FieldNav
}
