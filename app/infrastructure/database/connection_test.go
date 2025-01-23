package database

import (
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/yanosea/jrp/v2/pkg/proxy"

	"go.uber.org/mock/gomock"
)

func Test_dbConnection_Close(t *testing.T) {
	type fields struct {
		sql            proxy.Sql
		db             proxy.DB
		driverName     string
		dataSourceName string
		mutex          *sync.RWMutex
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing (db is nil)",
			fields: fields{
				sql:            proxy.NewSql(),
				db:             nil,
				driverName:     "sqlite",
				dataSourceName: "test.db",
				mutex:          &sync.RWMutex{},
			},
			wantErr: false,
			setup:   nil,
		},
		{
			name: "positive testing (db is not nil)",
			fields: fields{
				sql:            proxy.NewSql(),
				db:             nil,
				driverName:     "sqlite",
				dataSourceName: "test.db",
				mutex:          &sync.RWMutex{},
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().Close().Return(nil)
				tt.db = mockDB
			},
		},
		{
			name: "negative testing (c.db.Close() failed)",
			fields: fields{
				sql:            proxy.NewSql(),
				db:             nil,
				driverName:     "sqlite",
				dataSourceName: "test.db",
				mutex:          &sync.RWMutex{},
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().Close().Return(errors.New("proxy.DB.Close() failed"))
				tt.db = mockDB
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
			c := &dbConnection{
				sql:            tt.fields.sql,
				db:             tt.fields.db,
				driverName:     tt.fields.driverName,
				dataSourceName: tt.fields.dataSourceName,
				mutex:          tt.fields.mutex,
			}
			if err := c.Close(); (err != nil) != tt.wantErr {
				mockCtrl := gomock.NewController(t)
				defer mockCtrl.Finish()
				t.Errorf("dbConnection.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_dbConnection_Open(t *testing.T) {
	type fields struct {
		sql            proxy.Sql
		db             proxy.DB
		driverName     string
		dataSourceName string
		mutex          *sync.RWMutex
		want           proxy.DB
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing (db is not nil)",
			fields: fields{
				sql:            proxy.NewSql(),
				db:             nil,
				driverName:     "sqlite",
				dataSourceName: "test.db",
				mutex:          &sync.RWMutex{},
				want:           nil,
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDB := proxy.NewMockDB(mockCtrl)
				tt.db = mockDB
				tt.want = mockDB
			},
		},
		{
			name: "positive testing (db is nil)",
			fields: fields{
				sql:            proxy.NewSql(),
				db:             nil,
				driverName:     "sqlite",
				dataSourceName: "test.db",
				mutex:          &sync.RWMutex{},
				want:           nil,
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDB := proxy.NewMockDB(mockCtrl)
				mockSql := proxy.NewMockSql(mockCtrl)
				mockSql.EXPECT().Open(tt.driverName, tt.dataSourceName).Return(mockDB, nil)
				tt.sql = mockSql
				tt.want = mockDB
			},
		},
		{
			name: "negative testing (c.sql.Open() failed)",
			fields: fields{
				sql:            nil,
				db:             nil,
				driverName:     "sqlite",
				dataSourceName: "test.db",
				mutex:          &sync.RWMutex{},
				want:           nil,
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockSql := proxy.NewMockSql(mockCtrl)
				mockSql.EXPECT().Open(tt.driverName, tt.dataSourceName).Return(nil, errors.New("proxy.Sql.Open() failed"))
				tt.sql = mockSql
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
			c := &dbConnection{
				sql:            tt.fields.sql,
				db:             tt.fields.db,
				driverName:     tt.fields.driverName,
				dataSourceName: tt.fields.dataSourceName,
				mutex:          tt.fields.mutex,
			}
			got, err := c.Open()
			if (err != nil) != tt.wantErr {
				t.Errorf("dbConnection.Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.fields.want) {
				t.Errorf("dbConnection.Open() = %v, want %v", got, tt.fields.want)
			}
		})
	}
}
