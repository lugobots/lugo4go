// Code generated by MockGen. DO NOT EDIT.
// Source: ./contracts.go

// Package lugo4go_test is a generated GoMock package.
package lugo4go_test

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v2 "github.com/lugobots/lugo4go/v2"
	proto "github.com/lugobots/lugo4go/v2/proto"
)

// MockTurnHandler is a mock of TurnHandler interface.
type MockTurnHandler struct {
	ctrl     *gomock.Controller
	recorder *MockTurnHandlerMockRecorder
}

// MockTurnHandlerMockRecorder is the mock recorder for MockTurnHandler.
type MockTurnHandlerMockRecorder struct {
	mock *MockTurnHandler
}

// NewMockTurnHandler creates a new mock instance.
func NewMockTurnHandler(ctrl *gomock.Controller) *MockTurnHandler {
	mock := &MockTurnHandler{ctrl: ctrl}
	mock.recorder = &MockTurnHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTurnHandler) EXPECT() *MockTurnHandlerMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MockTurnHandler) Handle(ctx context.Context, snapshot *proto.GameSnapshot) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Handle", ctx, snapshot)
}

// Handle indicates an expected call of Handle.
func (mr *MockTurnHandlerMockRecorder) Handle(ctx, snapshot interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockTurnHandler)(nil).Handle), ctx, snapshot)
}

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
func (m *MockOrderSender) Send(ctx context.Context, turn uint32, orders []proto.PlayerOrder, debugMsg string) (*proto.OrderResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", ctx, turn, orders, debugMsg)
	ret0, _ := ret[0].(*proto.OrderResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Send indicates an expected call of Send.
func (mr *MockOrderSenderMockRecorder) Send(ctx, turn, orders, debugMsg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockOrderSender)(nil).Send), ctx, turn, orders, debugMsg)
}

// MockTurnOrdersSender is a mock of TurnOrdersSender interface.
type MockTurnOrdersSender struct {
	ctrl     *gomock.Controller
	recorder *MockTurnOrdersSenderMockRecorder
}

// MockTurnOrdersSenderMockRecorder is the mock recorder for MockTurnOrdersSender.
type MockTurnOrdersSenderMockRecorder struct {
	mock *MockTurnOrdersSender
}

// NewMockTurnOrdersSender creates a new mock instance.
func NewMockTurnOrdersSender(ctrl *gomock.Controller) *MockTurnOrdersSender {
	mock := &MockTurnOrdersSender{ctrl: ctrl}
	mock.recorder = &MockTurnOrdersSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTurnOrdersSender) EXPECT() *MockTurnOrdersSenderMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockTurnOrdersSender) Send(ctx context.Context, orders []proto.PlayerOrder, debugMsg string) (*proto.OrderResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", ctx, orders, debugMsg)
	ret0, _ := ret[0].(*proto.OrderResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Send indicates an expected call of Send.
func (mr *MockTurnOrdersSenderMockRecorder) Send(ctx, orders, debugMsg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockTurnOrdersSender)(nil).Send), ctx, orders, debugMsg)
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

// AsGoalkeeper mocks base method.
func (m *MockBot) AsGoalkeeper(ctx context.Context, sender v2.TurnOrdersSender, snapshot *proto.GameSnapshot, state v2.PlayerState) error {
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

// OnDefending mocks base method.
func (m *MockBot) OnDefending(ctx context.Context, sender v2.TurnOrdersSender, snapshot *proto.GameSnapshot) error {
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

// OnDisputing mocks base method.
func (m *MockBot) OnDisputing(ctx context.Context, sender v2.TurnOrdersSender, snapshot *proto.GameSnapshot) error {
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

// OnHolding mocks base method.
func (m *MockBot) OnHolding(ctx context.Context, sender v2.TurnOrdersSender, snapshot *proto.GameSnapshot) error {
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
func (m *MockBot) OnSupporting(ctx context.Context, sender v2.TurnOrdersSender, snapshot *proto.GameSnapshot) error {
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

// MockLogger is a mock of Logger interface.
type MockLogger struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerMockRecorder
}

// MockLoggerMockRecorder is the mock recorder for MockLogger.
type MockLoggerMockRecorder struct {
	mock *MockLogger
}

// NewMockLogger creates a new mock instance.
func NewMockLogger(ctrl *gomock.Controller) *MockLogger {
	mock := &MockLogger{ctrl: ctrl}
	mock.recorder = &MockLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogger) EXPECT() *MockLoggerMockRecorder {
	return m.recorder
}

// Debug mocks base method.
func (m *MockLogger) Debug(args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debug", varargs...)
}

// Debug indicates an expected call of Debug.
func (mr *MockLoggerMockRecorder) Debug(args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debug", reflect.TypeOf((*MockLogger)(nil).Debug), args...)
}

// Debugf mocks base method.
func (m *MockLogger) Debugf(template string, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{template}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debugf", varargs...)
}

// Debugf indicates an expected call of Debugf.
func (mr *MockLoggerMockRecorder) Debugf(template interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{template}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debugf", reflect.TypeOf((*MockLogger)(nil).Debugf), varargs...)
}

// Errorf mocks base method.
func (m *MockLogger) Errorf(template string, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{template}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Errorf", varargs...)
}

// Errorf indicates an expected call of Errorf.
func (mr *MockLoggerMockRecorder) Errorf(template interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{template}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Errorf", reflect.TypeOf((*MockLogger)(nil).Errorf), varargs...)
}

// Fatalf mocks base method.
func (m *MockLogger) Fatalf(template string, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{template}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Fatalf", varargs...)
}

// Fatalf indicates an expected call of Fatalf.
func (mr *MockLoggerMockRecorder) Fatalf(template interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{template}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fatalf", reflect.TypeOf((*MockLogger)(nil).Fatalf), varargs...)
}

// Infof mocks base method.
func (m *MockLogger) Infof(template string, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{template}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Infof", varargs...)
}

// Infof indicates an expected call of Infof.
func (mr *MockLoggerMockRecorder) Infof(template interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{template}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Infof", reflect.TypeOf((*MockLogger)(nil).Infof), varargs...)
}

// Warnf mocks base method.
func (m *MockLogger) Warnf(template string, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{template}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Warnf", varargs...)
}

// Warnf indicates an expected call of Warnf.
func (mr *MockLoggerMockRecorder) Warnf(template interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{template}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warnf", reflect.TypeOf((*MockLogger)(nil).Warnf), varargs...)
}
