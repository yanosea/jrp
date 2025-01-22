package utility

import (
	"reflect"
	"testing"

	"github.com/yanosea/jrp/pkg/proxy"
)

func TestNewRandUtil(t *testing.T) {
	rand := proxy.NewRand()

	type args struct {
		rand proxy.Rand
	}
	tests := []struct {
		name string
		args args
		want RandUtil
	}{
		{
			name: "positive testing",
			args: args{
				rand: rand,
			},
			want: &randUtil{
				rand: rand,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRandUtil(tt.args.rand); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRandUtil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_randomUtil_GenerateRandNumber(t *testing.T) {
	rand := proxy.NewRand()

	type fields struct {
		Rand proxy.Rand
	}
	type args struct {
		max int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		isRand bool
	}{
		{
			name: "positive testing (max < 0)",
			fields: fields{
				Rand: rand,
			},
			args: args{
				max: -1,
			},
			want:   0,
			isRand: false,
		},
		{
			name: "positive testing (max = 0)",
			fields: fields{
				Rand: rand,
			},
			args: args{
				max: 0,
			},
			want:   0,
			isRand: false,
		},
		{
			name: "positive testing (max = 1)",
			fields: fields{
				Rand: rand,
			},
			args: args{
				max: 1,
			},
			want:   0,
			isRand: true,
		},
		{
			name: "positive testing (max = 10)",
			fields: fields{
				Rand: rand,
			},
			args: args{
				max: 10,
			},
			want:   0,
			isRand: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ru := &randUtil{
				rand: tt.fields.Rand,
			}
			got := ru.GenerateRandomNumber(tt.args.max)
			if !tt.isRand {
				if got != tt.want {
					t.Errorf("randomUtil.GenerateRandNumber() = %v, want %v", got, tt.want)
				}
			} else if got < 0 || got > tt.args.max {
				t.Errorf("randomUtil.GenerateRandNumber() = %v, max %v", got, tt.args.max)
			}
		})
	}
}
