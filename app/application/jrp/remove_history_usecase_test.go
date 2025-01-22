package jrp

import (
	"context"
	"errors"
	"reflect"
	"testing"

	historyDomain "github.com/yanosea/jrp/app/domain/jrp/history"
	"go.uber.org/mock/gomock"
)

func TestNewRemoveHistoryUseCase(t *testing.T) {
	type args struct {
		historyRepo historyDomain.HistoryRepository
	}
	tests := []struct {
		name  string
		args  args
		want  *removeHistoryUseCase
		setup func(mockCtrl *gomock.Controller, tt *args) *removeHistoryUseCase
	}{
		{
			name: "positive testing",
			args: args{
				historyRepo: nil,
			},
			want: nil,
			setup: func(mockCtrl *gomock.Controller, tt *args) *removeHistoryUseCase {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				tt.historyRepo = mockHistoryRepo
				return &removeHistoryUseCase{
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
			if got := NewRemoveHistoryUseCase(tt.args.historyRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRemoveHistoryUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeHistoryUseCase_Run(t *testing.T) {
	type fields struct {
		historyRepo historyDomain.HistoryRepository
	}
	type args struct {
		ctx   context.Context
		ids   []int
		all   bool
		force bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing (all and force)",
			args: args{
				ctx:   context.Background(),
				ids:   nil,
				all:   true,
				force: true,
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				mockHistoryRepo.EXPECT().DeleteAll(gomock.Any()).Return(1, nil)
				tt.historyRepo = mockHistoryRepo
			},
		},
		{
			name: "positive testing (all and not force)",
			args: args{
				ctx:   context.Background(),
				ids:   nil,
				all:   true,
				force: false,
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				mockHistoryRepo.EXPECT().DeleteByIsFavoritedIs(gomock.Any(), 0).Return(1, nil)
				tt.historyRepo = mockHistoryRepo
			},
		},
		{
			name: "positive testing (not all and force)",
			args: args{
				ctx:   context.Background(),
				ids:   []int{1, 2},
				all:   false,
				force: true,
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				mockHistoryRepo.EXPECT().DeleteByIdIn(gomock.Any(), []int{1, 2}).Return(1, nil)
				tt.historyRepo = mockHistoryRepo
			},
		},
		{
			name: "positive testing (not all and not force)",
			args: args{
				ctx:   context.Background(),
				ids:   []int{1, 2},
				all:   false,
				force: false,
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				mockHistoryRepo.EXPECT().DeleteByIdInAndIsFavoritedIs(gomock.Any(), []int{1, 2}, 0).Return(1, nil)
				tt.historyRepo = mockHistoryRepo
			},
		},
		{
			name: "negative testing (err != nil)",
			args: args{
				ctx:   context.Background(),
				ids:   nil,
				all:   true,
				force: true,
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				mockHistoryRepo.EXPECT().DeleteAll(gomock.Any()).Return(0, errors.New("HistoryRepository.DeleteAll() failed"))
				tt.historyRepo = mockHistoryRepo
			},
		},
		{
			name: "negative testing (rows affected == 0)",
			args: args{
				ctx:   context.Background(),
				ids:   nil,
				all:   true,
				force: true,
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockHistoryRepo := historyDomain.NewMockHistoryRepository(mockCtrl)
				mockHistoryRepo.EXPECT().DeleteAll(gomock.Any()).Return(0, nil)
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
			uc := &removeHistoryUseCase{
				historyRepo: tt.fields.historyRepo,
			}
			if err := uc.Run(tt.args.ctx, tt.args.ids, tt.args.all, tt.args.force); (err != nil) != tt.wantErr {
				t.Errorf("removeHistoryUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
