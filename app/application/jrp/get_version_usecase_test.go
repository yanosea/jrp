package jrp

import (
	"reflect"
	"testing"
)

func TestNewGetVersionUseCase(t *testing.T) {
	tests := []struct {
		name string
		want *getVersionUseCase
	}{
		{
			name: "positive testing",
			want: &getVersionUseCase{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGetVersionUseCase(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGetVersionUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getVersionUseCase_Run(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name string
		uc   *getVersionUseCase
		args args
		want *GetVersionUseCaseOutputDto
	}{
		{
			name: "positive testing",
			uc:   &getVersionUseCase{},
			args: args{
				version: "0.0.0",
			},
			want: &GetVersionUseCaseOutputDto{
				Version: "0.0.0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &getVersionUseCase{}
			if got := uc.Run(tt.args.version); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getVersionUseCase.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
