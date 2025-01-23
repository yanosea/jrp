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

func TestNewSaveHistoryUseCase(t *testing.T) {
	type args struct {
		historyRepo historyDomain.HistoryRepository
	}
	tests := []struct {
		name  string
		args  args
		want  *saveHistoryUseCase
		setup func(mockCtrl *gomock.Controller, tt *args) *saveHistoryUseCase
	}{
		{
			name: "positive testing",
			args: args{
				historyRepo: nil,
			},
			want: nil,
			setup: func(mockCtrl *gomock.Controller, tt *args) *saveHistoryUseCase {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				tt.historyRepo = mockHistoryRepo
				return &saveHistoryUseCase{
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
			if got := NewSaveHistoryUseCase(tt.args.historyRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSaveHistoryUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_saveHistoryUseCase_Run(t *testing.T) {
	type fields struct {
		historyRepo historyDomain.HistoryRepository
	}
	type args struct {
		ctx       context.Context
		inputDtos []*SaveHistoryUseCaseInputDto
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*SaveHistoryUseCaseOutputDto
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing",
			fields: fields{
				historyRepo: nil,
			},
			args: args{
				ctx: context.Background(),
				inputDtos: []*SaveHistoryUseCaseInputDto{
					{
						Phrase:      "test",
						Prefix:      "",
						Suffix:      "",
						IsFavorited: 0,
					},
					{
						Phrase:      "test2",
						Prefix:      "",
						Suffix:      "",
						IsFavorited: 0,
					},
				},
			},
			want: []*SaveHistoryUseCaseOutputDto{
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
				tt.historyRepo = mockHistoryRepo
				mockHistoryRepo.EXPECT().SaveAll(gomock.Any(), gomock.Any()).Return([]*historyDomain.History{
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
			},
		},
		{
			name: "negative testing",
			fields: fields{
				historyRepo: nil,
			},
			args: args{
				ctx: context.Background(),
				inputDtos: []*SaveHistoryUseCaseInputDto{
					{
						Phrase:      "test",
						Prefix:      "",
						Suffix:      "",
						IsFavorited: 0,
					},
				},
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				tt.historyRepo = mockHistoryRepo
				mockHistoryRepo.EXPECT().SaveAll(gomock.Any(), gomock.Any()).Return(nil, errors.New("HistoryRepository.SaveAll() error"))
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
			uc := &saveHistoryUseCase{
				historyRepo: tt.fields.historyRepo,
			}
			got, err := uc.Run(tt.args.ctx, tt.args.inputDtos)
			if (err != nil) != tt.wantErr {
				t.Errorf("saveHistoryUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && tt.want != nil {
				for i := range got {
					if got[i].ID != tt.want[i].ID {
						t.Errorf("saveHistoryUseCase.Run() = %v, want %v", got[i].ID, tt.want[i].ID)
					}
					if got[i].Phrase != tt.want[i].Phrase {
						t.Errorf("saveHistoryUseCase.Run() = %v, want %v", got[i].Phrase, tt.want[i].Phrase)
					}
					if got[i].Prefix != tt.want[i].Prefix {
						t.Errorf("saveHistoryUseCase.Run() = %v, want %v", got[i].Prefix, tt.want[i].Prefix)
					}
					if got[i].Suffix != tt.want[i].Suffix {
						t.Errorf("saveHistoryUseCase.Run() = %v, want %v", got[i].Suffix, tt.want[i].Suffix)
					}
					if got[i].IsFavorited != tt.want[i].IsFavorited {
						t.Errorf("saveHistoryUseCase.Run() = %v, want %v", got[i].IsFavorited, tt.want[i].IsFavorited)
					}
				}
			}
		})
	}
}
