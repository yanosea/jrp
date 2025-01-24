package repository

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	historyDomain "github.com/yanosea/jrp/v2/app/domain/jrp/history"
	"github.com/yanosea/jrp/v2/app/infrastructure/database"

	"github.com/yanosea/jrp/v2/pkg/proxy"

	"go.uber.org/mock/gomock"
)

var (
	now = time.Now()
)

func TestNewHistoryRepository(t *testing.T) {
	cm := database.NewConnectionManager(proxy.NewSql())

	tests := []struct {
		name string
		want historyDomain.HistoryRepository
	}{
		{
			name: "positive testing",
			want: &historyRepository{
				connManager: cm,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHistoryRepository(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHistoryRepository() = %v, want %v", got, tt.want)
			}
		})
	}
	if err := database.ResetConnectionManager(); err != nil {
		t.Errorf("Failed to reset connection manager: %v", err)
	}
}

func Test_historyRepository_DeleteAll(t *testing.T) {
	type fields struct {
		connManager database.ConnectionManager
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		testData []*historyDomain.History
		want     int
		wantErr  bool
		setup    func(mockCtrl *gomock.Controller, tt *fields)
		cleanup  func()
	}{
		{
			name: "positive testing (no histories in the database)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
			},
			testData: nil,
			want:     0,
			wantErr:  false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (1 history in the database)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    1,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (2 history in the database)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    2,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "negative testing (getJrpDB(ctx, h.connManager) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
			},
			testData: nil,
			want:     0,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(nil, errors.New("ConnectionManager.GetConnection() failed"))
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: false}) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
			},
			testData: nil,
			want:     0,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(nil, errors.New("DB.BeginTx() failed"))
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (tx.ExecContext(ctx, DeleteAllQuery) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
			},
			testData: nil,
			want:     0,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockTx := proxy.NewMockTx(mockCtrl)
				mockTx.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("proxy.Tx.ExecContext() failed"))
				mockTx.EXPECT().Rollback().Return(nil)
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(mockTx, nil)
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (tx.ExecContext(ctx, DeleteSequenceQuery) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
			},
			testData: nil,
			want:     0,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockTx := proxy.NewMockTx(mockCtrl)
				mockTx.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockTx.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("proxy.Tx.ExecContext() failed"))
				mockTx.EXPECT().Rollback().Return(nil)
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(mockTx, nil)
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (result.RowsAffected() failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
			},
			testData: nil,
			want:     0,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockResult := proxy.NewMockResult(mockCtrl)
				mockResult.EXPECT().RowsAffected().Return(int64(0), errors.New("proxy.Result.RowsAffected() failed"))
				mockTx := proxy.NewMockTx(mockCtrl)
				mockTx.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockResult, nil)
				mockTx.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockTx.EXPECT().Rollback().Return(nil)
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(mockTx, nil)
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (tx.Commit() failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
			},
			testData: nil,
			want:     0,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockResult := proxy.NewMockResult(mockCtrl)
				mockResult.EXPECT().RowsAffected().Return(int64(0), nil)
				mockTx := proxy.NewMockTx(mockCtrl)
				mockTx.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockResult, nil)
				mockTx.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockTx.EXPECT().Commit().Return(errors.New("proxy.Tx.Commit() failed"))
				mockTx.EXPECT().Rollback().Return(nil)
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(mockTx, nil)
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			h := &historyRepository{
				connManager: tt.fields.connManager,
			}
			if len(tt.testData) > 0 {
				if _, err := h.SaveAll(tt.args.ctx, tt.testData); err != nil {
					t.Errorf("Failed to save test data: %v", err)
				}
			}
			got, err := h.DeleteAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("historyRepository.DeleteAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("historyRepository.DeleteAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_historyRepository_DeleteByIdIn(t *testing.T) {
	type fields struct {
		connManager database.ConnectionManager
	}
	type args struct {
		ctx context.Context
		ids []int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		testData []*historyDomain.History
		want     int
		wantErr  bool
		setup    func(mockCtrl *gomock.Controller, tt *fields)
		cleanup  func()
	}{
		{
			name: "positive testing (ids of args are empty)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
				ids: []int{},
			},
			testData: nil,
			want:     0,
			wantErr:  false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (no histories in the database)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
				ids: []int{1},
			},
			testData: nil,
			want:     0,
			wantErr:  false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (1 history in the database, id of args does not exist)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
				ids: []int{2},
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    0,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (1 history in the database, id of args exists)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
				ids: []int{1},
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    1,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (2 histories in the database, ids of args does not exist)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
				ids: []int{3, 4},
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    0,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (2 histories in the database, one of the id of args exists)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
				ids: []int{2, 3},
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    1,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (2 histories in the database, al of the id of args exists)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
				ids: []int{1, 2},
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    2,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "negative testing (getJrpDB(ctx, h.connManager) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
				ids: []int{1},
			},
			testData: nil,
			want:     0,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(nil, errors.New("ConnectionManager.GetConnection() failed"))
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (db.ExecContext(ctx, DeleteByIdInQuery, args...) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
				ids: []int{1},
			},
			testData: nil,
			want:     0,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("proxy.DB.ExecContext() failed"))
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (result.RowsAffected() failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
				ids: []int{1},
			},
			testData: nil,
			want:     0,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockResult := proxy.NewMockResult(mockCtrl)
				mockResult.EXPECT().RowsAffected().Return(int64(0), errors.New("proxy.Result.RowsAffected() failed"))
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockResult, nil)
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			h := &historyRepository{
				connManager: tt.fields.connManager,
			}
			if len(tt.testData) > 0 {
				if _, err := h.SaveAll(tt.args.ctx, tt.testData); err != nil {
					t.Errorf("Failed to save test data: %v", err)
				}
			}
			got, err := h.DeleteByIdIn(tt.args.ctx, tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("historyRepository.DeleteByIdIn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("historyRepository.DeleteByIdIn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_historyRepository_DeleteByIdInAndIsFavoritedIs(t *testing.T) {
	type fields struct {
		connManager database.ConnectionManager
	}
	type args struct {
		ctx         context.Context
		ids         []int
		isFavorited int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		testData []*historyDomain.History
		want     int
		wantErr  bool
		setup    func(mockCtrl *gomock.Controller, tt *fields)
		cleanup  func()
	}{
		{
			name: "positive testing (ids of args are empty)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				ids:         []int{},
				isFavorited: 0,
			},
			testData: nil,
			want:     0,
			wantErr:  false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (no histories in the database)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				ids:         []int{1},
				isFavorited: 0,
			},
			testData: nil,
			want:     0,
			wantErr:  false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (1 history in the database, id of args does not exist)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				ids:         []int{2},
				isFavorited: 0,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    0,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (1 history in the database (isFavorited = 0), id of args exists but isFavorited does not match)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				ids:         []int{1},
				isFavorited: 1,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    0,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (1 history in the database (isFavorited = 0), id of args exists and isFavorited matches)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				ids:         []int{1},
				isFavorited: 0,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    1,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (2 histories in the database (isFavorited = 0, 1), ids of args do not exist)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				ids:         []int{3, 4},
				isFavorited: 0,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    0,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (2 histories in the database, (isFavorited = 1) all of the id of args exists)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				ids:         []int{1, 2},
				isFavorited: 1,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    2,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (2 histories in the database (isFavorited = 1), one of the id of args exists and isFavorited matches)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				ids:         []int{2, 3},
				isFavorited: 1,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    1,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (2 histories in the database (isFavorited = 1), all of the id of args exists and isFavorited matches)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				ids:         []int{1, 2},
				isFavorited: 1,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    2,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (2 histories in the database (isFavorited = 0, 1), all of the id of args exists but isFavorited does not match)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				ids:         []int{1, 2},
				isFavorited: 1,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    0,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (2 histories in the database (isFavorited = 0, 1), one id exists but isFavorited does not match)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				ids:         []int{1, 2},
				isFavorited: 1,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    1,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (2 histories in the database (isFavorited = 0), all of the id of args exists and isFavorited matches)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				ids:         []int{1, 2},
				isFavorited: 0,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    2,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (3 histories in the database (isFavorited = 0,0,1), some ids exist and some isFavorited matches)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				ids:         []int{1, 2, 3},
				isFavorited: 0,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test3",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    2,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "negative testing (getJrpDB(ctx, h.connManager) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
				ids: []int{1},
			},
			testData: nil,
			want:     0,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(nil, errors.New("ConnectionManager.GetConnection() failed"))
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (db.ExecContext(ctx, DeleteByIdInAndIsFavoritedIsQuery, args.ids, args.isFavorited) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
				ids: []int{1},
			},
			testData: nil,
			want:     0,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("proxy.DB.ExecContext() failed"))
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (result.RowsAffected() failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
				ids: []int{1},
			},
			testData: nil,
			want:     0,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockResult := proxy.NewMockResult(mockCtrl)
				mockResult.EXPECT().RowsAffected().Return(int64(0), errors.New("proxy.Result.RowsAffected() failed"))
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockResult, nil)
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			h := &historyRepository{
				connManager: tt.fields.connManager,
			}
			if len(tt.testData) > 0 {
				if _, err := h.SaveAll(tt.args.ctx, tt.testData); err != nil {
					t.Errorf("Failed to save test data: %v", err)
				}
			}
			got, err := h.DeleteByIdInAndIsFavoritedIs(tt.args.ctx, tt.args.ids, tt.args.isFavorited)
			if (err != nil) != tt.wantErr {
				t.Errorf("historyRepository.DeleteByIdInAndIsFavoritedIs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("historyRepository.DeleteByIdInAndIsFavoritedIs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_historyRepository_DeleteByIsFavoritedIs(t *testing.T) {
	type fields struct {
		connManager database.ConnectionManager
	}
	type args struct {
		ctx         context.Context
		isFavorited int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		testData []*historyDomain.History
		want     int
		wantErr  bool
		setup    func(mockCtrl *gomock.Controller, tt *fields)
		cleanup  func()
	}{
		{
			name: "positive testing (no histories in the database)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				isFavorited: 0,
			},
			testData: nil,
			want:     0,
			wantErr:  false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (1 history in the database (isFavorited = 0), isFavorited matches)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				isFavorited: 0,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    1,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (1 history in the database (isFavorited = 0), isFavorited does not match)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				isFavorited: 1,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    0,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (2 histories in the database (isFavorited = 0, 1), isFavorited matches one)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				isFavorited: 1,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    1,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (2 histories in the database (isFavorited = 1), isFavorited matches all)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				isFavorited: 1,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    2,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "negative testing (getJrpDB(ctx, h.connManager) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				isFavorited: 0,
			},
			testData: nil,
			want:     0,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(nil, errors.New("ConnectionManager.GetConnection() failed"))
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (db.ExecContext(ctx, DeleteByIsFavoritedIsQuery, args.isFavorited) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				isFavorited: 0,
			},
			testData: nil,
			want:     0,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("proxy.DB.ExecContext() failed"))
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (result.RowsAffected() failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				isFavorited: 0,
			},
			testData: nil,
			want:     0,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockResult := proxy.NewMockResult(mockCtrl)
				mockResult.EXPECT().RowsAffected().Return(int64(0), errors.New("proxy.Result.RowsAffected() failed"))
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockResult, nil)
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			h := &historyRepository{
				connManager: tt.fields.connManager,
			}
			if len(tt.testData) > 0 {
				if _, err := h.SaveAll(tt.args.ctx, tt.testData); err != nil {
					t.Errorf("Failed to save test data: %v", err)
				}
			}
			got, err := h.DeleteByIsFavoritedIs(tt.args.ctx, tt.args.isFavorited)
			if (err != nil) != tt.wantErr {
				t.Errorf("historyRepository.DeleteByIsFavoritedIs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("historyRepository.DeleteByIsFavoritedIs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_historyRepository_FindAll(t *testing.T) {
	type fields struct {
		connManager database.ConnectionManager
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		testData []*historyDomain.History
		want     []*historyDomain.History
		wantErr  bool
		setup    func(mockCtrl *gomock.Controller, tt *fields)
		cleanup  func()
	}{
		{
			name: "positive testing (no histories in the database)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
			},
			testData: nil,
			want:     nil,
			wantErr:  false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (1 history in the database)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want: []*historyDomain.History{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (2 histories in the database)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix2",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix2",
						Valid:  true,
					},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want: []*historyDomain.History{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sql.NullString{String: "prefix2", Valid: true},
					Suffix:      sql.NullString{String: "suffix2", Valid: true},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "negative testing (getJrpDB(ctx, h.connManager) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(nil, errors.New("ConnectionManager.GetConnection() failed"))
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (db.QueryContext(ctx, FindAllQuery) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().QueryContext(gomock.Any(), gomock.Any()).Return(nil, errors.New("DB.QueryContext() failed"))
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (rows.Scan() failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockRows := proxy.NewMockRows(mockCtrl)
				mockRows.EXPECT().Next().Return(true)
				mockRows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("Rows.Scan() failed"))
				mockRows.EXPECT().Close().Return(nil)
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().QueryContext(gomock.Any(), gomock.Any()).Return(mockRows, nil)
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			h := &historyRepository{
				connManager: tt.fields.connManager,
			}
			if len(tt.testData) > 0 {
				if _, err := h.SaveAll(tt.args.ctx, tt.testData); err != nil {
					t.Errorf("Failed to save test data: %v", err)
				}
			}
			got, err := h.FindAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("historyRepository.FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("historyRepository.FindAll() returned %d items, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i].ID != tt.want[i].ID {
					t.Errorf("historyRepository.FindAll()[%d].ID = %v, want %v", i, got[i].ID, tt.want[i].ID)
				}
				if got[i].Phrase != tt.want[i].Phrase {
					t.Errorf("historyRepository.FindAll()[%d].Phrase = %v, want %v", i, got[i].Phrase, tt.want[i].Phrase)
				}
				if got[i].Prefix != tt.want[i].Prefix {
					t.Errorf("historyRepository.FindAll()[%d].Prefix = %v, want %v", i, got[i].Prefix, tt.want[i].Prefix)
				}
				if got[i].Suffix != tt.want[i].Suffix {
					t.Errorf("historyRepository.FindAll()[%d].Suffix = %v, want %v", i, got[i].Suffix, tt.want[i].Suffix)
				}
				if got[i].IsFavorited != tt.want[i].IsFavorited {
					t.Errorf("historyRepository.FindAll()[%d].IsFavorited = %v, want %v", i, got[i].IsFavorited, tt.want[i].IsFavorited)
				}
			}
		})
	}
}

func Test_historyRepository_FindByIsFavoritedIs(t *testing.T) {
	type fields struct {
		connManager database.ConnectionManager
	}
	type args struct {
		ctx         context.Context
		isFavorited int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		testData []*historyDomain.History
		want     []*historyDomain.History
		wantErr  bool
		setup    func(mockCtrl *gomock.Controller, tt *fields)
		cleanup  func()
	}{
		{
			name: "positive testing (no histories in the database)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				isFavorited: 0,
			},
			testData: nil,
			want:     nil,
			wantErr:  false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (1 history in the database (isFavorited = 0), isFavorited matches)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				isFavorited: 0,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want: []*historyDomain.History{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (1 history in the database (isFavorited = 0), isFavorited does not match)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				isFavorited: 1,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (2 histories in the database (isFavorited = 0, 1), isFavorited matches one)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				isFavorited: 1,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want: []*historyDomain.History{
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (1 history in the database (isFavorited = 0), isFavorited does not match)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				isFavorited: 1,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    nil,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (2 histories in the database (isFavorited = 0, 1), isFavorited matches one)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				isFavorited: 1,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want: []*historyDomain.History{
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "negative testing (getJrpDB(ctx, h.connManager) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				isFavorited: 0,
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(nil, errors.New("ConnectionManager.GetConnection() failed"))
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (db.QueryContext(ctx, FindByIsFavoritedIsQuery, isFavorited) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				isFavorited: 0,
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().QueryContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("proxy.DB.QueryContext() failed"))
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (rows.Scan() failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				isFavorited: 0,
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockRows := proxy.NewMockRows(mockCtrl)
				mockRows.EXPECT().Next().Return(true)
				mockRows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("Rows.Scan() failed"))
				mockRows.EXPECT().Close().Return(nil)
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().QueryContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockRows, nil)
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			h := &historyRepository{
				connManager: tt.fields.connManager,
			}
			if len(tt.testData) > 0 {
				if _, err := h.SaveAll(tt.args.ctx, tt.testData); err != nil {
					t.Errorf("Failed to save test data: %v", err)
				}
			}
			got, err := h.FindByIsFavoritedIs(tt.args.ctx, tt.args.isFavorited)
			if (err != nil) != tt.wantErr {
				t.Errorf("historyRepository.FindByIsFavoritedIs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("historyRepository.FindByIsFavoritedIs() returned %d items, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i].ID != tt.want[i].ID {
					t.Errorf("historyRepository.FindByIsFavoritedIs()[%d].ID = %v, want %v", i, got[i].ID, tt.want[i].ID)
				}
				if got[i].Phrase != tt.want[i].Phrase {
					t.Errorf("historyRepository.FindByIsFavoritedIs()[%d].Phrase = %v, want %v", i, got[i].Phrase, tt.want[i].Phrase)
				}
				if got[i].Prefix != tt.want[i].Prefix {
					t.Errorf("historyRepository.FindByIsFavoritedIs()[%d].Prefix = %v, want %v", i, got[i].Prefix, tt.want[i].Prefix)
				}
				if got[i].Suffix != tt.want[i].Suffix {
					t.Errorf("historyRepository.FindByIsFavoritedIs()[%d].Suffix = %v, want %v", i, got[i].Suffix, tt.want[i].Suffix)
				}
				if got[i].IsFavorited != tt.want[i].IsFavorited {
					t.Errorf("historyRepository.FindByIsFavoritedIs()[%d].IsFavorited = %v, want %v", i, got[i].IsFavorited, tt.want[i].IsFavorited)
				}
			}
		})
	}
}

func Test_historyRepository_FindByIsFavoritedIsAndPhraseContains(t *testing.T) {
	type fields struct {
		connManager database.ConnectionManager
	}
	type args struct {
		ctx         context.Context
		keywords    []string
		and         bool
		isFavorited int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		testData []*historyDomain.History
		want     []*historyDomain.History
		wantErr  bool
		setup    func(mockCtrl *gomock.Controller, tt *fields)
		cleanup  func()
	}{
		{
			name: "positive testing (no histories in the database)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				keywords:    []string{"test"},
				and:         false,
				isFavorited: 0,
			},
			testData: nil,
			want:     nil,
			wantErr:  false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (1 history in the database (isFavorited = 0), keyword and isFavorited match)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				keywords:    []string{"test"},
				and:         false,
				isFavorited: 0,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want: []*historyDomain.History{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (2 histories in the database (isFavorited = 0, 1), keyword matches both but isFavorited matches one)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				keywords:    []string{"test"},
				and:         false,
				isFavorited: 1,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test1",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want: []*historyDomain.History{
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (2 histories in the database (isFavorited = 0), AND search with multiple keywords)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				keywords:    []string{"test", "match"},
				and:         true,
				isFavorited: 0,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test only",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test match both",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want: []*historyDomain.History{
				{
					ID:          2,
					Phrase:      "test match both",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (3 histories in the database (isFavorited = 0), OR search with multiple keywords)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				keywords:    []string{"test", "new"},
				and:         false,
				isFavorited: 0,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test only",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "new content",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "other content",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want: []*historyDomain.History{
				{
					ID:          1,
					Phrase:      "test only",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:          2,
					Phrase:      "new content",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (2 histories in the database (isFavorited = 0), AND search with multiple keywords)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				keywords:    []string{"test", "match"},
				and:         true,
				isFavorited: 0,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test only",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test match both",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want: []*historyDomain.History{
				{
					ID:          2,
					Phrase:      "test match both",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "negative testing (getJrpDB(ctx, h.connManager) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				keywords:    []string{"test"},
				and:         false,
				isFavorited: 0,
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(nil, errors.New("ConnectionManager.GetConnection() failed"))
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (db.QueryContext(ctx, FindByIsFavoritedIsAndPhraseContainsQuery, args...) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				keywords:    []string{"test"},
				and:         false,
				isFavorited: 0,
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().QueryContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("proxy.DB.QueryContext() failed"))
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (rows.Scan() failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				keywords:    []string{"test"},
				and:         false,
				isFavorited: 0,
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockRows := proxy.NewMockRows(mockCtrl)
				mockRows.EXPECT().Next().Return(true)
				mockRows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("Rows.Scan() failed"))
				mockRows.EXPECT().Close().Return(nil)
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().QueryContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockRows, nil)
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			h := &historyRepository{
				connManager: tt.fields.connManager,
			}
			if len(tt.testData) > 0 {
				if _, err := h.SaveAll(tt.args.ctx, tt.testData); err != nil {
					t.Errorf("Failed to save test data: %v", err)
				}
			}
			got, err := h.FindByIsFavoritedIsAndPhraseContains(tt.args.ctx, tt.args.keywords, tt.args.and, tt.args.isFavorited)
			if (err != nil) != tt.wantErr {
				t.Errorf("historyRepository.FindByIsFavoritedIsAndPhraseContainsAndPhraseContains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("historyRepository.FindByIsFavoritedIsAndPhraseContains() returned %d items, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i].ID != tt.want[i].ID {
					t.Errorf("historyRepository.FindByIsFavoritedIsAndPhraseContains()[%d].ID = %v, want %v", i, got[i].ID, tt.want[i].ID)
				}
				if got[i].Phrase != tt.want[i].Phrase {
					t.Errorf("historyRepository.FindByIsFavoritedIsAndPhraseContains()[%d].Phrase = %v, want %v", i, got[i].Phrase, tt.want[i].Phrase)
				}
				if got[i].Prefix != tt.want[i].Prefix {
					t.Errorf("historyRepository.FindByIsFavoritedIsAndPhraseContains()[%d].Prefix = %v, want %v", i, got[i].Prefix, tt.want[i].Prefix)
				}
				if got[i].Suffix != tt.want[i].Suffix {
					t.Errorf("historyRepository.FindByIsFavoritedIsAndPhraseContains()[%d].Suffix = %v, want %v", i, got[i].Suffix, tt.want[i].Suffix)
				}
				if got[i].IsFavorited != tt.want[i].IsFavorited {
					t.Errorf("historyRepository.FindByIsFavoritedIsAndPhraseContains()[%d].IsFavorited = %v, want %v", i, got[i].IsFavorited, tt.want[i].IsFavorited)
				}
			}
		})
	}
}

func Test_historyRepository_FindByPhraseContains(t *testing.T) {
	type fields struct {
		connManager database.ConnectionManager
	}
	type args struct {
		ctx      context.Context
		keywords []string
		and      bool
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		testData []*historyDomain.History
		want     []*historyDomain.History
		wantErr  bool
		setup    func(mockCtrl *gomock.Controller, tt *fields)
		cleanup  func()
	}{
		{
			name: "positive testing (no histories in the database)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:      context.Background(),
				keywords: []string{"test"},
				and:      false,
			},
			testData: nil,
			want:     nil,
			wantErr:  false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (1 history in the database, keyword matches)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:      context.Background(),
				keywords: []string{"test"},
				and:      false,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want: []*historyDomain.History{
				{
					ID:     1,
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (2 histories in the database, keyword matches one)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:      context.Background(),
				keywords: []string{"test1"},
				and:      false,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test1",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want: []*historyDomain.History{
				{
					ID:     1,
					Phrase: "test1",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (2 histories in the database, AND search with multiple keywords)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:      context.Background(),
				keywords: []string{"test", "match"},
				and:      true,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test only",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test match both",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want: []*historyDomain.History{
				{
					ID:     2,
					Phrase: "test match both",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (3 histories in the database, OR search with multiple keywords)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:      context.Background(),
				keywords: []string{"test", "new"},
				and:      false,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test only",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "new content",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "other content",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want: []*historyDomain.History{
				{
					ID:     1,
					Phrase: "test only",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:     2,
					Phrase: "new content",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "negative testing (getJrpDB(ctx, h.connManager) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:      context.Background(),
				keywords: []string{"test"},
				and:      false,
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(nil, errors.New("ConnectionManager.GetConnection() failed"))
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (db.QueryContext(ctx, FindByPhraseContainsQuery, args...) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:      context.Background(),
				keywords: []string{"test"},
				and:      false,
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().QueryContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("proxy.DB.QueryContext() failed"))
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (rows.Scan() failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:      context.Background(),
				keywords: []string{"test"},
				and:      false,
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockRows := proxy.NewMockRows(mockCtrl)
				mockRows.EXPECT().Next().Return(true)
				mockRows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("Rows.Scan() failed"))
				mockRows.EXPECT().Close().Return(nil)
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().QueryContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockRows, nil)
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			h := &historyRepository{
				connManager: tt.fields.connManager,
			}
			if len(tt.testData) > 0 {
				if _, err := h.SaveAll(tt.args.ctx, tt.testData); err != nil {
					t.Errorf("Failed to save test data: %v", err)
				}
			}
			got, err := h.FindByPhraseContains(tt.args.ctx, tt.args.keywords, tt.args.and)
			if (err != nil) != tt.wantErr {
				t.Errorf("historyRepository.FindByPhraseContains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("historyRepository.FindByPhraseContains() returned %d items, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i].ID != tt.want[i].ID {
					t.Errorf("historyRepository.FindByPhraseContains()[%d].ID = %v, want %v", i, got[i].ID, tt.want[i].ID)
				}
				if got[i].Phrase != tt.want[i].Phrase {
					t.Errorf("historyRepository.FindByPhraseContains()[%d].Phrase = %v, want %v", i, got[i].Phrase, tt.want[i].Phrase)
				}
				if got[i].Prefix != tt.want[i].Prefix {
					t.Errorf("historyRepository.FindByPhraseContains()[%d].Prefix = %v, want %v", i, got[i].Prefix, tt.want[i].Prefix)
				}
				if got[i].Suffix != tt.want[i].Suffix {
					t.Errorf("historyRepository.FindByPhraseContains()[%d].Suffix = %v, want %v", i, got[i].Suffix, tt.want[i].Suffix)
				}
				if got[i].IsFavorited != tt.want[i].IsFavorited {
					t.Errorf("historyRepository.FindByPhraseContains()[%d].IsFavorited = %v, want %v", i, got[i].IsFavorited, tt.want[i].IsFavorited)
				}
			}
		})
	}
}

func Test_historyRepository_FindTopNByIsFavoritedIsAndByOrderByIdAsc(t *testing.T) {
	type fields struct {
		connManager database.ConnectionManager
	}
	type args struct {
		ctx         context.Context
		number      int
		isFavorited int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		testData []*historyDomain.History
		want     []*historyDomain.History
		wantErr  bool
		setup    func(mockCtrl *gomock.Controller, tt *fields)
		cleanup  func()
	}{
		{
			name: "positive testing (no histories in the database)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				number:      2,
				isFavorited: 0,
			},
			testData: nil,
			want:     nil,
			wantErr:  false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (1 history in the database (isFavorited = 0), isFavorited matches)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				number:      2,
				isFavorited: 0,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want: []*historyDomain.History{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (3 histories in the database (isFavorited = 0, 1), isFavorited matches two and number limits to one)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				number:      1,
				isFavorited: 0,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test1",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test3",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want: []*historyDomain.History{
				{
					ID:          3,
					Phrase:      "test3",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (3 histories in the database (isFavorited = 0, 1), isFavorited matches two and number matches)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				number:      2,
				isFavorited: 0,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test1",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test3",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want: []*historyDomain.History{
				{
					ID:          1,
					Phrase:      "test1",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:          3,
					Phrase:      "test3",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "negative testing (getJrpDB(ctx, h.connManager) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				number:      2,
				isFavorited: 0,
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(nil, errors.New("ConnectionManager.GetConnection() failed"))
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (db.QueryContext(ctx, FindTopNByIsFavoritedIsAndByOrderByIdAscQuery, isFavorited, number) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				number:      2,
				isFavorited: 0,
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().QueryContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("DB.QueryContext() failed"))
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (rows.Scan() failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				number:      2,
				isFavorited: 0,
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockRows := proxy.NewMockRows(mockCtrl)
				mockRows.EXPECT().Next().Return(true)
				mockRows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("Rows.Scan() failed"))
				mockRows.EXPECT().Close().Return(nil)
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().QueryContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockRows, nil)
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			h := &historyRepository{
				connManager: tt.fields.connManager,
			}
			if len(tt.testData) > 0 {
				if _, err := h.SaveAll(tt.args.ctx, tt.testData); err != nil {
					t.Errorf("Failed to save test data: %v", err)
				}
			}
			got, err := h.FindTopNByIsFavoritedIsAndByOrderByIdAsc(tt.args.ctx, tt.args.number, tt.args.isFavorited)
			if (err != nil) != tt.wantErr {
				t.Errorf("historyRepository.FindTopNByIsFavoritedIsAndByOrderByIdAsc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("historyRepository.FindTopNByIsFavoritedIsAndByOrderByIdAsc() returned %d items, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i].ID != tt.want[i].ID {
					t.Errorf("historyRepository.FindTopNByIsFavoritedIsAndByOrderByIdAsc()[%d].ID = %v, want %v", i, got[i].ID, tt.want[i].ID)
				}
				if got[i].Phrase != tt.want[i].Phrase {
					t.Errorf("historyRepository.FindTopNByIsFavoritedIsAndByOrderByIdAsc()[%d].Phrase = %v, want %v", i, got[i].Phrase, tt.want[i].Phrase)
				}
				if got[i].Prefix != tt.want[i].Prefix {
					t.Errorf("historyRepository.FindTopNByIsFavoritedIsAndByOrderByIdAsc()[%d].Prefix = %v, want %v", i, got[i].Prefix, tt.want[i].Prefix)
				}
				if got[i].Suffix != tt.want[i].Suffix {
					t.Errorf("historyRepository.FindTopNByIsFavoritedIsAndByOrderByIdAsc()[%d].Suffix = %v, want %v", i, got[i].Suffix, tt.want[i].Suffix)
				}
				if got[i].IsFavorited != tt.want[i].IsFavorited {
					t.Errorf("historyRepository.FindTopNByIsFavoritedIsAndByOrderByIdAsc()[%d].IsFavorited = %v, want %v", i, got[i].IsFavorited, tt.want[i].IsFavorited)
				}
			}
		})
	}
}

func Test_historyRepository_FindTopNByIsFavoritedIsAndByPhraseContainsOrderByIdAsc(t *testing.T) {
	type fields struct {
		connManager database.ConnectionManager
	}
	type args struct {
		ctx         context.Context
		keywords    []string
		and         bool
		number      int
		isFavorited int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		testData []*historyDomain.History
		want     []*historyDomain.History
		wantErr  bool
		setup    func(mockCtrl *gomock.Controller, tt *fields)
		cleanup  func()
	}{
		{
			name: "positive testing (no histories in the database)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				keywords:    []string{"test"},
				and:         false,
				number:      2,
				isFavorited: 0,
			},
			testData: nil,
			want:     nil,
			wantErr:  false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (3 histories in the database, match 2 by keywords and favorited, limit to 1)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				keywords:    []string{"test"},
				and:         false,
				number:      1,
				isFavorited: 1,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test1",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "other",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want: []*historyDomain.History{
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (3 histories, AND search with multiple keywords, limit 2)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				keywords:    []string{"test", "match"},
				and:         true,
				number:      2,
				isFavorited: 0,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test only",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test match both",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test match also",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want: []*historyDomain.History{
				{
					ID:          2,
					Phrase:      "test match both",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:          3,
					Phrase:      "test match also",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (3 histories in the database, OR search with multiple keywords, limit 2)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				keywords:    []string{"test", "new"},
				and:         false,
				number:      2,
				isFavorited: 0,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test first",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "new second",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "other content",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want: []*historyDomain.History{
				{
					ID:          1,
					Phrase:      "test first",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:          2,
					Phrase:      "new second",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "negative testing (getJrpDB(ctx, h.connManager) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				keywords:    []string{"test"},
				and:         false,
				number:      1,
				isFavorited: 0,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(nil, errors.New("ConnectionManager.GetConnection() failed"))
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (db.QueryContext(ctx, FindTopNByIsFavoritedIsAndByPhraseContainsOrderByIdAscQuery, args...) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				keywords:    []string{"test"},
				and:         false,
				number:      1,
				isFavorited: 0,
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().QueryContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("DB.QueryContext() failed"))
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (rows.Scan() failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				keywords:    []string{"test"},
				and:         false,
				number:      1,
				isFavorited: 0,
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockRows := proxy.NewMockRows(mockCtrl)
				mockRows.EXPECT().Next().Return(true)
				mockRows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("Rows.Scan() failed"))
				mockRows.EXPECT().Close().Return(nil)
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().QueryContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockRows, nil)
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			h := &historyRepository{
				connManager: tt.fields.connManager,
			}
			if len(tt.testData) > 0 {
				if _, err := h.SaveAll(tt.args.ctx, tt.testData); err != nil {
					t.Errorf("Failed to save test data: %v", err)
				}
			}
			got, err := h.FindTopNByIsFavoritedIsAndByPhraseContainsOrderByIdAsc(tt.args.ctx, tt.args.keywords, tt.args.and, tt.args.number, tt.args.isFavorited)
			if (err != nil) != tt.wantErr {
				t.Errorf("historyRepository.FindTopNByIsFavoritedIsAndByPhraseContainsOrderByIdAsc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("historyRepository.FindTopNByIsFavoritedIsAndByPhraseContainsOrderByIdAsc() returned %d items, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i].ID != tt.want[i].ID {
					t.Errorf("historyRepository.FindTopNByIsFavoritedIsAndByPhraseContainsOrderByIdAsc()[%d].ID = %v, want %v", i, got[i].ID, tt.want[i].ID)
				}
				if got[i].Phrase != tt.want[i].Phrase {
					t.Errorf("historyRepository.FindTopNByIsFavoritedIsAndByPhraseContainsOrderByIdAsc()[%d].Phrase = %v, want %v", i, got[i].Phrase, tt.want[i].Phrase)
				}
				if got[i].Prefix != tt.want[i].Prefix {
					t.Errorf("historyRepository.FindTopNByIsFavoritedIsAndByPhraseContainsOrderByIdAsc()[%d].Prefix = %v, want %v", i, got[i].Prefix, tt.want[i].Prefix)
				}
				if got[i].Suffix != tt.want[i].Suffix {
					t.Errorf("historyRepository.FindTopNByIsFavoritedIsAndByPhraseContainsOrderByIdAsc()[%d].Suffix = %v, want %v", i, got[i].Suffix, tt.want[i].Suffix)
				}
				if got[i].IsFavorited != tt.want[i].IsFavorited {
					t.Errorf("historyRepository.FindTopNByIsFavoritedIsAndByPhraseContainsOrderByIdAsc()[%d].IsFavorited = %v, want %v", i, got[i].IsFavorited, tt.want[i].IsFavorited)
				}
			}
		})
	}
}

func Test_historyRepository_FindTopNByOrderByIdAsc(t *testing.T) {
	type fields struct {
		connManager database.ConnectionManager
	}
	type args struct {
		ctx    context.Context
		number int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		testData []*historyDomain.History
		want     []*historyDomain.History
		wantErr  bool
		setup    func(mockCtrl *gomock.Controller, tt *fields)
		cleanup  func()
	}{
		{
			name: "positive testing (no histories in the database)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:    context.Background(),
				number: 2,
			},
			testData: nil,
			want:     nil,
			wantErr:  false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (1 history in the database, limit=2)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:    context.Background(),
				number: 2,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want: []*historyDomain.History{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (3 histories in the database, limit=2)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:    context.Background(),
				number: 2,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test1",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test3",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want: []*historyDomain.History{
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:          3,
					Phrase:      "test3",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "negative testing (getJrpDB(ctx, h.connManager) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:    context.Background(),
				number: 2,
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(nil, errors.New("ConnectionManager.GetConnection() failed"))
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (db.QueryContext() failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:    context.Background(),
				number: 2,
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().QueryContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("DB.QueryContext() failed"))
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (rows.Scan() failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:    context.Background(),
				number: 2,
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockRows := proxy.NewMockRows(mockCtrl)
				mockRows.EXPECT().Next().Return(true)
				mockRows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("Rows.Scan() failed"))
				mockRows.EXPECT().Close().Return(nil)
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().QueryContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockRows, nil)
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			h := &historyRepository{
				connManager: tt.fields.connManager,
			}
			if len(tt.testData) > 0 {
				if _, err := h.SaveAll(tt.args.ctx, tt.testData); err != nil {
					t.Errorf("Failed to save test data: %v", err)
				}
			}
			got, err := h.FindTopNByOrderByIdAsc(tt.args.ctx, tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("historyRepository.FindTopNByOrderByIdAs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("historyRepository.FindTopNByOrderByIdAs() returned %d items, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i].ID != tt.want[i].ID {
					t.Errorf("historyRepository.FindTopNByOrderByIdAs()[%d].ID = %v, want %v", i, got[i].ID, tt.want[i].ID)
				}
				if got[i].Phrase != tt.want[i].Phrase {
					t.Errorf("historyRepository.FindTopNByOrderByIdAs()[%d].Phrase = %v, want %v", i, got[i].Phrase, tt.want[i].Phrase)
				}
				if got[i].Prefix != tt.want[i].Prefix {
					t.Errorf("historyRepository.FindTopNByOrderByIdAs()[%d].Prefix = %v, want %v", i, got[i].Prefix, tt.want[i].Prefix)
				}
				if got[i].Suffix != tt.want[i].Suffix {
					t.Errorf("historyRepository.FindTopNByOrderByIdAs()[%d].Suffix = %v, want %v", i, got[i].Suffix, tt.want[i].Suffix)
				}
				if got[i].IsFavorited != tt.want[i].IsFavorited {
					t.Errorf("historyRepository.FindTopNByOrderByIdAs()[%d].IsFavorited = %v, want %v", i, got[i].IsFavorited, tt.want[i].IsFavorited)
				}
			}
		})
	}
}

func Test_historyRepository_FindTopNByPhraseContainsOrderByIdAsc(t *testing.T) {
	type fields struct {
		connManager database.ConnectionManager
	}
	type args struct {
		ctx      context.Context
		keywords []string
		and      bool
		number   int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		testData []*historyDomain.History
		want     []*historyDomain.History
		wantErr  bool
		setup    func(mockCtrl *gomock.Controller, tt *fields)
		cleanup  func()
	}{
		{
			name: "positive testing (no histories in the database)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:      context.Background(),
				keywords: []string{"test"},
				and:      false,
				number:   2,
			},
			testData: nil,
			want:     nil,
			wantErr:  false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (3 histories in the database, match 2 by keywords, limit to 1)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:      context.Background(),
				keywords: []string{"test"},
				and:      false,
				number:   1,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test1",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "other",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want: []*historyDomain.History{
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (3 histories, AND search with multiple keywords, limit 2)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:      context.Background(),
				keywords: []string{"test", "match"},
				and:      true,
				number:   2,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test only",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test match both",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test match also",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want: []*historyDomain.History{
				{
					ID:          2,
					Phrase:      "test match both",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:          3,
					Phrase:      "test match also",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (3 histories in the database, OR search with multiple keywords, limit 2)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:      context.Background(),
				keywords: []string{"test", "new"},
				and:      false,
				number:   2,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test first",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "new second",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "other content",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want: []*historyDomain.History{
				{
					ID:          1,
					Phrase:      "test first",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:          2,
					Phrase:      "new second",
					Prefix:      sql.NullString{String: "prefix", Valid: true},
					Suffix:      sql.NullString{String: "suffix", Valid: true},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "negative testing (getJrpDB(ctx, h.connManager) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:      context.Background(),
				keywords: []string{"test"},
				and:      false,
				number:   1,
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(nil, errors.New("ConnectionManager.GetConnection() failed"))
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (db.QueryContext(ctx, FindTopNByPhraseContainsOrderByIdAscQuery, args...) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:      context.Background(),
				keywords: []string{"test"},
				and:      false,
				number:   1,
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().QueryContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("proxy.DB.QueryContext() failed"))
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (rows.Scan() failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:      context.Background(),
				keywords: []string{"test"},
				and:      false,
				number:   1,
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockRows := proxy.NewMockRows(mockCtrl)
				mockRows.EXPECT().Next().Return(true)
				mockRows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("Rows.Scan() failed"))
				mockRows.EXPECT().Close().Return(nil)
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().QueryContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockRows, nil)
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			h := &historyRepository{
				connManager: tt.fields.connManager,
			}
			if len(tt.testData) > 0 {
				if _, err := h.SaveAll(tt.args.ctx, tt.testData); err != nil {
					t.Errorf("Failed to save test data: %v", err)
				}
			}
			got, err := h.FindTopNByPhraseContainsOrderByIdAsc(tt.args.ctx, tt.args.keywords, tt.args.and, tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("historyRepository.FindTopNByPhraseContainsOrderByIdAsc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("historyRepository.FindTopNByPhraseContainsOrderByIdAsc() returned %d items, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i].ID != tt.want[i].ID {
					t.Errorf("historyRepository.FindTopNByPhraseContainsOrderByIdAsc()[%d].ID = %v, want %v", i, got[i].ID, tt.want[i].ID)
				}
				if got[i].Phrase != tt.want[i].Phrase {
					t.Errorf("historyRepository.FindTopNByPhraseContainsOrderByIdAsc()[%d].Phrase = %v, want %v", i, got[i].Phrase, tt.want[i].Phrase)
				}
				if got[i].Prefix != tt.want[i].Prefix {
					t.Errorf("historyRepository.FindTopNByPhraseContainsOrderByIdAsc()[%d].Prefix = %v, want %v", i, got[i].Prefix, tt.want[i].Prefix)
				}
				if got[i].Suffix != tt.want[i].Suffix {
					t.Errorf("historyRepository.FindTopNByPhraseContainsOrderByIdAsc()[%d].Suffix = %v, want %v", i, got[i].Suffix, tt.want[i].Suffix)
				}
				if got[i].IsFavorited != tt.want[i].IsFavorited {
					t.Errorf("historyRepository.FindTopNByPhraseContainsOrderByIdAsc()[%d].IsFavorited = %v, want %v", i, got[i].IsFavorited, tt.want[i].IsFavorited)
				}
			}
		})
	}
}

func Test_historyRepository_SaveAll(t *testing.T) {
	type fields struct {
		connManager database.ConnectionManager
	}
	type args struct {
		ctx  context.Context
		jrps []*historyDomain.History
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		testData []*historyDomain.History
		want     []*historyDomain.History
		wantErr  bool
		setup    func(mockCtrl *gomock.Controller, tt *fields)
		cleanup  func()
	}{
		{
			name: "positive testing (save empty slice)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:  context.Background(),
				jrps: []*historyDomain.History{},
			},
			testData: nil,
			want:     []*historyDomain.History{},
			wantErr:  false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (save 1 history)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
				jrps: []*historyDomain.History{
					{
						Phrase: "test",
						Prefix: sql.NullString{
							String: "prefix",
							Valid:  true,
						},
						Suffix: sql.NullString{
							String: "suffix",
							Valid:  true,
						},
						IsFavorited: 0,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				},
			},
			testData: nil,
			want: []*historyDomain.History{
				{
					ID:     1,
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (save multiple histories)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
				jrps: []*historyDomain.History{
					{
						Phrase: "test1",
						Prefix: sql.NullString{
							String: "prefix1",
							Valid:  true,
						},
						Suffix: sql.NullString{
							String: "suffix1",
							Valid:  true,
						},
						IsFavorited: 0,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
					{
						Phrase: "test2",
						Prefix: sql.NullString{
							String: "prefix2",
							Valid:  true,
						},
						Suffix: sql.NullString{
							String: "suffix2",
							Valid:  true,
						},
						IsFavorited: 1,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				},
			},
			testData: nil,
			want: []*historyDomain.History{
				{
					ID:     1,
					Phrase: "test1",
					Prefix: sql.NullString{
						String: "prefix1",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix1",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:     2,
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix2",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix2",
						Valid:  true,
					},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "negative testing (getJrpDB(ctx, h.connManager) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
				jrps: []*historyDomain.History{
					{
						Phrase: "test",
						Prefix: sql.NullString{
							String: "prefix",
							Valid:  true,
						},
						Suffix: sql.NullString{
							String: "suffix",
							Valid:  true,
						},
						IsFavorited: 0,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				},
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(nil, errors.New("ConnectionManager.GetConnection() failed"))
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (db.BeginTx() failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
				jrps: []*historyDomain.History{
					{
						Phrase: "test",
						Prefix: sql.NullString{
							String: "prefix",
							Valid:  true,
						},
						Suffix: sql.NullString{
							String: "suffix",
							Valid:  true,
						},
						IsFavorited: 0,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				},
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(nil, errors.New("DB.BeginTx() failed"))
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (tx.ExecContext() failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
				jrps: []*historyDomain.History{
					{
						Phrase: "test",
						Prefix: sql.NullString{
							String: "prefix",
							Valid:  true,
						},
						Suffix: sql.NullString{
							String: "suffix",
							Valid:  true,
						},
						IsFavorited: 0,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				},
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockTx := proxy.NewMockTx(mockCtrl)
				mockTx.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("Tx.ExecContext() failed"))
				mockTx.EXPECT().Rollback().Return(nil)
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(mockTx, nil)
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (result.LastInsertId() failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
				jrps: []*historyDomain.History{
					{
						Phrase: "test",
						Prefix: sql.NullString{
							String: "prefix",
							Valid:  true,
						},
						Suffix: sql.NullString{
							String: "suffix",
							Valid:  true,
						},
						IsFavorited: 0,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				},
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockResult := proxy.NewMockResult(mockCtrl)
				mockResult.EXPECT().LastInsertId().Return(int64(0), errors.New("Result.LastInsertId() failed"))
				mockTx := proxy.NewMockTx(mockCtrl)
				mockTx.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockResult, nil)
				mockTx.EXPECT().Rollback().Return(nil)
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(mockTx, nil)
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (tx.Commit() failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx: context.Background(),
				jrps: []*historyDomain.History{
					{
						Phrase: "test",
						Prefix: sql.NullString{
							String: "prefix",
							Valid:  true,
						},
						Suffix: sql.NullString{
							String: "suffix",
							Valid:  true,
						},
						IsFavorited: 0,
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				},
			},
			testData: nil,
			want:     nil,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockResult := proxy.NewMockResult(mockCtrl)
				mockResult.EXPECT().LastInsertId().Return(int64(1), nil)
				mockTx := proxy.NewMockTx(mockCtrl)
				mockTx.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockResult, nil)
				mockTx.EXPECT().Commit().Return(errors.New("Tx.Commit() failed"))
				mockTx.EXPECT().Rollback().Return(nil)
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(mockTx, nil)
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			h := &historyRepository{
				connManager: tt.fields.connManager,
			}
			got, err := h.SaveAll(tt.args.ctx, tt.args.jrps)
			if (err != nil) != tt.wantErr {
				t.Errorf("historyRepository.SaveAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("historyRepository.FindTopNByPhraseContainsOrderByIdAsc() returned %d items, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i].ID != tt.want[i].ID {
					t.Errorf("historyRepository.SaveAll()[%d].ID = %v, want %v", i, got[i].ID, tt.want[i].ID)
				}
				if got[i].Phrase != tt.want[i].Phrase {
					t.Errorf("historyRepository.SaveAll()[%d].Phrase = %v, want %v", i, got[i].Phrase, tt.want[i].Phrase)
				}
				if got[i].Prefix != tt.want[i].Prefix {
					t.Errorf("historyRepository.SaveAll()[%d].Prefix = %v, want %v", i, got[i].Prefix, tt.want[i].Prefix)
				}
				if got[i].Suffix != tt.want[i].Suffix {
					t.Errorf("historyRepository.SaveAll()[%d].Suffix = %v, want %v", i, got[i].Suffix, tt.want[i].Suffix)
				}
				if got[i].IsFavorited != tt.want[i].IsFavorited {
					t.Errorf("historyRepository.SaveAll()[%d].IsFavorited = %v, want %v", i, got[i].IsFavorited, tt.want[i].IsFavorited)
				}
			}
		})
	}
}

func Test_historyRepository_UpdateIsFavoritedByIdIn(t *testing.T) {
	type fields struct {
		connManager database.ConnectionManager
	}
	type args struct {
		ctx         context.Context
		isFavorited int
		ids         []int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		testData []*historyDomain.History
		want     int
		wantErr  bool
		setup    func(mockCtrl *gomock.Controller, tt *fields)
		cleanup  func()
	}{
		{
			name: "positive testing (no histories in the database)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				isFavorited: 1,
				ids:         []int{1},
			},
			testData: nil,
			want:     0,
			wantErr:  false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (1 history in the database, id exists)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				isFavorited: 1,
				ids:         []int{1},
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    1,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (2 histories in the database, one id exists)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				isFavorited: 1,
				ids:         []int{1},
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test1",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    1,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (2 histories in the database, both ids exist)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				isFavorited: 1,
				ids:         []int{1, 2},
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test1",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    2,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "negative testing (getJrpDB(ctx, h.connManager) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				isFavorited: 1,
				ids:         []int{1},
			},
			testData: nil,
			want:     0,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(nil, errors.New("ConnectionManager.GetConnection() failed"))
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (db.ExecContext(ctx, UpdateIsFavoritedByIdInQuery, args...) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				isFavorited: 1,
				ids:         []int{1},
			},
			testData: nil,
			want:     0,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("DB.ExecContext() failed"))
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (result.RowsAffected() failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:         context.Background(),
				isFavorited: 1,
				ids:         []int{1},
			},
			testData: nil,
			want:     0,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockResult := proxy.NewMockResult(mockCtrl)
				mockResult.EXPECT().RowsAffected().Return(int64(0), errors.New("Result.RowsAffected() failed"))
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockResult, nil)
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			h := &historyRepository{
				connManager: tt.fields.connManager,
			}
			if len(tt.testData) > 0 {
				if _, err := h.SaveAll(tt.args.ctx, tt.testData); err != nil {
					t.Errorf("Failed to save test data: %v", err)
				}
			}
			got, err := h.UpdateIsFavoritedByIdIn(tt.args.ctx, tt.args.isFavorited, tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("historyRepository.UpdateIsFavoritedByIdIn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("historyRepository.UpdateIsFavoritedByIdIn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_historyRepository_UpdateIsFavoritedByIsFavoritedIs(t *testing.T) {
	type fields struct {
		connManager database.ConnectionManager
	}
	type args struct {
		ctx           context.Context
		isFavorited   int
		isFavoritedIs int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		testData []*historyDomain.History
		want     int
		wantErr  bool
		setup    func(mockCtrl *gomock.Controller, tt *fields)
		cleanup  func()
	}{
		{
			name: "positive testing (no histories in the database)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:           context.Background(),
				isFavorited:   1,
				isFavoritedIs: 0,
			},
			testData: nil,
			want:     0,
			wantErr:  false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (1 history in the database (isFavorited = 0), isFavoritedIs matches)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:           context.Background(),
				isFavorited:   1,
				isFavoritedIs: 0,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    1,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (1 history in the database (isFavorited = 0), isFavoritedIs does not match)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:           context.Background(),
				isFavorited:   1,
				isFavoritedIs: 1,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    0,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (2 histories in the database (isFavorited = 0), isFavoritedIs matches both)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:           context.Background(),
				isFavorited:   1,
				isFavoritedIs: 0,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test1",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    2,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "positive testing (2 histories in the database (isFavorited = 0, 1), isFavoritedIs matches one)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:           context.Background(),
				isFavorited:   1,
				isFavoritedIs: 0,
			},
			testData: []*historyDomain.History{
				{
					Phrase: "test1",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 0,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					Phrase: "test2",
					Prefix: sql.NullString{
						String: "prefix",
						Valid:  true,
					},
					Suffix: sql.NullString{
						String: "suffix",
						Valid:  true,
					},
					IsFavorited: 1,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			want:    1,
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "negative testing (getJrpDB(ctx, h.connManager) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:           context.Background(),
				isFavorited:   1,
				isFavoritedIs: 0,
			},
			testData: nil,
			want:     0,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(nil, errors.New("ConnectionManager.GetConnection() failed"))
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (db.ExecContext(ctx, UpdateIsFavoritedByIsFavoritedIsQuery, isFavorited, isFavoritedIs) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:           context.Background(),
				isFavorited:   1,
				isFavoritedIs: 0,
			},
			testData: nil,
			want:     0,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("DB.ExecContext() failed"))
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (result.RowsAffected() failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:           context.Background(),
				isFavorited:   1,
				isFavoritedIs: 0,
			},
			testData: nil,
			want:     0,
			wantErr:  true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockResult := proxy.NewMockResult(mockCtrl)
				mockResult.EXPECT().RowsAffected().Return(int64(0), errors.New("Result.RowsAffected() failed"))
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockResult, nil)
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			h := &historyRepository{
				connManager: tt.fields.connManager,
			}
			if len(tt.testData) > 0 {
				if _, err := h.SaveAll(tt.args.ctx, tt.testData); err != nil {
					t.Errorf("Failed to save test data: %v", err)
				}
			}
			got, err := h.UpdateIsFavoritedByIsFavoritedIs(tt.args.ctx, tt.args.isFavorited, tt.args.isFavoritedIs)
			if (err != nil) != tt.wantErr {
				t.Errorf("historyRepository.UpdateIsFavoritedByIsFavoritedIs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("historyRepository.UpdateIsFavoritedByIsFavoritedIs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getJrpDB(t *testing.T) {
	type args struct {
		ctx         context.Context
		connManager database.ConnectionManager
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *args)
		cleanup func()
	}{
		{
			name: "positive testing",
			args: args{
				ctx:         context.Background(),
				connManager: nil,
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *args) {
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.JrpDB,
					DSN:    filepath.Join(os.TempDir(), "jrp.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
				if err := os.Remove(filepath.Join(os.TempDir(), "jrp.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
			},
		},
		{
			name: "negative testing (GetConnection() failed)",
			args: args{
				ctx:         context.Background(),
				connManager: nil,
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(nil, errors.New("ConnectionManager.GetConnection() failed"))
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (conn.Open() failed)",
			args: args{
				ctx:         context.Background(),
				connManager: nil,
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(nil, errors.New("DBConnection.Open() failed"))
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (db.ExecContext() failed)",
			args: args{
				ctx:         context.Background(),
				connManager: nil,
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().ExecContext(gomock.Any(), gomock.Any()).Return(nil, errors.New("DB.ExecContext() failed"))
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.JrpDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.args)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			_, err := getJrpDB(tt.args.ctx, tt.args.connManager)
			if (err != nil) != tt.wantErr {
				t.Errorf("getJrpDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
