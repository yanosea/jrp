package constant

const (
	HISTORY_CLEAR_HELP_TEMPLATE = `📜✨ Clear the history of the "generate" command.

You can clear all history.
Also, you can clear the history even if it is favorited by using the "-f" or ""--force" flag.

Usage:
  jrp history clear [flag]
  jrp history cl    [flag]
  jrp history c     [flag]

Flags:
  -f, --force  🚫 force clear the history even if it is favorited
  -h, --help   🤝 help for clear
`
	HISTORY_CLEAR_USE   = "clear"
	HISTORY_CLEAR_SHORT = "📜✨ Clear the history of the \"generate\" command."
	HISTORY_CLEAR_LONG  = `📜✨ Clear the history of the "generate" command.

You can clear all history.
Also, you can clear the history even if it is favorited by using the "-f" or ""--force" flag.
`
	HISTORY_CLEAR_FLAG_FORCE             = "force"
	HISTORY_CLEAR_FLAG_FORCE_SHORTHAND   = "f"
	HISTORY_CLEAR_FLAG_FORCE_DEFAULT     = false
	HISTORY_CLEAR_FLAG_FORCE_DESCRIPTION = "force clear the history even if it is favorited"

	HISTORY_CLEAR_MESSAGE_CLEARED_SUCCESSFULLY = "✅ Cleared history successfully!"
	HISTORY_CLEAR_MESSAGE_CLEARED_NONE         = "⚡ No history to clear..."
	HISTORY_CLEAR_MESSAGE_CLEARED_FAIRULE      = "❌ Failed to clear history..."
)

func GetHistoryClearAliases() []string {
	return []string{"cl", "c"}
}
