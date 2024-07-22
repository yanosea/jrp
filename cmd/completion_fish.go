package cmd

import (
	"github.com/spf13/cobra"
)

const (
	completion_fish_help_template = `🔧🐟 Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

  jrp completion fish | source

To load completions for every new session, execute once:

  jrp completion fish > ~/.config/fish/completions/jrp.fish

You will need to start a new shell for this setup to take effect.

Usage:
  jrp completion fish [flags]

Flags:
  -h, --help   help for fish
`
	completion_fish_use   = "fish"
	completion_fish_short = "🔧🐟 Generate the autocompletion script for the fish shell."
	completion_fish_long  = `🔧🐟 Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

  jrp completion fish | source

To load completions for every new session, execute once:

  jrp completion fish > ~/.config/fish/completions/jrp.fish

You will need to start a new shell for this setup to take effect.`
)

func newCompletionFishCommand(globalOption *GlobalOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   completion_fish_use,
		Short: completion_fish_short,
		Long:  completion_fish_long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.GenFishCompletion(globalOption.Out, false)
		},
	}

	cmd.SetOut(globalOption.Out)
	cmd.SetErr(globalOption.ErrOut)

	cmd.SetHelpTemplate(completion_fish_help_template)

	return cmd
}
