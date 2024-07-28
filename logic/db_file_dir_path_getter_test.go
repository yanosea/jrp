package logic

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/yanosea/jrp/constant"
)

type MockUserProvider struct{}

func (m MockUserProvider) Current() (*user.User, error) {
	return nil, errors.New("mock error : Current() failed")
}

func TestGetDBFileDirPath(t *testing.T) {
	var testUser, _ = user.Current()
	tests := []struct {
		name         string
		wordNetJpDir string
		want         string
		wantErr      bool
	}{
		{
			name:         "positive testing (no env)",
			wordNetJpDir: "",
			want:         filepath.Join(testUser.HomeDir, ".local", "share", "jrp"),
			wantErr:      false,
		}, {
			name:         "positive testing (with env)",
			wordNetJpDir: filepath.Join(testUser.HomeDir, "jrp"),
			want:         filepath.Join(testUser.HomeDir, "jrp"),
			wantErr:      false,
		}, {
			name:         "negative testing (user.Current() fails)",
			wordNetJpDir: "",
			want:         "",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		provider := DefaultUserProvider{}
		os.Unsetenv(constant.JRP_ENV_WORDNETJP_DIR)
		t.Run(tt.name, func(t *testing.T) {
			if tt.wordNetJpDir != "" {
				os.Setenv(constant.JRP_ENV_WORDNETJP_DIR, tt.wordNetJpDir)
				defer os.Unsetenv(constant.JRP_ENV_WORDNETJP_DIR)
			}
			got, err := GetDBFileDirPath(provider)
			if err != nil && !tt.wantErr && got != tt.want {
				t.Errorf("GetDBFileDirPath() = %v, want %v", got, tt.want)
				return
			}
			if tt.wantErr {
				mockProvider := MockUserProvider{}
				_, err := GetDBFileDirPath(mockProvider)
				if err == nil {
					t.Error("Expected error when user.Current() fails, but got nil")
				}
			}
		})
	}
}
