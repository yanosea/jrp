package generate

import (
	"strconv"

	c "github.com/spf13/cobra"

	jrpApp "github.com/yanosea/jrp/app/application/jrp"
	wnjpnApp "github.com/yanosea/jrp/app/application/wnjpn"
	"github.com/yanosea/jrp/app/infrastructure/database"
	"github.com/yanosea/jrp/app/infrastructure/jrp/repository"
	"github.com/yanosea/jrp/app/infrastructure/wnjpn/query_service"
	"github.com/yanosea/jrp/app/presentation/cli/jrp/formatter"

	"github.com/yanosea/jrp/pkg/proxy"
)

// GenerateOptions provides the options for the generate command.
type GenerateOptions struct {
	// Number is a flag to specify the number of phrases to generate.
	Number int
	// Prefix is a flag to specify the prefix of the phrases to generate.
	Prefix string
	// Suffix is a flag to specify the suffix of the phrases to generate.
	Suffix string
	// DryRun is a flag to generate phrases without saving to the history.
	DryRun bool
	// Format is a flag to specify the format of the output.
	Format string
	// Interactive is a flag to generate Japanese random phrases interactively.
	Interactive bool
	// Timeout is a flag to specify the timeout in seconds for the interactive mode.
	Timeout int
}

var (
	// GenerateOps is a variable to store the generate options with the default values for injecting the dependencies in testing.
	GenerateOps = GenerateOptions{
		Number:      1,
		Prefix:      "",
		Suffix:      "",
		DryRun:      false,
		Format:      "table",
		Interactive: false,
		Timeout:     30,
	}
)

// NewGenerateCommand returns a new instance of the generate command.
func NewGenerateCommand(
	cobra proxy.Cobra,
	interactiveCmd proxy.Command,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("generate")
	cmd.SetAliases([]string{"gen", "g"})
	cmd.SetUsageTemplate(generateUsageTemplate)
	cmd.SetHelpTemplate(generateHelpTemplate)
	cmd.SetArgs(cobra.MaximumNArgs(1))
	cmd.SetSilenceErrors(true)
	cmd.Flags().IntVarP(
		&GenerateOps.Number,
		"number",
		"n",
		1,
		"🔢 number of phrases to generate (default 1, e.g. : 10)",
	)
	cmd.Flags().StringVarP(
		&GenerateOps.Prefix,
		"prefix",
		"p",
		"",
		"🔡 prefix of phrases to generate",
	)
	cmd.Flags().StringVarP(
		&GenerateOps.Suffix,
		"suffix",
		"s",
		"",
		"🔡 suffix of phrases to generate",
	)
	cmd.Flags().BoolVarP(
		&GenerateOps.DryRun,
		"dry-run",
		"d",
		false,
		"🧪 generate phrases without saving to the history",
	)
	cmd.Flags().StringVarP(
		&GenerateOps.Format,
		"format",
		"f",
		"table",
		"📝 format of the output (default \"table\", e.g. : \"plain\")",
	)
	cmd.Flags().BoolVarP(
		&GenerateOps.Interactive,
		"interactive",
		"i",
		false,
		"💬 generate Japanese random phrases interactively",
	)
	cmd.Flags().IntVarP(
		&GenerateOps.Timeout,
		"timeout",
		"t",
		30,
		"⌛ timeout in seconds for the interactive mode (default 30, e.g. : 10)",
	)
	cmd.AddCommand(interactiveCmd)
	cmd.SetRunE(
		func(cmd *c.Command, args []string) error {
			return runGenerate(
				cmd,
				args,
				interactiveCmd,
				output,
			)
		},
	)

	return cmd
}

