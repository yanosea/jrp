// Code generated by MockGen. DO NOT EDIT.
// Source: ../app/library/jrpwriter/jrpwriter.go
//
// Generated by this command:
//
//	mockgen -source=../app/library/jrpwriter/jrpwriter.go -destination=../mock/app/library/jrpwriter/jrpwriter.go -package=mockjrpwriter
//

// Package mockjrpwriter is a generated GoMock package.
package mockjrpwriter

import (
	reflect "reflect"

	model "github.com/yanosea/jrp/app/database/jrp/model"
	ioproxy "github.com/yanosea/jrp/app/proxy/io"
	gomock "go.uber.org/mock/gomock"
)

// MockJrpWritable is a mock of JrpWritable interface.
type MockJrpWritable struct {
	ctrl     *gomock.Controller
	recorder *MockJrpWritableMockRecorder
}

// MockJrpWritableMockRecorder is the mock recorder for MockJrpWritable.
type MockJrpWritableMockRecorder struct {
	mock *MockJrpWritable
}

// NewMockJrpWritable creates a new mock instance.
func NewMockJrpWritable(ctrl *gomock.Controller) *MockJrpWritable {
	mock := &MockJrpWritable{ctrl: ctrl}
	mock.recorder = &MockJrpWritableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJrpWritable) EXPECT() *MockJrpWritableMockRecorder {
	return m.recorder
}

// WriteAsTable mocks base method.
func (m *MockJrpWritable) WriteAsTable(writer ioproxy.WriterInstanceInterface, jrps []model.Jrp) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteAsTable", writer, jrps)
}

// WriteAsTable indicates an expected call of WriteAsTable.
func (mr *MockJrpWritableMockRecorder) WriteAsTable(writer, jrps any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteAsTable", reflect.TypeOf((*MockJrpWritable)(nil).WriteAsTable), writer, jrps)
}

// WriteGenerateResultAsTable mocks base method.
func (m *MockJrpWritable) WriteGenerateResultAsTable(writer ioproxy.WriterInstanceInterface, jrps []model.Jrp) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteGenerateResultAsTable", writer, jrps)
}

// WriteGenerateResultAsTable indicates an expected call of WriteGenerateResultAsTable.
func (mr *MockJrpWritableMockRecorder) WriteGenerateResultAsTable(writer, jrps any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteGenerateResultAsTable", reflect.TypeOf((*MockJrpWritable)(nil).WriteGenerateResultAsTable), writer, jrps)
}