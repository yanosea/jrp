package jrp

import (
	"errors"
	"reflect"
	"testing"

	"github.com/yanosea/jrp/pkg/proxy"
	"github.com/yanosea/jrp/pkg/utility"

	"go.uber.org/mock/gomock"
)

func TestNewDownloadUseCase(t *testing.T) {
	tests := []struct {
		name string
		want *downloadUseCase
	}{
		{
			name: "positive testing",
			want: &downloadUseCase{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDownloadUseCase(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDownloadUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_downloadUseCase_Run(t *testing.T) {
	origDu := Du
	origFu := Fu

	type args struct {
		wnJpnDBPath string
	}
	tests := []struct {
		name    string
		uc      *downloadUseCase
		args    args
		wantErr bool
		setup   func(mockCtrl *gomock.Controller)
		clear   func()
	}{
		{
			name: "positive testing",
			uc:   &downloadUseCase{},
			args: args{
				wnJpnDBPath: "/tmp/wnjpn.db",
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller) {
				mockReadCloser := proxy.NewMockReadCloser(mockCtrl)
				mockResp := proxy.NewMockResponse(mockCtrl)
				mockResp.EXPECT().GetBody().Return(mockReadCloser)
				mockResp.EXPECT().Close().Return(nil)
				mockFu := utility.NewMockFileUtil(mockCtrl)
				mockFu.EXPECT().IsExist("/tmp/wnjpn.db").Return(false)
				mockFu.EXPECT().SaveToTempFile(mockReadCloser, "wnjpn.db.gz").Return("/tmp/wjpn.db.gz", nil)
				mockFu.EXPECT().ExtractGzFile("/tmp/wjpn.db.gz", "/tmp/wnjpn.db").Return(nil)
				mockFu.EXPECT().RemoveAll("/tmp/wjpn.db.gz").Return(nil)
				mockDu := utility.NewMockDownloadUtil(mockCtrl)
				mockDu.EXPECT().Download("https://github.com/bond-lab/wnja/releases/download/v1.1/wnjpn.db.gz").Return(mockResp, nil)
				Du = mockDu
				Fu = mockFu
			},
			clear: func() {
				Du = origDu
				Fu = origFu
			},
		},
		{
			name: "negative testing (fu.IsExist(wnJpnDBPath) returns true)",
			uc:   &downloadUseCase{},
			args: args{
				wnJpnDBPath: "/tmp/wnjpn.db",
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller) {
				mockFu := utility.NewMockFileUtil(mockCtrl)
				mockFu.EXPECT().IsExist("/tmp/wnjpn.db").Return(true)
				Fu = mockFu
			},
			clear: func() {
				Fu = origFu
			},
		},
		{
			name: "negative testing (du.Download() failed)",
			uc:   &downloadUseCase{},
			args: args{
				wnJpnDBPath: "/tmp/wnjpn.db",
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller) {
				mockFu := utility.NewMockFileUtil(mockCtrl)
				mockFu.EXPECT().IsExist("/tmp/wnjpn.db").Return(false)
				mockDu := utility.NewMockDownloadUtil(mockCtrl)
				mockDu.EXPECT().Download("https://github.com/bond-lab/wnja/releases/download/v1.1/wnjpn.db.gz").Return(nil, errors.New("DownloadUtil.Download() failed"))
				Du = mockDu
				Fu = mockFu
			},
			clear: func() {
				Du = origDu
				Fu = origFu
			},
		},
		{
			name: "negative testing (fu.SaveToTempFile(resp.GetBody(), \"wnjpn.db.gz\") failed)",
			uc:   &downloadUseCase{},
			args: args{
				wnJpnDBPath: "/tmp/wnjpn.db",
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller) {
				mockReadCloser := proxy.NewMockReadCloser(mockCtrl)
				mockResp := proxy.NewMockResponse(mockCtrl)
				mockResp.EXPECT().GetBody().Return(mockReadCloser)
				mockResp.EXPECT().Close().Return(nil)
				mockFu := utility.NewMockFileUtil(mockCtrl)
				mockFu.EXPECT().IsExist("/tmp/wnjpn.db").Return(false)
				mockFu.EXPECT().SaveToTempFile(mockReadCloser, "wnjpn.db.gz").Return("", errors.New("FileUtil.SaveToTempFile() failed"))
				mockDu := utility.NewMockDownloadUtil(mockCtrl)
				mockDu.EXPECT().Download("https://github.com/bond-lab/wnja/releases/download/v1.1/wnjpn.db.gz").Return(mockResp, nil)
				Du = mockDu
				Fu = mockFu
			},
			clear: func() {
				Du = origDu
				Fu = origFu
			},
		},
		{
			name: "negative testing (fu.RemoveAll(tempFilePath) failed)",
			uc:   &downloadUseCase{},
			args: args{
				wnJpnDBPath: "/tmp/wnjpn.db",
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller) {
				mockReadCloser := proxy.NewMockReadCloser(mockCtrl)
				mockResp := proxy.NewMockResponse(mockCtrl)
				mockResp.EXPECT().GetBody().Return(mockReadCloser)
				mockResp.EXPECT().Close().Return(nil)
				mockFu := utility.NewMockFileUtil(mockCtrl)
				mockFu.EXPECT().IsExist("/tmp/wnjpn.db").Return(false)
				mockFu.EXPECT().SaveToTempFile(mockReadCloser, "wnjpn.db.gz").Return("/tmp/wjpn.db.gz", nil)
				mockFu.EXPECT().ExtractGzFile("/tmp/wjpn.db.gz", "/tmp/wnjpn.db").Return(errors.New("FileUtil.ExtractGzFile() failed"))
				mockFu.EXPECT().RemoveAll("/tmp/wjpn.db.gz").Return(nil)
				mockDu := utility.NewMockDownloadUtil(mockCtrl)
				mockDu.EXPECT().Download("https://github.com/bond-lab/wnja/releases/download/v1.1/wnjpn.db.gz").Return(mockResp, nil)
				Du = mockDu
				Fu = mockFu
			},
			clear: func() {
				Du = origDu
				Fu = origFu
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
				if tt.clear != nil {
					tt.clear()
				}
			}()
			uc := &downloadUseCase{}
			if err := uc.Run(tt.args.wnJpnDBPath); (err != nil) != tt.wantErr {
				t.Errorf("downloadUseCase.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
