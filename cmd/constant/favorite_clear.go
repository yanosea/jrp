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
  -h, --help  🤝 help for clear
`
	FAVORITE_CLEAR_USE   = "clear"
	FAVORITE_CLEAR_SHORT = "⭐✨ Clear the favorited phrase(s)."
	FAVORITE_CLEAR_LONG  = `⭐✨ Clear the favorited phrase(s).

You can clear all favorited phrase(s).
This is the same as the "favorite remove -a" command.
This does not remove the history of the "generate" command, just remove the favorite.
`
	FAVORITE_CLEAR_MESSAGE_CLEARED_SUCCESSFULLY = "✅ Cleared favorited phrase(s) successfully!"
	FAVORITE_CLEAR_MESSAGE_CLEARED_NONE         = "⚡ No favorited phrase(s) to clear..."
	FAVORITE_CLEAR_MESSAGE_CLEARED_FAIRULE      = "❌ Failed to clear favorited phraase(s)..."
)

func GetFavoriteClearAliases() []string {
	return []string{"cl", "c"}
}
