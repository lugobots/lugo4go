package lugo4go

type Error string

func (e Error) Error() string { return string(e) }

const (
	ErrGRPCConnectionClosed = Error("grpc connection closed by the server")
	ErrGRPCConnectionLost   = Error("grpc stream error")
)
