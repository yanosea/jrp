package jrp

import (
	"context"
	"errors"
	"reflect"
	"testing"

	historyDomain "github.com/yanosea/jrp/v2/app/domain/jrp/history"

	"go.uber.org/mock/gomock"
)

func TestNewUnfavoriteUseCase(t *testing.T) {
	type args struct {
		historyRepo historyDomain.HistoryRepository
	}
	tests := []struct {
		name  string
		args  args
		want  *unfavoriteUseCase
		setup func(mockCtrl *gomock.Controller, tt *args) *unfavoriteUseCase
	}{
		{
			name: "positive testing",
			args: args{
				historyRepo: nil,
			},
			want: nil,
			setup: func(mockCtrl *gomock.Controller, tt *args) *unfavoriteUseCase {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				tt.historyRepo = mockHistoryRepo
				return &unfavoriteUseCase{
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
			if got := NewUnfavoriteUseCase(tt.args.historyRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUnfavoriteUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_unfavoriteUseCase_Run(t *testing.T) {
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
				mockHistoryRepo.EXPECT().UpdateIsFavoritedByIsFavoritedIs(gomock.Any(), 0, 1).Return(1, nil)
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
				mockHistoryRepo.EXPECT().UpdateIsFavoritedByIdIn(gomock.Any(), 0, []int{1, 2}).Return(1, nil)
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
				ids: nil,
				all: true,
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				mockHistoryRepo.EXPECT().UpdateIsFavoritedByIsFavoritedIs(gomock.Any(), 0, 1).Return(0, errors.New("error"))
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
				ids: nil,
				all: true,
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				mockHistoryRepo.EXPECT().UpdateIsFavoritedByIsFavoritedIs(gomock.Any(), 0, 1).Return(0, nil)
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
			uc := &unfavoriteUseCase{
				historyRepo: tt.fields.historyRepo,
			}
			if err := uc.Run(tt.args.ctx, tt.args.ids, tt.args.all); (err != nil) != tt.wantErr {
				t.Errorf("unfavoriteUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
