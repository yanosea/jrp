package presenter

import (
	"testing"

	"github.com/yanosea/jrp/v2/pkg/proxy"
	"github.com/yanosea/jrp/v2/pkg/utility"

	"go.uber.org/mock/gomock"
)

func TestRunPrompt(t *testing.T) {
	origPu := Pu

	type args struct {
		label string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		setup   func(mockCtrl *gomock.Controller)
		cleanup func()
	}{
		{
			name: "positive testing",
			args: args{
				label: "test label",
			},
			want:    "test answer",
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller) {
				mockPrompt := proxy.NewMockPrompt(mockCtrl)
				mockPrompt.EXPECT().SetLabel("test label")
				mockPrompt.EXPECT().Run().Return("test answer", nil)
				mockPromptui := proxy.NewMockPromptui(mockCtrl)
				mockPromptui.EXPECT().NewPrompt().Return(mockPrompt)
				Pu = utility.NewPromptUtil(mockPromptui)
			},
			cleanup: func() {
				Pu = origPu
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
			got, err := RunPrompt(tt.args.label)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunPrompt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RunPrompt() = %v, want %v", got, tt.want)
			}
		})
	}
}
