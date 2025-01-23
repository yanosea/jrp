package presenter

import (
	"github.com/yanosea/jrp/v2/pkg/proxy"
	"github.com/yanosea/jrp/v2/pkg/utility"
)

var (
	// Su is a variable that contains the SpinnerUtil struct for injecting dependencies in testing.
	Su = utility.NewSpinnerUtil(proxy.NewSpinners())
	// spinner is a variable that contains the Spinner struct.
	spinner proxy.Spinner
)

// StartSpinner gets and starts the spinner.
func StartSpinner(isRversed bool, color string, suffix string) error {
	sp, err := Su.GetSpinner(isRversed, color, suffix)
	if err != nil {
		return err
	}

	spinner = sp
	sp.Start()

	return nil
}

// StopSpinner stops the spinner.
func StopSpinner() {
	if spinner != nil {
		spinner.Stop()
	}
}
