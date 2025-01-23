package presenter

import (
	"errors"
	"testing"

	"github.com/yanosea/jrp/v2/pkg/proxy"
	"github.com/yanosea/jrp/v2/pkg/utility"

	"go.uber.org/mock/gomock"
)

func TestStartSpinner(t *testing.T) {
	origSu := Su

	type args struct {
		isRversed bool
		color     string
		suffix    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		setup   func(mockCtrl *gomock.Controller)
		cleanup func()
	}{
		{
			name: "positive testing",
			args: args{
				isRversed: false,
				color:     "red",
				suffix:    "test suffix",
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller) {
				mockSpinner := proxy.NewMockSpinner(mockCtrl)
				mockSpinner.EXPECT().SetColor("red").Return(nil)
				mockSpinner.EXPECT().SetSuffix("test suffix")
				mockSpinner.EXPECT().Start()
				mockSpinners := proxy.NewMockSpinners(mockCtrl)
				mockSpinners.EXPECT().NewSpinner().Return(mockSpinner)
				Su = utility.NewSpinnerUtil(mockSpinners)
			},
			cleanup: func() {
				Su = origSu
			},
		},
		{
			name: "negative testing (s.GetSpinner(isRversed, color, suffix) failed)",
			args: args{
				isRversed: false,
				color:     "invalid color",
				suffix:    "test suffix",
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller) {
				mockSpinner := proxy.NewMockSpinner(mockCtrl)
				mockSpinner.EXPECT().SetColor("invalid color").Return(errors.New("SpinnerProxy.SetColor() failed"))
				mockSpinners := proxy.NewMockSpinners(mockCtrl)
				mockSpinners.EXPECT().NewSpinner().Return(mockSpinner)
				Su = utility.NewSpinnerUtil(mockSpinners)
			},
			cleanup: func() {
				Su = origSu
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
				if tt.cleanup != nil {
					tt.cleanup()
				}
			}()
			if err := StartSpinner(tt.args.isRversed, tt.args.color, tt.args.suffix); (err != nil) != tt.wantErr {
				t.Errorf("StartSpinner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStopSpinner(t *testing.T) {
	origSpinner := spinner

	tests := []struct {
		name    string
		setup   func(mockCtrl *gomock.Controller)
		cleanup func()
	}{
		{
			name: "positive testing",
			setup: func(mockCtrl *gomock.Controller) {
				mockSpinner := proxy.NewMockSpinner(mockCtrl)
				mockSpinner.EXPECT().Stop()
				spinner = mockSpinner
			},
			cleanup: func() {
				spinner = origSpinner
			},
		},
	}
	for _, tt := range tests {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		if tt.setup != nil {
			tt.setup(mockCtrl)
		}
		defer func() {
			if tt.cleanup != nil {
				tt.cleanup()
			}
		}()
		t.Run(tt.name, func(t *testing.T) {
			StopSpinner()
		})
	}
}
