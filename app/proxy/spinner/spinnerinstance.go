package spinnerproxy

import (
	"github.com/briandowns/spinner"
)

// SpinnerInstanceInterface is an interface for spinner.Spinner
type SpinnerInstanceInterface interface {
	Reverse()
	SetColor(colors ...string) error
	SetSuffix(suffix string)
	Start()
	Stop()
}

// SpinnerInstance is a struct that implements SpinnerInstanceInterface.
type SpinnerInstance struct {
	FieldSpinner *spinner.Spinner
}

// Reverse is a proxy for spinner.Spinner.Reverse.
func (s *SpinnerInstance) Reverse() {
	s.FieldSpinner.Reverse()
}

// SetColor is a proxy for spinner.Spinner.Color.
func (s *SpinnerInstance) SetColor(colors ...string) error {
	return s.FieldSpinner.Color(colors...)
}

// SetSuffix is a proxy for spinner.Spinner.Suffix.
func (s *SpinnerInstance) SetSuffix(suffix string) {
	s.FieldSpinner.Suffix = suffix
}

// Start is a proxy for spinner.Spinner.Start.
func (s *SpinnerInstance) Start() {
	s.FieldSpinner.Start()
}

// Stop is a proxy for spinner.Spinner.Stop.
func (s *SpinnerInstance) Stop() {
	s.FieldSpinner.Stop()
}
