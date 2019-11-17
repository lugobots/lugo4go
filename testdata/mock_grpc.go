// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/makeitplay/client-player-go (interfaces: Client)

// Package testdata is a generated GoMock package.
package testdata

import (
	gomock "github.com/golang/mock/gomock"
	client_player_go "github.com/makeitplay/client-player-go"
	lugo "github.com/makeitplay/client-player-go/lugo"
	reflect "reflect"
)

// MockClient is a mock of Client interface
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// OnNewTurn mocks base method
func (m *MockClient) OnNewTurn(arg0 func(*lugo.GameSnapshot) client_player_go.DecisionMaker) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnNewTurn", arg0)
}

// OnNewTurn indicates an expected call of OnNewTurn
func (mr *MockClientMockRecorder) OnNewTurn(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnNewTurn", reflect.TypeOf((*MockClient)(nil).OnNewTurn), arg0)
}

// Stop mocks base method
func (m *MockClient) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop
func (mr *MockClientMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockClient)(nil).Stop))
}
