package constant

const (
	GENARETE_HELP_TEMPLATE = `✨ Generate Japanese random phrase(s).

You can specify how many phrases to generate by flag "-n" or "--number" or a number argument.
If both are provided, the larger number takes precedence.

And you can specify the prefix or suffix of the phrase(s) to generate
by the flag "-p" or "--prefix" and "-s" or "--suffix".

Usage:
  jrp generate [flags]
  jrp gen      [flags]
  jrp g        [flags]

Flags:
  -n, --number   🔢 number of phrases to generate (default 1, e.g: 10).
  -p  --prefix   💬 prefix of phrase(s) to generate.
  -s  --suffix   💬 suffix of phrase(s) to generate.
  -d  --dry-run  🧪 generate phrase(s) without saving to the history.
  -P, --plain    📝 plain text output instead of table output.
  -h, --help     🤝 help for generate

Arguments:
  number  🔢 number of phrases to generate (default 1, e.g: 10).
`
	GENERATE_USE   = "generate"
	GENERATE_SHORT = "✨ Generate Japanese random phrase(s)."
	GENERATE_LONG  = `✨ Generate Japanese random phrase(s).

You can specify how many phrases to generate by flag "-n" or "--number" or a number argument.
If both are provided, the larger number takes precedence.

And you can specify the prefix or suffix of the phrase(s) to generate
by the flag "-p" or "--prefix" and "-s" or "--suffix".
`
	GENERATE_FLAG_NUMBER              = "number"
	GENERATE_FLAG_NUMBER_SHORTHAND    = "n"
	GENERATE_FLAG_NUMBER_DEFAULT      = 1
	GENERATE_FLAG_NUMBER_DESCRIPTION  = "number of phrases to generate"
	GENERATE_FLAG_PREFIX              = "prifix"
	GENERATE_FLAG_PREFIX_SHORTHAND    = "p"
	GENERATE_FLAG_PREFIX_DEFAULT      = ""
	GENERATE_FLAG_PREFIX_DESCRIPTION  = "prefix of phrase(s) to generate"
	GENERATE_FLAG_SUFFIX              = "suffix"
	GENERATE_FLAG_SUFFIX_SHORTHAND    = "s"
	GENERATE_FLAG_SUFFIX_DEFAULT      = ""
	GENERATE_FLAG_SUFFIX_DESCRIPTION  = "suffix of phrase(s) to generate"
	GENERATE_FLAG_DRY_RUN             = "dry-run"
	GENERATE_FLAG_DRY_RUN_SHORTHAND   = "d"
	GENERATE_FLAG_DRY_RUN_DEFAULT     = false
	GENERATE_FLAG_DRY_RUN_DESCRIPTION = "generate phrase(s) without saving to the history"
	GENERATE_FLAG_PLAIN               = "plain"
	GENERATE_FLAG_PLAIN_SHORTHAND     = "P"
	GENERATE_FLAG_PLAIN_DEFAULT       = false
	GENERATE_FLAG_PLAIN_DESCRIPTION   = "plain text output instead of table output"

	GENERATE_MESSAGE_GENERATE_FAILURE         = "❌ Failed to generate the phrase(s)..."
	GENERATE_MESSAGE_NOTIFY_DOWNLOAD_REQUIRED = "⚡ You have to execute \"download\" to use jrp..."
	GENERATE_MESSAGE_NOTIFY_USE_ONLY_ONE      = "⚡ You can use only one of prefix or suffix..."
	GENERATE_MESSAGE_SAVED_FAILURE            = "❌ Failed to save the history..."
	GENERATE_MESSAGE_SAVED_NONE               = "⚡ No phrase(s) to save to the history..."
	GENERATE_MESSAGE_SAVED_NOT_ALL            = "⚡ Some phrase(s) are not saved to the history..."
)

func GetGenerateAliases() []string {
	return []string{"gen", "g"}
}