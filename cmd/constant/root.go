package constant

const (
	ROOT_HELP_TEMPLATE = `🎲 jrp is the CLI tool to generate Japanese random phrase(s).

You can generate Japanese random phrase(s).

You can specify how many phrases to generate by flag "-n" or "--number" or a number argument.
If both are provided, the larger number takes precedence.

And you can specify the prefix or suffix of the phrase(s) to generate
by the flag "-p" or "--prefix" and "-s" or "--suffix".

Usage:
  jrp [flags]
  jrp [command]

Available Subcommands:
  download, dl,   d  📦 Download WordNet Japan sqlite3 database file from the official site.
  generate, gen,  g  ✨ Generate Japanese random phrase(s). You can abbreviate "generate" sub command. ("jrp" and "jrp generate" are the same.)
  history,  hist, h  📜 Manage the history of the "generate" command.
  favorite, fav,  f  ⭐ Manage the favorited phrase(s) of the history of "generate" command.
  help               🤝 Help of jrp.
  completion         🔧 Generate the autocompletion script for the specified shell.
  version            🔖 Show the version of jrp.

Flags:
  -n, --number   🔢 number of phrases to generate (default 1, e.g: 10).
  -p  --prefix   💬 prefix of phrase(s) to generate.
  -s  --suffix   💬 suffix of phrase(s) to generate.
  -d  --dry-run  🧪 generate phrase(s) without saving to the history.
  -P, --plain    📝 plain text output instead of table output.
  -h, --help     🤝 help for jrp
  -v, --version  🔖 version for jrp

Arguments:
  number  🔢 number of phrases to generate (e.g: 10).

Use "jrp [command] --help" for more information about a command.
`
	ROOT_USE   = "jrp"
	ROOT_SHORT = "🎲 jrp is the CLI tool to generate Japanese random phrase(s)."
	ROOT_LONG  = `🎲 jrp is the CLI tool to generate Japanese random phrase(s).

You can generate Japanese random phrase(s).

You can specify how many phrases to generate by flag "-n" or "--number" or a number argument.
If both are provided, the larger number takes precedence.

And you can specify the prefix or suffix of the phrase(s) to generate
by the flag "-p" or "--prefix" and "-s" or "--suffix".
`
	ROOT_FLAG_NUMBER              = "number"
	ROOT_FLAG_NUMBER_SHORTHAND    = "n"
	ROOT_FLAG_NUMBER_DEFAULT      = 1
	ROOT_FLAG_NUMBER_DESCRIPTION  = "number of phrases to generate"
	ROOT_FLAG_PREFIX              = "prifix"
	ROOT_FLAG_PREFIX_SHORTHAND    = "p"
	ROOT_FLAG_PREFIX_DEFAULT      = ""
	ROOT_FLAG_PREFIX_DESCRIPTION  = "prefix of phrase(s) to generate"
	ROOT_FLAG_SUFFIX              = "suffix"
	ROOT_FLAG_SUFFIX_SHORTHAND    = "s"
	ROOT_FLAG_SUFFIX_DEFAULT      = ""
	ROOT_FLAG_SUFFIX_DESCRIPTION  = "suffix of phrase(s) to generate"
	ROOT_FLAG_DRY_RUN             = "dry-run"
	ROOT_FLAG_DRY_RUN_SHORTHAND   = "d"
	ROOT_FLAG_DRY_RUN_DEFAULT     = false
	ROOT_FLAG_DRY_RUN_DESCRIPTION = "generate phrase(s) without saving to the history"
	ROOT_FLAG_PLAIN               = "plain"
	ROOT_FLAG_PLAIN_SHORTHAND     = "P"
	ROOT_FLAG_PLAIN_DEFAULT       = false
	ROOT_FLAG_PLAIN_DESCRIPTION   = "plain text output instead of table output"

	ROOT_MESSAGE_GENERATE_FAILURE         = "❌ Failed to generate the phrase(s)..."
	ROOT_MESSAGE_NOTIFY_DOWNLOAD_REQUIRED = "⚡ You have to execute \"download\" to use jrp..."
	ROOT_MESSAGE_NOTIFY_USE_ONLY_ONE      = "⚡ You can use only one of prefix or suffix..."
	ROOT_MESSAGE_SAVED_FAILURE            = "❌ Failed to save the history..."
	ROOT_MESSAGE_SAVED_NONE               = "⚡ No phrase(s) to save to the history..."
	ROOT_MESSAGE_SAVED_NOT_ALL            = "⚡ Some phrase(s) are not saved to the history..."
)