package constant

const (
	COMPLETION_USE           = "completion"
	COMPLETION_HELP_TEMPLATE = `🔧 Generate the autocompletion script for the specified shell.

Usage:
  jrp completion [flags]
  jrp completion [command]

Available Subommands:
  bash        🔧🐚 Generate the autocompletion script for the bash shell.
  fish        🔧🐟 Generate the autocompletion script for the fish shell.
  powershell  🔧🪟 Generate the autocompletion script for the powershell shell.
  zsh         🔧🧙 Generate the autocompletion script for the zsh shell.

Flags:
  -h, --help  🤝 help for completion

Use "jrp completion [command] --help" for more information about a command.
`
	COMPLETION_MESSAGE_NO_SUB_COMMAND = `Use sub command below...

  - 🐚 bash
  - 🐟 fish
  - 🪟 powershell
  - 🧙 zsh
`
)
