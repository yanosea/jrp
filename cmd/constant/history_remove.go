package constant

const (
	HISTORY_REMOVE_USE           = "remove"
	HISTORY_REMOVE_HELP_TEMPLATE = `📜🧹 Remove the history of the "generate" command.

You can specify the history to remove with ID argument(s).
You have to get ID from the "history" command.
Multiple ID's can be specified separated by spaces.

You can remove all history by flag "-a" or "--all".
This is the same as the "history clear" command.

Also, you can remove the history even if it is favorited by using the "-f" or ""--force" flag.

Usage:
  jrp history remove [flag]
  jrp history rm     [flag]
  jrp history r      [flag]

Flags:
  -a, --all    ✨ remove all history
  -f, --force  💪 remove the history even if it is favorited
  -no-confirm  🚫 do not confirm before removing all the history
  -h, --help   🤝 help for remove

Arguments:
  ID  🆔 remove the history by the ID (e.g: 1 2 3)
`
	HISTORY_REMOVE_FLAG_ALL                    = "all"
	HISTORY_REMOVE_FLAG_ALL_SHORTHAND          = "a"
	HISTORY_REMOVE_FLAG_ALL_DEFAULT            = false
	HISTORY_REMOVE_FLAG_ALL_DESCRIPTION        = "remove all history"
	HISTORY_REMOVE_FLAG_FORCE                  = "force"
	HISTORY_REMOVE_FLAG_FORCE_SHORTHAND        = "f"
	HISTORY_REMOVE_FLAG_FORCE_DEFAULT          = false
	HISTORY_REMOVE_FLAG_FORCE_DESCRIPTION      = "force remove the history even if it is favorited"
	HISTORY_REMOVE_FLAG_NO_CONFIRM             = "no-confirm"
	HISTORY_REMOVE_FLAG_NO_CONFIRM_SHORTHAND   = ""
	HISTORY_REMOVE_FLAG_NO_CONFIRM_DEFAULT     = false
	HISTORY_REMOVE_FLAG_NO_CONFIRM_DESCRIPTION = "do not confirm before removing all the history"

	HISTORY_REMOVE_MESSAGE_NO_ID_SPECIFIED      = "⚡ No ID argument(s) specified..."
	HISTORY_REMOVE_MESSAGE_REMOVED_SUCCESSFULLY = "✅ Removed the history successfully!"
	HISTORY_REMOVE_MESSAGE_REMOVED_FAILURE      = "❌ Failed to remove the history..."
	HISTORY_REMOVE_MESSAGE_REMOVED_NONE         = "⚡ No history to remove..."
	HISTORY_REMOVE_MESSAGE_REMOVED_NOT_ALL      = "⚡ Some phtase(s) was not removed because the id does not exist or have already favorited..."
	HISTORY_REMOVE_MESSAGE_REMOVE_ALL_CANCELED  = "🚫 Cancelled removing all the history."
	HISTORY_REMOVE_PROMPT_REMOVE_ALL_LABEL      = "Proceed with removing all the history? [y/N]"
)

func GetHistoryRemoveAliases() []string {
	return []string{"rm", "r"}
}
