package history

import (
	"strconv"

	c "github.com/spf13/cobra"

	jrpApp "github.com/yanosea/jrp/app/application/jrp"
	"github.com/yanosea/jrp/app/infrastructure/jrp/repository"
	"github.com/yanosea/jrp/app/presentation/cli/jrp/formatter"
	"github.com/yanosea/jrp/app/presentation/cli/jrp/presenter"

	"github.com/yanosea/jrp/pkg/proxy"
)

// RemoveOptions provides the options for the remove command.
type RemoveOptions struct {
	// All is a flag to remove all phrases.
	All bool
	// Force is a flag to remove the histories even if it is favorited.
	Force bool
	// NoConfirm is a flag to not confirm before removing all the histories.
	NoConfirm bool
}

var (
	// removeOps is a variable to store the remove options with the default values for injecting the dependencies in testing.
	removeOps = RemoveOptions{
		All:       false,
		Force:     false,
		NoConfirm: false,
	}
)

// NewRemoveCommand returns a new instance of the remove command.
func NewRemoveCommand(
	cobra proxy.Cobra,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("remove")
	cmd.SetAliases([]string{"rm", "r"})
	cmd.SetUsageTemplate(removeUsageTemplate)
	cmd.SetHelpTemplate(removeHelpTemplate)
	cmd.SetSilenceErrors(true)
	cmd.Flags().BoolVarP(
		&removeOps.All,
		"all",
		"a",
		false,
		"✨ remove all history",
	)
	cmd.Flags().BoolVarP(
		&removeOps.Force,
		"force",
		"f",
		false,
		"💪 remove the histories even if it is favorited",
	)
	cmd.Flags().BoolVarP(
		&removeOps.NoConfirm,
		"no-confirm",
		"",
		false,
		"🚫 do not confirm before removing all the histories",
	)

	cmd.SetRunE(
		func(cmd *c.Command, args []string) error {
			return runRemove(
				cmd,
				args,
				output,
			)
		},
	)

	return cmd
}

// runRemove runs the remove command.
func runRemove(
	cmd *c.Command,
	args []string,
	output *string,
) error {
	if len(args) == 0 && !removeOps.All {
		o := formatter.Yellow("⚡ No ID arguments specified...")
		*output = o
		return nil
	}

	var ids []int
	for _, arg := range args {
		id, err := strconv.Atoi(arg)
		if err != nil {
			o := formatter.Red("🚨 The ID argument must be an integer...")
			*output = o
			return err
		}
		ids = append(ids, id)
	}

	historyRepo := repository.NewHistoryRepository()
	rhuc := jrpApp.NewRemoveHistoryUseCase(historyRepo)

	if removeOps.All && !removeOps.NoConfirm {
		if answer, err := presenter.RunPrompt(
			"Proceed with removing all the histories? [y/N]",
		); err != nil {
			return err
		} else if answer != "y" && answer != "Y" {
			o := formatter.Yellow("🚫 Cancelled removing all the histories.")
			*output = o
			return nil
		}
	}

	if err := rhuc.Run(
		cmd.Context(),
		ids,
		removeOps.All,
		removeOps.Force,
	); err != nil && err.Error() == "no histories to remove" {
		o := formatter.Yellow("⚡ No histories to remove...")
		*output = o
		return nil
	} else if err != nil {
		return err
	}

	o := formatter.Green("✅ Removed successfully!")
	*output = o

	return nil
}

const (
	// removeHelpTemplate is the help template of the remove command.
	removeHelpTemplate = `📜🧹 Remove the histories of the "generate" command.

You can specify the histories to remove with ID arguments.
You have to get ID from the "history" command.
Multiple ID's can be specified separated by spaces.

You can remove all the histories by flag "-a" or "--all".
This is the same as the "history clear" command.

Also, you can remove the histories even if it is favorited by using the "-f" or ""--force" flag.

` + removeUsageTemplate
	// removeUsageTemplate is the usage template of the remove command.
	removeUsageTemplate = `Usage:
  jrp history remove [flag]
  jrp history rm     [flag]
  jrp history r      [flag]

Flags:
  -a, --all    ✨ remove all histories
  -f, --force  💪 remove the histories even if it is favorited
  -no-confirm  🚫 do not confirm before removing all the histories
  -h, --help   🤝 help for remove

Arguments:
  ID  🆔 remove the history by the ID (e.g: 1 2 3)
`
)
