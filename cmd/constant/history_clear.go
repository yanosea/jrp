package constant

const (
	HISTORY_CLEAR_USE           = "clear"
	HISTORY_CLEAR_HELP_TEMPLATE = `📜✨ Clear the history of the "generate" command.

You can clear all history.
Also, you can clear the history even if it is favorited by using the "-f" or ""--force" flag.

Usage:
  jrp history clear [flag]
  jrp history cl    [flag]
  jrp history c     [flag]

Flags:
  -f, --force  💪 clear all the history even if it is favorited
  -no-confirm  🚫 do not confirm before clearing the history
  -h, --help   🤝 help for clear
`
	HISTORY_CLEAR_FLAG_FORCE                  = "force"
	HISTORY_CLEAR_FLAG_FORCE_SHORTHAND        = "f"
	HISTORY_CLEAR_FLAG_FORCE_DEFAULT          = false
	HISTORY_CLEAR_FLAG_FORCE_DESCRIPTION      = "clear all the history even if it is favorited"
	HISTORY_CLEAR_FLAG_NO_CONFIRM             = "no-confirm"
	HISTORY_CLEAR_FLAG_NO_CONFIRM_SHORTHAND   = ""
	HISTORY_CLEAR_FLAG_NO_CONFIRM_DEFAULT     = false
	HISTORY_CLEAR_FLAG_NO_CONFIRM_DESCRIPTION = "do not confirm before clearing the history"

	HISTORY_CLEAR_MESSAGE_CLEARED_SUCCESSFULLY = "✅ Cleared history successfully!"
	HISTORY_CLEAR_MESSAGE_CLEARED_NONE         = "⚡ No history to clear..."
	HISTORY_CLEAR_MESSAGE_CLEARED_FAIRULE      = "❌ Failed to clear history..."
	HISTORY_CLEAR_MESSAGE_CLEAR_CANCELED       = "🚫 Cancelled clearing the history."
	HISTORY_CLEAR_PROMPT_LABEL                 = "Proceed with clearing the history? [y/N]"
)

func GetHistoryClearAliases() []string {
	return []string{"cl", "c"}
}
