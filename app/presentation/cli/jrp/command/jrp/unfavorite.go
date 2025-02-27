package jrp

import (
	"strconv"

	c "github.com/spf13/cobra"

	jrpApp "github.com/yanosea/jrp/v2/app/application/jrp"
	"github.com/yanosea/jrp/v2/app/infrastructure/jrp/repository"
	"github.com/yanosea/jrp/v2/app/presentation/cli/jrp/formatter"
	"github.com/yanosea/jrp/v2/app/presentation/cli/jrp/presenter"

	"github.com/yanosea/jrp/v2/pkg/proxy"
)

// UnfavoriteOptions provides the options for the unfavorite command.
type UnfavoriteOptions struct {
	// All is a flag to unfavorite all the histories.
	All bool
	// NoConfirm is a flag to not confirm before unfavoriting all the historyies.
	NoConfirm bool
}

var (
	// unfavoriteOps is a variable to store the unfavorite options with the default values for injecting the dependencies in testing.
	unfavoriteOps = UnfavoriteOptions{
		All:       false,
		NoConfirm: false,
	}
)

// NewUnfavoriteCommand returns a new instance of the unfavorite command.
func NewUnfavoriteCommand(
	cobra proxy.Cobra,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("unfavorite")
	cmd.SetAliases([]string{"unf", "u"})
	cmd.SetUsageTemplate(unfavoriteUsageTemplate)
	cmd.SetHelpTemplate(unfavoriteHelpTemplatep)
	cmd.SetSilenceErrors(true)
	cmd.Flags().BoolVarP(
		&unfavoriteOps.All,
		"all",
		"a",
		false,
		"✨ remove all favorited phrases",
	)
	cmd.Flags().BoolVarP(
		&unfavoriteOps.NoConfirm,
		"no-confirm",
		"",
		false,
		"🚫 do not confirm before removing all the favorited phrases",
	)

	cmd.SetRunE(
		func(cmd *c.Command, args []string) error {
			return runUnfavorite(
				cmd,
				args,
				output,
			)
		},
	)

	return cmd
}

// runUnfavorite runs the unfavorite command.
func runUnfavorite(
	cmd *c.Command,
	args []string,
	output *string,
) error {
	if len(args) == 0 && !unfavoriteOps.All {
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
	uuc := jrpApp.NewUnfavoriteUseCase(historyRepo)

	if unfavoriteOps.All && !unfavoriteOps.NoConfirm {
		if answer, err := presenter.RunPrompt(
			"Proceed with unfavoriting all the histories? [y/N]",
		); err != nil {
			return err
		} else if answer != "y" && answer != "Y" {
			o := formatter.Yellow("🚫 Cancelled unfavoriting all the favorited histories.")
			*output = o
			return nil
		}
	}

	if err := uuc.Run(
		cmd.Context(),
		ids,
		unfavoriteOps.All,
	); err != nil && err.Error() == "no favorited histories to unfavorite" {
		o := formatter.Yellow("⚡ No favorited histories to unfavorite...")
		*output = o
		return nil
	} else if err != nil {
		return err
	}

	o := formatter.Green("✅ Unfavorited successfully!")
	*output = o

	return nil
}

const (
	// unfavoriteHelpTemplatep is the help template of the unfavorite command.
	unfavoriteHelpTemplatep = `⭐🧹 Unfavorite the favorited histories with the "favorite" command.

You can specify the favorited histories to unfavorite with ID arguments.
You have to get ID from the "history" command.
Multiple ID's can be specified separated by spaces.

This does not remove the history of the "generate" command, just unfavorite.

Also, you can unfavorite all the favorited histories with the "-a" or "--all" flag.

` + unfavoriteUsageTemplate
	// unfavoriteUsageTemplate is the usage template of the unfavorite command.
	unfavoriteUsageTemplate = `Usage:
  jrp unfavorite [flag] [arguments]
  jrp unf        [flag] [arguments]
  jrp u          [flag] [arguments]

Flags:
  -a, --all    ✨ unfavorite all the favorited histories
  -no-confirm  🚫 do not confirm before unfavoriting all the favorited histories
  -h, --help   🤝 help for unfavorite

Arguments:
  ID  🆔 unfavorite with the the ID of the favorited history (e.g. : 1 2 3)
`
)
