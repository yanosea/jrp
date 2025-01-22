package wnjpn

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"
)

func Test_newFetchWordsUseCase(t *testing.T) {
	type args struct {
		wordQueryService WordQueryService
	}
	tests := []struct {
		name  string
		args  args
		want  *FetchWordsUseCaseStruct
		setup func(mockCtrl *gomock.Controller, tt *args) *FetchWordsUseCaseStruct
	}{
		{
			name: "positive testing",
			args: args{
				wordQueryService: nil,
			},
			want: &FetchWordsUseCaseStruct{
				wordQueryService: nil,
			},
			setup: func(mockCtrl *gomock.Controller, tt *args) *FetchWordsUseCaseStruct {
				mockWordQueryService := NewMockWordQueryService(mockCtrl)
				tt.wordQueryService = mockWordQueryService
				return newFetchWordsUseCase(mockWordQueryService)
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
			if got := newFetchWordsUseCase(tt.args.wordQueryService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFetchWordsUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFetchWordsUseCaseStruct_Run(t *testing.T) {
	type fields struct {
		wordQueryService WordQueryService
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
		want    []*FetchWordsUseCaseOutputDto
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing",
			fields: fields{
				wordQueryService: nil,
			},
			args: args{
				ctx:  context.Background(),
				lang: "jpn",
				pos:  []string{"a", "v", "n"},
			},
			want: []*FetchWordsUseCaseOutputDto{
				{
					WordID: 1,
					Lang:   "jpn",
					Lemma:  "test",
					Pron:   "test",
					Pos:    "a",
				},
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockWordQueryService := NewMockWordQueryService(mockCtrl)
				mockWordQueryService.EXPECT().FindByLangIsAndPosIn(gomock.Any(), "jpn", []string{"a", "v", "n"}).Return([]*FetchWordsDto{
					{
						WordID: 1,
						Lang:   sql.NullString{String: "jpn", Valid: true},
						Lemma:  sql.NullString{String: "test", Valid: true},
						Pron:   sql.NullString{String: "test", Valid: true},
						Pos:    sql.NullString{String: "a", Valid: true},
					},
				}, nil)
				tt.wordQueryService = mockWordQueryService
			},
		},
		{
			name: "negative testing (uc.wordQueryService.FindByLangIsAndPosIn(ctx, lang, pos) failed)",
			fields: fields{
				wordQueryService: nil,
			},
			args: args{
				ctx:  context.Background(),
				lang: "jpn",
				pos:  []string{"a", "v", "n"},
			},
			want:    nil,
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockWordQueryService := NewMockWordQueryService(mockCtrl)
				mockWordQueryService.EXPECT().FindByLangIsAndPosIn(gomock.Any(), "jpn", []string{"a", "v", "n"}).Return(nil, errors.New("WordQueryService.FindByLangIsAndPosIn() failed"))
				tt.wordQueryService = mockWordQueryService
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
			uc := &FetchWordsUseCaseStruct{
				wordQueryService: tt.fields.wordQueryService,
			}
			got, err := uc.Run(tt.args.ctx, tt.args.lang, tt.args.pos)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchWordsUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fetchWordsUseCase.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
