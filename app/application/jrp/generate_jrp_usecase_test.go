package jrp

import (
	"reflect"
	"testing"

	"github.com/yanosea/jrp/v2/pkg/utility"

	"go.uber.org/mock/gomock"
)

func TestNewGenerateJrpUseCase(t *testing.T) {
	tests := []struct {
		name string
		want *generateJrpUseCase
	}{
		{
			name: "positive testing",
			want: &generateJrpUseCase{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGenerateJrpUseCase(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGenerateJrpUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateJrpUseCase_RunWithPrefix(t *testing.T) {
	origRu := ru

	type args struct {
		dtos   []*GenerateJrpUseCaseInputDto
		prefix string
	}
	tests := []struct {
		name    string
		uc      *generateJrpUseCase
		args    args
		want    *GenerateJrpUseCaseOutputDto
		setup   func(mockCtrl *gomock.Controller)
		cleanup func()
	}{
		{
			name: "positive testing (len(dtos) == 0)",
			args: args{
				dtos:   nil,
				prefix: "prefix",
			},
			want:    nil,
			setup:   nil,
			cleanup: nil,
		},
		{
			name: "positive testing",
			args: args{
				dtos: []*GenerateJrpUseCaseInputDto{
					{
						WordID: 1,
						Lang:   "jpn",
						Lemma:  "testa",
						Pron:   "test",
						Pos:    "a",
					},
					{
						WordID: 2,
						Lang:   "jpn",
						Lemma:  "testv",
						Pron:   "test",
						Pos:    "v",
					},
					{
						WordID: 3,
						Lang:   "jpn",
						Lemma:  "testn",
						Pron:   "test",
						Pos:    "n",
					},
				},
				prefix: "prefix",
			},
			want: &GenerateJrpUseCaseOutputDto{
				ID:          0,
				Phrase:      "prefixtestn",
				Prefix:      "prefix",
				Suffix:      "",
				IsFavorited: 0,
			},
			setup: func(mockCtrl *gomock.Controller) {
				mockRu := utility.NewMockRandUtil(mockCtrl)
				mockRu.EXPECT().GenerateRandomNumber(3).Return(0)
				mockRu.EXPECT().GenerateRandomNumber(3).Return(1)
				mockRu.EXPECT().GenerateRandomNumber(3).Return(2)
				ru = mockRu
			},
			cleanup: func() {
				ru = origRu
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
			uc := &generateJrpUseCase{}
			got := uc.RunWithPrefix(tt.args.dtos, tt.args.prefix)
			if (got == nil) != (tt.want == nil) {
				t.Errorf("generateJrpUseCase.RunWithPrefix() nil check failed: got %v, want %v", got, tt.want)
				return
			}
			if got != nil && tt.want != nil {
				if got.ID != tt.want.ID {
					t.Errorf("generateJrpUseCase.RunWithPrefix() = %v, want %v", got.ID, tt.want.ID)
				}
				if got.Phrase != tt.want.Phrase {
					t.Errorf("generateJrpUseCase.RunWithPrefix() = %v, want %v", got.Phrase, tt.want.Phrase)
				}
				if got.Prefix != tt.want.Prefix {
					t.Errorf("generateJrpUseCase.RunWithPrefix() = %v, want %v", got.Prefix, tt.want.Prefix)
				}
				if got.Suffix != tt.want.Suffix {
					t.Errorf("generateJrpUseCase.RunWithPrefix() = %v, want %v", got.Suffix, tt.want.Suffix)
				}
				if got.IsFavorited != tt.want.IsFavorited {
					t.Errorf("generateJrpUseCase.RunWithPrefix() = %v, want %v", got.IsFavorited, tt.want.IsFavorited)
				}
			}
		})
	}
}

func Test_generateJrpUseCase_RunWithSuffix(t *testing.T) {
	origRu := ru

	type args struct {
		dtos   []*GenerateJrpUseCaseInputDto
		suffix string
	}
	tests := []struct {
		name    string
		uc      *generateJrpUseCase
		args    args
		want    *GenerateJrpUseCaseOutputDto
		setup   func(mockCtrl *gomock.Controller)
		cleanup func()
	}{
		{
			name: "positive testing (len(dtos) == 0)",
			args: args{
				dtos:   nil,
				suffix: "suffix",
			},
			want:    nil,
			setup:   nil,
			cleanup: nil,
		},
		{
			name: "positive testing",
			args: args{
				dtos: []*GenerateJrpUseCaseInputDto{
					{
						WordID: 1,
						Lang:   "jpn",
						Lemma:  "testn",
						Pron:   "test",
						Pos:    "n",
					},
					{
						WordID: 2,
						Lang:   "jpn",
						Lemma:  "testa",
						Pron:   "test",
						Pos:    "a",
					},
				},
				suffix: "suffix",
			},
			want: &GenerateJrpUseCaseOutputDto{
				ID:          0,
				Phrase:      "testasuffix",
				Prefix:      "",
				Suffix:      "suffix",
				IsFavorited: 0,
			},
			setup: func(mockCtrl *gomock.Controller) {
				mockRu := utility.NewMockRandUtil(mockCtrl)
				mockRu.EXPECT().GenerateRandomNumber(2).Return(0)
				mockRu.EXPECT().GenerateRandomNumber(2).Return(1)
				ru = mockRu
			},
			cleanup: func() {
				ru = origRu
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
			uc := &generateJrpUseCase{}
			got := uc.RunWithSuffix(tt.args.dtos, tt.args.suffix)
			if (got == nil) != (tt.want == nil) {
				t.Errorf("generateJrpUseCase.RunWithPrefix() = %v, want not nil", got)
			}
			if got != nil && tt.want != nil {
				if got.ID != tt.want.ID {
					t.Errorf("generateJrpUseCase.RunWithSuffix() = %v, want %v", got.ID, tt.want.ID)
				}
				if got.Phrase != tt.want.Phrase {
					t.Errorf("generateJrpUseCase.RunWithSuffix() = %v, want %v", got.Phrase, tt.want.Phrase)
				}
				if got.Suffix != tt.want.Suffix {
					t.Errorf("generateJrpUseCase.RunWithSuffix() = %v, want %v", got.Suffix, tt.want.Suffix)
				}
				if got.Suffix != tt.want.Suffix {
					t.Errorf("generateJrpUseCase.RunWithSuffix() = %v, want %v", got.Suffix, tt.want.Suffix)
				}
				if got.IsFavorited != tt.want.IsFavorited {
					t.Errorf("generateJrpUseCase.RunWithSuffix() = %v, want %v", got.IsFavorited, tt.want.IsFavorited)
				}
			}
		})
	}
}

func Test_generateJrpUseCase_RunWithRandom(t *testing.T) {
	origRu := ru

	type args struct {
		dtos []*GenerateJrpUseCaseInputDto
	}
	tests := []struct {
		name    string
		uc      *generateJrpUseCase
		args    args
		want    *GenerateJrpUseCaseOutputDto
		setup   func(mockCtrl *gomock.Controller)
		cleanup func()
	}{
		{
			name: "positive testing (len(dtos) == 0)",
			args: args{
				dtos: nil,
			},
			want:    nil,
			setup:   nil,
			cleanup: nil,
		},
		{
			name: "positive testing",
			args: args{
				dtos: []*GenerateJrpUseCaseInputDto{
					{
						WordID: 1,
						Lang:   "jpn",
						Lemma:  "testn",
						Pron:   "test",
						Pos:    "n",
					},
					{
						WordID: 2,
						Lang:   "jpn",
						Lemma:  "testa",
						Pron:   "test",
						Pos:    "a",
					},
					{
						WordID: 3,
						Lang:   "jpn",
						Lemma:  "testv",
						Pron:   "test",
						Pos:    "v",
					},
					{
						WordID: 4,
						Lang:   "jpn",
						Lemma:  "testa",
						Pron:   "test",
						Pos:    "a",
					},
					{
						WordID: 5,
						Lang:   "jpn",
						Lemma:  "testa",
						Pron:   "test",
						Pos:    "a",
					},
					{
						WordID: 6,
						Lang:   "jpn",
						Lemma:  "testn",
						Pron:   "test",
						Pos:    "n",
					},
				},
			},
			want: &GenerateJrpUseCaseOutputDto{
				ID:          0,
				Phrase:      "testatestn",
				Prefix:      "",
				Suffix:      "",
				IsFavorited: 0,
			},
			setup: func(mockCtrl *gomock.Controller) {
				mockRu := utility.NewMockRandUtil(mockCtrl)
				mockRu.EXPECT().GenerateRandomNumber(6).Return(0)
				mockRu.EXPECT().GenerateRandomNumber(6).Return(1)
				mockRu.EXPECT().GenerateRandomNumber(6).Return(2)
				mockRu.EXPECT().GenerateRandomNumber(6).Return(3)
				mockRu.EXPECT().GenerateRandomNumber(6).Return(4)
				mockRu.EXPECT().GenerateRandomNumber(6).Return(5)
				ru = mockRu
			},
			cleanup: func() {
				ru = origRu
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
			uc := &generateJrpUseCase{}
			got := uc.RunWithRandom(tt.args.dtos)
			if (got == nil) != (tt.want == nil) {
				t.Errorf("generateJrpUseCase.RunWithPrefix() = %v, want not nil", got)
			}
			if got != nil && tt.want != nil {
				if got.ID != tt.want.ID {
					t.Errorf("generateJrpUseCase.RunWithRandom() = %v, want %v", got.ID, tt.want.ID)
				}
				if got.Phrase != tt.want.Phrase {
					t.Errorf("generateJrpUseCase.RunWithRandom() = %v, want %v", got.Phrase, tt.want.Phrase)
				}
				if got.Prefix != tt.want.Prefix {
					t.Errorf("generateJrpUseCase.RunWithRandom() = %v, want %v", got.Prefix, tt.want.Prefix)
				}
				if got.Suffix != tt.want.Suffix {
					t.Errorf("generateJrpUseCase.RunWithRandom() = %v, want %v", got.Suffix, tt.want.Suffix)
				}
				if got.IsFavorited != tt.want.IsFavorited {
					t.Errorf("generateJrpUseCase.RunWithRandom() = %v, want %v", got.IsFavorited, tt.want.IsFavorited)
				}
			}
		})
	}
}
