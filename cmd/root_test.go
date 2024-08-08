package cmd_test

import (
	"errors"
	"io"
	"os"
	"testing"

	gomock "go.uber.org/mock/gomock"

	"github.com/yanosea/jrp/cmd"
	"github.com/yanosea/jrp/internal/cmdwrapper"
	"github.com/yanosea/jrp/mock/cmdwrapper"
)

func TestExecute(t *testing.T) {
	type args struct {
		globalOption *cmd.GlobalOption
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
		setup   func(mockCmd *mock_cmdwrapper.MockICommand)
	}{
		{
			name:    "positive testing",
			args:    args{globalOption: cmd.NewGlobalOption(os.Stdout, os.Stderr)},
			want:    0,
			wantErr: false,
			setup: func(mockCmd *mock_cmdwrapper.MockICommand) {
				mockCmd.EXPECT().Execute().Return(nil)
			},
		},
		{
			name:    "negative testing (rootCmd.Execute() fails)",
			args:    args{globalOption: cmd.NewGlobalOption(os.Stdout, os.Stderr)},
			want:    1,
			wantErr: true,
			setup: func(mockCmd *mock_cmdwrapper.MockICommand) {
				mockCmd.EXPECT().Execute().Return(errors.New("failed to execute command"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCmd := mock_cmdwrapper.NewMockICommand(ctrl)
			if tt.setup != nil {
				tt.setup(mockCmd)
			}

			tt.args.globalOption.NewRootCommand = func(ow, ew io.Writer) cmdwrapper.ICommand {
				return mockCmd
			}

			if got := tt.args.globalOption.Execute(); (got != 0) != tt.wantErr {
				t.Errorf("Execute() = %v, want = %v", got, tt.want)
			}
		})
	}
}
