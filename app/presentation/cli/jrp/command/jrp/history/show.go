package history

import (
	"strconv"

	c "github.com/spf13/cobra"

	jrpApp "github.com/yanosea/jrp/app/application/jrp"
	"github.com/yanosea/jrp/app/infrastructure/jrp/repository"
	"github.com/yanosea/jrp/app/presentation/cli/jrp/formatter"

	"github.com/yanosea/jrp/pkg/proxy"
)

// ShowOptions provides the options for the show command.
type ShowOptions struct {
	// Number is a flag to specify the number of histories to show.
	Number int
	// All is a flag to show all the histories.
	All bool
	// Favorited is a flag to show only favorited histories.
	Favorited bool
	// Format is a flag to specify the format of the output.
	Format string
}

var (
	// showOps is a variable to store the show options with the default values for injecting the dependencies in testing.
	showOps = ShowOptions{
		Number:    1,
		All:       false,
		Favorited: false,
		Format:    "table",
	}
)

// NewShowCommand returns a new instance of the show command.
func NewShowCommand(
	cobra proxy.Cobra,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("show")
	cmd.SetAliases([]string{"sh", "s"})
	cmd.SetUsageTemplate(showUsageTemplate)
	cmd.SetHelpTemplate(showHelpTemplate)
	cmd.SetArgs(cobra.MaximumNArgs(1))
	cmd.SetSilenceErrors(true)
	cmd.Flags().IntVarP(
		&showOps.Number,
		"number",
		"n",
		10,
		"🔢 number how many histories to show (default 10, e.g. : 50)",
	)
	cmd.Flags().BoolVarP(
		&showOps.All,
		"all",
		"a",
		false,
		"📁 show all the history",
	)
	cmd.Flags().BoolVarP(
		&showOps.Favorited,
		"favorited",
		"F",
		false,
		"🌟 show only favorited histories",
	)
	cmd.Flags().StringVarP(
		&showOps.Format,
		"format",
		"f",
		"table",
		"📝 format of the output (default \"table\", e.g. : \"plain\")",
	)

	cmd.SetRunE(
		func(cmd *c.Command, args []string) error {
			return runShow(
				cmd,
				args,
				output,
			)
		},
	)

	return cmd
}

// runShow runs the show command.
func runShow(
	cmd *c.Command,
	args []string,
	output *string,
) error {
	var number int = showOps.Number
	isDefaultNumber := number == 10
	if len(args) > 0 {
		argNumber, err := strconv.Atoi(args[0])
		if err != nil {
			o := formatter.Red("🚨 The number argument must be an integer...")
			*output = o
			return err
		}

		if isDefaultNumber {
			number = argNumber
		} else {
			if argNumber > number {
				number = argNumber
			}
		}
	}

	historyRepo := repository.NewHistoryRepository()
	ghuc := jrpApp.NewGetHistoryUseCase(historyRepo)

	ghoDtos, err := ghuc.Run(
		cmd.Context(),
		showOps.All,
		showOps.Favorited,
		number,
	)
	if err != nil {
		return err
	}

	if len(ghoDtos) == 0 {
		o := formatter.Yellow("⚡ No histories found...")
		*output = o
		return nil
	}

	f, err := formatter.NewFormatter(showOps.Format)
	if err != nil {
		o := formatter.Red("❌ Failed to create a formatter...")
		*output = o
		return err
	}
	o := f.Format(ghoDtos)
	*output = o

	return nil
}

const (
	// showHelpTemplate is the help template of the show command.
	showHelpTemplate = `📜📖 Show the histories of the "generate" command.

You can specify how many histories to show by flag "-n" or "--number" or a number argument.
jrp will get the most recent histories from the histories.
If you don't specify the number of histories, jrp will show the most recent 10 histories by default.
If both are provided, the larger number takes precedence.

Also, you can show all the histories by flag "-a" or "--all".
If you use the flag, the number flag or argument will be ignored.

` + showUsageTemplate
	// showUsageTemplate is the usage template of the show command.
	showUsageTemplate = `Usage:
  jrp history show [flag]
  jrp history sh   [flag]
  jrp history s    [flag]

Flags:
  -n, --number     🔢 number how many histories to show (default 10, e.g. : 50)
  -a, --all        📁 show all the histories
  -F, --favorited  🌟 show only favorited histories
  -f, --format     📝 format of the output (default "table", e.g. : "plain")
  -h, --help       🤝 help for show

Argument:
  number  🔢 number how many histories to show (default 10, e.g. : 50)
`
)
