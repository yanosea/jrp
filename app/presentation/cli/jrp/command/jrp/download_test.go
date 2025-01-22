package jrp

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/fatih/color"

	jrpApp "github.com/yanosea/jrp/app/application/jrp"
	baseConfig "github.com/yanosea/jrp/app/config"
	"github.com/yanosea/jrp/app/presentation/cli/jrp/config"
	"github.com/yanosea/jrp/app/presentation/cli/jrp/presenter"

	"github.com/yanosea/jrp/pkg/proxy"
	"github.com/yanosea/jrp/pkg/utility"

	"go.uber.org/mock/gomock"
)

func TestNewDownloadCommand(t *testing.T) {
	type args struct {
		cobra  proxy.Cobra
		conf   *config.JrpCliConfig
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
				conf: &config.JrpCliConfig{
					JrpConfig: baseConfig.JrpConfig{
						JrpDBType:   "sqlite",
						JrpDBDsn:    filepath.Join(os.TempDir(), "jrp.db"),
						WNJpnDBType: "sqlite",
						WNJpnDBDsn:  filepath.Join(os.TempDir(), "wnjpn.db"),
					},
				},
				output: new(string),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDownloadCommand(tt.args.cobra, tt.args.conf, tt.args.output)
			if got == nil {
				t.Errorf("NewDownloadCommand() = %v, want not nil", got)
			} else {
				if err := got.RunE(nil, []string{}); err != nil {
					t.Errorf("Failed to run the download command: %v", err)
				}
			}
		})
	}
}

func Test_runDownload(t *testing.T) {
	var output string
	var hiddenFilePath string
	fu := utility.NewFileUtil(proxy.NewGzip(), proxy.NewIo(), proxy.NewOs())
	origSu := presenter.Su
	origDu := jrpApp.Du

	type args struct {
		conf   *config.JrpCliConfig
		output *string
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
			name: "positive testing (download executed)",
			args: args{
				conf: &config.JrpCliConfig{
					JrpConfig: baseConfig.JrpConfig{
						JrpDBType:   "sqlite",
						JrpDBDsn:    filepath.Join(os.TempDir(), "jrp.db"),
						WNJpnDBType: "sqlite",
						WNJpnDBDsn:  filepath.Join(os.TempDir(), "wnjpn.db"),
					},
				},
				output: &output,
			},
			want:    color.GreenString("✅ Downloaded successfully! Now, you are ready to use jrp!"),
			wantErr: false,
			setup: func(_ *gomock.Controller) {
				if err := os.Remove(filepath.Join(os.TempDir(), "wnjpn.db")); err != nil && !os.IsNotExist(err) {
					t.Errorf("Failed to remove test database: %v", err)
				}
				output = ""
			},
			cleanup: nil,
		},
		{
			name: "positive testing (already downloaded)",
			args: args{
				conf: &config.JrpCliConfig{
					JrpConfig: baseConfig.JrpConfig{
						JrpDBType:   "sqlite",
						JrpDBDsn:    filepath.Join(os.TempDir(), "jrp.db"),
						WNJpnDBType: "sqlite",
						WNJpnDBDsn:  filepath.Join(os.TempDir(), "wnjpn.db"),
					},
				},
				output: &output,
			},
			want:    color.GreenString("✅ You are already ready to use jrp!"),
			wantErr: false,
			setup: func(_ *gomock.Controller) {
				duc := jrpApp.NewDownloadUseCase()
				if err := duc.Run(filepath.Join(os.TempDir(), "wnjpn.db")); err != nil && err.Error() != "wnjpn.db already exists" {
					t.Errorf("Failed to download WordNet Japan DB file: %v", err)
				}
				output = ""
			},
		},
		{
			name: "negative testing (conf.WNJpnDBType != database.SQLite)",
			args: args{
				conf: &config.JrpCliConfig{
					JrpConfig: baseConfig.JrpConfig{
						JrpDBType:   "sqlite",
						JrpDBDsn:    filepath.Join(os.TempDir(), "jrp.db"),
						WNJpnDBType: "test",
						WNJpnDBDsn:  filepath.Join(os.TempDir(), "wnjpn.db"),
					},
				},
				output: &output,
			},
			want:    color.RedString("❌ The type of WordNet Japan database is not sqlite..."),
			wantErr: false,
			setup: func(_ *gomock.Controller) {
				output = ""
			},
			cleanup: nil,
		},
		{
			name: "negative testing (presenter.StartSpinner() failed)",
			args: args{
				conf: &config.JrpCliConfig{
					JrpConfig: baseConfig.JrpConfig{
						JrpDBType:   "sqlite",
						JrpDBDsn:    filepath.Join(os.TempDir(), "jrp.db"),
						WNJpnDBType: "sqlite",
						WNJpnDBDsn:  filepath.Join(os.TempDir(), "wnjpn.db"),
					},
				},
				output: &output,
			},
			want:    color.RedString("❌ Failed to start spinner..."),
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller) {
				mockSu := utility.NewMockSpinnerUtil(mockCtrl)
				mockSu.EXPECT().GetSpinner(
					true,
					"yellow",
					color.YellowString("  📦 Downloading WordNet Japan sqlite database file from the official web site..."),
				).Return(nil, errors.New("SpinnerUtil.GetSpinner() failed"))
				presenter.Su = mockSu
				output = ""
			},
			cleanup: func() {
				presenter.Su = origSu
			},
		},
		{
			name: "negative testing (duc.Run(conf.WNJpnDBDsn) failed)",
			args: args{
				conf: &config.JrpCliConfig{
					JrpConfig: baseConfig.JrpConfig{
						JrpDBType:   "sqlite",
						JrpDBDsn:    filepath.Join(os.TempDir(), "jrp.db"),
						WNJpnDBType: "sqlite",
						WNJpnDBDsn:  filepath.Join(os.TempDir(), "wnjpn.db"),
					},
				},
				output: &output,
			},
			want:    color.RedString("❌ Failed to download WordNet Japan sqlite database file..."),
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller) {
				if fu.IsExist(filepath.Join(os.TempDir(), "wnjpn.db")) {
					if hfp, err := fu.HideFile(filepath.Join(os.TempDir(), "wnjpn.db")); err != nil {
						t.Errorf("Failed to hide the test database: %v", err)
					} else {
						hiddenFilePath = hfp
					}
				}
				mockDu := utility.NewMockDownloadUtil(mockCtrl)
				mockDu.EXPECT().Download("https://github.com/bond-lab/wnja/releases/download/v1.1/wnjpn.db.gz").Return(nil, errors.New("DownloadUtil.Download() failed"))
				jrpApp.Du = mockDu
				output = ""
			},
			cleanup: func() {
				jrpApp.Du = origDu
				if hiddenFilePath != "" {
					if err := fu.UnhideFile(hiddenFilePath); err != nil {
						t.Errorf("Failed to unhide the test database: %v", err)
					}
				}
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
			if err := runDownload(tt.args.conf, tt.args.output); (err != nil) != tt.wantErr {
				t.Errorf("runDownload() error = %v, wantErr %v", err, tt.wantErr)
			}
			if *tt.args.output != tt.want {
				t.Errorf("runDownload() = %v, want %v", *tt.args.output, tt.want)
			}
		})
	}
}
