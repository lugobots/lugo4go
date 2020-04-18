package lugo4go

// Error is used to identify internal errors
type Error string

// Error implements the native golang error interface
func (e Error) Error() string { return string(e) }

const (
	// ErrGRPCConnectionClosed identifies when the error returned is cased by the connection has been closed
	ErrGRPCConnectionClosed = Error("grpc connection closed by the server")

	// ErrGRPCConnectionLost identifies that something unexpected broke the gRPC connection
	ErrGRPCConnectionLost = Error("grpc stream error")
)
