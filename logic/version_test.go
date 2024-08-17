package logic

import (
	"runtime/debug"
	"testing"

	gomock "go.uber.org/mock/gomock"

	"github.com/yanosea/jrp/internal/buildinfo"

	mock_buildinfo "github.com/yanosea/jrp/mock/buildinfo"
)

func TestNewJrpVersionGetter(t *testing.T) {
	type args struct {
		v buildinfo.BuildInfoProvider
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive testing",
			args: args{v: buildinfo.RealBuildInfoProvider{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewJrpVersionGetter(tt.args.v)
			if u == nil {
				t.Errorf("NewJrpVersionGetter() : returned nil")
			}
		})
	}
}

func TestGetVersion(t *testing.T) {
	type args struct {
		version          string
		jrpVersionGetter *JrpVersionGetter
	}
	tests := []struct {
		name  string
		args  args
		want  string
		setup func(mockCtrl *gomock.Controller, tt *args)
	}{
		{
			name: "positive testing (version is embedded)",
			args: args{version: "vx.x.x", jrpVersionGetter: nil},
			want: "vx.x.x",
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				jrpVersionGetter := NewJrpVersionGetter(buildinfo.RealBuildInfoProvider{})
				tt.jrpVersionGetter = jrpVersionGetter
			},
		}, {
			name: "positive testing (version is not embedded and ReadBuildInfo() returns not ok)",
			args: args{version: "", jrpVersionGetter: nil},
			want: "unknown",
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mb := mock_buildinfo.NewMockBuildInfoProvider(mockCtrl)
				buildInfo, ok := debug.ReadBuildInfo()
				if !ok {
					t.Errorf("debug.BuildInfo is not found")
				}
				mb.EXPECT().ReadBuildInfo().Return(buildInfo, false)
				jrpVersionGetter := NewJrpVersionGetter(mb)
				tt.jrpVersionGetter = jrpVersionGetter
			},
		}, {
			name: "positive testing (version is not embedded and ReadBuildInfo() returns ok)",
			args: args{version: "", jrpVersionGetter: nil},
			want: "vy.y.y",
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mb := mock_buildinfo.NewMockBuildInfoProvider(mockCtrl)
				mockBuildInfo := &debug.BuildInfo{
					Main: debug.Module{Version: "vy.y.y"},
				}
				mb.EXPECT().ReadBuildInfo().Return(mockBuildInfo, true)
				jrpVersionGetter := NewJrpVersionGetter(mb)
				tt.jrpVersionGetter = jrpVersionGetter
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			if tt.setup != nil {
				tt.setup(ctrl, &tt.args)
			}

			if got := tt.args.jrpVersionGetter.GetVersion(tt.args.version); got != tt.want {
				t.Errorf("GetVersion() : got = %v, want = %v", got, tt.want)
			}
		})
	}
}
