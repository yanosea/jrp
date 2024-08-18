package constant

const (
	GENARETE_HELP_TEMPLATE = `✨ Generate Japanese random phrase(s).

You can specify how many phrases to generate by flag "-n" or "--number" or a number argument.
If both are provided, the larger number takes precedence.

And you can specify the prefix and suffix of the phrase to generate
by the flag "-p" or "--prefix" and "-s" or "--suffix".

Usage:
  jrp generate [flags]
  jrp gen [flags]
  jrp g [flags]

Flags:
  -n, --number    🔢 number of phrases to generate (default 1).
  -p  --prefix    💬 prefix of phrases to generate.
  -s  --suffix    💬 suffix of phrases to generate.
  -h, --help      🤝 help for generate

Arguments:
  number  🔢 number of phrases to generate (e.g: 10).
`
	GENERATE_USE   = "generate"
	GENERATE_SHORT = "✨ Generate Japanese random phrase(s)."
	GENERATE_LONG  = `✨ Generate Japanese random phrase(s).

You can generate Japanese random phrase.
`
	GENERATE_FLAG_NUMBER             = "number"
	GENERATE_FLAG_NUMBER_SHORTHAND   = "n"
	GENERATE_FLAG_NUMBER_DESCRIPTION = "number of phrases to generate"
	GENERATE_FLAG_PREFIX             = "prifix"
	GENERATE_FLAG_PREFIX_SHORTHAND   = "p"
	GENERATE_FLAG_PREFIX_DESCRIPTION = "prefix of phrases to generate"
	GENERATE_FLAG_SUFFIX             = "suffix"
	GENERATE_FLAG_SUFFIX_SHORTHAND   = "s"
	GENERATE_FLAG_SUFFIX_DESCRIPTION = "suffix of phrases to generate"

	GENERATE_MESSAGE_NOTIFY_DOWNLOAD_REQUIRED = "⚡ You have to execute 'download' to use jrp..."
	GENERATE_MESSAGE_NOTIFY_USE_ONLY_ONE      = "⚡ You can use only one of prefix or suffix..."

	GENERATE_SQL_GET_ALL_JAPANESE_AVN_WORDS = "SELECT word.Lemma, word.Pos FROM word WHERE word.Lang = 'jpn' AND word.Pos in ('a', 'v', 'n');"
	GENERATE_SQL_GET_ALL_JAPANESE_AV_WORDS  = "SELECT word.Lemma, word.Pos FROM word WHERE word.Lang = 'jpn' AND word.Pos in ('a', 'v');"
	GENERATE_SQL_GET_ALL_JAPANESE_N_WORDS   = "SELECT word.Lemma, word.Pos FROM word WHERE word.Lang = 'jpn' AND word.Pos = 'n';"
)

func GetGenerateAliases() []string {
	return []string{"gen", "g"}
}
