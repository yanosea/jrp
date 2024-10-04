package constant

const (
	HISTORY_REMOVE_USE           = "remove"
	HISTORY_REMOVE_HELP_TEMPLATE = `📜🧹 Remove the history of the "generate" command.

You can specify the history to remove with ID argument(s).
You have to get ID from the "history" command.
Multiple ID's can be specified separated by spaces.

Also, you can remove the history even if it is favorited by using the "-f" or ""--force" flag.

Usage:
  jrp history remove [flag]
  jrp history rm     [flag]
  jrp history r      [flag]

Flags:
  -f, --force  🧹 remove the history even if it is favorited
  -h, --help   🤝 help for remove

Arguments:
  ID  🆔 remove the history by the ID (e.g: 1 2 3).
`
	HISTORY_REMOVE_FLAG_FORCE             = "force"
	HISTORY_REMOVE_FLAG_FORCE_SHORTHAND   = "f"
	HISTORY_REMOVE_FLAG_FORCE_DEFAULT     = false
	HISTORY_REMOVE_FLAG_FORCE_DESCRIPTION = "force remove the history even if it is favorited"

	HISTORY_REMOVE_MESSAGE_NO_ID_SPECIFIED      = "⚡ No ID argument(s) specified..."
	HISTORY_REMOVE_MESSAGE_REMOVED_SUCCESSFULLY = "✅ Removed the history successfully!"
	HISTORY_REMOVE_MESSAGE_REMOVED_FAILURE      = "❌ Failed to remove the history..."
	HISTORY_REMOVE_MESSAGE_REMOVED_NONE         = "⚡ No history to remove..."
	HISTORY_REMOVE_MESSAGE_REMOVED_NOT_ALL      = "⚡ Some phtase(s) was not removed because the id does not exist or have already favorited..."
)

func GetHistoryRemoveAliases() []string {
	return []string{"rm", "r"}
}
