package logic

import "testing"

func TestDownload(t *testing.T) {
	type args struct {
		e  Env
		u  User
		fs FileSystem
		hc HttpClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Download(tt.args.e, tt.args.u, tt.args.fs, tt.args.hc); (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
