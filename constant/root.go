package constant

const (
	ROOT_HELP_TEMPLATE = `🎲 jrp is the CLI tool to generate random Japanese phrases.

You can specify how many phrases to generate by flag "-n" or "--number" or a number argument.

Usage:
  jrp [flags]
  jrp [command]

Available Subcommands:
  download    📥 Download Japanese Wordnet sqlite3 database file from the official site.
  generate    ✨ Generate Japanese random phrases.
  help        🤝 Help of jrp.
  completion  🔧 Generate the autocompletion script for the specified shell.
  version     🔖 Show the version of jrp.

Flags:
  -n, --number    🔢 number of phrases to generate (default 1). You can abbreviate "generate" sub command.
  -h, --help      🤝 help for jrp
  -v, --version   🔖 version for jrp

Arguments:
  number  🔢 number of phrases to generate (e.g: 10).

Use "jrp [command] --help" for more information about a command.
`
	ROOT_USE   = "jrp"
	ROOT_SHORT = "🎲 jrp is the CLI tool to generate random japanese phrases."
	ROOT_LONG  = `🎲 jrp is the CLI tool to generate random japanese phrases.

You can specify how many phrases to generate.`
	ROOT_FLAG_NUMBER             = "number"
	ROOT_FLAG_NUMBER_SHORTHAND   = "n"
	ROOT_FLAG_NUMBER_DESCRIPTION = "number of phrases to generate"
)
