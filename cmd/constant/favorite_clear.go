package constant

const (
	FAVORITE_CLEAR_HELP_TEMPLATE = `⭐✨ Clear the favorited phrase(s).

You can clear all favorited phrase(s).
This is the same as the "favorite remove -a" command.
This does not clear the history of the "generate" command, just clear the favorite.

Usage:
  jrp favorite clear [flag]
  jrp favorite cl    [flag]
  jrp favorite c     [flag]

Flags:
  -no-confirm  🚫 do not confirm before clearing the history
  -h, --help   🤝 help for clear
`
	FAVORITE_CLEAR_USE   = "clear"
	FAVORITE_CLEAR_SHORT = "⭐✨ Clear the favorited phrase(s)."
	FAVORITE_CLEAR_LONG  = `⭐✨ Clear the favorited phrase(s).

You can clear all favorited phrase(s).
This is the same as the "favorite remove -a" command.
This does not remove the history of the "generate" command, just remove the favorite.
`
	FAVORITE_CLEAR_FLAG_NO_CONFIRM             = "no-confirm"
	FAVORITE_CLEAR_FLAG_NO_CONFIRM_SHORTHAND   = ""
	FAVORITE_CLEAR_FLAG_NO_CONFIRM_DEFAULT     = false
	FAVORITE_CLEAR_FLAG_NO_CONFIRM_DESCRIPTION = "do not confirm before clearing the favorited phrase(s)"

	FAVORITE_CLEAR_MESSAGE_CLEARED_SUCCESSFULLY = "✅ Cleared favorited phrase(s) successfully!"
	FAVORITE_CLEAR_MESSAGE_CLEARED_NONE         = "⚡ No favorited phrase(s) to clear..."
	FAVORITE_CLEAR_MESSAGE_CLEARED_FAIRULE      = "❌ Failed to clear favorited phraase(s)..."
	FAVORITE_CLEAR_MESSAGE_CLEAR_CANCELED       = "🚫 Cancelled clearing the favorited phrases(s)."
	FAVORITE_CLEAR_PROMPT_LABEL                 = "Proceed with clearing the favorited phrases(s)? [y/N]"
)

func GetFavoriteClearAliases() []string {
	return []string{"cl", "c"}
}
