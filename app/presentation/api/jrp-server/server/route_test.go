package server

import (
	"testing"

	"github.com/yanosea/jrp/v2/pkg/proxy"

	"go.uber.org/mock/gomock"
)

func TestBind(t *testing.T) {
	type args struct {
		e proxy.Echo
	}
	tests := []struct {
		name  string
		args  args
		setup func(mockCtrl *gomock.Controller, tt *args)
	}{
		{
			name: "positive testing",
			args: args{
				e: nil,
			},
			setup: func(mockCtrl *gomock.Controller, tt *args) {
				mockGroup := proxy.NewMockGroup(mockCtrl)
				mockGroup.EXPECT().GET("/jrp", gomock.Any())
				mockEcho := proxy.NewMockEcho(mockCtrl)
				mockEcho.EXPECT().Group("/api").Return(mockGroup)
				mockEcho.EXPECT().Get("/swagger/*", gomock.Any())
				tt.e = mockEcho
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.setup(mockCtrl, &tt.args)
			}
			Bind(tt.args.e)
		})
	}
}
