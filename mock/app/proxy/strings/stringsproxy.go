// Code generated by MockGen. DO NOT EDIT.
// Source: ../app/proxy/strings/stringsproxy.go
//
// Generated by this command:
//
//	mockgen -source=../app/proxy/strings/stringsproxy.go -destination=../mock/app/proxy/strings/stringsproxy.go -package=mockstringsproxy
//

// Package mockstringsproxy is a generated GoMock package.
package mockstringsproxy

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockStrings is a mock of Strings interface.
type MockStrings struct {
	ctrl     *gomock.Controller
	recorder *MockStringsMockRecorder
}

// MockStringsMockRecorder is the mock recorder for MockStrings.
type MockStringsMockRecorder struct {
	mock *MockStrings
}

// NewMockStrings creates a new mock instance.
func NewMockStrings(ctrl *gomock.Controller) *MockStrings {
	mock := &MockStrings{ctrl: ctrl}
	mock.recorder = &MockStringsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStrings) EXPECT() *MockStringsMockRecorder {
	return m.recorder
}

// Join mocks base method.
func (m *MockStrings) Join(elems []string, sep string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Join", elems, sep)
	ret0, _ := ret[0].(string)
	return ret0
}

// Join indicates an expected call of Join.
func (mr *MockStringsMockRecorder) Join(elems, sep any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Join", reflect.TypeOf((*MockStrings)(nil).Join), elems, sep)
}