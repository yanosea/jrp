package presenter

import (
	"os"
	"testing"

	"github.com/yanosea/jrp/v2/pkg/proxy"
	"github.com/yanosea/jrp/v2/pkg/utility"
)

func TestPrint(t *testing.T) {
	stdBuffer := proxy.NewBuffer()
	errBuffer := proxy.NewBuffer()

	type fields struct {
		StdBuffer proxy.Buffer
		ErrBuffer proxy.Buffer
	}
	type args struct {
		fnc func()
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStdOut string
		wantStdErr string
		wantErr    bool
	}{
		{
			name: "positive testing (stdout, not \\n)",
			fields: fields{
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					if err := Print(os.Stdout, "test stdout"); err != nil {
						t.Errorf("Print() error = %v", err)
					}
				},
			},
			wantStdOut: "test stdout\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (stdout, \\n)",
			fields: fields{
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					if err := Print(os.Stdout, ""); err != nil {
						t.Errorf("Print() error = %v", err)
					}
				},
			},
			wantStdOut: "\n",
			wantStdErr: "",
			wantErr:    false,
		},
		{
			name: "positive testing (stderr, not \\n)",
			fields: fields{
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					if err := Print(os.Stderr, "test stderr"); err != nil {
						t.Errorf("Print() error = %v", err)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "test stderr\n",
			wantErr:    false,
		},
		{
			name: "positive testing (stderr, \\n)",
			fields: fields{
				StdBuffer: stdBuffer,
				ErrBuffer: errBuffer,
			},
			args: args{
				fnc: func() {
					if err := Print(os.Stderr, ""); err != nil {
						t.Errorf("Print() error = %v", err)
					}
				},
			},
			wantStdOut: "",
			wantStdErr: "\n",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := utility.NewCapturer(tt.fields.StdBuffer, tt.fields.ErrBuffer)
			gotStdOut, gotStdErr, err := c.CaptureOutput(tt.args.fnc)
			if (err != nil) != tt.wantErr {
				t.Errorf("Capturer.CaptureOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStdOut != tt.wantStdOut {
				t.Errorf("Print() gotStdOut = %v, want %v", gotStdOut, tt.wantStdOut)
			}
			if gotStdErr != tt.wantStdErr {
				t.Errorf("Print() gotStdErr = %v, want %v", gotStdErr, tt.wantStdErr)
			}
		})
	}
}
