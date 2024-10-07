package promptuiproxy

import (
	"github.com/manifoldco/promptui"
)

// PromptInstanceInterface is an interface for promptui.Prompt.
type PromptInstanceInterface interface {
	Run() (string, error)
	SetLabel(label string)
}

// PromptInstance is a struct that implements PromptInstanceInterface.
type PromptInstance struct {
	FieldPrompt *promptui.Prompt
}

// Run is a proxy for promptui.Prompt.Run.
func (p *PromptInstance) Run() (string, error) {
	return p.FieldPrompt.Run()
}

// SetLabel is a proxy for promptui.Prompt.Label.
func (p *PromptInstance) SetLabel(label string) {
	p.FieldPrompt.Label = label
}
