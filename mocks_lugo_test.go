// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/lugobots/lugo4go/v2/lugo (interfaces: PlayerOrder,GameServer,GameClient,Game_JoinATeamClient,Game_JoinATeamServer,BroadcastClient,Broadcast_OnEventClient,BroadcastServer,Broadcast_OnEventServer)

// Package lugo4go_test is a generated GoMock package.
package lugo4go_test

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	empty "github.com/golang/protobuf/ptypes/empty"
	lugo "github.com/lugobots/lugo4go/v2/lugo"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
	reflect "reflect"
)

// MockPlayerOrder is a mock of PlayerOrder interface
type MockPlayerOrder struct {
	ctrl     *gomock.Controller
	recorder *MockPlayerOrderMockRecorder
}

// MockPlayerOrderMockRecorder is the mock recorder for MockPlayerOrder
type MockPlayerOrderMockRecorder struct {
	mock *MockPlayerOrder
}

// NewMockPlayerOrder creates a new mock instance
func NewMockPlayerOrder(ctrl *gomock.Controller) *MockPlayerOrder {
	mock := &MockPlayerOrder{ctrl: ctrl}
	mock.recorder = &MockPlayerOrderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPlayerOrder) EXPECT() *MockPlayerOrderMockRecorder {
	return m.recorder
}

// LugoOrdersUnifier mocks base method
func (m *MockPlayerOrder) LugoOrdersUnifier() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "LugoOrdersUnifier")
}

// LugoOrdersUnifier indicates an expected call of LugoOrdersUnifier
func (mr *MockPlayerOrderMockRecorder) LugoOrdersUnifier() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LugoOrdersUnifier", reflect.TypeOf((*MockPlayerOrder)(nil).LugoOrdersUnifier))
}

// isOrder_Action mocks base method
func (m *MockPlayerOrder) isOrder_Action() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "isOrder_Action")
}

// isOrder_Action indicates an expected call of isOrder_Action
func (mr *MockPlayerOrderMockRecorder) isOrder_Action() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "isOrder_Action", reflect.TypeOf((*MockPlayerOrder)(nil).isOrder_Action))
}

// MockGameServer is a mock of GameServer interface
type MockGameServer struct {
	ctrl     *gomock.Controller
	recorder *MockGameServerMockRecorder
}

// MockGameServerMockRecorder is the mock recorder for MockGameServer
type MockGameServerMockRecorder struct {
	mock *MockGameServer
}

// NewMockGameServer creates a new mock instance
func NewMockGameServer(ctrl *gomock.Controller) *MockGameServer {
	mock := &MockGameServer{ctrl: ctrl}
	mock.recorder = &MockGameServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGameServer) EXPECT() *MockGameServerMockRecorder {
	return m.recorder
}

// JoinATeam mocks base method
func (m *MockGameServer) JoinATeam(arg0 *lugo.JoinRequest, arg1 lugo.Game_JoinATeamServer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "JoinATeam", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// JoinATeam indicates an expected call of JoinATeam
func (mr *MockGameServerMockRecorder) JoinATeam(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "JoinATeam", reflect.TypeOf((*MockGameServer)(nil).JoinATeam), arg0, arg1)
}

