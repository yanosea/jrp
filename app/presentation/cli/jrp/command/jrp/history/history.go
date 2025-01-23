package history

import (
	c "github.com/spf13/cobra"

	"github.com/yanosea/jrp/v2/pkg/proxy"
)

// HistoryOptions provides the options for the generate command.
type HistoryOptions struct {
	ShowOptions ShowOptions
}

var (
	// historyOps is a variable to store the history options with the default values for injecting the dependencies in testing.
	historyOps = HistoryOptions{
		ShowOptions: ShowOptions{
			Number:    1,
			All:       false,
			Favorited: false,
			Format:    "table",
		},
	}
)

// NewHistoryCommand returns a new instance of the history command.
func NewHistoryCommand(
	cobra proxy.Cobra,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("history")
	cmd.SetAliases([]string{"hist", "h"})
	cmd.SetUsageTemplate(historyUsageTemplate)
	cmd.SetHelpTemplate(historyHelpTemplate)
	cmd.SetArgs(cobra.MaximumNArgs(1))
	cmd.SetSilenceErrors(true)
	cmd.Flags().IntVarP(
		&historyOps.ShowOptions.Number,
		"number",
		"n",
		10,
		"🔢 number how many histories to show (default 10, e.g. : 50)",
	)
	cmd.Flags().BoolVarP(
		&historyOps.ShowOptions.All,
		"all",
		"a",
		false,
		"📁 show all the history",
	)
	cmd.Flags().BoolVarP(
		&historyOps.ShowOptions.Favorited,
		"favorited",
		"F",
		false,
		"🌟 show only favorited histories",
	)
	cmd.Flags().StringVarP(
		&historyOps.ShowOptions.Format,
		"format",
		"f",
		"table",
		"📝 format of the output (default \"table\", e.g. : \"plain\")",
	)

	showCmd := NewShowCommand(
		cobra,
		output,
	)
	cmd.AddCommand(
		NewClearCommand(
			cobra,
			output,
		),
		NewRemoveCommand(
			cobra,
			output,
		),
		NewSearchCommand(
			cobra,
			output,
		),
		showCmd,
	)

	cmd.SetRunE(
		func(cmd *c.Command, args []string) error {
			return runHistory(
				cmd,
				showCmd,
				args,
			)
		},
	)

	return cmd
}

// runHistory runs the history command.
func runHistory(
	cmd *c.Command,
	showCmd proxy.Command,
	args []string,
) error {
	showOps = historyOps.ShowOptions
	return showCmd.RunE(cmd, args)
}

const (
	// historyHelpTemplate is the help template of the history command.
	historyHelpTemplate = `📜 Manage the histories of the "generate" command.

You can show, search, remove and clear the histories of the "generate" command.

You can specify how many histories to show by flag "-n" or "--number" or a number argument.
jrp will get the most recent histories from the histories.
If you don't specify the number of histories, jrp will show the most recent 10 histories by default.
If both are provided, the larger number takes precedence.

Also, you can show all the histories the history by flag "-a" or "--all".
If you use the flag, the number flag or argument will be ignored.

` + historyUsageTemplate
	// historyUsageTemplate is the usage template of the history command.
	historyUsageTemplate = `Usage:
  jrp history [flag]
  jrp hist    [flag]
  jrp h       [flag]
  jrp history [command]
  jrp hist    [command]
  jrp h       [command]

Available Subommands:
  show,   sh, s  📜📖 Show the histories of the "generate" command.
                      You can abbreviate "show" sub command. ("jrp history" and "jrp history show" are the same.)
  search, se, S  📜🔍 Search the histories of the "generate" command.
  remove, rm, r  📜🧹 Remove the histories of the "generate" command.
  clear,  cl, c  📜✨ Clear the histories of the "generate" command.

Flags:
  -n, --number     🔢 number how many histories to show (default 10, e.g. : 50)
  -a, --all        📁 show all the histories
  -F, --favorited  🌟 show only favorited histories
  -f, --format     📝 format of the output (default "table", e.g. : "plain")
  -h, --help       🤝 help for history

Argument:
  number  🔢 number how many histories to show (default 10, e.g. : 50)

Use "jrp history [command] --help" for more information about a command.
`
)
