package spinnerproxy

import (
	"github.com/briandowns/spinner"

	"github.com/yanosea/jrp/app/proxy/time"
)

// Spinner is an interface for spinner.
type Spinner interface {
	NewSpinner() SpinnerInstanceInterface
}

// SpinnerProxy is a struct that implements Spinner.
type SpinnerProxy struct{}

// New is a constructor for SpinnerProxy.
func New() Spinner {
	return &SpinnerProxy{}
}

// New is a proxy for spinner.New.
func (*SpinnerProxy) NewSpinner() SpinnerInstanceInterface {
	return &SpinnerInstance{FieldSpinner: spinner.New(spinner.CharSets[11], 100*timeproxy.Millisecond)}
}