// SendOrders mocks base method
func (m *MockGameServer) SendOrders(arg0 context.Context, arg1 *lugo.OrderSet) (*lugo.OrderResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendOrders", arg0, arg1)
	ret0, _ := ret[0].(*lugo.OrderResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendOrders indicates an expected call of SendOrders
func (mr *MockGameServerMockRecorder) SendOrders(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendOrders", reflect.TypeOf((*MockGameServer)(nil).SendOrders), arg0, arg1)
}

// MockGameClient is a mock of GameClient interface
type MockGameClient struct {
	ctrl     *gomock.Controller
	recorder *MockGameClientMockRecorder
}

// MockGameClientMockRecorder is the mock recorder for MockGameClient
type MockGameClientMockRecorder struct {
	mock *MockGameClient
}

// NewMockGameClient creates a new mock instance
func NewMockGameClient(ctrl *gomock.Controller) *MockGameClient {
	mock := &MockGameClient{ctrl: ctrl}
	mock.recorder = &MockGameClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGameClient) EXPECT() *MockGameClientMockRecorder {
	return m.recorder
}

// JoinATeam mocks base method
func (m *MockGameClient) JoinATeam(arg0 context.Context, arg1 *lugo.JoinRequest, arg2 ...grpc.CallOption) (lugo.Game_JoinATeamClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "JoinATeam", varargs...)
	ret0, _ := ret[0].(lugo.Game_JoinATeamClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// JoinATeam indicates an expected call of JoinATeam
func (mr *MockGameClientMockRecorder) JoinATeam(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "JoinATeam", reflect.TypeOf((*MockGameClient)(nil).JoinATeam), varargs...)
}

// SendOrders mocks base method
func (m *MockGameClient) SendOrders(arg0 context.Context, arg1 *lugo.OrderSet, arg2 ...grpc.CallOption) (*lugo.OrderResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SendOrders", varargs...)
	ret0, _ := ret[0].(*lugo.OrderResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendOrders indicates an expected call of SendOrders
func (mr *MockGameClientMockRecorder) SendOrders(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendOrders", reflect.TypeOf((*MockGameClient)(nil).SendOrders), varargs...)
}

// MockGame_JoinATeamClient is a mock of Game_JoinATeamClient interface
type MockGame_JoinATeamClient struct {
	ctrl     *gomock.Controller
	recorder *MockGame_JoinATeamClientMockRecorder
}

// MockGame_JoinATeamClientMockRecorder is the mock recorder for MockGame_JoinATeamClient
type MockGame_JoinATeamClientMockRecorder struct {
	mock *MockGame_JoinATeamClient
}

// NewMockGame_JoinATeamClient creates a new mock instance
func NewMockGame_JoinATeamClient(ctrl *gomock.Controller) *MockGame_JoinATeamClient {
	mock := &MockGame_JoinATeamClient{ctrl: ctrl}
	mock.recorder = &MockGame_JoinATeamClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGame_JoinATeamClient) EXPECT() *MockGame_JoinATeamClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method
func (m *MockGame_JoinATeamClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend
func (mr *MockGame_JoinATeamClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockGame_JoinATeamClient)(nil).CloseSend))
}

// Context mocks base method
func (m *MockGame_JoinATeamClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context
func (mr *MockGame_JoinATeamClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockGame_JoinATeamClient)(nil).Context))
}

// Header mocks base method
func (m *MockGame_JoinATeamClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header
func (mr *MockGame_JoinATeamClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockGame_JoinATeamClient)(nil).Header))
}

// Recv mocks base method
func (m *MockGame_JoinATeamClient) Recv() (*lugo.GameSnapshot, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*lugo.GameSnapshot)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv
func (mr *MockGame_JoinATeamClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockGame_JoinATeamClient)(nil).Recv))
}

// RecvMsg mocks base method
func (m *MockGame_JoinATeamClient) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg
func (mr *MockGame_JoinATeamClientMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockGame_JoinATeamClient)(nil).RecvMsg), arg0)
}

// SendMsg mocks base method
func (m *MockGame_JoinATeamClient) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg
func (mr *MockGame_JoinATeamClientMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockGame_JoinATeamClient)(nil).SendMsg), arg0)
}

// Trailer mocks base method
func (m *MockGame_JoinATeamClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer
func (mr *MockGame_JoinATeamClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockGame_JoinATeamClient)(nil).Trailer))
}

// MockGame_JoinATeamServer is a mock of Game_JoinATeamServer interface
type MockGame_JoinATeamServer struct {
	ctrl     *gomock.Controller
	recorder *MockGame_JoinATeamServerMockRecorder
}

// MockGame_JoinATeamServerMockRecorder is the mock recorder for MockGame_JoinATeamServer
type MockGame_JoinATeamServerMockRecorder struct {
	mock *MockGame_JoinATeamServer
}

// NewMockGame_JoinATeamServer creates a new mock instance
func NewMockGame_JoinATeamServer(ctrl *gomock.Controller) *MockGame_JoinATeamServer {
	mock := &MockGame_JoinATeamServer{ctrl: ctrl}
	mock.recorder = &MockGame_JoinATeamServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGame_JoinATeamServer) EXPECT() *MockGame_JoinATeamServerMockRecorder {
	return m.recorder
}

