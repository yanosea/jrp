package history

import (
	c "github.com/spf13/cobra"

	jrpApp "github.com/yanosea/jrp/v2/app/application/jrp"
	"github.com/yanosea/jrp/v2/app/infrastructure/jrp/repository"
	"github.com/yanosea/jrp/v2/app/presentation/cli/jrp/formatter"
	"github.com/yanosea/jrp/v2/app/presentation/cli/jrp/presenter"

	"github.com/yanosea/jrp/v2/pkg/proxy"
)

// ClearOptions provides the options for the clear command.
type ClearOptions struct {
	// Force is a flag to clear the histories even if it is favorited.
	Force bool
	// NoConfirm is a flag to not confirm before removing all the histories.
	NoConfirm bool
}

var (
	// clearOps is a variable to store the clear options with the default values for injecting the dependencies in testing.
	clearOps = ClearOptions{
		Force:     false,
		NoConfirm: false,
	}
)

// NewClearCommand returns a new instance of the clear command.
func NewClearCommand(
	cobra proxy.Cobra,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("clear")
	cmd.SetAliases([]string{"cl", "c"})
	cmd.SetUsageTemplate(clearUsageTemplate)
	cmd.SetHelpTemplate(clearHelpTemplate)
	cmd.SetArgs(cobra.ExactArgs(0))
	cmd.SetSilenceErrors(true)
	cmd.Flags().BoolVarP(
		&clearOps.Force,
		"force",
		"f",
		false,
		"💪 clear the histories even if it is favorited",
	)
	cmd.Flags().BoolVarP(
		&clearOps.NoConfirm,
		"no-confirm",
		"",
		false,
		"🚫 do not confirm before clearing the histories",
	)

	cmd.SetRunE(
		func(cmd *c.Command, _ []string) error {
			return runClear(
				cmd,
				output,
			)
		},
	)

	return cmd
}

// runClear runs the clear command.
func runClear(
	cmd *c.Command,
	output *string,
) error {
	var ids []int

	historyRepo := repository.NewHistoryRepository()
	rhuc := jrpApp.NewRemoveHistoryUseCase(historyRepo)

	if !clearOps.NoConfirm {
		if answer, err := presenter.RunPrompt(
			"Proceed with clearing the histories? [y/N]",
		); err != nil {
			return err
		} else if answer != "y" && answer != "Y" {
			o := formatter.Yellow("🚫 Cancelled clearing the histories.")
			*output = o
			return nil
		}
	}

	if err := rhuc.Run(
		cmd.Context(),
		ids,
		true,
		clearOps.Force,
	); err != nil && err.Error() == "no histories to remove" {
		o := formatter.Yellow("⚡ No histories to clear...")
		*output = o
		return nil
	} else if err != nil {
		return err
	}

	o := formatter.Green("✅ Cleared successfully!")
	*output = o

	return nil
}

const (
	// clearHelpTemplate is the help template of the clear command.
	clearHelpTemplate = `📜✨ Clear the histories of the "generate" command.

You can clear the histories.
This is the same as the "history remove -a" command.
Also, you can clear the histories even if it is favorited by using the "-f" or ""--force" flag.

` + clearUsageTemplate
	// clearUsageTemplate is the usage template of the clear command.
	clearUsageTemplate = `Usage:
  jrp history clear [flag]
  jrp history cl    [flag]
  jrp history c     [flag]

Flags:
  -f, --force  💪 clear the histories even if it is favorited
  -no-confirm  🚫 do not confirm before clearing the histories
  -h, --help   🤝 help for clear
`
)
