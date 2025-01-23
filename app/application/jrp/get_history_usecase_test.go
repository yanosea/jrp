package jrp

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"

	historyDomain "github.com/yanosea/jrp/v2/app/domain/jrp/history"

	"go.uber.org/mock/gomock"
)

func TestNewGetHistoryUseCase(t *testing.T) {
	type args struct {
		historyRepo historyDomain.HistoryRepository
	}
	tests := []struct {
		name  string
		args  args
		want  *getHistoryUseCase
		setup func(mockCtrl *gomock.Controller, tt *args) *getHistoryUseCase
	}{
		{
			name: "positive testing",
			args: args{
				historyRepo: nil,
			},
			want: nil,
			setup: func(mockCtrl *gomock.Controller, tt *args) *getHistoryUseCase {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				tt.historyRepo = mockHistoryRepo
				return &getHistoryUseCase{
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
			if got := NewGetHistoryUseCase(tt.args.historyRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGetHistoryUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getHistoryUseCase_Run(t *testing.T) {
	type fields struct {
		historyRepo historyDomain.HistoryRepository
	}
	type args struct {
		ctx       context.Context
		all       bool
		favorited bool
		number    int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*GetHistoryUseCaseOutputDto
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
				all:       true,
				favorited: true,
				number:    0,
			},
			want: []*GetHistoryUseCaseOutputDto{
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
				mockHistoryRepo.EXPECT().FindByIsFavoritedIs(gomock.Any(), 1).Return([]*historyDomain.History{
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
				all:       true,
				favorited: false,
				number:    0,
			},
			want: []*GetHistoryUseCaseOutputDto{
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
				mockHistoryRepo.EXPECT().FindAll(gomock.Any()).Return([]*historyDomain.History{
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
				all:       false,
				favorited: true,
				number:    1,
			},
			want: []*GetHistoryUseCaseOutputDto{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      "",
					Suffix:      "",
					IsFavorited: 1,
				},
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				mockHistoryRepo.EXPECT().FindTopNByIsFavoritedIsAndByOrderByIdAsc(gomock.Any(), 1, 1).Return([]*historyDomain.History{
					{
						ID:          1,
						Phrase:      "test",
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
				all:       false,
				favorited: false,
				number:    1,
			},
			want: []*GetHistoryUseCaseOutputDto{
				{
					ID:          1,
					Phrase:      "test",
					Prefix:      "",
					Suffix:      "",
					IsFavorited: 0,
				},
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				mockHistoryRepo.EXPECT().FindTopNByOrderByIdAsc(gomock.Any(), 1).Return([]*historyDomain.History{
					{
						ID:          1,
						Phrase:      "test",
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
				all:       true,
				favorited: true,
				number:    0,
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				mockHistoryRepo.EXPECT().FindByIsFavoritedIs(gomock.Any(), 1).Return(nil, errors.New("HistoryRepository.FindByIsFavoritedIs() failed"))
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
			uc := &getHistoryUseCase{
				historyRepo: tt.fields.historyRepo,
			}
			got, err := uc.Run(tt.args.ctx, tt.args.all, tt.args.favorited, tt.args.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("getHistoryUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
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
