package utility

import (
	"errors"
	"reflect"
	"testing"

	"github.com/yanosea/jrp/pkg/proxy"

	"go.uber.org/mock/gomock"
)

func TestNewSpinnerUtil(t *testing.T) {
	spinners := proxy.NewSpinners()

	type args struct {
		spinners proxy.Spinners
	}
	tests := []struct {
		name string
		args args
		want SpinnerUtil
	}{
		{
			name: "positive testing",
			args: args{
				spinners: spinners,
			},
			want: &spinnerUtil{
				spinners: spinners,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSpinnerUtil(tt.args.spinners); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSpinnerUtil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_spinnerUtil_GetSpinner(t *testing.T) {
	type fields struct {
		spinners proxy.Spinners
	}
	type args struct {
		isReversed bool
		color      string
		suffix     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    proxy.Spinner
		wantErr bool
		setup   func(mockCtrl *gomock.Controller, tt *fields) proxy.Spinner
	}{
		{
			name: "positive testing (isReversed = false)",
			fields: fields{
				spinners: proxy.NewSpinners(),
			},
			args: args{
				isReversed: false,
				color:      "red",
				suffix:     "suffix",
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) proxy.Spinner {
				mockSpinners := proxy.NewMockSpinners(mockCtrl)
				mockSpinner := proxy.NewMockSpinner(mockCtrl)
				mockSpinner.EXPECT().SetColor("red").Return(nil)
				mockSpinner.EXPECT().SetSuffix("suffix")
				mockSpinners.EXPECT().NewSpinner().Return(mockSpinner)
				tt.spinners = mockSpinners
				return mockSpinner
			},
		},
		{
			name: "positive testing (isReversed = true)",
			fields: fields{
				spinners: proxy.NewSpinners(),
			},
			args: args{
				isReversed: true,
				color:      "blue",
				suffix:     "test",
			},
			wantErr: false,
			setup: func(mockCtrl *gomock.Controller, tt *fields) proxy.Spinner {
				mockSpinners := proxy.NewMockSpinners(mockCtrl)
				mockSpinner := proxy.NewMockSpinner(mockCtrl)
				mockSpinner.EXPECT().Reverse()
				mockSpinner.EXPECT().SetColor("blue").Return(nil)
				mockSpinner.EXPECT().SetSuffix("test")
				mockSpinners.EXPECT().NewSpinner().Return(mockSpinner)
				tt.spinners = mockSpinners
				return mockSpinner
			},
		},
		{
			name: "negative testing (spinner.SetColor() failed)",
			fields: fields{
				spinners: proxy.NewSpinners(),
			},
			args: args{
				isReversed: false,
				color:      "invalid",
				suffix:     "test",
			},
			wantErr: true,
			setup: func(mockCtrl *gomock.Controller, tt *fields) proxy.Spinner {
				mockSpinners := proxy.NewMockSpinners(mockCtrl)
				mockSpinner := proxy.NewMockSpinner(mockCtrl)
				mockSpinner.EXPECT().SetColor("invalid").Return(errors.New("invalid color"))
				mockSpinners.EXPECT().NewSpinner().Return(mockSpinner)
				tt.spinners = mockSpinners
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			if tt.setup != nil {
				tt.want = tt.setup(mockCtrl, &tt.fields)
			}
			s := &spinnerUtil{
				spinners: tt.fields.spinners,
			}
			got, err := s.GetSpinner(tt.args.isReversed, tt.args.color, tt.args.suffix)
			if (err != nil) != tt.wantErr {
				t.Errorf("spinnerUtil.GetSpinner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("spinnerUtil.GetSpinner() = %v, want %v", got, tt.want)
			}
		})
	}
}
