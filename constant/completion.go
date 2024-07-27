package constant

const (
	COMPLETION_HELP_TEMPLATE = `🔧 Generate the autocompletion script for the specified shell.

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
	COMPLETION_USE   = "completion"
	COMPLETION_SHORT = "🔧 Generate the autocompletion script for the specified shell."
	COMPLETION_LONG  = `🔧 Generate the autocompletion script for the specified shell.

See each sub-command's help for details on how to use the generated script.
You must use sub command below...

  * 🐚 bash
  * 🐟 fish
  * 🪟 powershell
  * 🧙 zsh`
	COMPLETION_MESSAGE_NO_SUB_COMMAND = `Use sub command below...

  * 🐚 bash
  * 🐟 fish
  * 🪟 powershell
  * 🧙 zsh`
)
