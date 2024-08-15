package cmd_test

import (
	"errors"
	"os"
	"testing"

	gomock "go.uber.org/mock/gomock"

	"github.com/yanosea/jrp/cmd"
	"github.com/yanosea/jrp/internal/database"
	"github.com/yanosea/jrp/internal/fs"
	"github.com/yanosea/jrp/internal/usermanager"
	"github.com/yanosea/jrp/logic"

	mock_generator "github.com/yanosea/jrp/mock/generator"
)

func TestNewGenerateCommand(t *testing.T) {
	type args struct {
		globalOption *cmd.GlobalOption
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive testing",
			args: args{globalOption: &cmd.GlobalOption{Out: os.Stdout, ErrOut: os.Stderr, Args: []string{"generate", "0"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cmd.NewGenerateCommand(tt.args.globalOption)
			if got == nil {
				t.Errorf("NewDownloadCommand() : returned nil")
			}
		})
	}
}

func TestGenerateRunE(t *testing.T) {
	type args struct {
		o cmd.GenerateOption
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive testing",
			args: args{o: cmd.GenerateOption{Out: os.Stdout, ErrOut: os.Stderr, Args: []string{"generate", "0"}, Number: 0, Generator: logic.NewJapaneseRandomPhraseGenerator(usermanager.OSUserProvider{}, database.SQLiteProvider{}, fs.OsFileManager{})}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			err := tt.args.o.GenerateRunE(nil, nil)
			if err != nil {
				t.Errorf("GenerateRunE() : error = %v", err)
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	type args struct {
		o cmd.GenerateOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *args)
	}{
		{
			name:    "positive testing",
			args:    args{o: cmd.GenerateOption{Out: os.Stdout, ErrOut: os.Stderr, Args: []string{"generate", "0"}, Number: 1, Generator: logic.NewJapaneseRandomPhraseGenerator(usermanager.OSUserProvider{}, database.SQLiteProvider{}, fs.OsFileManager{})}},
			wantErr: false,
			setup:   nil,
		}, {
			name:    "negative testing (Generate() fails)",
			args:    args{o: cmd.GenerateOption{Out: os.Stdout, ErrOut: os.Stderr, Args: []string{"generate", "0"}, Number: 1, Generator: nil}},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mg := mock_generator.NewMockGenerator(mockCtrl)
				mg.EXPECT().Generate(tt.o.Number).Return(nil, errors.New("failed to generate japanese random phrase"))
				tt.o.Generator = mg
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			if tt.setup != nil {
				tt.setup(ctrl, &tt.args)
			}

			err := tt.args.o.Generate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() : error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
