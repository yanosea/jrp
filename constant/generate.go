package constant

const (
	GENARETE_HELP_TEMPLATE = `✨ Generate Japanese random phrase(s).

You can specify how many phrases to generate by flag "-n" or "--number" or a number argument.
If both are provided, the larger number takes precedence.

Usage:
  jrp generate [flags]
  jrp gen [flags]
  jrp g [flags]

Flags:
  -n, --number    🔢 number of phrases to generate (default 1). You can abbreviate "generate" sub command such ah (jrp -n 10, jrp 10).
  -h, --help      🤝 help for generate

Arguments:
  number  🔢 number of phrases to generate (e.g: 10).
`
	GENERATE_USE   = "generate"
	GENERATE_SHORT = "✨ Generate Japanese random phrase(s)."
	GENERATE_LONG  = `✨ Generate Japanese random phrase(s).

You can generate Japanese random phrase.
You can specify the number of phrases to generate by the flag "-n" or "--number".
`
	GENERATE_FLAG_NUMBER             = "number"
	GENERATE_FLAG_NUMBER_SHORTHAND   = "n"
	GENERATE_FLAG_NUMBER_DESCRIPTION = "number of phrases to generate"

	GENERATE_MESSAGE_NOTIFY_DOWNLOAD_REQUIRED = "⚡ You have to execute 'download' to use jrp..."
	GENERATE_SQL_GET_ALL_JAPANESE_AVN_WORDS   = "SELECT word.Lemma, word.Pos FROM word WHERE word.Lang = 'jpn' AND word.Pos in ('a', 'v', 'n')"
)

func GetGenerateAliases() []string {
	return []string{"gen", "g"}
}
