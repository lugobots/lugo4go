// Code generated by MockGen. DO NOT EDIT.
// Source: team/interfaces.go

// Package team_test is a generated GoMock package.
package team_test

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	lugo "github.com/lugobots/lugo4go/v2/lugo"
	team "github.com/lugobots/lugo4go/v2/team"
	reflect "reflect"
)

// MockOrderSender is a mock of OrderSender interface.
type MockOrderSender struct {
	ctrl     *gomock.Controller
	recorder *MockOrderSenderMockRecorder
}

// MockOrderSenderMockRecorder is the mock recorder for MockOrderSender.
type MockOrderSenderMockRecorder struct {
	mock *MockOrderSender
}

// NewMockOrderSender creates a new mock instance.
func NewMockOrderSender(ctrl *gomock.Controller) *MockOrderSender {
	mock := &MockOrderSender{ctrl: ctrl}
	mock.recorder = &MockOrderSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderSender) EXPECT() *MockOrderSenderMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockOrderSender) Send(ctx context.Context, turn uint32, orders []lugo.PlayerOrder, debugMsg string) (*lugo.OrderResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", ctx, turn, orders, debugMsg)
	ret0, _ := ret[0].(*lugo.OrderResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Send indicates an expected call of Send.
func (mr *MockOrderSenderMockRecorder) Send(ctx, turn, orders, debugMsg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockOrderSender)(nil).Send), ctx, turn, orders, debugMsg)
}

// MockBot is a mock of Bot interface.
type MockBot struct {
	ctrl     *gomock.Controller
	recorder *MockBotMockRecorder
}

// MockBotMockRecorder is the mock recorder for MockBot.
type MockBotMockRecorder struct {
	mock *MockBot
}

// NewMockBot creates a new mock instance.
func NewMockBot(ctrl *gomock.Controller) *MockBot {
	mock := &MockBot{ctrl: ctrl}
	mock.recorder = &MockBotMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBot) EXPECT() *MockBotMockRecorder {
	return m.recorder
}

// OnDisputing mocks base method.
func (m *MockBot) OnDisputing(ctx context.Context, sender team.OrderSender, snapshot *lugo.GameSnapshot) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OnDisputing", ctx, sender, snapshot)
	ret0, _ := ret[0].(error)
	return ret0
}

// OnDisputing indicates an expected call of OnDisputing.
func (mr *MockBotMockRecorder) OnDisputing(ctx, sender, snapshot interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnDisputing", reflect.TypeOf((*MockBot)(nil).OnDisputing), ctx, sender, snapshot)
}

// OnDefending mocks base method.
func (m *MockBot) OnDefending(ctx context.Context, sender team.OrderSender, snapshot *lugo.GameSnapshot) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OnDefending", ctx, sender, snapshot)
	ret0, _ := ret[0].(error)
	return ret0
}

// OnDefending indicates an expected call of OnDefending.
func (mr *MockBotMockRecorder) OnDefending(ctx, sender, snapshot interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnDefending", reflect.TypeOf((*MockBot)(nil).OnDefending), ctx, sender, snapshot)
}

// OnHolding mocks base method.
func (m *MockBot) OnHolding(ctx context.Context, sender team.OrderSender, snapshot *lugo.GameSnapshot) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OnHolding", ctx, sender, snapshot)
	ret0, _ := ret[0].(error)
	return ret0
}

// OnHolding indicates an expected call of OnHolding.
func (mr *MockBotMockRecorder) OnHolding(ctx, sender, snapshot interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnHolding", reflect.TypeOf((*MockBot)(nil).OnHolding), ctx, sender, snapshot)
}

// OnSupporting mocks base method.
func (m *MockBot) OnSupporting(ctx context.Context, sender team.OrderSender, snapshot *lugo.GameSnapshot) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OnSupporting", ctx, sender, snapshot)
	ret0, _ := ret[0].(error)
	return ret0
}

// OnSupporting indicates an expected call of OnSupporting.
func (mr *MockBotMockRecorder) OnSupporting(ctx, sender, snapshot interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnSupporting", reflect.TypeOf((*MockBot)(nil).OnSupporting), ctx, sender, snapshot)
}

// AsGoalkeeper mocks base method.
func (m *MockBot) AsGoalkeeper(ctx context.Context, sender team.OrderSender, snapshot *lugo.GameSnapshot, state team.PlayerState) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AsGoalkeeper", ctx, sender, snapshot, state)
	ret0, _ := ret[0].(error)
	return ret0
}

// AsGoalkeeper indicates an expected call of AsGoalkeeper.
func (mr *MockBotMockRecorder) AsGoalkeeper(ctx, sender, snapshot, state interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AsGoalkeeper", reflect.TypeOf((*MockBot)(nil).AsGoalkeeper), ctx, sender, snapshot, state)
}

