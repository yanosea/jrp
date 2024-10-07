package versionprovider

import (
	"reflect"
	"runtime/debug"
	"testing"

	"github.com/yanosea/jrp/app/proxy/debug"

	"github.com/yanosea/jrp/mock/app/proxy/debug"
	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	debugProxy := debugproxy.New()

	type args struct {
		debugProxy debugproxy.Debug
	}
	tests := []struct {
		name string
		args args
		want *VersionProvider
	}{
		{
			name: "positive testing",
			args: args{
				debugProxy: debugProxy,
			},
			want: &VersionProvider{
				DebugProxy: debugProxy,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.debugProxy); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() : got =\n%v, want =\n%v", got, tt.want)
			}
		})
	}
}

func TestVersionProvider_GetVersion(t *testing.T) {
	debugProxy := debugproxy.New()

	type fields struct {
		versionProvider *VersionProvider
	}
	type args struct {
		embeddedVersion string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		setup  func(mockCtrl *gomock.Controller, tt *fields)
	}{
		{
			name:   "positive testing (version is embedded)",
			fields: fields{versionProvider: nil},
			args: args{
				embeddedVersion: "vx.x.x",
			},
			want: "vx.x.x",
			setup: func(_ *gomock.Controller, tt *fields) {
				versionProvider := New(debugProxy)
				tt.versionProvider = versionProvider
			},
		},
		{
			name:   "positive testing (version is not embedded and DebugProxy.ReadBuildInfo() returns not ok)",
			fields: fields{versionProvider: nil},
			args: args{
				embeddedVersion: "",
			},
			want: "unknown",
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDebugProxy := mockdebugproxy.NewMockDebug(mockCtrl)
				buildInfo, ok := debugProxy.ReadBuildInfo()
				if !ok {
					t.Errorf("DebugProxy.ReadBuildInfo() failed")
				}
				mockDebugProxy.EXPECT().ReadBuildInfo().Return(buildInfo, false)
				versionProvider := New(mockDebugProxy)
				tt.versionProvider = versionProvider
			},
		},
		{
			name:   "positive testing (version is not embedded and DebugProxy.ReadBuildInfo() returns ok, but version is empty)",
			fields: fields{versionProvider: nil},
			args: args{
				embeddedVersion: "",
			},
			want: "devel",
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDebugProxy := mockdebugproxy.NewMockDebug(mockCtrl)
				mockBuildInfo := &debugproxy.BuildInfoInstance{
					FieldBuildInfo: &debug.BuildInfo{
						Main: debug.Module{
							Version: "",
						},
					},
				}
				mockDebugProxy.EXPECT().ReadBuildInfo().Return(mockBuildInfo, true)
				versionProvider := New(mockDebugProxy)
				tt.versionProvider = versionProvider
			},
		},
		{
			name:   "positive testing (version is not embedded and DebugProxy.ReadBuildInfo() returns ok, but version is (devel))",
			fields: fields{versionProvider: nil},
			args: args{
				embeddedVersion: "",
			},
			want: "devel",
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDebugProxy := mockdebugproxy.NewMockDebug(mockCtrl)
				mockBuildInfo := &debugproxy.BuildInfoInstance{
					FieldBuildInfo: &debug.BuildInfo{
						Main: debug.Module{
							Version: "(devel)",
						},
					},
				}
				mockDebugProxy.EXPECT().ReadBuildInfo().Return(mockBuildInfo, true)
				versionProvider := New(mockDebugProxy)
				tt.versionProvider = versionProvider
			},
		},
		{
			name:   "positive testing (version is not embedded and DebugProxy.ReadBuildInfo() returns ok, version is not empty, and version is not (devel))",
			fields: fields{versionProvider: nil},
			args: args{
				embeddedVersion: "",
			},
			want: "vy.y.y",
			setup: func(mockCtrl *gomock.Controller, tt *fields) {
				mockDebugProxy := mockdebugproxy.NewMockDebug(mockCtrl)
				mockBuildInfo := &debugproxy.BuildInfoInstance{
					FieldBuildInfo: &debug.BuildInfo{
						Main: debug.Module{
							Version: "vy.y.y",
						},
					},
				}
				mockDebugProxy.EXPECT().ReadBuildInfo().Return(mockBuildInfo, true)
				versionProvider := New(mockDebugProxy)
				tt.versionProvider = versionProvider
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				mockCtrl := gomock.NewController(t)
				defer mockCtrl.Finish()
				tt.setup(mockCtrl, &tt.fields)
			}
			if got := tt.fields.versionProvider.GetVersion(tt.args.embeddedVersion); got != tt.want {
				t.Errorf("VersionProvider.GetVersion() : got =\n%v, want =\n%v", got, tt.want)
			}
		})
	}
}
