// Code generated by MockGen. DO NOT EDIT.
// Source: ../app/proxy/spinner/spinnerproxy.go
//
// Generated by this command:
//
//	mockgen -source=../app/proxy/spinner/spinnerproxy.go -destination=../mock/app/proxy/spinner/spinnerproxy.go -package=mockspinnerproxy
//

// Package mockspinnerproxy is a generated GoMock package.
package mockspinnerproxy

import (
	reflect "reflect"

	spinnerproxy "github.com/yanosea/jrp/app/proxy/spinner"
	gomock "go.uber.org/mock/gomock"
)

// MockSpinner is a mock of Spinner interface.
type MockSpinner struct {
	ctrl     *gomock.Controller
	recorder *MockSpinnerMockRecorder
}

// MockSpinnerMockRecorder is the mock recorder for MockSpinner.
type MockSpinnerMockRecorder struct {
	mock *MockSpinner
}

// NewMockSpinner creates a new mock instance.
func NewMockSpinner(ctrl *gomock.Controller) *MockSpinner {
	mock := &MockSpinner{ctrl: ctrl}
	mock.recorder = &MockSpinnerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSpinner) EXPECT() *MockSpinnerMockRecorder {
	return m.recorder
}

// NewSpinner mocks base method.
func (m *MockSpinner) NewSpinner() spinnerproxy.SpinnerInstanceInterface {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewSpinner")
	ret0, _ := ret[0].(spinnerproxy.SpinnerInstanceInterface)
	return ret0
}

// NewSpinner indicates an expected call of NewSpinner.
func (mr *MockSpinnerMockRecorder) NewSpinner() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewSpinner", reflect.TypeOf((*MockSpinner)(nil).NewSpinner))
}