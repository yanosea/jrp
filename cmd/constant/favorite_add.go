package constant

const (
	FAVORITE_ADD_USE           = "add"
	FAVORITE_ADD_HELP_TEMPLATE = `⭐📌 Favorite (add) phrase(s) in the history of the "generate" command.

You can specify the phrase(s) to favorite with ID argument(s).
You have to get ID from the "history" command.
Multiple ID's can be specified separated by spaces.

Usage:
  jrp favorite add [flag]
  jrp favorite ad  [flag]
  jrp favorite a   [flag]

Flags:
  -h, --help  🤝 help for add

Arguments:
  ID  🆔 add the favorite by the ID (e.g: 1 2 3).
`
	FAVORITE_ADD_MESSAGE_NO_ID_SPECIFIED    = "⚡ No ID argument(s) specified..."
	FAVORITE_ADD_MESSAGE_ADDED_SUCCESSFULLY = "✅ Favorite successfully!"
	FAVORITE_ADD_MESSAGE_ADDED_FAILURE      = "❌ Failed favorite..."
	FAVORITE_ADD_MESSAGE_ADDED_NONE         = "⚡ No phrase(s) to favorite..."
	FAVORITE_ADD_MESSAGE_ADDED_NOT_ALL      = "⚡ Some phrase(s) are not favorited because the id does not exist or have already favorited..."
)

func GetFavoriteAddAliases() []string {
	return []string{"ad", "a"}
}
