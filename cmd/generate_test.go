package cmd

import (
	"os"
	"testing"
)

func TestNewGenerateCommand(t *testing.T) {
	type args struct {
		globalOption *GlobalOption
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			// TODO : check the output
			name:    "positive testing",
			args:    args{globalOption: &GlobalOption{Out: os.Stdout, ErrOut: os.Stderr}},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = newGenerateCommand(tt.args.globalOption)
		})
	}
}

func TestGenerate(t *testing.T) {
	type args struct {
		generateOption *generateOption
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			// TODO : check the output
			name:    "positive testing",
			args:    args{},
			want:    0,
			wantErr: false,
		}, {
			// TODO : negative testing
			name:    "negative testing (japaneseRandomPhraseGenaretaer.Generate() fails)",
			args:    args{},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}
