// Code generated by MockGen. DO NOT EDIT.
// Source: ../app/proxy/sql/nullstringinstance.go
//
// Generated by this command:
//
//	mockgen -source=../app/proxy/sql/nullstringinstance.go -destination=../mock/app/proxy/sql/nullstringinstance.go -package=mocksqlproxy
//

// Package mocksqlproxy is a generated GoMock package.
package mocksqlproxy

import (
	driver "database/sql/driver"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockNullStringInstanceInterface is a mock of NullStringInstanceInterface interface.
type MockNullStringInstanceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockNullStringInstanceInterfaceMockRecorder
}

// MockNullStringInstanceInterfaceMockRecorder is the mock recorder for MockNullStringInstanceInterface.
type MockNullStringInstanceInterfaceMockRecorder struct {
	mock *MockNullStringInstanceInterface
}

// NewMockNullStringInstanceInterface creates a new mock instance.
func NewMockNullStringInstanceInterface(ctrl *gomock.Controller) *MockNullStringInstanceInterface {
	mock := &MockNullStringInstanceInterface{ctrl: ctrl}
	mock.recorder = &MockNullStringInstanceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNullStringInstanceInterface) EXPECT() *MockNullStringInstanceInterfaceMockRecorder {
	return m.recorder
}

// Scan mocks base method.
func (m *MockNullStringInstanceInterface) Scan(value any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Scan", value)
	ret0, _ := ret[0].(error)
	return ret0
}

// Scan indicates an expected call of Scan.
func (mr *MockNullStringInstanceInterfaceMockRecorder) Scan(value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scan", reflect.TypeOf((*MockNullStringInstanceInterface)(nil).Scan), value)
}

// Value mocks base method.
func (m *MockNullStringInstanceInterface) Value() (driver.Value, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Value")
	ret0, _ := ret[0].(driver.Value)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Value indicates an expected call of Value.
func (mr *MockNullStringInstanceInterfaceMockRecorder) Value() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Value", reflect.TypeOf((*MockNullStringInstanceInterface)(nil).Value))
}
