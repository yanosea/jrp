package database

import (
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/yanosea/jrp/v2/pkg/proxy"
	"go.uber.org/mock/gomock"
)

func TestNewConnectionManager(t *testing.T) {
	origGcm := gcm
	sql := proxy.NewSql()

	type args struct {
		sql proxy.Sql
	}
	tests := []struct {
		name    string
		args    args
		want    ConnectionManager
		setup   func()
		cleanup func()
	}{
		{
			name: "positive testing (gcm is nil)",
			args: args{
				sql: sql,
			},
			want: &connectionManager{
				sql:         sql,
				connections: make(map[DBName]DBConnection),
				mutex:       &sync.RWMutex{},
			},
			setup: func() {
				gcm = nil
			},
			cleanup: func() {
				gcm = origGcm
			},
		},
		{
			name: "positive testing (gcm is not nil)",
			args: args{
				sql: proxy.NewSql(),
			},
			want: &connectionManager{
				sql:         sql,
				connections: make(map[DBName]DBConnection),
				mutex:       &sync.RWMutex{},
			},
			setup: func() {
				gcm = &connectionManager{
					sql:         sql,
					connections: make(map[DBName]DBConnection),
					mutex:       &sync.RWMutex{},
				}
			},
			cleanup: func() {
				gcm = origGcm
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			if got := NewConnectionManager(tt.args.sql); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConnectionManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetConnectionManager(t *testing.T) {
	origGcm := gcm
	sql := proxy.NewSql()

	tests := []struct {
		name    string
		want    ConnectionManager
		setup   func()
		cleanup func()
	}{
		{
			name: "positive testing (gcm is nil)",
			want: nil,
			setup: func() {
				gcm = nil
			},
			cleanup: func() {
				gcm = origGcm
			},
		},
		{
			name: "positive testing (gcm is not nil)",
			want: &connectionManager{
				sql:         sql,
				connections: make(map[DBName]DBConnection),
				mutex:       &sync.RWMutex{},
			},
			setup: func() {
				gcm = &connectionManager{
					sql:         sql,
					connections: make(map[DBName]DBConnection),
					mutex:       &sync.RWMutex{},
				}
			},
			cleanup: func() {
				gcm = origGcm
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			if got := GetConnectionManager(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConnectionManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getConnectionManager(t *testing.T) {
	origGcm := gcm
	sql := proxy.NewSql()

	tests := []struct {
		name    string
		want    ConnectionManager
		setup   func()
		cleanup func()
	}{
		{
			name: "positive testing (gcm is nil)",
			want: nil,
			setup: func() {
				gcm = nil
			},
			cleanup: func() {
				gcm = origGcm
			},
		},
		{
			name: "positive testing (gcm is not nil)",
			want: &connectionManager{
				sql:         sql,
				connections: make(map[DBName]DBConnection),
				mutex:       &sync.RWMutex{},
			},
			setup: func() {
				gcm = &connectionManager{
					sql:         sql,
					connections: make(map[DBName]DBConnection),
					mutex:       &sync.RWMutex{},
				}
			},
			cleanup: func() {
				gcm = origGcm
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			if got := getConnectionManager(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getConnectionManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResetConnectionManager(t *testing.T) {
	origGcm := gcm
	sql := proxy.NewSql()

	tests := []struct {
		name    string
		isNil   bool
		wantErr bool
		setup   func(mockCtrl *gomock.Controller)
		cleanup func()
	}{
		{
			name:    "positive testing (gcm is nil)",
			isNil:   true,
			wantErr: false,
			setup: func(_ *gomock.Controller) {
				gcm = nil
			},
			cleanup: func() {
				gcm = origGcm
			},
		},
		{
			name:    "positive testing (gcm is not nil)",
			isNil:   false,
			wantErr: false,
			setup: func(_ *gomock.Controller) {
				gcm = &connectionManager{
					sql:         sql,
					connections: make(map[DBName]DBConnection),
					mutex:       &sync.RWMutex{},
				}
			},
			cleanup: func() {
				gcm = origGcm
			},
		},
		{
			name:    "negative testing (gcm.CloseAllConnections() failed)",
			isNil:   false,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller) {
				mockGcm := NewMockConnectionManager(mockCtrl)
				mockGcm.EXPECT().CloseAllConnections().Return(errors.New("CloseAllConnections() failed"))
				gcm = mockGcm
			},
			cleanup: func() {
				gcm = origGcm
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(gomock.NewController(t))
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			if err := ResetConnectionManager(); (err != nil) != tt.wantErr {
				t.Errorf("ResetConnectionManager() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if gcm == nil {
					t.Error("ResetConnectionManager() should not set gcm to nil when error occurs")
				}
			} else {
				if gcm != nil {
					t.Error("ResetConnectionManager() should set gcm to nil when no error occurs")
				}
			}
		})
	}
}

func Test_connectionManager_CloseAllConnections(t *testing.T) {
	type fields struct {
		sql         proxy.Sql
		connections map[DBName]DBConnection
		mu          *sync.RWMutex
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing",
			fields: fields{
				sql:         proxy.NewSql(),
				connections: make(map[DBName]DBConnection),
				mu:          &sync.RWMutex{},
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockConn := NewMockDBConnection(mockCtrl)
				mockConn.EXPECT().Close().Return(nil)
				tt.connections[DBName("test")] = mockConn
			},
		},
		{
			name: "negative testing (conn.Close() failed)",
			fields: fields{
				sql:         proxy.NewSql(),
				connections: make(map[DBName]DBConnection),
				mu:          &sync.RWMutex{},
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockConn := NewMockDBConnection(mockCtrl)
				mockConn.EXPECT().Close().Return(errors.New("DBConnection.Close() failed"))
				tt.connections[DBName("test")] = mockConn
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			gcm := &connectionManager{
				sql:         tt.fields.sql,
				connections: tt.fields.connections,
				mutex:       tt.fields.mu,
			}
			if err := gcm.CloseAllConnections(); (err != nil) != tt.wantErr {
				t.Errorf("connectionManager.CloseAllConnections() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_connectionManager_CloseConnection(t *testing.T) {
	type fields struct {
		sql         proxy.Sql
		connections map[DBName]DBConnection
		mu          *sync.RWMutex
	}
	type args struct {
		dbType DBName
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing (gcm.connections[dbType] exists)",
			fields: fields{
				sql:         proxy.NewSql(),
				connections: make(map[DBName]DBConnection),
				mu:          &sync.RWMutex{},
			},
			args: args{
				dbType: DBName("test"),
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockConn := NewMockDBConnection(mockCtrl)
				mockConn.EXPECT().Close().Return(nil)
				tt.connections[DBName("test")] = mockConn
			},
		},
		{
			name: "positive testing (gcm.connections[dbType] does not exist)",
			fields: fields{
				sql:         proxy.NewSql(),
				connections: make(map[DBName]DBConnection),
				mu:          &sync.RWMutex{},
			},
			args: args{
				dbType: DBName("not_exist"),
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				tt.connections[DBName("test")] = nil
			},
		},
		{
			name: "negative testing (conn.Close() failed)",
			fields: fields{
				sql:         proxy.NewSql(),
				connections: make(map[DBName]DBConnection),
				mu:          &sync.RWMutex{},
			},
			args: args{
				dbType: DBName("test"),
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockConn := NewMockDBConnection(mockCtrl)
				mockConn.EXPECT().Close().Return(errors.New("DBConnection.Close() failed"))
				tt.connections[DBName("test")] = mockConn
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			gcm := &connectionManager{
				sql:         tt.fields.sql,
				connections: tt.fields.connections,
				mutex:       tt.fields.mu,
			}
			if err := gcm.CloseConnection(tt.args.dbType); (err != nil) != tt.wantErr {
				t.Errorf("connectionManager.CloseConnection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_connectionManager_GetConnection(t *testing.T) {
	type fields struct {
		sql         proxy.Sql
		connections map[DBName]DBConnection
		mu          *sync.RWMutex
		want        DBConnection
	}
	type args struct {
		dbType DBName
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing (gcm.connections[dbType] exists)",
			fields: fields{
				sql:         proxy.NewSql(),
				connections: make(map[DBName]DBConnection),
				mu:          &sync.RWMutex{},
				want:        nil,
			},
			args: args{
				dbType: DBName("test"),
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockConn := NewMockDBConnection(mockCtrl)
				tt.connections[DBName("test")] = mockConn
				tt.want = mockConn
			},
		},
		{
			name: "positive testing (gcm.connections[dbType] does not exist)",
			fields: fields{
				sql:         proxy.NewSql(),
				connections: make(map[DBName]DBConnection),
				mu:          &sync.RWMutex{},
				want:        nil,
			},
			args: args{
				dbType: DBName("not_exist"),
			},
			wantErr: true,
			setup:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			gcm := &connectionManager{
				sql:         tt.fields.sql,
				connections: tt.fields.connections,
				mutex:       tt.fields.mu,
			}
			got, err := gcm.GetConnection(tt.args.dbType)
			if (err != nil) != tt.wantErr {
				t.Errorf("connectionManager.GetConnection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.fields.want) {
				t.Errorf("connectionManager.GetConnection() = %v, want %v", got, tt.fields.want)
			}
		})
	}
}

func Test_connectionManager_InitializeConnection(t *testing.T) {
	type fields struct {
		sql         proxy.Sql
		connections map[DBName]DBConnection
		mu          *sync.RWMutex
	}
	type args struct {
		config ConnectionConfig
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func(tt *fields)
	}{
		{
			name: "positive testing (gcm.connections[config.DBName] does not exist)",
			fields: fields{
				sql:         proxy.NewSql(),
				connections: make(map[DBName]DBConnection),
				mu:          &sync.RWMutex{},
			},
			args: args{
				config: ConnectionConfig{
					DBName: DBName("test"),
					DBType: DBType("sqlite"),
					DSN:    "test.db",
				},
			},
			wantErr: false,
			setup:   nil,
		},
		{
			name: "negative testing (gcm.connections[config.DBName] already exists)",
			fields: fields{
				sql:         proxy.NewSql(),
				connections: make(map[DBName]DBConnection),
				mu:          &sync.RWMutex{},
			},
			args: args{
				config: ConnectionConfig{
					DBName: DBName("test"),
					DBType: DBType("sqlite"),
					DSN:    "test.db",
				},
			},
			wantErr: true,
			setup: func(tt *fields) {
				tt.connections[DBName("test")] = &dbConnection{
					sql:            tt.sql,
					driverName:     "sqlite",
					dataSourceName: "test.db",
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(&tt.fields)
			}
			gcm := &connectionManager{
				sql:         tt.fields.sql,
				connections: tt.fields.connections,
				mutex:       tt.fields.mu,
			}
			if err := gcm.InitializeConnection(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("connectionManager.InitializeConnection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
