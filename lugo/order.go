package lugo

// PlayerOrder is a hack to create a single type that includes all orders structs.
// The structs and interfaces generated from the protobuf files does not creates a common type that allows us
// to set a type of argument. So this interface make it possible
type PlayerOrder interface {
	UnifierLugoOrders()
	isOrder_Action()
}

func (*Order_Move) UnifierLugoOrders()  {}
func (*Order_Catch) UnifierLugoOrders() {}
func (*Order_Kick) UnifierLugoOrders()  {}
func (*Order_Jump) UnifierLugoOrders()  {}