// Context mocks base method
func (m *MockGame_JoinATeamServer) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context
func (mr *MockGame_JoinATeamServerMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockGame_JoinATeamServer)(nil).Context))
}

// RecvMsg mocks base method
func (m *MockGame_JoinATeamServer) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg
func (mr *MockGame_JoinATeamServerMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockGame_JoinATeamServer)(nil).RecvMsg), arg0)
}

// Send mocks base method
func (m *MockGame_JoinATeamServer) Send(arg0 *lugo.GameSnapshot) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send
func (mr *MockGame_JoinATeamServerMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockGame_JoinATeamServer)(nil).Send), arg0)
}

// SendHeader mocks base method
func (m *MockGame_JoinATeamServer) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader
func (mr *MockGame_JoinATeamServerMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockGame_JoinATeamServer)(nil).SendHeader), arg0)
}

// SendMsg mocks base method
func (m *MockGame_JoinATeamServer) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg
func (mr *MockGame_JoinATeamServerMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockGame_JoinATeamServer)(nil).SendMsg), arg0)
}

// SetHeader mocks base method
func (m *MockGame_JoinATeamServer) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader
func (mr *MockGame_JoinATeamServerMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockGame_JoinATeamServer)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method
func (m *MockGame_JoinATeamServer) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer
func (mr *MockGame_JoinATeamServerMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockGame_JoinATeamServer)(nil).SetTrailer), arg0)
}

// MockBroadcastClient is a mock of BroadcastClient interface
type MockBroadcastClient struct {
	ctrl     *gomock.Controller
	recorder *MockBroadcastClientMockRecorder
}

// MockBroadcastClientMockRecorder is the mock recorder for MockBroadcastClient
type MockBroadcastClientMockRecorder struct {
	mock *MockBroadcastClient
}

// NewMockBroadcastClient creates a new mock instance
func NewMockBroadcastClient(ctrl *gomock.Controller) *MockBroadcastClient {
	mock := &MockBroadcastClient{ctrl: ctrl}
	mock.recorder = &MockBroadcastClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBroadcastClient) EXPECT() *MockBroadcastClientMockRecorder {
	return m.recorder
}

// OnEvent mocks base method
func (m *MockBroadcastClient) OnEvent(arg0 context.Context, arg1 *empty.Empty, arg2 ...grpc.CallOption) (lugo.Broadcast_OnEventClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "OnEvent", varargs...)
	ret0, _ := ret[0].(lugo.Broadcast_OnEventClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OnEvent indicates an expected call of OnEvent
func (mr *MockBroadcastClientMockRecorder) OnEvent(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnEvent", reflect.TypeOf((*MockBroadcastClient)(nil).OnEvent), varargs...)
}

// MockBroadcast_OnEventClient is a mock of Broadcast_OnEventClient interface
type MockBroadcast_OnEventClient struct {
	ctrl     *gomock.Controller
	recorder *MockBroadcast_OnEventClientMockRecorder
}

// MockBroadcast_OnEventClientMockRecorder is the mock recorder for MockBroadcast_OnEventClient
type MockBroadcast_OnEventClientMockRecorder struct {
	mock *MockBroadcast_OnEventClient
}

// NewMockBroadcast_OnEventClient creates a new mock instance
func NewMockBroadcast_OnEventClient(ctrl *gomock.Controller) *MockBroadcast_OnEventClient {
	mock := &MockBroadcast_OnEventClient{ctrl: ctrl}
	mock.recorder = &MockBroadcast_OnEventClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBroadcast_OnEventClient) EXPECT() *MockBroadcast_OnEventClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method
func (m *MockBroadcast_OnEventClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend
func (mr *MockBroadcast_OnEventClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockBroadcast_OnEventClient)(nil).CloseSend))
}

// Context mocks base method
func (m *MockBroadcast_OnEventClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context
func (mr *MockBroadcast_OnEventClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockBroadcast_OnEventClient)(nil).Context))
}

// Header mocks base method
func (m *MockBroadcast_OnEventClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header
func (mr *MockBroadcast_OnEventClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockBroadcast_OnEventClient)(nil).Header))
}

// Recv mocks base method
func (m *MockBroadcast_OnEventClient) Recv() (*lugo.GameEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*lugo.GameEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv
func (mr *MockBroadcast_OnEventClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockBroadcast_OnEventClient)(nil).Recv))
}

