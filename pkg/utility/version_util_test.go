package utility

import (
	"reflect"
	d "runtime/debug"
	"testing"

	"github.com/yanosea/jrp/pkg/proxy"

	"go.uber.org/mock/gomock"
)

func TestNewVersionUtil(t *testing.T) {
	debug := proxy.NewDebug()

	type args struct {
		debug proxy.Debug
	}
	tests := []struct {
		name string
		args args
		want VersionUtil
	}{
		{
			name: "positive testing",
			args: args{
				debug: debug,
			},
			want: &versionUtil{
				debug: debug,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewVersionUtil(tt.args.debug); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVersionUtil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_versionUtil_GetVersion(t *testing.T) {
	debug := proxy.NewDebug()

	type fields struct {
		Debug proxy.Debug
	}
	type args struct {
		version string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		setup  func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name: "positive testing (version is not empty)",
			fields: fields{
				Debug: debug,
			},
			args: args{
				version: "0.0.0",
			},
			want:  "0.0.0",
			setup: nil,
		},
		{
			name: "positive testing (version is empty, debug.ReadBuildInfo() returns false)",
			fields: fields{
				Debug: nil,
			},
			args: args{
				version: "",
			},
			want: "unknown",
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDebug := proxy.NewMockDebug(mockCtrl)
				mockDebug.EXPECT().ReadBuildInfo().Return(nil, false)
				tt.Debug = mockDebug
			},
		},
		{
			name: "positive testing (version is empty, debug.ReadBuildInfo() returns false)",
			fields: fields{
				Debug: nil,
			},
			args: args{
				version: "",
			},
			want: "unknown",
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDebug := proxy.NewMockDebug(mockCtrl)
				mockDebug.EXPECT().ReadBuildInfo().Return(
					&d.BuildInfo{
						Main: d.Module{
							Version: "unknown",
						},
					},
					false,
				)
				tt.Debug = mockDebug
			},
		},
		{
			name: "positive testing (version is empty, debug.ReadBuildInfo() returns true, i.Main.Version is not empty)",
			fields: fields{
				Debug: nil,
			},
			args: args{
				version: "",
			},
			want: "0.0.0",
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDebug := proxy.NewMockDebug(mockCtrl)
				mockDebug.EXPECT().ReadBuildInfo().Return(
					&d.BuildInfo{
						Main: d.Module{
							Version: "0.0.0",
						},
					},
					true,
				)
				tt.Debug = mockDebug
			},
		},
		{
			name: "positive testing (version is empty, debug.ReadBuildInfo() returns true, i.Main.Version is empty)",
			fields: fields{
				Debug: nil,
			},
			args: args{
				version: "",
			},
			want: "dev",
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDebug := proxy.NewMockDebug(mockCtrl)
				mockDebug.EXPECT().ReadBuildInfo().Return(
					&d.BuildInfo{
						Main: d.Module{
							Version: "",
						},
					},
					true,
				)
				tt.Debug = mockDebug
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
			v := &versionUtil{
				debug: tt.fields.Debug,
			}
			if got := v.GetVersion(tt.args.version); got != tt.want {
				t.Errorf("versionUtil.GetVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