// runGenerate runs the generate command.
func runGenerate(
	cmd *c.Command,
	args []string,
	interactiveCmd proxy.Command,
	output *string,
) error {
	if GenerateOps.Interactive {
		interactiveOps.Prefix = GenerateOps.Prefix
		interactiveOps.Suffix = GenerateOps.Suffix
		interactiveOps.Format = GenerateOps.Format
		interactiveOps.Timeout = GenerateOps.Timeout
		return interactiveCmd.RunE(cmd, args)
	}

	connManager := database.GetConnectionManager()
	if connManager == nil {
		o := formatter.Red("❌ Connection manager is not initialized...")
		*output = o
		return nil
	}

	_, err := connManager.GetConnection(database.WNJpnDB)
	if err != nil && err.Error() == "connection not initialized" {
		o := formatter.Yellow("⚡ You have to execute \"download\" to use jrp...")
		*output = o
		return nil
	} else if err != nil {
		return err
	}

	needRandomPrefix := GenerateOps.Prefix == ""
	needRandomSuffix := GenerateOps.Suffix == ""
	if !needRandomPrefix && !needRandomSuffix {
		o := formatter.Yellow("⚡ You can't specify both prefix and suffix at the same time...")
		*output = o
		return nil
	}

	var pos []string
	if needRandomPrefix {
		pos = append(pos, "a", "v")
	}
	if needRandomSuffix {
		pos = append(pos, "n")
	}

	wordQueryService := query_service.NewWordQueryService()
	fwuc := wnjpnApp.NewFetchWordsUseCase(wordQueryService)

	fwoDtos, err := fwuc.Run(
		cmd.Context(),
		"jpn",
		pos,
	)
	if err != nil {
		return err
	}

	var number int = GenerateOps.Number
	if len(args) > 0 {
		argNumber, err := strconv.Atoi(args[0])
		if err != nil {
			o := formatter.Red("🚨 The number argument must be an integer...")
			*output = o
			return err
		}
		if argNumber > number {
			number = argNumber
		}
	}

	var gjiDtos []*jrpApp.GenerateJrpUseCaseInputDto
	for _, fwoDto := range fwoDtos {
		gjiDto := &jrpApp.GenerateJrpUseCaseInputDto{
			WordID: fwoDto.WordID,
			Lang:   fwoDto.Lang,
			Lemma:  fwoDto.Lemma,
			Pron:   fwoDto.Pron,
			Pos:    fwoDto.Pos,
		}
		gjiDtos = append(gjiDtos, gjiDto)
	}

	gjuc := jrpApp.NewGenerateJrpUseCase()
	var gjoDtos []*jrpApp.GenerateJrpUseCaseOutputDto
	for i := 0; i < number; i++ {
		var gjoDto *jrpApp.GenerateJrpUseCaseOutputDto
		if needRandomPrefix && needRandomSuffix {
			gjoDto = gjuc.RunWithRandom(gjiDtos)
		} else if needRandomPrefix {
			gjoDto = gjuc.RunWithSuffix(gjiDtos, GenerateOps.Suffix)
		} else {
			gjoDto = gjuc.RunWithPrefix(gjiDtos, GenerateOps.Prefix)
		}
		gjoDtos = append(gjoDtos, gjoDto)
	}

	if !GenerateOps.DryRun {
		var shiDtos []*jrpApp.SaveHistoryUseCaseInputDto
		for _, gjoDto := range gjoDtos {
			shiDto := &jrpApp.SaveHistoryUseCaseInputDto{
				Phrase:      gjoDto.Phrase,
				Prefix:      gjoDto.Prefix,
				Suffix:      gjoDto.Suffix,
				IsFavorited: gjoDto.IsFavorited,
				CreatedAt:   gjoDto.CreatedAt,
				UpdatedAt:   gjoDto.UpdatedAt,
			}
			shiDtos = append(shiDtos, shiDto)
		}

		historyRepo := repository.NewHistoryRepository()
		shuc := jrpApp.NewSaveHistoryUseCase(historyRepo)

		shoDtos, err := shuc.Run(cmd.Context(), shiDtos)
		if err != nil {
			return err
		}

		if len(shoDtos) > 0 {
			for i, gjoDto := range gjoDtos {
				gjoDto.ID = shoDtos[i].ID
			}
		}
	}

	f, err := formatter.NewFormatter(GenerateOps.Format)
	if err != nil {
		o := formatter.Red("❌ Failed to create a formatter...")
		*output = o
		return err
	}
	o := f.Format(gjoDtos)
	*output = o

	return nil
}

const (
	// generateHelpTemplate is the help template of the generate command.
	generateHelpTemplate = `✨ Generate Japanese random phrases.

You can specify how many phrases to generate by flag "-n" or "--number" or a number argument.
If both are provided, the larger number takes precedence.

And you can specify the prefix or suffix of the phrases to generate
by the flag "-p" or "--prefix" and "-s" or "--suffix".

Those commands below are the same.
  "jrp" : "jrp generate"
  "jrp interactive" : "jrp --interactive" : "jrp generate interactive" : "jrp generate --interactive"

` + generateUsageTemplate
	// generateUsageTemplate is the usage template of the generate command.
	generateUsageTemplate = `Usage:
  jrp generate [flags]
  jrp gen      [flags]
  jrp g        [flags]

Available Subcommands:
  interactive, int, i  💬 Generate Japanese random phrases interactively.

Flags:
  -n, --number       🔢 number of phrases to generate (default 1, e.g. : 10)
  -p, --prefix       🔡 prefix of phrases to generate
  -s, --suffix       🔡 suffix of phrases to generate
  -d, --dry-run      🧪 generate phrases without saving to the history
  -f, --format       📝 format of the output (default "table", e.g. : "plain")
  -i, --interactive  💬 generate Japanese random phrases interactively
  -t, --timeout      ⌛ timeout in seconds for the interactive mode (default 30, e.g. : 10)
  -h, --help         🤝 help for generate

Argument:
  number  🔢 number of phrases to generate (default 1, e.g. : 10)

Use "jrp generate [command] --help" for more information about a command.
`
)
