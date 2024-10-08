// Code generated by MockGen. DO NOT EDIT.
// Source: ../app/proxy/time/timeinstance.go
//
// Generated by this command:
//
//	mockgen -source=../app/proxy/time/timeinstance.go -destination=../mock/app/proxy/time/timeinstance.go -package=mocktimeproxy
//

// Package mocktimeproxy is a generated GoMock package.
package mocktimeproxy

import (
	driver "database/sql/driver"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockTimeInstanceInterface is a mock of TimeInstanceInterface interface.
type MockTimeInstanceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockTimeInstanceInterfaceMockRecorder
}

// MockTimeInstanceInterfaceMockRecorder is the mock recorder for MockTimeInstanceInterface.
type MockTimeInstanceInterfaceMockRecorder struct {
	mock *MockTimeInstanceInterface
}

// NewMockTimeInstanceInterface creates a new mock instance.
func NewMockTimeInstanceInterface(ctrl *gomock.Controller) *MockTimeInstanceInterface {
	mock := &MockTimeInstanceInterface{ctrl: ctrl}
	mock.recorder = &MockTimeInstanceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTimeInstanceInterface) EXPECT() *MockTimeInstanceInterfaceMockRecorder {
	return m.recorder
}

// Format mocks base method.
func (m *MockTimeInstanceInterface) Format(layout string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Format", layout)
	ret0, _ := ret[0].(string)
	return ret0
}

// Format indicates an expected call of Format.
func (mr *MockTimeInstanceInterfaceMockRecorder) Format(layout any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Format", reflect.TypeOf((*MockTimeInstanceInterface)(nil).Format), layout)
}

// Scan mocks base method.
func (m *MockTimeInstanceInterface) Scan(value any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Scan", value)
	ret0, _ := ret[0].(error)
	return ret0
}

// Scan indicates an expected call of Scan.
func (mr *MockTimeInstanceInterfaceMockRecorder) Scan(value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scan", reflect.TypeOf((*MockTimeInstanceInterface)(nil).Scan), value)
}

// Value mocks base method.
func (m *MockTimeInstanceInterface) Value() (driver.Value, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Value")
	ret0, _ := ret[0].(driver.Value)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Value indicates an expected call of Value.
func (mr *MockTimeInstanceInterfaceMockRecorder) Value() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Value", reflect.TypeOf((*MockTimeInstanceInterface)(nil).Value))
}
