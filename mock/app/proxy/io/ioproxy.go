// Code generated by MockGen. DO NOT EDIT.
// Source: ../app/proxy/io/ioproxy.go
//
// Generated by this command:
//
//	mockgen -source=../app/proxy/io/ioproxy.go -destination=../mock/app/proxy/io/ioproxy.go -package=mockioproxy
//

// Package mockioproxy is a generated GoMock package.
package mockioproxy

import (
	reflect "reflect"

	ioproxy "github.com/yanosea/jrp/app/proxy/io"
	gomock "go.uber.org/mock/gomock"
)

// MockIo is a mock of Io interface.
type MockIo struct {
	ctrl     *gomock.Controller
	recorder *MockIoMockRecorder
}

// MockIoMockRecorder is the mock recorder for MockIo.
type MockIoMockRecorder struct {
	mock *MockIo
}

// NewMockIo creates a new mock instance.
func NewMockIo(ctrl *gomock.Controller) *MockIo {
	mock := &MockIo{ctrl: ctrl}
	mock.recorder = &MockIoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIo) EXPECT() *MockIoMockRecorder {
	return m.recorder
}

// Copy mocks base method.
func (m *MockIo) Copy(dst ioproxy.WriterInstanceInterface, src ioproxy.ReaderInstanceInterface) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Copy", dst, src)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Copy indicates an expected call of Copy.
func (mr *MockIoMockRecorder) Copy(dst, src any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Copy", reflect.TypeOf((*MockIo)(nil).Copy), dst, src)
}
