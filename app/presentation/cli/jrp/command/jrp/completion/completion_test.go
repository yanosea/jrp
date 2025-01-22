package completion

import (
	"testing"

	"github.com/yanosea/jrp/pkg/proxy"
)

func TestNewCompletionCommand(t *testing.T) {
	type args struct {
		cobra  proxy.Cobra
		output *string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive testing",
			args: args{
				cobra: proxy.NewCobra(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCompletionCommand(tt.args.cobra, tt.args.output)
			if got == nil {
				t.Errorf("NewCompletionCommand() = %v, want not nil", got)
			}
		})
	}
}
