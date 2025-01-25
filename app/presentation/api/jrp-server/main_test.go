package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yanosea/jrp/v2/app/presentation/api/jrp-server/server"

	"github.com/yanosea/jrp/v2/pkg/proxy"

	"go.uber.org/mock/gomock"
)

func Test_main(t *testing.T) {
	os.Setenv("JRP_SERVER_WNJPN_DB_TYPE", "sqlite")
	os.Setenv("JRP_SERVER_WNJPN_DB", filepath.Join(os.TempDir(), "wnjpn.db"))
	origExit := exit
	exit = func(code int) {}
	defer func() {
		exit = origExit
	}()
	origJrpApiServerParams := jrpApiServerParams

	tests := []struct {
		name  string
		setup func(mockCtrl *gomock.Controller)
		clear func()
	}{
		{
			name: "positive testing",
			setup: func(mockCtrl *gomock.Controller) {
				mockGroup := proxy.NewMockGroup(mockCtrl)
				mockGroup.EXPECT().GET("/jrp", gomock.Any())
				mockEcho := proxy.NewMockEcho(mockCtrl)
				mockEcho.EXPECT().Use(gomock.Any())
				mockEcho.EXPECT().Use(gomock.Any())
				mockEcho.EXPECT().Group("/api").Return(mockGroup)
				mockEcho.EXPECT().Start(":8080")
				mockEcho.EXPECT().Get("/swagger/*", gomock.Any())
				mockLogger := proxy.NewMockLogger(mockCtrl)
				mockEchos := proxy.NewMockEchos(mockCtrl)
				mockEchos.EXPECT().NewEcho().Return(mockEcho, mockLogger)
				jrpApiServerParams.Echos = mockEchos
			},
			clear: func() {
				jrpApiServerParams = origJrpApiServerParams
			},
		},
		{
			name: "negative testing (serv.Init() failed)",
			setup: func(mockCtrl *gomock.Controller) {
				mockServ := server.NewMockServer(mockCtrl)
				mockServ.EXPECT().Init(gomock.Any(), gomock.Any(), gomock.Any()).Return(1)
				mockServ.EXPECT().Run()
				origNewServer := server.NewServer
				server.NewServer = func(echo proxy.Echos) server.Server {
					return mockServ
				}
				t.Cleanup(func() {
					server.NewServer = origNewServer
				})
			},
			clear: func() {
				jrpApiServerParams = origJrpApiServerParams
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
			main()
		})
	}
}
