package constant

const (
	COMPLETION_ZSH_HELP_TEMPLATE = `🔧🧙 Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need to enable it.

You can execute the following once:

  echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

  source <(jrp completion zsh)

To load completions for every new session, execute once:

  - 🐧 Linux:

    jrp completion zsh > "${fpath[1]}/_jrp"

  - 🍎 macOS:

    jrp completion zsh > $(brew --prefix)/share/zsh/site-functions/_jrp

You will need to start a new shell for this setup to take effect.

Usage:
  jrp completion zsh [flags]

Flags:
  -h, --help  🤝 help for zsh
`
	COMPLETION_ZSH_USE   = "zsh"
	COMPLETION_ZSH_SHORT = "🔧🧙 Generate the autocompletion script for the zsh shell."
	COMPLETION_ZSH_LONG  = `🔧🧙 Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need to enable it.

You can execute the following once:

  echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

  source <(jrp completion zsh)

To load completions for every new session, execute once:

  - 🐧 Linux:

    jrp completion zsh > "${fpath[1]}/_jrp"

  - 🍎 macOS:

    jrp completion zsh > $(brew --prefix)/share/zsh/site-functions/_jrp

You will need to start a new shell for this setup to take effect.
`
)
