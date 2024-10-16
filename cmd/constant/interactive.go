package constant

const (
	INTERACTIVE_USE           = "interactive"
	INTERACTIVE_HELP_TEMPLATE = `💬 Generate Japanese random phrase(s) interactively.

You can specify the prefix or suffix of the phrase(s) to generate
by the flag "-p" or "--prefix" and "-s" or "--suffix".

And you can choose to save or favorite the phrase(s) generated interactively.

Press the key for your action:
  "u"   : Favorite, continue.
  "i"   : Favorite, exit.
  "j"   : Save, continue.
  "k"   : Save, exit.
  "m"   : Skip, continue.
  other : Skip, exit.

Usage:
  jrp interactive [flags]
  jrp int         [flags]
  jrp i           [flags]

Flags:
  -p  --prefix   💬 prefix of phrase(s) to generate
  -s  --suffix   💬 suffix of phrase(s) to generate
  -P, --plain    📝 plain text output instead of table output
  -t  --timeout  ⏱️  timeout second for the interactive mode (default 30, e.g: 10)
  -h, --help     🤝 help for interactive
`
	INTERACTIVE_FLAG_PREFIX              = "prifix"
	INTERACTIVE_FLAG_PREFIX_SHORTHAND    = "p"
	INTERACTIVE_FLAG_PREFIX_DEFAULT      = ""
	INTERACTIVE_FLAG_PREFIX_DESCRIPTION  = "prefix of phrase(s) to generate"
	INTERACTIVE_FLAG_SUFFIX              = "suffix"
	INTERACTIVE_FLAG_SUFFIX_SHORTHAND    = "s"
	INTERACTIVE_FLAG_SUFFIX_DEFAULT      = ""
	INTERACTIVE_FLAG_SUFFIX_DESCRIPTION  = "suffix of phrase(s) to generate"
	INTERACTIVE_FLAG_PLAIN               = "plain"
	INTERACTIVE_FLAG_PLAIN_SHORTHAND     = "P"
	INTERACTIVE_FLAG_PLAIN_DEFAULT       = false
	INTERACTIVE_FLAG_PLAIN_DESCRIPTION   = "plain text output instead of table output"
	INTERACTIVE_FLAG_TIMEOUT             = "timeout"
	INTERACTIVE_FLAG_TIMEOUT_SHORTHAND   = "t"
	INTERACTIVE_FLAG_TIMEOUT_DEFAULT     = 30
	INTERACTIVE_FLAG_TIMEOUT_DESCRIPTION = "timeout second for the interactive mode (default 30, e.g: 10)"

	INTERACTIVE_MESSAGE_GENERATE_FAILURE         = "❌ Failed to generate the phrase(s)..."
	INTERACTIVE_MESSAGE_NOTIFY_DOWNLOAD_REQUIRED = "⚡ You have to execute \"download\" to use jrp..."
	INTERACTIVE_MESSAGE_NOTIFY_USE_ONLY_ONE      = "⚡ You can use only one of prefix or suffix..."
	INTERACTIVE_MESSAGE_SAVED_SUCCESSFULLY       = "✅ Saved successfully!"
	INTERACTIVE_MESSAGE_SAVED_FAILURE            = "❌ Failed to save the history..."
	INTERACTIVE_MESSAGE_SAVED_NONE               = "⚡ No phrase(s) to save to the history..."
	INTERACTIVE_MESSAGE_SAVED_NOT_ALL            = "⚡ Some phrase(s) are not saved to the history..."
	INTERACTIVE_MESSAGE_FAVORITED_SUCCESSFULLY   = "✅ Favorite successfully!"
	INTERACTIVE_MESSAGE_FAVORITED_FAILURE        = "❌ Failed favorite..."
	INTERACTIVE_MESSAGE_FAVORITED_NONE           = "⚡ No phrase(s) to favorite..."
	INTERACTIVE_MESSAGE_FAVORITED_NOT_ALL        = "⚡ Some phrase(s) are not favorited because the id does not exist or have already favorited..."
	INTERACTIVE_MESSAGE_EXIT                     = "🚪Exit!"
	INTERACTIVE_MESSAGE_SKIP                     = "⏭️  Skip!"
	INTERACTIVE_MESSAGE_PHASE                    = "🔄 Phase : "
	INTERACTIVE_PROMPT_LABEL                     = `🔽 Press the key for your action:
  "u"   : Favorite, continue.
  "i"   : Favorite, exit.
  "j"   : Save, continue.
  "k"   : Save, exit.
  "m"   : Skip, continue.
  other : Skip, exit.
`
)

func GetInteractiveAliases() []string {
	return []string{"int", "i"}
}

// InteractiveAnswer is a type for interactive answer.
type InteractiveAnswer int

const (
	// InteractiveAnswerSaveAndFavoriteAndContinue is a constant for save, favorite, and continue.
	InteractiveAnswerSaveAndFavoriteAndContinue InteractiveAnswer = iota
	// InteractiveAnswerSaveAndFavoriteAndExit is a constant for save, favorite, and exit.
	InteractiveAnswerSaveAndFavoriteAndExit
	// InteractiveAnswerSaveAndContinue is a constant for save and continue.
	InteractiveAnswerSaveAndContinue
	// InteractiveAnswerSaveAndExit is a constant for save and exit.
	InteractiveAnswerSaveAndExit
	// InteractiveAnswerSkipAndContinue is a constant for skip and continue.
	InteractiveAnswerSkipAndContinue
	// InteractiveAnswerSkipAndExit is a constant for skip and exit.
	InteractiveAnswerSkipAndExit
)
