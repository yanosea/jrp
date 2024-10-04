package promptuiproxy

import (
	"testing"

	"github.com/manifoldco/promptui"
)

func TestPromptInstance_Run(t *testing.T) {
	type fields struct {
		FieldPrompt *promptui.Prompt
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "positive testing",
			fields: fields{
				FieldPrompt: &promptui.Prompt{},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PromptInstance{
				FieldPrompt: tt.fields.FieldPrompt,
			}
			got, err := p.Run()
			if (err != nil) != tt.wantErr {
				t.Errorf("PromptInstance.Run() : error =\n%v, wantErr =\n%v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PromptInstance.Run() : got =\n%v, want =\n%v", got, tt.want)
			}
		})
	}
}

func TestPromptInstance_SetLabel(t *testing.T) {
	type fields struct {
		FieldPrompt *promptui.Prompt
	}
	type args struct {
		label string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "positive testing",
			fields: fields{
				FieldPrompt: &promptui.Prompt{},
			},
			args: args{
				label: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PromptInstance{
				FieldPrompt: tt.fields.FieldPrompt,
			}
			p.SetLabel(tt.args.label)
		})
	}
}