// RecvMsg mocks base method
func (m *MockBroadcast_OnEventClient) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg
func (mr *MockBroadcast_OnEventClientMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockBroadcast_OnEventClient)(nil).RecvMsg), arg0)
}

// SendMsg mocks base method
func (m *MockBroadcast_OnEventClient) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg
func (mr *MockBroadcast_OnEventClientMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockBroadcast_OnEventClient)(nil).SendMsg), arg0)
}

// Trailer mocks base method
func (m *MockBroadcast_OnEventClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer
func (mr *MockBroadcast_OnEventClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockBroadcast_OnEventClient)(nil).Trailer))
}

// MockBroadcastServer is a mock of BroadcastServer interface
type MockBroadcastServer struct {
	ctrl     *gomock.Controller
	recorder *MockBroadcastServerMockRecorder
}

// MockBroadcastServerMockRecorder is the mock recorder for MockBroadcastServer
type MockBroadcastServerMockRecorder struct {
	mock *MockBroadcastServer
}

// NewMockBroadcastServer creates a new mock instance
func NewMockBroadcastServer(ctrl *gomock.Controller) *MockBroadcastServer {
	mock := &MockBroadcastServer{ctrl: ctrl}
	mock.recorder = &MockBroadcastServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBroadcastServer) EXPECT() *MockBroadcastServerMockRecorder {
	return m.recorder
}

// OnEvent mocks base method
func (m *MockBroadcastServer) OnEvent(arg0 *empty.Empty, arg1 lugo.Broadcast_OnEventServer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OnEvent", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// OnEvent indicates an expected call of OnEvent
func (mr *MockBroadcastServerMockRecorder) OnEvent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnEvent", reflect.TypeOf((*MockBroadcastServer)(nil).OnEvent), arg0, arg1)
}

// MockBroadcast_OnEventServer is a mock of Broadcast_OnEventServer interface
type MockBroadcast_OnEventServer struct {
	ctrl     *gomock.Controller
	recorder *MockBroadcast_OnEventServerMockRecorder
}

// MockBroadcast_OnEventServerMockRecorder is the mock recorder for MockBroadcast_OnEventServer
type MockBroadcast_OnEventServerMockRecorder struct {
	mock *MockBroadcast_OnEventServer
}

// NewMockBroadcast_OnEventServer creates a new mock instance
func NewMockBroadcast_OnEventServer(ctrl *gomock.Controller) *MockBroadcast_OnEventServer {
	mock := &MockBroadcast_OnEventServer{ctrl: ctrl}
	mock.recorder = &MockBroadcast_OnEventServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBroadcast_OnEventServer) EXPECT() *MockBroadcast_OnEventServerMockRecorder {
	return m.recorder
}

// Context mocks base method
func (m *MockBroadcast_OnEventServer) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context
func (mr *MockBroadcast_OnEventServerMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockBroadcast_OnEventServer)(nil).Context))
}

// RecvMsg mocks base method
func (m *MockBroadcast_OnEventServer) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg
func (mr *MockBroadcast_OnEventServerMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockBroadcast_OnEventServer)(nil).RecvMsg), arg0)
}

// Send mocks base method
func (m *MockBroadcast_OnEventServer) Send(arg0 *lugo.GameEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send
func (mr *MockBroadcast_OnEventServerMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockBroadcast_OnEventServer)(nil).Send), arg0)
}

// SendHeader mocks base method
func (m *MockBroadcast_OnEventServer) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader
func (mr *MockBroadcast_OnEventServerMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockBroadcast_OnEventServer)(nil).SendHeader), arg0)
}

// SendMsg mocks base method
func (m *MockBroadcast_OnEventServer) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg
func (mr *MockBroadcast_OnEventServerMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockBroadcast_OnEventServer)(nil).SendMsg), arg0)
}

// SetHeader mocks base method
func (m *MockBroadcast_OnEventServer) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader
func (mr *MockBroadcast_OnEventServerMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockBroadcast_OnEventServer)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method
func (m *MockBroadcast_OnEventServer) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer
func (mr *MockBroadcast_OnEventServerMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockBroadcast_OnEventServer)(nil).SetTrailer), arg0)
}
