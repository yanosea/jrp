package history

import (
	c "github.com/spf13/cobra"

	jrpApp "github.com/yanosea/jrp/app/application/jrp"
	"github.com/yanosea/jrp/app/infrastructure/jrp/repository"
	"github.com/yanosea/jrp/app/presentation/cli/jrp/formatter"

	"github.com/yanosea/jrp/pkg/proxy"
)

// SearchOptions provides the options for the search command.
type SearchOptions struct {
	// Number is a flag to specify the number of histories to search.
	Number int
	// And is a flag to search histories by AND condition.
	And bool
	// All is a flag to search all histories.
	All bool
	// Favorited is a flag to show only favorited histories.
	Favorited bool
	// Format is a flag to specify the format of the output.
	Format string
}

var (
	// searchOps is a variable to store the search options with the default values for injecting the dependencies in testing.
	searchOps = SearchOptions{
		Number:    1,
		And:       false,
		All:       false,
		Favorited: false,
		Format:    "table",
	}
)

// NewSearchCommand returns a new instance of the search command.
func NewSearchCommand(
	cobra proxy.Cobra,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("search")
	cmd.SetAliases([]string{"se", "S"})
	cmd.SetUsageTemplate(searchUsageTemplate)
	cmd.SetHelpTemplate(searchHelpTemplate)
	cmd.SetSilenceErrors(true)
	cmd.Flags().IntVarP(
		&searchOps.Number,
		"number",
		"n",
		10,
		"🔢 number how many histories to search (default 10, e.g: 50)",
	)
	cmd.Flags().BoolVarP(
		&searchOps.And,
		"and",
		"A",
		false,
		"🧠 search histories by AND condition",
	)
	cmd.Flags().BoolVarP(
		&searchOps.All,
		"all",
		"a",
		false,
		"📁 search all histories",
	)
	cmd.Flags().BoolVarP(
		&searchOps.Favorited,
		"favorited",
		"F",
		false,
		"🌟 show only favorited histories",
	)
	cmd.Flags().StringVarP(
		&searchOps.Format,
		"format",
		"f",
		"table",
		"📝 format of the output (default \"table\", e.g: \"plain\")",
	)

	cmd.SetRunE(
		func(cmd *c.Command, args []string) error {
			return runSearch(
				cmd,
				args,
				output,
			)
		},
	)

	return cmd
}

// runSearch runs the search command.
func runSearch(
	cmd *c.Command,
	args []string,
	output *string,
) error {
	if len(args) == 0 {
		o := formatter.Yellow("⚡ No keywords provided...")
		*output = o
		return nil
	}

	historyRepo := repository.NewHistoryRepository()
	shuc := jrpApp.NewSearchHistoryUseCase(historyRepo)

	shoDtos, err := shuc.Run(
		cmd.Context(),
		args,
		searchOps.And,
		searchOps.All,
		searchOps.Favorited,
		searchOps.Number,
	)
	if err != nil {
		return err
	}

	if len(shoDtos) == 0 {
		o := formatter.Yellow("⚡ No histories found...")
		*output = o
		return nil
	}

	f, err := formatter.NewFormatter(searchOps.Format)
	if err != nil {
		o := formatter.Red("❌ Failed to create a formatter...")
		*output = o
		return err
	}
	o := f.Format(shoDtos)
	*output = o

	return nil
}

const (
	// searchHelpTemplate is the help template of the search command.
	searchHelpTemplate = `📜🔍 Search the histories of the "generate" command.

You can search histories with keyword arguments.
Multiple keywords are separated by a space.

If you want to search histories by AND condition, you can use flag "-A" or "--and".
OR condition is by default.

You can specify how many histories to show with flag "-n" or "--number".
If you don't specify the number of histories, jrp will show the most recent 10 histories by default.

Also, you can show all histories by flag "-a" or "--all".
If you use the flag, the number flag will be ignored.

` + searchUsageTemplate
	// searchUsageTemplate is the usage template of the search command.
	searchUsageTemplate = `Usage:
  jrp history search [flag]
  jrp history se     [flag]
  jrp history S      [flag]

Flags:
  -A, --and        🧠 search histories by AND condition
  -n, --number     🔢 number how many histories to show (default 10, e.g: 50)
  -a, --all        📁 show all histories
  -F, --favorited  🌟 show only favorited histories
  -f, --format     📝 format of the output (default "table", e.g: "plain")
  -h, --help       🤝 help for search

Arguments:
  keywords  🔡 search histories by keywords (multiple keywords are separated by space)
`
)
