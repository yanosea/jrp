// Code generated by MockGen. DO NOT EDIT.
// Source: ../app/proxy/tablewriter/tableinstance.go
//
// Generated by this command:
//
//	mockgen -source=../app/proxy/tablewriter/tableinstance.go -destination=../mock/app/proxy/tablewriter/tableinstance.go -package=mocktablewriterproxy
//

// Package mocktablewriterproxy is a generated GoMock package.
package mocktablewriterproxy

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockTableInstanceInterface is a mock of TableInstanceInterface interface.
type MockTableInstanceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockTableInstanceInterfaceMockRecorder
}

// MockTableInstanceInterfaceMockRecorder is the mock recorder for MockTableInstanceInterface.
type MockTableInstanceInterfaceMockRecorder struct {
	mock *MockTableInstanceInterface
}

// NewMockTableInstanceInterface creates a new mock instance.
func NewMockTableInstanceInterface(ctrl *gomock.Controller) *MockTableInstanceInterface {
	mock := &MockTableInstanceInterface{ctrl: ctrl}
	mock.recorder = &MockTableInstanceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTableInstanceInterface) EXPECT() *MockTableInstanceInterfaceMockRecorder {
	return m.recorder
}

// AppendBulk mocks base method.
func (m *MockTableInstanceInterface) AppendBulk(rows [][]string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AppendBulk", rows)
}

// AppendBulk indicates an expected call of AppendBulk.
func (mr *MockTableInstanceInterfaceMockRecorder) AppendBulk(rows any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AppendBulk", reflect.TypeOf((*MockTableInstanceInterface)(nil).AppendBulk), rows)
}

// Render mocks base method.
func (m *MockTableInstanceInterface) Render() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Render")
}

// Render indicates an expected call of Render.
func (mr *MockTableInstanceInterfaceMockRecorder) Render() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Render", reflect.TypeOf((*MockTableInstanceInterface)(nil).Render))
}

// SetAlignment mocks base method.
func (m *MockTableInstanceInterface) SetAlignment(align int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetAlignment", align)
}

// SetAlignment indicates an expected call of SetAlignment.
func (mr *MockTableInstanceInterfaceMockRecorder) SetAlignment(align any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAlignment", reflect.TypeOf((*MockTableInstanceInterface)(nil).SetAlignment), align)
}

// SetAutoFormatHeaders mocks base method.
func (m *MockTableInstanceInterface) SetAutoFormatHeaders(auto bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetAutoFormatHeaders", auto)
}

// SetAutoFormatHeaders indicates an expected call of SetAutoFormatHeaders.
func (mr *MockTableInstanceInterfaceMockRecorder) SetAutoFormatHeaders(auto any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAutoFormatHeaders", reflect.TypeOf((*MockTableInstanceInterface)(nil).SetAutoFormatHeaders), auto)
}

// SetAutoWrapText mocks base method.
func (m *MockTableInstanceInterface) SetAutoWrapText(auto bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetAutoWrapText", auto)
}

// SetAutoWrapText indicates an expected call of SetAutoWrapText.
func (mr *MockTableInstanceInterfaceMockRecorder) SetAutoWrapText(auto any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAutoWrapText", reflect.TypeOf((*MockTableInstanceInterface)(nil).SetAutoWrapText), auto)
}

// SetBorder mocks base method.
func (m *MockTableInstanceInterface) SetBorder(border bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetBorder", border)
}

// SetBorder indicates an expected call of SetBorder.
func (mr *MockTableInstanceInterfaceMockRecorder) SetBorder(border any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBorder", reflect.TypeOf((*MockTableInstanceInterface)(nil).SetBorder), border)
}

// SetCenterSeparator mocks base method.
func (m *MockTableInstanceInterface) SetCenterSeparator(sep string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetCenterSeparator", sep)
}

// SetCenterSeparator indicates an expected call of SetCenterSeparator.
func (mr *MockTableInstanceInterfaceMockRecorder) SetCenterSeparator(sep any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCenterSeparator", reflect.TypeOf((*MockTableInstanceInterface)(nil).SetCenterSeparator), sep)
}

// SetColumnSeparator mocks base method.
func (m *MockTableInstanceInterface) SetColumnSeparator(sep string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetColumnSeparator", sep)
}

// SetColumnSeparator indicates an expected call of SetColumnSeparator.
func (mr *MockTableInstanceInterfaceMockRecorder) SetColumnSeparator(sep any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetColumnSeparator", reflect.TypeOf((*MockTableInstanceInterface)(nil).SetColumnSeparator), sep)
}

// SetHeader mocks base method.
func (m *MockTableInstanceInterface) SetHeader(keys []string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetHeader", keys)
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockTableInstanceInterfaceMockRecorder) SetHeader(keys any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockTableInstanceInterface)(nil).SetHeader), keys)
}

// SetHeaderAlignment mocks base method.
func (m *MockTableInstanceInterface) SetHeaderAlignment(hAlign int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetHeaderAlignment", hAlign)
}

// SetHeaderAlignment indicates an expected call of SetHeaderAlignment.
func (mr *MockTableInstanceInterfaceMockRecorder) SetHeaderAlignment(hAlign any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeaderAlignment", reflect.TypeOf((*MockTableInstanceInterface)(nil).SetHeaderAlignment), hAlign)
}

// SetHeaderLine mocks base method.
func (m *MockTableInstanceInterface) SetHeaderLine(line bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetHeaderLine", line)
}

// SetHeaderLine indicates an expected call of SetHeaderLine.
func (mr *MockTableInstanceInterfaceMockRecorder) SetHeaderLine(line any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeaderLine", reflect.TypeOf((*MockTableInstanceInterface)(nil).SetHeaderLine), line)
}

// SetNoWhiteSpace mocks base method.
func (m *MockTableInstanceInterface) SetNoWhiteSpace(allow bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetNoWhiteSpace", allow)
}

// SetNoWhiteSpace indicates an expected call of SetNoWhiteSpace.
func (mr *MockTableInstanceInterfaceMockRecorder) SetNoWhiteSpace(allow any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetNoWhiteSpace", reflect.TypeOf((*MockTableInstanceInterface)(nil).SetNoWhiteSpace), allow)
}

// SetRowSeparator mocks base method.
func (m *MockTableInstanceInterface) SetRowSeparator(sep string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetRowSeparator", sep)
}

// SetRowSeparator indicates an expected call of SetRowSeparator.
func (mr *MockTableInstanceInterfaceMockRecorder) SetRowSeparator(sep any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRowSeparator", reflect.TypeOf((*MockTableInstanceInterface)(nil).SetRowSeparator), sep)
}

// SetTablePadding mocks base method.
func (m *MockTableInstanceInterface) SetTablePadding(padding string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTablePadding", padding)
}

// SetTablePadding indicates an expected call of SetTablePadding.
func (mr *MockTableInstanceInterfaceMockRecorder) SetTablePadding(padding any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTablePadding", reflect.TypeOf((*MockTableInstanceInterface)(nil).SetTablePadding), padding)
}