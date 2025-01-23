package completion

import (
	"github.com/yanosea/jrp/v2/pkg/proxy"
)

// NewCompletionCommand returns a new instance of the completion command.
func NewCompletionCommand(
	cobra proxy.Cobra,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("completion")
	cmd.SetAliases([]string{"comp", "c"})
	cmd.SetUsageTemplate(completionUsageTemplate)
	cmd.SetHelpTemplate(completionHelpTemplate)
	cmd.SetArgs(cobra.ExactArgs(0))
	cmd.SetSilenceErrors(true)
	cmd.AddCommand(
		NewCompletionBashCommand(cobra, output),
		NewCompletionFishCommand(cobra, output),
		NewCompletionPowerShellCommand(cobra, output),
		NewCompletionZshCommand(cobra, output),
	)

	return cmd
}

const (
	// completionHelpTemplate is the help template of the completion command.
	completionHelpTemplate = `🔧 Generate the autocompletion script for the specified shell.

` + completionUsageTemplate
	// compleitonUsageTemplate is the usage template of the completion command.
	completionUsageTemplate = `Usage:
  jrp completion [flags]
  jrp completion [command]

Available Subommands:
  bash        🔧🐚 Generate the autocompletion script for the bash shell.
  fish        🔧🐟 Generate the autocompletion script for the fish shell.
  powershell  🔧🪟 Generate the autocompletion script for the powershell shell.
  zsh         🔧🧙 Generate the autocompletion script for the zsh shell.

Flags:
  -h, --help  🤝 help for completion
`
)
