package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/util"
)

const (
	completion_help_template = `🔧 Generate the autocompletion script for the specified shell.

Usage:
  jrp completion [flags]
  jrp completion [command]

Available Commands:
  bash        🔧🐚 Generate the autocompletion script for the bash shell.
  fish        🔧🐟 Generate the autocompletion script for the fish shell.
  powershell  🔧🪟 Generate the autocompletion script for the powershell shell.
  zsh         🔧🧙 Generate the autocompletion script for the zsh shell.

Flags:
  -h, --help   help for completion

Use "jrp completion [command] --help" for more information about a command.
`
	completion_use   = "completion"
	completion_short = "🔧 Generate the autocompletion script for the specified shell."
	completion_long  = `🔧 Generate the autocompletion script for the specified shell.

See each sub-command's help for details on how to use the generated script.
You must use sub command below...

  * 🐚 bash
  * 🐟 fish
  * 🪟 powershell
  * 🧙 zsh`
	completion_message_no_sub_command = `Use sub command below...

  * 🐚 bash
  * 🐟 fish
  * 🪟 powershell
  * 🧙 zsh`
)

func newCompletionCommand(globalOption *GlobalOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   completion_use,
		Short: completion_short,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If no sub command is specified, print the message and return nil.
			util.PrintlnWithWriter(globalOption.Out, completion_message_no_sub_command)

			return nil
		},
	}

	cmd.SetOut(globalOption.Out)
	cmd.SetErr(globalOption.ErrOut)

	cmd.SetHelpTemplate(completion_help_template)

	cmd.AddCommand(
		newCompletionBashCommand(globalOption),
		newCompletionFishCommand(globalOption),
		newCompletionPowerShellCommand(globalOption),
		newCompletionZshCommand(globalOption),
	)

	return cmd
}
