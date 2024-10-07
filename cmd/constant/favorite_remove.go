package constant

const (
	FAVORITE_REMOVE_USE           = "remove"
	FAVORITE_REMOVE_HELP_TEMPLATE = `⭐🧹 Remove the favorited phrase(s).

You can specify the favorited phrase(s) to remove with ID argument(s).
You have to get ID from the "history" command.
Multiple ID's can be specified separated by spaces.
This does not remove the history of the "generate" command, just remove the favorited status.

Also, you can remove all favorited phrase(s) by flag "-a" or "--all".
This is the same as the "favorite clear" command.

Usage:
  jrp favorite remove [flag]
  jrp favorite rm     [flag]
  jrp favorite r      [flag]

Flags:
  -a, --all    ✨ remove all favorited phrase(s)
  -no-confirm  🚫 do not confirm before removing all the favorited phrase(s)
  -h, --help   🤝 help for remove

Arguments:
  ID  🆔 remove the favorite by the ID (e.g: 1 2 3)
`
	FAVORITE_REMOVE_FLAG_ALL                    = "all"
	FAVORITE_REMOVE_FLAG_ALL_SHORTHAND          = "a"
	FAVORITE_REMOVE_FLAG_ALL_DEFAULT            = false
	FAVORITE_REMOVE_FLAG_ALL_DESCRIPTION        = "remove all favorited phrase(s)"
	FAVORITE_REMOVE_FLAG_NO_CONFIRM             = "no-confirm"
	FAVORITE_REMOVE_FLAG_NO_CONFIRM_SHORTHAND   = ""
	FAVORITE_REMOVE_FLAG_NO_CONFIRM_DEFAULT     = false
	FAVORITE_REMOVE_FLAG_NO_CONFIRM_DESCRIPTION = "do not confirm before removing all the favorited phrase(s)"

	FAVORITE_REMOVE_MESSAGE_NO_ID_SPECIFIED      = "⚡ No ID argument(s) specified..."
	FAVORITE_REMOVE_MESSAGE_REMOVED_SUCCESSFULLY = "✅ Removed favorited phrase(s) successfully!"
	FAVORITE_REMOVE_MESSAGE_REMOVED_FAILURE      = "❌ Failed to remove favorited phrase(s)..."
	FAVORITE_REMOVE_MESSAGE_REMOVED_NONE         = "⚡ No favorited phrase(s) to remove..."
	FAVORITE_REMOVE_MESSAGE_REMOVED_NOT_ALL      = "⚡ Some favorited phrase(s) was not removed because the id does not exist or have not favorited..."
	FAVORITE_REMOVE_MESSAGE_REMOVE_ALL_CANCELED  = "🚫 Cancelled removing all the favorited phrase(s)."
	FAVORITE_REMOVE_PROMPT_REMOVE_ALL_LABEL      = "Proceed with removing all the favorited phrase(s)? [y/N]"
)

func GetFavoriteRemoveAliases() []string {
	return []string{"rm", "r"}
}
