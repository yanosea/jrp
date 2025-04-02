package generate

import (
	"os"
	"strconv"

	c "github.com/spf13/cobra"

	jrpApp "github.com/yanosea/jrp/v2/app/application/jrp"
	wnjpnApp "github.com/yanosea/jrp/v2/app/application/wnjpn"
	"github.com/yanosea/jrp/v2/app/infrastructure/database"
	"github.com/yanosea/jrp/v2/app/infrastructure/jrp/repository"
	"github.com/yanosea/jrp/v2/app/infrastructure/wnjpn/query_service"
	"github.com/yanosea/jrp/v2/app/presentation/cli/jrp/formatter"
	"github.com/yanosea/jrp/v2/app/presentation/cli/jrp/presenter"

	"github.com/yanosea/jrp/v2/pkg/proxy"
)

// InteractiveOptions provides the options for the interactive command.
type InteractiveOptions struct {
	// Prefix is a flag to specify the prefix of the phrases to generate.
	Prefix string
	// Suffix is a flag to specify the suffix of the phrases to generate.
	Suffix string
	// Format is a flag to specify the format of the output.
	Format string
	// Timeout is a flag to specify the timeout in seconds for the interactive mode.
	Timeout int
}

var (
	// interactiveOps is a variable to store the interactive options with the default values for injecting the dependencies in testing.
	interactiveOps = InteractiveOptions{
		Prefix:  "",
		Suffix:  "",
		Format:  "table",
		Timeout: 30,
	}
)

// NewInteractiveCommand returns a new instance of the interactive command.
func NewInteractiveCommand(
	cobra proxy.Cobra,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("interactive")
	cmd.SetAliases([]string{"int", "i"})
	cmd.SetUsageTemplate(interactiveUsageTemplate)
	cmd.SetHelpTemplate(interactiveHelpTemplate)
	cmd.SetArgs(cobra.MaximumNArgs(1))
	cmd.SetSilenceErrors(true)
	cmd.PersistentFlags().StringVarP(
		&interactiveOps.Prefix,
		"prefix",
		"p",
		"",
		"🔡 prefix of phrases to generate",
	)
	cmd.PersistentFlags().StringVarP(
		&interactiveOps.Suffix,
		"suffix",
		"s",
		"",
		"🔡 suffix of phrases to generate",
	)
	cmd.PersistentFlags().StringVarP(
		&interactiveOps.Format,
		"format",
		"f",
		"table",
		"📝 format of the output (default \"table\", e.g: \"plain\")",
	)
	cmd.PersistentFlags().IntVarP(
		&interactiveOps.Timeout,
		"timeout",
		"t",
		30,
		"⌛ timeout in seconds for the interactive mode (default 30, e.g: 10)",
	)

	cmd.SetRunE(
		func(cmd *c.Command, _ []string) error {
			return runInteractive(
				cmd,
				output,
			)
		},
	)

	return cmd
}

