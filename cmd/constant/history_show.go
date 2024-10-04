package constant

const (
	HISTORY_SHOW_USE           = "show"
	HISTORY_SHOW_HELP_TEMPLATE = `📜📖 Show the history of the "generate" command.

You can specify how many phrases to show by flag "-n" or "--number" or a number argument.
jrp will get the most recent phrase(s) from the history.
If you don't specify the number of phrases, jrp will show the most recent 10 phrases by default.
If both are provided, the larger number takes precedence.

Also, you can show all phrases in the history by flag "-a" or "--all".
If you use the flag, the number flag or argument will be ignored.

Usage:
  jrp history show [flag]
  jrp history sh   [flag]
  jrp history s    [flag]

Flags:
  -n, --number  📏 number how many phrases to show (default 10, e.g: 50)
  -a, --all     📁 show all history
  -P, --plain   📝 plain text output instead of table output
  -h, --help    🤝 help for show

Arguments:
  number  📏 number how many phrases to show (default 10, e.g: 50)
`
	HISTORY_SHOW_FLAG_NUMBER             = "number"
	HISTORY_SHOW_FLAG_NUMBER_SHORTHAND   = "n"
	HISTORY_SHOW_FLAG_NUMBER_DEFAULT     = 10
	HISTORY_SHOW_FLAG_NUMBER_DESCRIPTION = "number how many phrases to show"
	HISTORY_SHOW_FLAG_ALL                = "all"
	HISTORY_SHOW_FLAG_ALL_SHORTHAND      = "a"
	HISTORY_SHOW_FLAG_ALL_DEFAULT        = false
	HISTORY_SHOW_FLAG_ALL_DESCRIPTION    = "show all phrases in the history"
	HISTORY_SHOW_FLAG_PLAIN              = "plain"
	HISTORY_SHOW_FLAG_PLAIN_SHORTHAND    = "P"
	HISTORY_SHOW_FLAG_PLAIN_DEFAULT      = false
	HISTORY_SHOW_FLAG_PLAIN_DESCRIPTION  = "plain text output instead of table output"

	HISTORY_SHOW_MESSAGE_NO_HISTORY_FOUND = "⚡ No history found..."
)

func GetHistoryShowAliases() []string {
	return []string{"sh", "s"}
}
