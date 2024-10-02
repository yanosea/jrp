package constant

const (
	HISTORY_SEARCH_HELP_TEMPLATE = `📜🔍 Search the history of the "generate" command.

You can search phrase(s) with keyword argument(s).
Multiple keywords are separated by a space.

If you want to search phrase(s) by AND condition, you can use flag "-A" or "--and".
OR condition is by default.

You can specify how many results to show with flag "-n" or "--number".
If you don't specify the number of phrases, jrp will show the most recent 10 results by default.

Also, you can show all results in the history by flag "-a" or "--all".
If you use the flag, the number flag will be ignored.

Usage:
  jrp history search [flag]
  jrp history se     [flag]
  jrp history S      [flag]

Flags:
  -A, --and     🧠 search phrase(s) by AND condition.
  -n, --number  📏 number how many results to show (default 10, e.g: 50).
  -a, --all     📁 show all results in the history.
  -P, --plain   📝 plain text output instead of table output.
  -h, --help    🤝 help for search

Arguments:
  keywords  💬 search phrase(s) by keywords. Multiple keywords are separated by space.
`
	HISTORY_SEARCH_USE   = "search"
	HISTORY_SEARCH_SHORT = "📜🔍 Search the history of the \"generate\" command."
	HISTORY_SEARCH_LONG  = `📜🔍 Search the history of the "generate" command.

You can search phrase(s) with keyword argument(s).
Multiple keywords are separated by a space.

If you want to search phrase(s) by AND condition, you can use flag "-A" or "--and".
OR condition is by default.

You can specify how many results to show with flag "-n" or "--number".
If you don't specify the number of phrases, jrp will show the most recent 10 results by default.

Also, you can show all results in the history by flag "-a" or "--all".
If you use the flag, the number flag will be ignored.
`
	HISTORY_SEARCH_FLAG_AND                = "and"
	HISTORY_SEARCH_FLAG_AND_SHORTHAND      = "A"
	HISTORY_SEARCH_FLAG_AND_DEFAULT        = false
	HISTORY_SEARCH_FLAG_AND_DESCRIPTION    = "search phrase(s) by AND condition"
	HISTORY_SEARCH_FLAG_NUMBER             = "number"
	HISTORY_SEARCH_FLAG_NUMBER_SHORTHAND   = "n"
	HISTORY_SEARCH_FLAG_NUMBER_DEFAULT     = 10
	HISTORY_SEARCH_FLAG_NUMBER_DESCRIPTION = "number how many results to show"
	HISTORY_SEARCH_FLAG_ALL                = "all"
	HISTORY_SEARCH_FLAG_ALL_SHORTHAND      = "a"
	HISTORY_SEARCH_FLAG_ALL_DEFAULT        = false
	HISTORY_SEARCH_FLAG_ALL_DESCRIPTION    = "show all phrases in the history"
	HISTORY_SEARCH_FLAG_PLAIN              = "plain"
	HISTORY_SEARCH_FLAG_PLAIN_SHORTHAND    = "P"
	HISTORY_SEARCH_FLAG_PLAIN_DEFAULT      = false
	HISTORY_SEARCH_FLAG_PLAIN_DESCRIPTION  = "plain text output instead of table output"

	HISTORY_SEARCH_MESSAGE_NO_KEYWORDS_PROVIDED = "⚡ No keyword(s) provided..."
	HISTORY_SEARCH_MESSAGE_NO_RESULT_FOUND      = "⚡ No results found..."
)

func GetHistorySearchAliases() []string {
	return []string{"se", "S"}
}
