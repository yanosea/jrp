package cmd_test

import (
	"errors"
	"fmt"
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
		cmdArgs []string
		want    int
		wantErr bool
		setup   func(mockCmd *mock_cmdwrapper.MockICommand)
	}{
		{
			name:    "positive testing",
			args:    args{globalOption: cmd.NewGlobalOption(os.Stdout, os.Stderr, []string{})},
			cmdArgs: []string{},
			want:    0,
			wantErr: false,
			setup:   nil,
		},
		{
			name:    "positive testing with args",
			args:    args{globalOption: cmd.NewGlobalOption(os.Stdout, os.Stderr, []string{"testArg"})},
			cmdArgs: []string{"testArg"},
			want:    0,
			wantErr: false,
			setup:   nil,
		},
		{
			name:    "negative testing (rootCmd.Execute() fails)",
			args:    args{globalOption: cmd.NewGlobalOption(os.Stdout, os.Stderr, []string{})},
			cmdArgs: []string{},
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

			if tt.setup != nil {
				mockCmd := mock_cmdwrapper.NewMockICommand(ctrl)
				tt.setup(mockCmd)
				tt.args.globalOption.NewRootCommand = func(ow, ew io.Writer, cmdArgs []string) cmdwrapper.ICommand {
					return mockCmd
				}
			}

			fmt.Println("globalOption: args")
			fmt.Println(tt.args.globalOption.Args)

			if got := tt.args.globalOption.Execute(); (got != 0) != tt.wantErr {
				t.Errorf("Execute() = %v, want = %v", got, tt.want)
			}
		})
	}
}
