package completion

import (
	"bytes"

	c "github.com/spf13/cobra"

	"github.com/yanosea/jrp/pkg/proxy"
)

// NewCompletionFishCommand returns a new instance of the completion fish command.
func NewCompletionFishCommand(
	cobra proxy.Cobra,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("fish")
	cmd.SetAliases([]string{"fi", "f"})
	cmd.SetUsageTemplate(completionFishUsageTemplate)
	cmd.SetHelpTemplate(completionFishHelpTemplate)
	cmd.SetArgs(cobra.ExactArgs(0))
	cmd.SetSilenceErrors(true)
	cmd.SetRunE(
		func(cmd *c.Command, args []string) error {
			return runCompletionFish(cmd, output)
		},
	)

	return cmd
}

// runCompletionFish generates the autocompletion script for the fish shell.
func runCompletionFish(cmd *c.Command, output *string) error {
	buf := new(bytes.Buffer)
	err := cmd.Root().GenFishCompletion(buf, true)
	*output = buf.String()
	return err
}

const (
	// completionFishHelpTemplate is the help template of the completion fish command.
	completionFishHelpTemplate = `🔧🐟 Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

  jrp completion fish | source

To load completions for every new session, execute once:

  jrp completion fish > ~/.config/fish/completions/jrp.fish

You will need to start a new shell for this setup to take effect.

` + completionFishUsageTemplate
	// compleitonUsageTemplate is the usage template of the completion fish command.
	completionFishUsageTemplate = `Usage:
  jrp completion fish [flags]

Flags:
  -h, --help  🤝 help for fish
`
)
