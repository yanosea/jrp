package command

import (
	c "github.com/spf13/cobra"

	"github.com/yanosea/jrp/v2/app/presentation/cli/jrp/command/jrp"
	"github.com/yanosea/jrp/v2/app/presentation/cli/jrp/command/jrp/completion"
	"github.com/yanosea/jrp/v2/app/presentation/cli/jrp/command/jrp/generate"
	"github.com/yanosea/jrp/v2/app/presentation/cli/jrp/command/jrp/history"
	"github.com/yanosea/jrp/v2/app/presentation/cli/jrp/config"

	"github.com/yanosea/jrp/v2/pkg/proxy"
)

// RootOptions provides the options for the root command.
type RootOptions struct {
	// Version is a flag to show the version of jrp.
	Version bool
	// GenerateOptions provides the options for the generate command.
	GenerateOptions generate.GenerateOptions
}

var (
	// rootOps is a variable to store the root options with the default values for injecting the dependencies in testing.
	rootOps = RootOptions{
		Version: false,
		GenerateOptions: generate.GenerateOptions{
			Number:      1,
			Prefix:      "",
			Suffix:      "",
			DryRun:      false,
			Format:      "table",
			Interactive: false,
			Timeout:     30,
		},
	}
)

// NewRootCommand returns a new instance of the root command.
func NewRootCommand(
	cobra proxy.Cobra,
	version string,
	conf *config.JrpCliConfig,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("jrp")
	cmd.SetUsageTemplate(rootUsageTemplate)
	cmd.SetHelpTemplate(rootHelpTemplate)
	cmd.SetArgs(cobra.MaximumNArgs(1))
	cmd.SetSilenceErrors(true)
	cmd.PersistentFlags().BoolVarP(
		&rootOps.Version,
		"version",
		"v",
		false,
		"🔖 show the version of jrp",
	)
	cmd.Flags().IntVarP(
		&rootOps.GenerateOptions.Number,
		"number",
		"n",
		1,
		"🔢 number of phrases to generate (default 1, e.g. : 10)",
	)
	cmd.Flags().StringVarP(
		&rootOps.GenerateOptions.Prefix,
		"prefix",
		"p",
		"",
		"🔡 prefix of phrases to generate",
	)
	cmd.Flags().StringVarP(
		&rootOps.GenerateOptions.Suffix,
		"suffix",
		"s",
		"",
		"🔡 suffix of phrases to generate",
	)
	cmd.Flags().BoolVarP(
		&rootOps.GenerateOptions.DryRun,
		"dry-run",
		"d",
		false,
		"🧪 generate phrases without saving to the history",
	)
	cmd.Flags().StringVarP(
		&rootOps.GenerateOptions.Format,
		"format",
		"f",
		"table",
		"📝 format of the output (default \"table\", e.g. : \"plain\")",
	)
	cmd.Flags().BoolVarP(
		&rootOps.GenerateOptions.Interactive,
		"interactive",
		"i",
		false,
		"💬 generate Japanese random phrases interactively",
	)
	cmd.Flags().IntVarP(
		&rootOps.GenerateOptions.Timeout,
		"timeout",
		"t",
		30,
		"⌛ timeout in seconds for the interactive mode (default 30, e.g. : 10)",
	)
	interactiveCmd := generate.NewInteractiveCommand(
		cobra,
		output,
	)
	generateCmd := generate.NewGenerateCommand(
		cobra,
		interactiveCmd,
		output,
	)
	versionCmd := jrp.NewVersionCommand(
		cobra,
		version,
		output,
	)
	cmd.AddCommand(
		completion.NewCompletionCommand(
			cobra,
			output,
		),
		jrp.NewDownloadCommand(
			cobra,
			conf,
			output,
		),
		jrp.NewFavoriteCommand(
			cobra,
			output,
		),
		generateCmd,
		history.NewHistoryCommand(
			cobra,
			output,
		),
		interactiveCmd,
		jrp.NewUnfavoriteCommand(
			cobra,
			output,
		),
		versionCmd,
	)

	cmd.SetRunE(
		func(cmd *c.Command, args []string) error {
			return runRoot(cmd, args, generateCmd, versionCmd)
		},
	)

	return cmd
}

// runRoot runs the root command.
func runRoot(
	cmd *c.Command,
	args []string,
	generateCmd proxy.Command,
	versionCmd proxy.Command,
) error {
	if rootOps.Version {
		return versionCmd.RunE(cmd, args)
	}
	generate.GenerateOps = rootOps.GenerateOptions
	return generateCmd.RunE(cmd, args)
}

const (
	// rootHelpTemplate is the help template of the root command.
	rootHelpTemplate = `🎲 jrp is the CLI jokeey tool to generate Japanese random phrases.

You can generate Japanese random phrases.

You can specify how many phrases to generate by flag "-n" or "--number" or a number argument.
If both are provided, the larger number takes precedence.

And you can specify the prefix or suffix of the phrases to generate
by the flag "-p" or "--prefix" and "-s" or "--suffix".

Those commands below are the same.
  "jrp" : "jrp generate"
  "jrp interactive" : "jrp --interactive" : "jrp generate interactive" : "jrp generate --interactive"

` + rootUsageTemplate
	// rootUsageTemplate is the usage template of the root command.
	rootUsageTemplate = `Usage:
  jrp [flags]
  jrp [command]
  jrp [argument]

Available Subcommands:
  download,    dl,   d  📦 Download WordNet Japan sqlite database file from the official web site.
  generate,    gen,  g  ✨ Generate Japanese random phrases.
                           You can abbreviate "generate" sub command. ("jrp" and "jrp generate" are the same.)
  interactive, int,  i  💬 Generate Japanese random phrases interactively.
  history,     hist, h  📜 Manage the histories of the "generate" command.
  favorite,    fav,  f  ⭐ Favorite the histories of the "generate" command.
  unfavorite,  unf,  u  ❌ Unfavorite the favorited histories of the "generate" command.
  completion   comp, c  🔧 Generate the autocompletion script for the specified shell.
  version      ver,  v  🔖 Show the version of jrp.
  help                  🤝 Help for jrp.

Flags:
  -n, --number       🔢 number of phrases to generate (default 1, e.g. : 10)
  -p, --prefix       🔡 prefix of phrases to generate
  -s, --suffix       🔡 suffix of phrases to generate
  -d, --dry-run      🧪 generate phrases without saving as the histories
  -f, --format       📝 format of the output (default "table", e.g. : "plain")
  -i, --interactive  💬 generate Japanese random phrases interactively
  -t, --timeout      ⌛ timeout in seconds for the interactive mode (default 30, e.g. : 10)
  -h, --help         🤝 help for jrp
  -v, --version      🔖 version for jrp

Argument:
  number  🔢 number of phrases to generate (e.g. : 10)

Use "jrp [command] --help" for more information about a command.
`
)
