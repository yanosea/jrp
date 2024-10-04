package constant

const (
	FAVORITE_SEARCH_USE           = "search"
	FAVORITE_SEARCH_HELP_TEMPLATE = `⭐🔍 Search the favorited phrase(s).

You can search favorited phrase(s) with keyword argument(s).
Multiple keywords are separated by a space.

If you want to search favorited phrase(s) by AND condition, you can use flag "-A" or "--and".
OR condition is by default.

You can specify how many results to show with flag "-n" or "--number".
If you don't specify the number of phrases, jrp will show the most recent 10 results by default.

Also, you can show all results by flag "-a" or "--all".
If you use the flag, the number flag will be ignored.

Usage:
  jrp favorite search [flag]
  jrp favorite se     [flag]
  jrp favorite S      [flag]

Flags:
  -A, --and     🧠 search phrase(s) by AND condition
  -n, --number  📏 number how many results to show (default 10, e.g: 50)
  -a, --all     📁 show all results
  -P, --plain   📝 plain text output instead of table output
  -h, --help    🤝 help for search

Arguments:
  keywords  💬 search phrase(s) by keywords (multiple keywords are separated by space)
`
	FAVORITE_SEARCH_FLAG_AND                = "and"
	FAVORITE_SEARCH_FLAG_AND_SHORTHAND      = "A"
	FAVORITE_SEARCH_FLAG_AND_DEFAULT        = false
	FAVORITE_SEARCH_FLAG_AND_DESCRIPTION    = "search phrase(s) by AND condition"
	FAVORITE_SEARCH_FLAG_NUMBER             = "number"
	FAVORITE_SEARCH_FLAG_NUMBER_SHORTHAND   = "n"
	FAVORITE_SEARCH_FLAG_NUMBER_DEFAULT     = 10
	FAVORITE_SEARCH_FLAG_NUMBER_DESCRIPTION = "number how many results to show"
	FAVORITE_SEARCH_FLAG_ALL                = "all"
	FAVORITE_SEARCH_FLAG_ALL_SHORTHAND      = "a"
	FAVORITE_SEARCH_FLAG_ALL_DEFAULT        = false
	FAVORITE_SEARCH_FLAG_ALL_DESCRIPTION    = "show all phrases in the favorite"
	FAVORITE_SEARCH_FLAG_PLAIN              = "plain"
	FAVORITE_SEARCH_FLAG_PLAIN_SHORTHAND    = "P"
	FAVORITE_SEARCH_FLAG_PLAIN_DEFAULT      = false
	FAVORITE_SEARCH_FLAG_PLAIN_DESCRIPTION  = "plain text output instead of table output"

	FAVORITE_SEARCH_MESSAGE_NO_KEYWORDS_PROVIDED = "⚡ No keyword(s) provided..."
	FAVORITE_SEARCH_MESSAGE_NO_RESULT_FOUND      = "⚡ No results found..."
)

func GetFavoriteSearchAliases() []string {
	return []string{"se", "S"}
}
