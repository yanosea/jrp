// Code generated by MockGen. DO NOT EDIT.
// Source: ../app/proxy/color/colorproxy.go
//
// Generated by this command:
//
//	mockgen -source=../app/proxy/color/colorproxy.go -destination=../mock/app/proxy/color/colorproxy.go -package=mockcolorproxy
//

// Package mockcolorproxy is a generated GoMock package.
package mockcolorproxy

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockColor is a mock of Color interface.
type MockColor struct {
	ctrl     *gomock.Controller
	recorder *MockColorMockRecorder
}

// MockColorMockRecorder is the mock recorder for MockColor.
type MockColorMockRecorder struct {
	mock *MockColor
}

// NewMockColor creates a new mock instance.
func NewMockColor(ctrl *gomock.Controller) *MockColor {
	mock := &MockColor{ctrl: ctrl}
	mock.recorder = &MockColorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockColor) EXPECT() *MockColorMockRecorder {
	return m.recorder
}

// BlueString mocks base method.
func (m *MockColor) BlueString(format string, a ...any) string {
	m.ctrl.T.Helper()
	varargs := []any{format}
	for _, a_2 := range a {
		varargs = append(varargs, a_2)
	}
	ret := m.ctrl.Call(m, "BlueString", varargs...)
	ret0, _ := ret[0].(string)
	return ret0
}

// BlueString indicates an expected call of BlueString.
func (mr *MockColorMockRecorder) BlueString(format any, a ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{format}, a...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlueString", reflect.TypeOf((*MockColor)(nil).BlueString), varargs...)
}

// GreenString mocks base method.
func (m *MockColor) GreenString(format string, a ...any) string {
	m.ctrl.T.Helper()
	varargs := []any{format}
	for _, a_2 := range a {
		varargs = append(varargs, a_2)
	}
	ret := m.ctrl.Call(m, "GreenString", varargs...)
	ret0, _ := ret[0].(string)
	return ret0
}

// GreenString indicates an expected call of GreenString.
func (mr *MockColorMockRecorder) GreenString(format any, a ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{format}, a...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GreenString", reflect.TypeOf((*MockColor)(nil).GreenString), varargs...)
}

// RedString mocks base method.
func (m *MockColor) RedString(format string, a ...any) string {
	m.ctrl.T.Helper()
	varargs := []any{format}
	for _, a_2 := range a {
		varargs = append(varargs, a_2)
	}
	ret := m.ctrl.Call(m, "RedString", varargs...)
	ret0, _ := ret[0].(string)
	return ret0
}

// RedString indicates an expected call of RedString.
func (mr *MockColorMockRecorder) RedString(format any, a ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{format}, a...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RedString", reflect.TypeOf((*MockColor)(nil).RedString), varargs...)
}

// YellowString mocks base method.
func (m *MockColor) YellowString(format string, a ...any) string {
	m.ctrl.T.Helper()
	varargs := []any{format}
	for _, a_2 := range a {
		varargs = append(varargs, a_2)
	}
	ret := m.ctrl.Call(m, "YellowString", varargs...)
	ret0, _ := ret[0].(string)
	return ret0
}

// YellowString indicates an expected call of YellowString.
func (mr *MockColorMockRecorder) YellowString(format any, a ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{format}, a...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "YellowString", reflect.TypeOf((*MockColor)(nil).YellowString), varargs...)
}
