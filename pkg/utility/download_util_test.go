package utility

import (
	"errors"
	"reflect"
	"testing"

	"github.com/yanosea/jrp/pkg/proxy"

	"go.uber.org/mock/gomock"
)

func TestNewDownloadUtil(t *testing.T) {
	http := proxy.NewHttp()

	type args struct {
		http proxy.Http
	}
	tests := []struct {
		name string
		args args
		want DownloadUtil
	}{
		{
			name: "positive testing",
			args: args{
				http: http,
			},
			want: &downloadUtil{
				http: http,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDownloadUtil(tt.args.http); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDownloadUtil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_downloadUtil_Download(t *testing.T) {
	type fields struct {
		Http proxy.Http
	}
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing",
			fields: fields{
				Http: nil,
			},
			args: args{
				url: "http://example.com",
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockHttp.EXPECT().Get(gomock.Any()).Return(nil, nil)
				tt.Http = mockHttp
			},
		},
		{
			name: "negative testing (d.Http.Get(url) failed)",
			fields: fields{
				Http: nil,
			},
			args: args{
				url: "http://example.com",
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockHttp := proxy.NewMockHttp(mockCtrl)
				mockHttp.EXPECT().Get(gomock.Any()).Return(nil, errors.New("HttpProxy.Get() failed"))
				tt.Http = mockHttp
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.fields)
			}
			d := &downloadUtil{
				http: tt.fields.Http,
			}
			_, err := d.Download(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("downloadUtil.Download() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
