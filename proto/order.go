package proto

// PlayerOrder is a hack to create a single type that includes all orders structs.
// The structs and interfaces generated from the protobuf files does not creates a common type that allows us
// to set a type of argument. So this interface make it possible
type PlayerOrder interface {
	LugoOrdersUnifier()
	isOrder_Action()
}

// TROCAR NOME da library pra lugogrpc ou algo do tipo, e trocar o nome de ops para lugo

func (*Order_Move) LugoOrdersUnifier()  {}
func (*Order_Catch) LugoOrdersUnifier() {}
func (*Order_Kick) LugoOrdersUnifier()  {}
func (*Order_Jump) LugoOrdersUnifier()  {}
