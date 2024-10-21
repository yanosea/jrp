package constant

const (
	FAVORITE_USE           = "favorite"
	FAVORITE_HELP_TEMPLATE = `⭐ Manage the favorited phrase(s) of the history of "generate" command.

You can favorite (add) generated phrase(s) with its ID(s).
Also, You can show, remove, search and clear the phrase(s) you favorited.

You can specify how many phrases to show by flag "-n" or "--number" or a number argument.
jrp will get the most recent favorited phrase(s).
If you don't specify the number of phrases, jrp will show the most recent 10 phrases by default.
If both are provided, the larger number takes precedence.

Also, you can show all phrases in the favorite by flag "-a" or "--all".
If you use the flag, the number flag or argument will be ignored.

Usage:
  jrp favorite [flag]
  jrp fav      [flag]
  jrp f        [flag]
  jrp favorite [command]
  jrp fav      [command]
  jrp f        [command]

Available Subommands:
  show,   sh, s  ⭐📖 Show the favorited phrase(s). You can abbreviate "show" sub command. ("jrp favorite" and "jrp favorite show" are the same.)
  add,    ad, a  ⭐📌 Favorite (add) phrase(s) in the history of the "generate" command.
  search, se, S  ⭐🔍 Search the favorited phrase(s).
  remove, rm, r  ⭐🧹 Remove the favorited phrase(s).
  clear,  cl, c  ⭐✨ Clear the favorited phrase(s).

Flags:
  -n, --number  🔢 number how many phrases to show (default 10, e.g: 50)
  -a, --all     📁 show all favorited phrase(s)
  -P, --plain   📝 plain text output instead of table output
  -h, --help    🤝 help for favorite

Arguments:
  number  🔢 number how many phrases to show (default 10, e.g: 50)

Use "jrp favorite [command] --help" for more information about a command.
`
	FAVORITE_FLAG_NUMBER             = "number"
	FAVORITE_FLAG_NUMBER_SHORTHAND   = "n"
	FAVORITE_FLAG_NUMBER_DEFAULT     = 10
	FAVORITE_FLAG_NUMBER_DESCRIPTION = "number how many phrases to show"
	FAVORITE_FLAG_ALL                = "all"
	FAVORITE_FLAG_ALL_SHORTHAND      = "a"
	FAVORITE_FLAG_ALL_DEFAULT        = false
	FAVORITE_FLAG_ALL_DESCRIPTION    = "show all phrases in the favorite"
	FAVORITE_FLAG_PLAIN              = "plain"
	FAVORITE_FLAG_PLAIN_SHORTHAND    = "P"
	FAVORITE_FLAG_PLAIN_DEFAULT      = false
	FAVORITE_FLAG_PLAIN_DESCRIPTION  = "plain text output instead of table output"

	FAVORITE_MESSAGE_NO_FAVORITE_FOUND = "⚡ No favorited phrase(s) found..."
)

func GetFavoriteAliases() []string {
	return []string{"fav", "f"}
}
