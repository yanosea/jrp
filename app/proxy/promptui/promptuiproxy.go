package promptuiproxy

import (
	"github.com/manifoldco/promptui"
)

// Promptui is an interface for promptui.
type Promptui interface {
	NewPrompt() PromptInstanceInterface
}

// PromptuiProxy is a struct that implements Promptui.
type PromptuiProxy struct{}

// New is a constructor for PromptuiProxy.
func New() Promptui {
	return &PromptuiProxy{}
}

// NewPrompt is a proxy for getting promptui.Prompt struct.
func (*PromptuiProxy) NewPrompt() PromptInstanceInterface {
	return &PromptInstance{FieldPrompt: &promptui.Prompt{}}
}
