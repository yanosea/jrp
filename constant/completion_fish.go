package constant

const (
	COMPLETION_FISH_HELP_TEMPLATE = `🔧🐟 Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

  jrp completion fish | source

To load completions for every new session, execute once:

  jrp completion fish > ~/.config/fish/completions/jrp.fish

You will need to start a new shell for this setup to take effect.

Usage:
  jrp completion fish [flags]

Flags:
  -h, --help  🤝 help for fish
`
	COMPLETION_FISH_USE   = "fish"
	COMPLETION_FISH_SHORT = "🔧🐟 Generate the autocompletion script for the fish shell."
	COMPLETION_FISH_LONG  = `🔧🐟 Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

  jrp completion fish | source

To load completions for every new session, execute once:

  jrp completion fish > ~/.config/fish/completions/jrp.fish

You will need to start a new shell for this setup to take effect.`
)
