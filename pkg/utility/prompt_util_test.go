package utility

import (
	"reflect"
	"testing"

	"github.com/yanosea/jrp/pkg/proxy"

	"go.uber.org/mock/gomock"
)

func TestNewPromptUtil(t *testing.T) {
	promptui := proxy.NewPromptui()

	type args struct {
		promptui proxy.Promptui
	}
	tests := []struct {
		name string
		args args
		want PromptUtil
	}{
		{
			name: "positive testing",
			args: args{
				promptui: promptui,
			},
			want: &promptUtil{
				promptui: promptui,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPromptUtil(tt.args.promptui); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPromptUtil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_promptUtil_GetPrompt(t *testing.T) {
	type fields struct {
		promptui proxy.Promptui
	}
	type args struct {
		label string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   proxy.Prompt
		setup  func(mockCtrl *gomock.Controller, tt *fields) proxy.Prompt
	}{
		{
			name: "positive testing",
			fields: fields{
				promptui: nil,
			},
			args: args{
				label: "test",
			},
			want: nil,
			setup: func(mockCtrl *gomock.Controller, tt *fields) proxy.Prompt {
				mockPromptui := proxy.NewMockPromptui(mockCtrl)
				mockPrompt := proxy.NewMockPrompt(mockCtrl)
				mockPromptui.EXPECT().NewPrompt().Return(mockPrompt)
				mockPrompt.EXPECT().SetLabel("test")
				tt.promptui = mockPromptui
				return mockPrompt
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.want = tt.setup(mockCtrl, &tt.fields)
			}
			p := &promptUtil{
				promptui: tt.fields.promptui,
			}
			if got := p.GetPrompt(tt.args.label); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("promptUtil.GetPrompt() = %v, want %v", got, tt.want)
			}
		})
	}
}
