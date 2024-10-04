package constant

const (
	HISTORY_USE           = "history"
	HISTORY_HELP_TEMPLATE = `📜 Manage the history of the "generate" command.

You can show, remove, search and clear the history of the "generate" command.

You can specify how many phrases to show by flag "-n" or "--number" or a number argument.
jrp will get the most recent phrase(s) from the history.
If you don't specify the number of phrases, jrp will show the most recent 10 phrases by default.
If both are provided, the larger number takes precedence.

Also, you can show all phrases in the history by flag "-a" or "--all".
If you use the flag, the number flag or argument will be ignored.

Usage:
  jrp history [flag]
  jrp hist    [flag]
  jrp h       [flag]
  jrp history [command]
  jrp hist    [command]
  jrp h       [command]

Available Subommands:
  show    📜📖 Show the history of the "generate" command. You can abbreviate "show" sub command. ("jrp history" and "jrp history show" are the same.)
  search  📜🔍 Search the history of the "generate" command.
  remove  📜🧹 Remove the history of the "generate" command.
  clear   📜✨ Clear the history of the "generate" command.

Flags:
  -n, --number  📏 number how many phrases to show (default 10, e.g: 50)
  -a, --all     📁 show all history
  -P, --plain   📝 plain text output instead of table output
  -h, --help    🤝 help for history

Arguments:
  number  📏 number how many phrases to show (default 10, e.g: 50)

Use "jrp history [command] --help" for more information about a command.
`
	HISTORY_FLAG_NUMBER             = "number"
	HISTORY_FLAG_NUMBER_SHORTHAND   = "n"
	HISTORY_FLAG_NUMBER_DEFAULT     = 10
	HISTORY_FLAG_NUMBER_DESCRIPTION = "number how many phrases to show"
	HISTORY_FLAG_ALL                = "all"
	HISTORY_FLAG_ALL_SHORTHAND      = "a"
	HISTORY_FLAG_ALL_DEFAULT        = false
	HISTORY_FLAG_ALL_DESCRIPTION    = "show all phrases in the history"
	HISTORY_FLAG_PLAIN              = "plain"
	HISTORY_FLAG_PLAIN_SHORTHAND    = "P"
	HISTORY_FLAG_PLAIN_DEFAULT      = false
	HISTORY_FLAG_PLAIN_DESCRIPTION  = "plain text output instead of table output"

	HISTORY_MESSAGE_NO_HISTORY_FOUND = "⚡ No history found..."
)

func GetHistoryAliases() []string {
	return []string{"hist", "h"}
}
