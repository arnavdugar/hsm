// Code generated by MockGen. DO NOT EDIT.
// Source: machine.go

// Package queue is a generated GoMock package.
package queue

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHandler is a mock of Handler interface.
type MockHandler struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerMockRecorder
}

// MockHandlerMockRecorder is the mock recorder for MockHandler.
type MockHandlerMockRecorder struct {
	mock *MockHandler
}

// NewMockHandler creates a new mock instance.
func NewMockHandler(ctrl *gomock.Controller) *MockHandler {
	mock := &MockHandler{ctrl: ctrl}
	mock.recorder = &MockHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHandler) EXPECT() *MockHandlerMockRecorder {
	return m.recorder
}

// HandleClose mocks base method.
func (m *MockHandler) HandleClose() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleClose")
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleClose indicates an expected call of HandleClose.
func (mr *MockHandlerMockRecorder) HandleClose() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleClose", reflect.TypeOf((*MockHandler)(nil).HandleClose))
}

// HandleConsumeElement mocks base method.
func (m *MockHandler) HandleConsumeElement() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleConsumeElement")
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleConsumeElement indicates an expected call of HandleConsumeElement.
func (mr *MockHandlerMockRecorder) HandleConsumeElement() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleConsumeElement", reflect.TypeOf((*MockHandler)(nil).HandleConsumeElement))
}

// HandlePushElement mocks base method.
func (m *MockHandler) HandlePushElement(data QueueElement) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandlePushElement", data)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandlePushElement indicates an expected call of HandlePushElement.
func (mr *MockHandlerMockRecorder) HandlePushElement(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandlePushElement", reflect.TypeOf((*MockHandler)(nil).HandlePushElement), data)
}

// HasSingleElement mocks base method.
func (m *MockHandler) HasSingleElement() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasSingleElement")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HasSingleElement indicates an expected call of HasSingleElement.
func (mr *MockHandlerMockRecorder) HasSingleElement() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasSingleElement", reflect.TypeOf((*MockHandler)(nil).HasSingleElement))
}