// MockPositioner is a mock of Positioner interface.
type MockPositioner struct {
	ctrl     *gomock.Controller
	recorder *MockPositionerMockRecorder
}

// MockPositionerMockRecorder is the mock recorder for MockPositioner.
type MockPositionerMockRecorder struct {
	mock *MockPositioner
}

// NewMockPositioner creates a new mock instance.
func NewMockPositioner(ctrl *gomock.Controller) *MockPositioner {
	mock := &MockPositioner{ctrl: ctrl}
	mock.recorder = &MockPositionerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPositioner) EXPECT() *MockPositionerMockRecorder {
	return m.recorder
}

// GetRegion mocks base method.
func (m *MockPositioner) GetRegion(col, row uint8) (team.FieldNav, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRegion", col, row)
	ret0, _ := ret[0].(team.FieldNav)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRegion indicates an expected call of GetRegion.
func (mr *MockPositionerMockRecorder) GetRegion(col, row interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRegion", reflect.TypeOf((*MockPositioner)(nil).GetRegion), col, row)
}

// GetPointRegion mocks base method.
func (m *MockPositioner) GetPointRegion(point lugo.Point) (team.FieldNav, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPointRegion", point)
	ret0, _ := ret[0].(team.FieldNav)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPointRegion indicates an expected call of GetPointRegion.
func (mr *MockPositionerMockRecorder) GetPointRegion(point interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPointRegion", reflect.TypeOf((*MockPositioner)(nil).GetPointRegion), point)
}

// MockFieldNav is a mock of FieldNav interface.
type MockFieldNav struct {
	ctrl     *gomock.Controller
	recorder *MockFieldNavMockRecorder
}

// MockFieldNavMockRecorder is the mock recorder for MockFieldNav.
type MockFieldNavMockRecorder struct {
	mock *MockFieldNav
}

// NewMockFieldNav creates a new mock instance.
func NewMockFieldNav(ctrl *gomock.Controller) *MockFieldNav {
	mock := &MockFieldNav{ctrl: ctrl}
	mock.recorder = &MockFieldNavMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFieldNav) EXPECT() *MockFieldNavMockRecorder {
	return m.recorder
}

// String mocks base method.
func (m *MockFieldNav) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String.
func (mr *MockFieldNavMockRecorder) String() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockFieldNav)(nil).String))
}

// Col mocks base method.
func (m *MockFieldNav) Col() uint8 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Col")
	ret0, _ := ret[0].(uint8)
	return ret0
}

// Col indicates an expected call of Col.
func (mr *MockFieldNavMockRecorder) Col() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Col", reflect.TypeOf((*MockFieldNav)(nil).Col))
}

// Row mocks base method.
func (m *MockFieldNav) Row() uint8 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Row")
	ret0, _ := ret[0].(uint8)
	return ret0
}

// Row indicates an expected call of Row.
func (mr *MockFieldNavMockRecorder) Row() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Row", reflect.TypeOf((*MockFieldNav)(nil).Row))
}

// Center mocks base method.
func (m *MockFieldNav) Center() lugo.Point {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Center")
	ret0, _ := ret[0].(lugo.Point)
	return ret0
}

// Center indicates an expected call of Center.
func (mr *MockFieldNavMockRecorder) Center() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Center", reflect.TypeOf((*MockFieldNav)(nil).Center))
}

// Front mocks base method.
func (m *MockFieldNav) Front() team.FieldNav {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Front")
	ret0, _ := ret[0].(team.FieldNav)
	return ret0
}

// Front indicates an expected call of Front.
func (mr *MockFieldNavMockRecorder) Front() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Front", reflect.TypeOf((*MockFieldNav)(nil).Front))
}

// Back mocks base method.
func (m *MockFieldNav) Back() team.FieldNav {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Back")
	ret0, _ := ret[0].(team.FieldNav)
	return ret0
}

// Back indicates an expected call of Back.
func (mr *MockFieldNavMockRecorder) Back() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Back", reflect.TypeOf((*MockFieldNav)(nil).Back))
}

// Left mocks base method.
func (m *MockFieldNav) Left() team.FieldNav {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Left")
	ret0, _ := ret[0].(team.FieldNav)
	return ret0
}

// Left indicates an expected call of Left.
func (mr *MockFieldNavMockRecorder) Left() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Left", reflect.TypeOf((*MockFieldNav)(nil).Left))
}

// Right mocks base method.
func (m *MockFieldNav) Right() team.FieldNav {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Right")
	ret0, _ := ret[0].(team.FieldNav)
	return ret0
}

// Right indicates an expected call of Right.
func (mr *MockFieldNavMockRecorder) Right() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Right", reflect.TypeOf((*MockFieldNav)(nil).Right))
}
