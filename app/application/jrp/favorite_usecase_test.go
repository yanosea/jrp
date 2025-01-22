package jrp

import (
	"context"
	"errors"
	"reflect"
	"testing"

	historyDomain "github.com/yanosea/jrp/app/domain/jrp/history"

	"go.uber.org/mock/gomock"
)

func TestNewFavoriteUseCase(t *testing.T) {
	type args struct {
		historyRepo historyDomain.HistoryRepository
	}
	tests := []struct {
		name  string
		args  args
		want  *favoriteUseCase
		setup func(mockCtrl *gomock.Controller, tt *args) *favoriteUseCase
	}{
		{
			name: "positive testing",
			args: args{
				historyRepo: nil,
			},
			want: nil,
			setup: func(mockCtrl *gomock.Controller, tt *args) *favoriteUseCase {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				tt.historyRepo = mockHistoryRepo
				return &favoriteUseCase{
					historyRepo: mockHistoryRepo,
				}
			},
		},
	}
	for _, tt := range tests {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		if tt.setup != nil {
			tt.want = tt.setup(mockCtrl, &tt.args)
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFavoriteUseCase(tt.args.historyRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFavoriteUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_favoriteUseCase_Run(t *testing.T) {
	type fields struct {
		historyRepo historyDomain.HistoryRepository
	}
	type args struct {
		ctx context.Context
		ids []int
		all bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing (all)",
			fields: fields{
				historyRepo: nil,
			},
			args: args{
				ctx: context.Background(),
				ids: nil,
				all: true,
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				mockHistoryRepo.EXPECT().UpdateIsFavoritedByIsFavoritedIs(gomock.Any(), 1, 0).Return(1, nil)
				tt.historyRepo = mockHistoryRepo
			},
		},
		{
			name: "positive testing (not all)",
			fields: fields{
				historyRepo: nil,
			},
			args: args{
				ctx: context.Background(),
				ids: []int{1, 2},
				all: false,
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				mockHistoryRepo.EXPECT().UpdateIsFavoritedByIdIn(gomock.Any(), 1, []int{1, 2}).Return(2, nil)
				tt.historyRepo = mockHistoryRepo
			},
		},
		{
			name: "negative testing (err != nil)",
			fields: fields{
				historyRepo: nil,
			},
			args: args{
				ctx: context.Background(),
				ids: []int{1, 2},
				all: false,
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				mockHistoryRepo.EXPECT().UpdateIsFavoritedByIdIn(gomock.Any(), 1, []int{1, 2}).Return(0, errors.New("HistoryRepository.UpdateIsFavoritedByIdIn() failed"))
				tt.historyRepo = mockHistoryRepo
			},
		},
		{
			name: "negative testing (rowsAffected == 0)",
			fields: fields{
				historyRepo: nil,
			},
			args: args{
				ctx: context.Background(),
				ids: []int{1, 2},
				all: false,
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				mockHistoryRepo.EXPECT().UpdateIsFavoritedByIdIn(gomock.Any(), 1, []int{1, 2}).Return(0, nil)
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
			uc := &favoriteUseCase{
				historyRepo: tt.fields.historyRepo,
			}
			if err := uc.Run(tt.args.ctx, tt.args.ids, tt.args.all); (err != nil) != tt.wantErr {
				t.Errorf("favoriteUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
