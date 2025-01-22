package presenter

import (
	"github.com/yanosea/jrp/pkg/proxy"
	"github.com/yanosea/jrp/pkg/utility"
)

var (
	// Pu is a variable that contains the PromptUtil struct for injecting dependencies in testing.
	Pu = utility.NewPromptUtil(proxy.NewPromptui())
)

// RunPrompt runs the prompt.
func RunPrompt(label string) (string, error) {
	prompt := Pu.GetPrompt(label)
	return prompt.Run()
}
