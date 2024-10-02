package constant

const (
	COMPLETION_BASH_HELP_TEMPLATE = `🔧🐚 Generate the autocompletion script for the bash shell.

This script depends on the "bash-completion" package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

  source <(jrp completion bash)

To load completions for every new session, execute once:

  - 🐧 Linux:

    jrp completion bash > /etc/bash_completion.d/jrp

  - 🍎 macOS:

    jrp completion bash > $(brew --prefix)/etc/bash_completion.d/jrp

You will need to start a new shell for this setup to take effect.

Usage:
  jrp completion bash [flags]

Flags:
  -h, --help  🤝 help for bash
`
	COMPLETION_BASH_USE   = "bash"
	COMPLETION_BASH_SHORT = "🔧🐚 Generate the autocompletion script for the bash shell."
	COMPLETION_BASH_LONG  = `🔧🐚 Generate the autocompletion script for the bash shell.

This script depends on the "bash-completion" package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

  source <(jrp completion bash)

To load completions for every new session, execute once:

  - 🐧 Linux:

    jrp completion bash > /etc/bash_completion.d/jrp

  - 🍎 macOS:

    jrp completion bash > $(brew --prefix)/etc/bash_completion.d/jrp

You will need to start a new shell for this setup to take effect.
`
)