// runInteractive runs the interactive command.
func runInteractive(
	cmd *c.Command,
	output *string,
) error {
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

	needRandomPrefix := interactiveOps.Prefix == ""
	needRandomSuffix := interactiveOps.Suffix == ""
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

	phase := 1
	for {
		if err := presenter.Print(os.Stdout, formatter.Blue("🔄 Phase : "+strconv.Itoa(phase))); err != nil {
			return err
		}

		var gjiDtos []*jrpApp.GenerateJrpUseCaseInputDto
		for _, fwDto := range fwoDtos {
			gjiDto := &jrpApp.GenerateJrpUseCaseInputDto{
				WordID: fwDto.WordID,
				Lang:   fwDto.Lang,
				Lemma:  fwDto.Lemma,
				Pron:   fwDto.Pron,
				Pos:    fwDto.Pos,
			}
			gjiDtos = append(gjiDtos, gjiDto)
		}

		gjuc := jrpApp.NewGenerateJrpUseCase()
		var gjoDtos []*jrpApp.GenerateJrpUseCaseOutputDto
		var gjoDto *jrpApp.GenerateJrpUseCaseOutputDto
		if needRandomPrefix && needRandomSuffix {
			gjoDto = gjuc.RunWithRandom(gjiDtos)
		} else if needRandomPrefix {
			gjoDto = gjuc.RunWithSuffix(gjiDtos, GenerateOps.Suffix)
		} else {
			gjoDto = gjuc.RunWithPrefix(gjiDtos, GenerateOps.Prefix)
		}
		gjoDtos = append(gjoDtos, gjoDto)

		f, err := formatter.NewFormatter(interactiveOps.Format)
		if err != nil {
			o := formatter.Red("❌ Failed to create a formatter...")
			*output = o
			return err
		}
		o := f.Format(gjoDtos)
		if err := presenter.Print(os.Stdout, "\n"); err != nil {
			return err
		}
		if err := presenter.Print(os.Stdout, o); err != nil {
			return err
		}
		if err := presenter.Print(os.Stdout, "\n"); err != nil {
			return err
		}
		if err := presenter.Print(os.Stdout, formatter.Yellow(interactivePromptLabel)); err != nil {
			return err
		}

		if err := presenter.OpenKeyboard(); err != nil {
			return err
		}
		answer, err := presenter.GetKey(interactiveOps.Timeout)
		if err != nil {
			return err
		}
		if err := presenter.CloseKeyboard(); err != nil {
			return err
		}

		var save bool
		var cont bool
		switch answer {
		case "u", "U":
			gjoDtos[0].IsFavorited = 1
			save = true
			cont = true
		case "i", "I":
			gjoDtos[0].IsFavorited = 1
			save = true
			cont = false
		case "j", "J":
			save = true
			cont = true
		case "k", "K":
			save = true
			cont = false
		case "m", "M":
			save = false
			cont = true
		default:
			save = false
			cont = false
		}

		if save {
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

			_, err = shuc.Run(cmd.Context(), shiDtos)
			if err != nil {
				return err
			}

			if answer == "u" || answer == "U" || answer == "i" || answer == "I" {
				if err := presenter.Print(os.Stdout, formatter.Green("✅ Favorited successfully!")); err != nil {
					return err
				}
				if err := presenter.Print(os.Stdout, "\n"); err != nil {
					return err
				}
			} else {
				if err := presenter.Print(os.Stdout, formatter.Green("✅ Saved successfully!")); err != nil {
					return err
				}
				if err := presenter.Print(os.Stdout, "\n"); err != nil {
					return err
				}
			}
		} else {
			if err := presenter.Print(os.Stdout, formatter.Yellow("⏩ Skip!")); err != nil {
				return err
			}
			if err := presenter.Print(os.Stdout, "\n"); err != nil {
				return err
			}
		}

		if !cont {
			if err := presenter.Print(os.Stdout, "🚪 Exit!"); err != nil {
				return err
			}
			break
		}

		phase++
	}

	return nil
}

const (
	// interactiveHelpTemplate is the help template of the interactive command.
	interactiveHelpTemplate = `💬 Generate Japanese random phrases interactively.

You can specify the prefix or suffix of the phrases to generate
by the flag "-p" or "--prefix" and "-s" or "--suffix".

And you can choose to save or favorite the phrases generated interactively.

Press either key below for your action:
  "u"   : Favorite, continue.
  "i"   : Favorite, exit.
  "j"   : Save, continue.
  "k"   : Save, exit.
  "m"   : Skip, continue.
  other : Skip, exit.

` + generateUsageTemplate
	// interactiveUsageTemplate is the usage template of the interactive command.
	interactiveUsageTemplate = `Usage:
  jrp interactive [flags]
  jrp int         [flags]
  jrp i           [flags]

Flags:
  -p, --prefix   🔡 prefix of phrases to generate
  -s, --suffix   🔡 suffix of phrases to generate
  -P, --plain    📝 plain text output instead of table output
  -t, --timeout  ⌛ timeout second for the interactive mode (default 30, e.g: 10)
  -h, --help     🤝 help for interactive
`
	// interactivePromptLabel is the prompt label of the interactive command.
	interactivePromptLabel = `🔽 Press either key below for your action:
  "u"   : Favorite, continue.
  "i"   : Favorite, exit.
  "j"   : Save, continue.
  "k"   : Save, exit.
  "m"   : Skip, continue.
  other : Skip, exit.
`
)
