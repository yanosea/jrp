package jrp

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"

	historyDomain "github.com/yanosea/jrp/app/domain/jrp/history"

	"go.uber.org/mock/gomock"
)

func TestNewSearchHistoryUseCase(t *testing.T) {
	type args struct {
		historyRepo historyDomain.HistoryRepository
	}
	tests := []struct {
		name  string
		args  args
		want  *searchHistoryUseCase
		setup func(mockCtrl *gomock.Controller, tt *args) *searchHistoryUseCase
	}{
		{
			name: "positive testing",
			args: args{
				historyRepo: nil,
			},
			want: nil,
			setup: func(mockCtrl *gomock.Controller, tt *args) *searchHistoryUseCase {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				tt.historyRepo = mockHistoryRepo
				return &searchHistoryUseCase{
					historyRepo: mockHistoryRepo,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.want = tt.setup(mockCtrl, &tt.args)
			}
			if got := NewSearchHistoryUseCase(tt.args.historyRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSearchHistoryUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_searchHistoryUseCase_Run(t *testing.T) {
	type fields struct {
		historyRepo historyDomain.HistoryRepository
	}
	type args struct {
		ctx       context.Context
		keywords  []string
		and       bool
		all       bool
		favorited bool
		number    int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*SearchHistoryUseCaseOutputDto
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing (all and favorited)",
			fields: fields{
				historyRepo: nil,
			},
			args: args{
				ctx:       context.Background(),
				keywords:  []string{"test"},
				and:       true,
				all:       true,
				favorited: true,
			},
			want: []*SearchHistoryUseCaseOutputDto{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      "",
					Suffix:      "",
					IsFavorited: 1,
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      "",
					Suffix:      "",
					IsFavorited: 1,
				},
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				mockHistoryRepo.EXPECT().FindByIsFavoritedIsAndPhraseContains(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*historyDomain.History{
					{
						ID:          1,
						Phrase:      "test",
						Prefix:      sql.NullString{String: "", Valid: false},
						Suffix:      sql.NullString{String: "", Valid: false},
						IsFavorited: 1,
					},
					{
						ID:          2,
						Phrase:      "test2",
						Prefix:      sql.NullString{String: "", Valid: false},
						Suffix:      sql.NullString{String: "", Valid: false},
						IsFavorited: 1,
					},
				}, nil)
				tt.historyRepo = mockHistoryRepo
			},
		},
		{
			name: "positive testing (all and not favorited)",
			fields: fields{
				historyRepo: nil,
			},
			args: args{
				ctx:       context.Background(),
				keywords:  []string{"test"},
				and:       true,
				all:       true,
				favorited: false,
			},
			want: []*SearchHistoryUseCaseOutputDto{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      "",
					Suffix:      "",
					IsFavorited: 0,
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      "",
					Suffix:      "",
					IsFavorited: 0,
				},
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				mockHistoryRepo.EXPECT().FindByPhraseContains(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*historyDomain.History{
					{
						ID:          1,
						Phrase:      "test",
						Prefix:      sql.NullString{String: "", Valid: false},
						Suffix:      sql.NullString{String: "", Valid: false},
						IsFavorited: 0,
					},
					{
						ID:          2,
						Phrase:      "test2",
						Prefix:      sql.NullString{String: "", Valid: false},
						Suffix:      sql.NullString{String: "", Valid: false},
						IsFavorited: 0,
					},
				}, nil)
				tt.historyRepo = mockHistoryRepo
			},
		},
		{
			name: "positive testing (not all and favorited)",
			fields: fields{
				historyRepo: nil,
			},
			args: args{
				ctx:       context.Background(),
				keywords:  []string{"test"},
				and:       true,
				all:       false,
				favorited: true,
			},
			want: []*SearchHistoryUseCaseOutputDto{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      "",
					Suffix:      "",
					IsFavorited: 1,
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      "",
					Suffix:      "",
					IsFavorited: 1,
				},
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				mockHistoryRepo.EXPECT().FindTopNByIsFavoritedIsAndByPhraseContainsOrderByIdAsc(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*historyDomain.History{
					{
						ID:          1,
						Phrase:      "test",
						Prefix:      sql.NullString{String: "", Valid: false},
						Suffix:      sql.NullString{String: "", Valid: false},
						IsFavorited: 1,
					},
					{
						ID:          2,
						Phrase:      "test2",
						Prefix:      sql.NullString{String: "", Valid: false},
						Suffix:      sql.NullString{String: "", Valid: false},
						IsFavorited: 1,
					},
				}, nil)
				tt.historyRepo = mockHistoryRepo
			},
		},
		{
			name: "positive testing (not all and not favorited)",
			fields: fields{
				historyRepo: nil,
			},
			args: args{
				ctx:       context.Background(),
				keywords:  []string{"test"},
				and:       true,
				all:       false,
				favorited: false,
			},
			want: []*SearchHistoryUseCaseOutputDto{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      "",
					Suffix:      "",
					IsFavorited: 0,
				},
				{
					ID:          2,
					Phrase:      "test2",
					Prefix:      "",
					Suffix:      "",
					IsFavorited: 0,
				},
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				mockHistoryRepo.EXPECT().FindTopNByPhraseContainsOrderByIdAsc(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*historyDomain.History{
					{
						ID:          1,
						Phrase:      "test",
						Prefix:      sql.NullString{String: "", Valid: false},
						Suffix:      sql.NullString{String: "", Valid: false},
						IsFavorited: 0,
					},
					{
						ID:          2,
						Phrase:      "test2",
						Prefix:      sql.NullString{String: "", Valid: false},
						Suffix:      sql.NullString{String: "", Valid: false},
						IsFavorited: 0,
					},
				}, nil)
				tt.historyRepo = mockHistoryRepo
			},
		},
		{
			name: "negative testing (err != nil)",
			fields: fields{
				historyRepo: nil,
			},
			args: args{
				ctx:       context.Background(),
				keywords:  []string{"test"},
				and:       true,
				all:       true,
				favorited: true,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				mockHistoryRepo.EXPECT().FindByIsFavoritedIsAndPhraseContains(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				tt.historyRepo = mockHistoryRepo
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
			uc := &searchHistoryUseCase{
				historyRepo: tt.fields.historyRepo,
			}
			got, err := uc.Run(tt.args.ctx, tt.args.keywords, tt.args.and, tt.args.all, tt.args.favorited, tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("searchHistoryUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil {
				for i := range got {
					if got[i].ID != tt.want[i].ID {
						t.Errorf("getHistoryUseCase.Run() = %v, want %v", got[i].ID, tt.want[i].ID)
					}
					if got[i].Phrase != tt.want[i].Phrase {
						t.Errorf("getHistoryUseCase.Run() = %v, want %v", got[i].Phrase, tt.want[i].Phrase)
					}
					if got[i].Prefix != tt.want[i].Prefix {
						t.Errorf("getHistoryUseCase.Run() = %v, want %v", got[i].Prefix, tt.want[i].Prefix)
					}
					if got[i].Suffix != tt.want[i].Suffix {
						t.Errorf("getHistoryUseCase.Run() = %v, want %v", got[i].Suffix, tt.want[i].Suffix)
					}
					if got[i].IsFavorited != tt.want[i].IsFavorited {
						t.Errorf("getHistoryUseCase.Run() = %v, want %v", got[i].IsFavorited, tt.want[i].IsFavorited)
					}
				}
			}
		})
	}
}
