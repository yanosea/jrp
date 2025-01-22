package utility

import (
	"github.com/yanosea/jrp/pkg/proxy"
)

type SpinnerUtil interface {
	GetSpinner(isReversed bool, color string, suffix string) (proxy.Spinner, error)
}

// spinnerUtil is a struct that implements the SpinnerUtil interface.
type spinnerUtil struct {
	spinners proxy.Spinners
}

// NewSpinnerUtil returns a new instance of the spinnerUtil struct.
func NewSpinnerUtil(
	spinners proxy.Spinners,
) SpinnerUtil {
	return &spinnerUtil{
		spinners: spinners,
	}
}

// GetSpinner returns a new instance of the spinner.Spinner.
func (s *spinnerUtil) GetSpinner(isReversed bool, color string, suffix string) (proxy.Spinner, error) {
	spinner := s.spinners.NewSpinner()
	if isReversed {
		spinner.Reverse()
	}
	if err := spinner.SetColor(color); err != nil {
		return nil, err
	}
	spinner.SetSuffix(suffix)

	return spinner, nil
}
