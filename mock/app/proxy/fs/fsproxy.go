// Code generated by MockGen. DO NOT EDIT.
// Source: ../app/proxy/fs/fsproxy.go
//
// Generated by this command:
//
//	mockgen -source=../app/proxy/fs/fsproxy.go -destination=../mock/app/proxy/fs/fsproxy.go -package=mockfsproxy
//

// Package mockfsproxy is a generated GoMock package.
package mockfsproxy

import (
	gomock "go.uber.org/mock/gomock"
)

// MockFs is a mock of Fs interface.
type MockFs struct {
	ctrl     *gomock.Controller
	recorder *MockFsMockRecorder
}

// MockFsMockRecorder is the mock recorder for MockFs.
type MockFsMockRecorder struct {
	mock *MockFs
}

// NewMockFs creates a new mock instance.
func NewMockFs(ctrl *gomock.Controller) *MockFs {
	mock := &MockFs{ctrl: ctrl}
	mock.recorder = &MockFsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFs) EXPECT() *MockFsMockRecorder {
	return m.recorder
}
