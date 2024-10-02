package constant

const (
	FAVORITE_SHOW_HELP_TEMPLATE = `⭐📖 Show the favorited phrase(s).

You can specify how many phrases to show by flag "-n" or "--number" or a number argument.
jrp will get the most recent favorited phrase(s).
If you don't specify the number of phrases, jrp will show the most recent 10 phrases by default.
If both are provided, the larger number takes precedence.

Also, you can show all phrases in the favorite by flag "-a" or "--all".
If you use the flag, the number flag or argument will be ignored.

Usage:
  jrp favorite show [flag]
  jrp favorite sh   [flag]
  jrp favorite s    [flag]

Flags:
  -n, --number  📏 number how many phrases to show (default 10, e.g: 50).
  -a, --all     📁 show all phrases in the favorite.
  -P, --plain   📝 plain text output instead of table output.
  -h, --help    🤝 help for show

Arguments:
  number  📏 number how many phrases to show (default 10, e.g: 50).
`
	FAVORITE_SHOW_USE   = "show"
	FAVORITE_SHOW_SHORT = "⭐📖 Show the favorited phrase(s)."
	FAVORITE_SHOW_LONG  = `⭐📖 Show the favorited phrase(s).

You can specify how many phrases to show by flag "-n" or "--number" or a number argument.
jrp will get the most recent favorited phrase(s).
If you don't specify the number of phrases, jrp will show the most recent 10 phrases by default.

Also, you can show all phrases in the favorite by flag "-a" or "--all".
If you use the flag, the number flag or argument will be ignored.
`
	FAVORITE_SHOW_FLAG_NUMBER             = "number"
	FAVORITE_SHOW_FLAG_NUMBER_SHORTHAND   = "n"
	FAVORITE_SHOW_FLAG_NUMBER_DEFAULT     = 10
	FAVORITE_SHOW_FLAG_NUMBER_DESCRIPTION = "number how many phrases to show"
	FAVORITE_SHOW_FLAG_ALL                = "all"
	FAVORITE_SHOW_FLAG_ALL_SHORTHAND      = "a"
	FAVORITE_SHOW_FLAG_ALL_DEFAULT        = false
	FAVORITE_SHOW_FLAG_ALL_DESCRIPTION    = "show all phrases in the favorite"
	FAVORITE_SHOW_FLAG_PLAIN              = "plain"
	FAVORITE_SHOW_FLAG_PLAIN_SHORTHAND    = "P"
	FAVORITE_SHOW_FLAG_PLAIN_DEFAULT      = false
	FAVORITE_SHOW_FLAG_PLAIN_DESCRIPTION  = "plain text output instead of table output"

	FAVORITE_SHOW_MESSAGE_NO_FAVORITE_FOUND = "⚡ No favorited phrase(s) found..."
)

func GetFavoriteShowAliases() []string {
	return []string{"sh", "s"}
}
