package formatter

import (
	"errors"
	"reflect"
	"testing"
	"time"

	jrpApp "github.com/yanosea/jrp/v2/app/application/jrp"

	"github.com/yanosea/jrp/v2/pkg/proxy"
	"github.com/yanosea/jrp/v2/pkg/utility"

	"go.uber.org/mock/gomock"
)

func TestNewTableFormatter(t *testing.T) {
	tests := []struct {
		name string
		want *TableFormatter
	}{
		{
			name: "positive testing",
			want: &TableFormatter{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTableFormatter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTableFormatter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTableFormatter_Format(t *testing.T) {
	ti := time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)
	su := utility.NewStringsUtil()

	type args struct {
		result any
	}
	tests := []struct {
		name    string
		f       *TableFormatter
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "positive testing (result is []*jrpApp.GenerateJrpUseCaseOutputDto)",
			f:    &TableFormatter{},
			args: args{
				result: []*jrpApp.GenerateJrpUseCaseOutputDto{
					{
						ID:          1,
						Phrase:      "phrase1",
						Prefix:      "prefix1",
						Suffix:      "suffix1",
						IsFavorited: 1,
						CreatedAt:   ti,
						UpdatedAt:   ti,
					},
					{
						ID:          2,
						Phrase:      "phrase2",
						Prefix:      "prefix2",
						Suffix:      "suffix2",
						IsFavorited: 1,
						CreatedAt:   ti,
						UpdatedAt:   ti,
					},
				},
			},
			want:    "IDPHRASEPREFIXSUFFIXCREATEDAT1phrase1prefix1suffix12006-01-0215:04:052phrase2prefix2suffix22006-01-0215:04:05TOTAL:2jrps!",
			wantErr: false,
		},
		{
			name: "positive testing (result is []*jrpApp.GetHistoryUseCaseOutputDto)",
			f:    &TableFormatter{},
			args: args{
				result: []*jrpApp.GetHistoryUseCaseOutputDto{
					{
						ID:          1,
						Phrase:      "phrase1",
						Prefix:      "prefix1",
						Suffix:      "suffix1",
						IsFavorited: 1,
						CreatedAt:   ti,
						UpdatedAt:   ti,
					},
					{
						ID:          2,
						Phrase:      "phrase2",
						Prefix:      "prefix2",
						Suffix:      "suffix2",
						IsFavorited: 1,
						CreatedAt:   ti,
						UpdatedAt:   ti,
					},
				},
			},
			want:    "IDPHRASEPREFIXSUFFIXISFAVORITEDCREATEDATUPDATEDAT1phrase1prefix1suffix1○2006-01-0215:04:052006-01-0215:04:052phrase2prefix2suffix2○2006-01-0215:04:052006-01-0215:04:05TOTAL:2jrps!",
			wantErr: false,
		},
		{
			name: "positive testing (result is []*jrpApp.SearchHistoryUseCaseOutputDto)",
			f:    &TableFormatter{},
			args: args{
				result: []*jrpApp.SearchHistoryUseCaseOutputDto{
					{
						ID:          1,
						Phrase:      "phrase1",
						Prefix:      "prefix1",
						Suffix:      "suffix1",
						IsFavorited: 1,
						CreatedAt:   ti,
						UpdatedAt:   ti,
					},
					{
						ID:          2,
						Phrase:      "phrase2",
						Prefix:      "prefix2",
						Suffix:      "suffix2",
						IsFavorited: 1,
						CreatedAt:   ti,
						UpdatedAt:   ti,
					},
				},
			},
			want:    "IDPHRASEPREFIXSUFFIXISFAVORITEDCREATEDATUPDATEDAT1phrase1prefix1suffix1○2006-01-0215:04:052006-01-0215:04:052phrase2prefix2suffix2○2006-01-0215:04:052006-01-0215:04:05TOTAL:2jrps!",
			wantErr: false,
		},
		{
			name: "negative testing (result is invalid)",
			f:    &TableFormatter{},
			args: args{
				result: "invalid",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &TableFormatter{}
			got, err := f.Format(tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("TableFormatter.Format() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && su.RemoveNewLines(su.RemoveSpaces(su.RemoveTabs(got))) != tt.want {
				t.Errorf("TableFormatter.Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTableFormatter_formatGenerateJrp(t *testing.T) {
	ti := time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)

	type args struct {
		items []*jrpApp.GenerateJrpUseCaseOutputDto
	}
	tests := []struct {
		name string
		f    *TableFormatter
		args args
		want tableData
	}{
		{
			name: "positive testing (no id)",
			f:    &TableFormatter{},
			args: args{
				items: []*jrpApp.GenerateJrpUseCaseOutputDto{
					{
						Phrase:    "phrase1",
						Prefix:    "prefix1",
						Suffix:    "suffix1",
						CreatedAt: ti,
					},
					{
						Phrase:    "phrase2",
						Prefix:    "prefix2",
						Suffix:    "suffix2",
						CreatedAt: ti,
					},
				},
			},
			want: tableData{
				header: []string{"phrase", "prefix", "suffix", "created_at"},
				rows: [][]string{
					{"phrase1", "prefix1", "suffix1", "2006-01-02 15:04:05"},
					{"phrase2", "prefix2", "suffix2", "2006-01-02 15:04:05"},
				},
			},
		},
		{
			name: "positive testing (with id)",
			f:    &TableFormatter{},
			args: args{
				items: []*jrpApp.GenerateJrpUseCaseOutputDto{
					{
						ID:        1,
						Phrase:    "phrase1",
						Prefix:    "prefix1",
						Suffix:    "suffix1",
						CreatedAt: ti,
					},
					{
						ID:        2,
						Phrase:    "phrase2",
						Prefix:    "prefix2",
						Suffix:    "suffix2",
						CreatedAt: ti,
					},
				},
			},
			want: tableData{
				header: []string{"id", "phrase", "prefix", "suffix", "created_at"},
				rows: [][]string{
					{"1", "phrase1", "prefix1", "suffix1", "2006-01-02 15:04:05"},
					{"2", "phrase2", "prefix2", "suffix2", "2006-01-02 15:04:05"},
					{"", "", "", "", ""},
					{"TOTAL : 2 jrps!", "", "", "", ""},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.formatGenerateJrp(tt.args.items); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TableFormatter.formatGenerateJrp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatHistory(t *testing.T) {
	ti := time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)

	type args struct {
		items   any
		getData func(any) (int, string, string, string, int, time.Time, time.Time)
	}
	tests := []struct {
		name string
		f    *TableFormatter
		args args
		want tableData
	}{
		{
			name: "positive testing (items is []*jrpApp.GetHistoryUseCaseOutputDto)",
			f:    &TableFormatter{},
			args: args{
				items: []*jrpApp.GetHistoryUseCaseOutputDto{
					{
						ID:          1,
						Phrase:      "phrase1",
						Prefix:      "prefix1",
						Suffix:      "suffix1",
						IsFavorited: 1,
						CreatedAt:   ti,
						UpdatedAt:   ti,
					},
					{
						ID:          2,
						Phrase:      "phrase2",
						Prefix:      "prefix2",
						Suffix:      "suffix2",
						IsFavorited: 1,
						CreatedAt:   ti,
						UpdatedAt:   ti,
					},
				},
				getData: func(v any) (int, string, string, string, int, time.Time, time.Time) {
					dto := v.(*jrpApp.GetHistoryUseCaseOutputDto)
					return dto.ID, dto.Phrase, dto.Prefix, dto.Suffix, dto.IsFavorited, dto.CreatedAt, dto.UpdatedAt
				},
			},
			want: tableData{
				header: []string{"id", "phrase", "prefix", "suffix", "is_favorited", "created_at", "updated_at"},
				rows: [][]string{
					{"1", "phrase1", "prefix1", "suffix1", "○", "2006-01-02 15:04:05", "2006-01-02 15:04:05"},
					{"2", "phrase2", "prefix2", "suffix2", "○", "2006-01-02 15:04:05", "2006-01-02 15:04:05"},
					{"", "", "", "", "", "", ""},
					{"TOTAL : 2 jrps!", "", "", "", "", "", ""},
				},
			},
		},
		{
			name: "positive testing (items is []*jrpApp.SearchHistoryUseCaseOutputDto)",
			f:    &TableFormatter{},
			args: args{
				items: []*jrpApp.SearchHistoryUseCaseOutputDto{
					{
						ID:          1,
						Phrase:      "phrase1",
						Prefix:      "prefix1",
						Suffix:      "suffix1",
						IsFavorited: 1,
						CreatedAt:   ti,
						UpdatedAt:   ti,
					},
					{
						ID:          2,
						Phrase:      "phrase2",
						Prefix:      "prefix2",
						Suffix:      "suffix2",
						IsFavorited: 1,
						CreatedAt:   ti,
						UpdatedAt:   ti,
					},
				},
				getData: func(v any) (int, string, string, string, int, time.Time, time.Time) {
					dto := v.(*jrpApp.SearchHistoryUseCaseOutputDto)
					return dto.ID, dto.Phrase, dto.Prefix, dto.Suffix, dto.IsFavorited, dto.CreatedAt, dto.UpdatedAt
				},
			},
			want: tableData{
				header: []string{"id", "phrase", "prefix", "suffix", "is_favorited", "created_at", "updated_at"},
				rows: [][]string{
					{"1", "phrase1", "prefix1", "suffix1", "○", "2006-01-02 15:04:05", "2006-01-02 15:04:05"},
					{"2", "phrase2", "prefix2", "suffix2", "○", "2006-01-02 15:04:05", "2006-01-02 15:04:05"},
					{"", "", "", "", "", "", ""},
					{"TOTAL : 2 jrps!", "", "", "", "", "", ""},
				},
			},
		},
		{
			name: "positive testing (items is invalid)",
			f:    &TableFormatter{},
			args: args{
				items: "invalid",
				getData: func(v any) (int, string, string, string, int, time.Time, time.Time) {
					return 0, "", "", "", 0, time.Time{}, time.Time{}
				},
			},
			want: tableData{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.formatHistory(tt.args.items, tt.args.getData); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TableFormatter.formatHistory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTableFormatter_addTotalRow(t *testing.T) {
	type args struct {
		rows [][]string
	}
	tests := []struct {
		name string
		f    *TableFormatter
		args args
		want [][]string
	}{
		{
			name: "positive testing (rows is empty)",
			f:    &TableFormatter{},
			args: args{
				rows: [][]string{},
			},
			want: [][]string{},
		},
		{
			name: "positive testing (rows is not empty)",
			f:    &TableFormatter{},
			args: args{
				rows: [][]string{
					{"1", "phrase1", "prefix1", "suffix1", "2006-01-02 15:04:05"},
					{"2", "phrase2", "prefix2", "suffix2", "2006-01-02 15:04:05"},
				},
			},
			want: [][]string{
				{"1", "phrase1", "prefix1", "suffix1", "2006-01-02 15:04:05"},
				{"2", "phrase2", "prefix2", "suffix2", "2006-01-02 15:04:05"},
				{"", "", "", "", ""},
				{"TOTAL : 2 jrps!", "", "", "", ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &TableFormatter{}
			if got := f.addTotalFooter(tt.args.rows); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TableFormatter.addTotalRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTableFormatter_getTableString(t *testing.T) {
	su := utility.NewStringsUtil()
	origTu := Tu

	type args struct {
		data tableData
	}
	tests := []struct {
		name    string
		f       *TableFormatter
		args    args
		want    string
		wantErr bool
		setup   func(mockCtrl *gomock.Controller)
		cleanup func()
	}{
		{
			name: "positive testing (data is empty)",
			f:    &TableFormatter{},
			args: args{
				data: tableData{},
			},
			want:    "",
			wantErr: false,
			setup:   nil,
			cleanup: nil,
		},
		{
			name: "positive testing (data.Header is empty, data.Rows is empty)",
			f:    &TableFormatter{},
			args: args{
				data: tableData{
					header: []string{},
					rows:   [][]string{},
				},
			},
			want:    "",
			wantErr: false,
			setup:   nil,
			cleanup: nil,
		},
		{
			name: "positive testing (data.Header is not empty, data.Rows is empty)",
			f:    &TableFormatter{},
			args: args{
				data: tableData{
					header: []string{"id", "phrase", "prefix", "suffix", "created_at"},
					rows:   [][]string{},
				},
			},
			want:    "",
			wantErr: false,
			setup:   nil,
			cleanup: nil,
		},
		{
			name: "positive testing (data.Header is empty, data.Rows is not empty)",
			f:    &TableFormatter{},
			args: args{
				data: tableData{
					header: []string{},
					rows: [][]string{
						{"1", "phrase1", "prefix1", "suffix1", "2006-01-02 15:04:05"},
						{"2", "phrase2", "prefix2", "suffix2", "2006-01-02 15:04:05"},
						{},
						{"TOTAL : 2 jrps!"},
					},
				},
			},
			want:    "",
			wantErr: false,
			setup:   nil,
			cleanup: nil,
		},
		{
			name: "positive testing (data.Header is not empty, data.Rows is not empty)",
			f:    &TableFormatter{},
			args: args{
				data: tableData{
					header: []string{"id", "phrase", "prefix", "suffix", "created_at"},
					rows: [][]string{
						{"1", "phrase1", "prefix1", "suffix1", "2006-01-02 15:04:05"},
						{"2", "phrase2", "prefix2", "suffix2", "2006-01-02 15:04:05"},
						{},
						{"TOTAL : 2 jrps!"},
					},
				},
			},
			want:    "IDPHRASEPREFIXSUFFIXCREATEDAT1phrase1prefix1suffix12006-01-0215:04:052phrase2prefix2suffix22006-01-0215:04:05TOTAL:2jrps!",
			wantErr: false,
			setup:   nil,
			cleanup: nil,
		},
		{
			name: "negative testing (table.Bulk(data.rows) failed)",
			f:    &TableFormatter{},
			args: args{
				data: tableData{
					header: []string{"id", "phrase", "prefix", "suffix", "created_at"},
					rows: [][]string{
						{"1", "phrase1", "prefix1", "suffix1", "2006-01-02 15:04:05"},
						{"2", "phrase2", "prefix2", "suffix2", "2006-01-02 15:04:05"},
						{},
						{"TOTAL : 2 jrps!"},
					},
				},
			},
			want:    "",
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller) {
				mockTable := proxy.NewMockTable(mockCtrl)
				mockTable.EXPECT().Header(gomock.Any())
				mockTable.EXPECT().Bulk(gomock.Any()).Return(errors.New("bulk failed"))
				mockTableWriter := proxy.NewMockTableWriter(mockCtrl)
				mockTableWriter.EXPECT().NewTable(gomock.Any())
				Tu = utility.NewTableWriterUtil(mockTableWriter)
			},
			cleanup: func() {
				Tu = origTu
			},
		},
		{
			name: "negative testing (table.Render() failed)",
			f:    &TableFormatter{},
			args: args{
				data: tableData{
					header: []string{"id", "phrase", "prefix", "suffix", "created_at"},
					rows: [][]string{
						{"1", "phrase1", "prefix1", "suffix1", "2006-01-02 15:04:05"},
						{"2", "phrase2", "prefix2", "suffix2", "2006-01-02 15:04:05"},
						{},
						{"TOTAL : 2 jrps!"},
					},
				},
			},
			want:    "",
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller) {
				mockTable := proxy.NewMockTable(mockCtrl)
				mockTable.EXPECT().Header(gomock.Any())
				mockTable.EXPECT().Bulk(gomock.Any()).Return(nil)
				mockTable.EXPECT().Render().Return(errors.New("render failed"))
				mockTableWriter := proxy.NewMockTableWriter(mockCtrl)
				mockTableWriter.EXPECT().NewTable(gomock.Any())
				Tu = utility.NewTableWriterUtil(mockTableWriter)
			},
			cleanup: func() {
				Tu = origTu
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl)
			}
			defer func() {
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			f := &TableFormatter{}
			got, err := f.getTableString(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("TableFormatter.getTableString() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != "" && su.RemoveNewLines(su.RemoveSpaces(su.RemoveTabs(got))) != tt.want {
				t.Errorf("TableFormatter.getTableString() = %v, want %v", got, tt.want)
			}
		})
	}
}
