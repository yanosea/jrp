package query_service

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	jrpApp "github.com/yanosea/jrp/app/application/jrp"
	wnjpnApp "github.com/yanosea/jrp/app/application/wnjpn"
	"github.com/yanosea/jrp/app/infrastructure/database"

	"github.com/yanosea/jrp/pkg/proxy"

	"go.uber.org/mock/gomock"
)

func TestNewWordQueryService(t *testing.T) {
	cm := database.NewConnectionManager(proxy.NewSql())

	tests := []struct {
		name string
		want wnjpnApp.WordQueryService
	}{
		{
			name: "positive testing",
			want: &wordQueryService{
				connManager: cm,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWordQueryService(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWordQueryService() = %v, want %v", got, tt.want)
			}
		})
	}
	if err := database.ResetConnectionManager(); err != nil {
		t.Errorf("Failed to reset connection manager: %v", err)
	}
}

func Test_wordQueryService_FindByLangIsAndPosIn(t *testing.T) {
	duc := jrpApp.NewDownloadUseCase()
	if err := duc.Run(filepath.Join(os.TempDir(), "wnjpn.db")); err != nil && err.Error() != "wnjpn.db already exists" {
		t.Errorf("Failed to download WordNet Japan DB file: %v", err)
	}

	type fields struct {
		connManager database.ConnectionManager
	}
	type args struct {
		ctx  context.Context
		lang string
		pos  []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
		cleanup func()
	}{
		{
			name: "positive testing",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:  context.Background(),
				lang: "jpn",
				pos:  []string{"a", "v", "n"},
			},
			wantErr: false,
			setup: func(_ *gomock.Controller, tt *fields) {
				tt.connManager = database.NewConnectionManager(proxy.NewSql())
				if err := tt.connManager.InitializeConnection(database.ConnectionConfig{
					DBType: database.SQLite,
					DBName: database.WNJpnDB,
					DSN:    filepath.Join(os.TempDir(), "wnjpn.db"),
				}); err != nil {
					t.Errorf("Failed to initialize connection: %v", err)
				}
			},
			cleanup: func() {
				if err := database.ResetConnectionManager(); err != nil {
					t.Errorf("Failed to reset connection manager: %v", err)
				}
			},
		},
		{
			name: "negative testing (w.connManager.GetConnection(database.WNJpnDB) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:  context.Background(),
				lang: "jpn",
				pos:  []string{"a", "v", "n"},
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.WNJpnDB).Return(nil, errors.New("ConnectionManager.GetConnection() failed"))
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (conn.Open() failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:  context.Background(),
				lang: "jpn",
				pos:  []string{"a", "v", "n"},
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(nil, errors.New("DBConnection.Open() failed"))
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.WNJpnDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (db.QueryContext(ctx, query, params...) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:  context.Background(),
				lang: "jpn",
				pos:  []string{"a", "v", "n"},
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().QueryContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("proxy.DB.QueryContext() failed"))
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.WNJpnDB).Return(mockConnection, nil)
				tt.connManager = mockConnManager
			},
			cleanup: nil,
		},
		{
			name: "negative testing (rows.Scan(&word.WordID, &word.Lang, &word.Lemma, &word.Pron, &word.Pos) failed)",
			fields: fields{
				connManager: nil,
			},
			args: args{
				ctx:  context.Background(),
				lang: "jpn",
				pos:  []string{"a", "v", "n"},
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockRows := proxy.NewMockRows(mockCtrl)
				mockRows.EXPECT().Next().Return(true)
				mockRows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("proxy.Rows.Scan() failed"))
				mockRows.EXPECT().Close().Return(nil)
				mockDB := proxy.NewMockDB(mockCtrl)
				mockDB.EXPECT().QueryContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockRows, nil)
				mockConnection := database.NewMockDBConnection(mockCtrl)
				mockConnection.EXPECT().Open().Return(mockDB, nil)
				mockConnManager := database.NewMockConnectionManager(mockCtrl)
				mockConnManager.EXPECT().GetConnection(database.WNJpnDB).Return(mockConnection, nil)
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
			w := &wordQueryService{
				connManager: tt.fields.connManager,
			}
			got, err := w.FindByLangIsAndPosIn(tt.args.ctx, tt.args.lang, tt.args.pos)
			if (err != nil) != tt.wantErr {
				t.Errorf("wordQueryService.FindByLangIsAndPosIn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) == 0 {
				t.Errorf("wordQueryService.FindByLangIsAndPosIn() got = %v, want not empty", got)
			}
		})
	}
}
