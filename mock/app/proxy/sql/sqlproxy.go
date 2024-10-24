// Code generated by MockGen. DO NOT EDIT.
// Source: ../app/proxy/sql/sqlproxy.go
//
// Generated by this command:
//
//	mockgen -source=../app/proxy/sql/sqlproxy.go -destination=../mock/app/proxy/sql/sqlproxy.go -package=mocksqlproxy
//

// Package mocksqlproxy is a generated GoMock package.
package mocksqlproxy

import (
	reflect "reflect"

	sqlproxy "github.com/yanosea/jrp/app/proxy/sql"
	gomock "go.uber.org/mock/gomock"
)

// MockSql is a mock of Sql interface.
type MockSql struct {
	ctrl     *gomock.Controller
	recorder *MockSqlMockRecorder
}

// MockSqlMockRecorder is the mock recorder for MockSql.
type MockSqlMockRecorder struct {
	mock *MockSql
}

// NewMockSql creates a new mock instance.
func NewMockSql(ctrl *gomock.Controller) *MockSql {
	mock := &MockSql{ctrl: ctrl}
	mock.recorder = &MockSqlMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSql) EXPECT() *MockSqlMockRecorder {
	return m.recorder
}

// Open mocks base method.
func (m *MockSql) Open(driverName, dataSourceName string) (sqlproxy.DBInstanceInterface, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Open", driverName, dataSourceName)
	ret0, _ := ret[0].(sqlproxy.DBInstanceInterface)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Open indicates an expected call of Open.
func (mr *MockSqlMockRecorder) Open(driverName, dataSourceName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Open", reflect.TypeOf((*MockSql)(nil).Open), driverName, dataSourceName)
}

// StringToNullString mocks base method.
func (m *MockSql) StringToNullString(s string) *sqlproxy.NullStringInstance {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StringToNullString", s)
	ret0, _ := ret[0].(*sqlproxy.NullStringInstance)
	return ret0
}

// StringToNullString indicates an expected call of StringToNullString.
func (mr *MockSqlMockRecorder) StringToNullString(s any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StringToNullString", reflect.TypeOf((*MockSql)(nil).StringToNullString), s)
}
