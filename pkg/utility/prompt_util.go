package utility

import (
	"github.com/yanosea/jrp/pkg/proxy"
)

type PromptUtil interface {
	GetPrompt(label string) proxy.Prompt
}

// promptUtil is a struct that implements the PromptUtil interface.
type promptUtil struct {
	promptui proxy.Promptui
}

// NewPromptUtil returns a new instance of the promptUtil struct.
func NewPromptUtil(
	promptui proxy.Promptui,
) PromptUtil {
	return &promptUtil{
		promptui: promptui,
	}
}

// GetPrompt returns a new instance of the promptui.Prompt.
func (p *promptUtil) GetPrompt(label string) proxy.Prompt {
	prompt := p.promptui.NewPrompt()
	prompt.SetLabel(label)

	return prompt
}
