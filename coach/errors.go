package coach

// Error helps to define internal errors
type Error string

func (e Error) Error() string { return string(e) }

const (
	// ErrMinCols defines an error for invalid number of cols
	ErrMinCols = Error("number of cols lower the minimum")
	// ErrMaxCols defines an error for invalid number of cols
	ErrMaxCols = Error("number of cols higher the maximum")
	// ErrMinRows defines an error for invalid number of rows
	ErrMinRows = Error("number of rows lower the minimum")
	// ErrMaxRows defines an error for invalid number of rows
	ErrMaxRows = Error("number of rows higher the maximum")
)

const (
	ErrPlayerNotFound = Error("player not found in the game Snapshot")
	ErrNoBall         = Error("no ball found in the snapshot")
)
