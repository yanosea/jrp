package config

import (
	"reflect"
	"testing"

	"github.com/yanosea/jrp/v2/pkg/proxy"
	"github.com/yanosea/jrp/v2/pkg/utility"
)

func TestNewConfigurator(t *testing.T) {
	envconfig := proxy.NewEnvconfig()
	fileUtil := utility.NewFileUtil(
		proxy.NewGzip(),
		proxy.NewIo(),
		proxy.NewOs(),
	)

	type args struct {
		envconfigProxy proxy.Envconfig
		fileUtil       utility.FileUtil
	}
	tests := []struct {
		name string
		args args
		want *BaseConfigurator
	}{
		{
			name: "positive testing",
			args: args{
				envconfigProxy: envconfig,
				fileUtil:       fileUtil,
			},
			want: &BaseConfigurator{
				Envconfig: envconfig,
				FileUtil:  fileUtil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfigurator(tt.args.envconfigProxy, tt.args.fileUtil); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfigurator() = %v, want %v", got, tt.want)
			}
		})
	}
}
