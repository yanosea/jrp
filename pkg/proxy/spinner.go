package proxy

import (
	"time"

	"github.com/briandowns/spinner"
)

// Spinners is an interface that provides a proxy of the methods of spinner.
type Spinners interface {
	NewSpinner() Spinner
}

// spinnersProxy is a proxy struct that implements the Spinner interface.
type spinnersProxy struct{}

// NewSpinners returns a new instance of the Spinners interface.
func NewSpinners() Spinners {
	return &spinnersProxy{}
}

// NewSpinner returns a new instance of the spinner.Spinner.
func (s *spinnersProxy) NewSpinner() Spinner {
	return &spinnerProxy{spinner: spinner.New(spinner.CharSets[11], 100*time.Millisecond)}
}

// Spinner is an interface that provides a proxy of the methods of spinner.Spinner.
type Spinner interface {
	Reverse()
	SetColor(colors ...string) error
	SetSuffix(suffix string)
	Start()
	Stop()
}

// spinnerProxy is a proxy struct that implements the Spinner interface.
type spinnerProxy struct {
	spinner *spinner.Spinner
}

// Reverse sets the direction of the spinner to reverse.
func (s *spinnerProxy) Reverse() {
	s.spinner.Reverse()
}

// SetColor sets the color of the spinner.
func (s *spinnerProxy) SetColor(colors ...string) error {
	return s.spinner.Color(colors...)
}

// SetSuffix sets the suffix of the spinner.
func (s *spinnerProxy) SetSuffix(suffix string) {
	s.spinner.Suffix = suffix
}

// Start starts the spinner.
func (s *spinnerProxy) Start() {
	s.spinner.Start()
}

// Stop stops the spinner.
func (s *spinnerProxy) Stop() {
	s.spinner.Stop()
}
