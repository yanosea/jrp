// Code generated by MockGen. DO NOT EDIT.
// Source: ../app/proxy/os/fileinstance.go
//
// Generated by this command:
//
//	mockgen -source=../app/proxy/os/fileinstance.go -destination=../mock/app/proxy/os/fileinstance.go -package=mockosproxy
//

// Package mockosproxy is a generated GoMock package.
package mockosproxy

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockFileInstanceInterface is a mock of FileInstanceInterface interface.
type MockFileInstanceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockFileInstanceInterfaceMockRecorder
}

// MockFileInstanceInterfaceMockRecorder is the mock recorder for MockFileInstanceInterface.
type MockFileInstanceInterfaceMockRecorder struct {
	mock *MockFileInstanceInterface
}

// NewMockFileInstanceInterface creates a new mock instance.
func NewMockFileInstanceInterface(ctrl *gomock.Controller) *MockFileInstanceInterface {
	mock := &MockFileInstanceInterface{ctrl: ctrl}
	mock.recorder = &MockFileInstanceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileInstanceInterface) EXPECT() *MockFileInstanceInterfaceMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockFileInstanceInterface) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockFileInstanceInterfaceMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockFileInstanceInterface)(nil).Close))
}

// Read mocks base method.
func (m *MockFileInstanceInterface) Read(p []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", p)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockFileInstanceInterfaceMockRecorder) Read(p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockFileInstanceInterface)(nil).Read), p)
}

// Seek mocks base method.
func (m *MockFileInstanceInterface) Seek(offset int64, whence int) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Seek", offset, whence)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Seek indicates an expected call of Seek.
func (mr *MockFileInstanceInterfaceMockRecorder) Seek(offset, whence any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Seek", reflect.TypeOf((*MockFileInstanceInterface)(nil).Seek), offset, whence)
}

// Write mocks base method.
func (m *MockFileInstanceInterface) Write(b []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", b)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write.
func (mr *MockFileInstanceInterfaceMockRecorder) Write(b any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockFileInstanceInterface)(nil).Write), b)
}