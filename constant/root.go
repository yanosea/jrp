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
  download, dl, d   📦 Download Japanese Wordnet sqlite3 database file from the official site.
  generate, gen, g  ✨ Generate Japanese random phrase(s). You can abbreviate "generate" sub command. ("jrp" and "jrp generate" are the same.)
  help              🤝 Help of jrp.
  completion        🔧 Generate the autocompletion script for the specified shell.
  version           🔖 Show the version of jrp.

Flags:
  -n, --number   🔢 number of phrases to generate (default 1, e.g: 10).
  -p  --prefix   💬 prefix of phrase(s) to generate.
  -s  --suffix   💬 suffix of phrase(s) to generate.
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
	ROOT_FLAG_NUMBER             = "number"
	ROOT_FLAG_NUMBER_SHORTHAND   = "n"
	ROOT_FLAG_NUMBER_DESCRIPTION = "number of phrases to generate"
	ROOT_FLAG_PREFIX             = "prifix"
	ROOT_FLAG_PREFIX_SHORTHAND   = "p"
	ROOT_FLAG_PREFIX_DESCRIPTION = "prefix of phrase(s) to generate"
	ROOT_FLAG_SUFFIX             = "suffix"
	ROOT_FLAG_SUFFIX_SHORTHAND   = "s"
	ROOT_FLAG_SUFFIX_DESCRIPTION = "suffix of phrase(s) to generate"
)
