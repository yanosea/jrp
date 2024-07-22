package cmd

import (
	"github.com/spf13/cobra"
)

const (
	completion_bash_help_template = `🔧🐚 Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

  source <(jrp completion bash)

To load completions for every new session, execute once:

  * 🐧 Linux:

    jrp completion bash > /etc/bash_completion.d/jrp

  * 🍎 macOS:

    jrp completion bash > $(brew --prefix)/etc/bash_completion.d/jrp

You will need to start a new shell for this setup to take effect.

Usage:
  jrp completion bash [flags]

Flags:
  -h, --help   help for bash
`
	completion_bash_use   = "bash"
	completion_bash_short = "🔧🐚 Generate the autocompletion script for the bash shell."
	completion_bash_long  = `🔧🐚 Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

  source <(jrp completion bash)

To load completions for every new session, execute once:

  * 🐧 Linux:

    jrp completion bash > /etc/bash_completion.d/jrp

  * 🍎 macOS:

    jrp completion bash > $(brew --prefix)/etc/bash_completion.d/jrp

You will need to start a new shell for this setup to take effect.`
)

func newCompletionBashCommand(globalOption *GlobalOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   completion_bash_use,
		Short: completion_bash_short,
		Long:  completion_bash_long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.GenBashCompletion(globalOption.Out)
		},
	}

	cmd.SetOut(globalOption.Out)
	cmd.SetErr(globalOption.ErrOut)

	cmd.SetHelpTemplate(completion_bash_help_template)

	return cmd
}
