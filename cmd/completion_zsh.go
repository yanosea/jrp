package cmd

import (
	"github.com/spf13/cobra"
)

const (
	completion_zsh_help_template = `🔧🧙 Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need to enable it.

You can execute the following once:

  echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

  source <(jrp completion zsh)

To load completions for every new session, execute once:

  * 🐧 Linux:

    jrp completion zsh > "${fpath[1]}/_jrp"

  * 🍎 macOS:

    jrp completion zsh > $(brew --prefix)/share/zsh/site-functions/_jrp

You will need to start a new shell for this setup to take effect.

Usage:
  jrp completion zsh [flags]

Flags:
  -h, --help   help for zsh
`
	completion_zsh_use   = "zsh"
	completion_zsh_short = "🔧🧙 Generate the autocompletion script for the zsh shell."
	completion_zsh_long  = `🔧🧙 Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need to enable it.

You can execute the following once:

  echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

  source <(jrp completion zsh)

To load completions for every new session, execute once:

  * 🐧 Linux:

    jrp completion zsh > "${fpath[1]}/_jrp"

  * 🍎 macOS:

    jrp completion zsh > $(brew --prefix)/share/zsh/site-functions/_jrp

You will need to start a new shell for this setup to take effect.`
)

func newCompletionZshCommand(globalOption *GlobalOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   completion_zsh_use,
		Short: completion_zsh_short,
		Long:  completion_zsh_long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.GenZshCompletion(globalOption.Out)
		},
	}

	cmd.SetOut(globalOption.Out)
	cmd.SetErr(globalOption.ErrOut)

	cmd.SetHelpTemplate(completion_zsh_help_template)

	return cmd
}
