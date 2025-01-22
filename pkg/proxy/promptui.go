package proxy

import (
	"github.com/manifoldco/promptui"
)

// Promptui is an interface that provides a proxy of the methods of promptui.
type Promptui interface {
	NewPrompt() Prompt
}

// promptuiProxy is a proxy struct that implements the Promptui interface.
type promptuiProxy struct{}

// NewPromptui returns a new instance of the Promptui interface.
func NewPromptui() Promptui {
	return &promptuiProxy{}
}

// NewPrompt returns a new instance of the promptui.Prompt.
func (p *promptuiProxy) NewPrompt() Prompt {
	return &promptProxy{prompt: &promptui.Prompt{}}
}

// Prompt is an interface that provides a proxy of the methods of promptui.Prompt.
type Prompt interface {
	Run() (string, error)
	SetLabel(label string)
}

// promptProxy is a proxy struct that implements the Prompt interface.
type promptProxy struct {
	prompt *promptui.Prompt
}

// Run runs the prompt.
func (p *promptProxy) Run() (string, error) {
	return p.prompt.Run()
}

// SetLabel sets the label of the prompt.
func (p *promptProxy) SetLabel(label string) {
	p.prompt.Label = label
}
