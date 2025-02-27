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

// FavoriteOptions provides the options for the favorite command.
type FavoriteOptions struct {
	// All is a flag to add all phrases.
	All bool
	// NoConfirm is a flag to not confirm before removing all the histories.
	NoConfirm bool
}

var (
	// favoriteOps is a variable to store the favorite options with the default values for injecting the dependencies in testing.
	favoriteOps = FavoriteOptions{
		All:       false,
		NoConfirm: false,
	}
)

// NewFavoriteCommand returns a new instance of the favorite command.
func NewFavoriteCommand(
	cobra proxy.Cobra,
	output *string,
) proxy.Command {
	cmd := cobra.NewCommand()
	cmd.SetUse("favorite")
	cmd.SetAliases([]string{"fav", "f"})
	cmd.SetUsageTemplate(favoriteUsageTemplate)
	cmd.SetHelpTemplate(favoriteHelpTemplate)
	cmd.SetSilenceErrors(true)
	cmd.Flags().BoolVarP(
		&favoriteOps.All,
		"all",
		"a",
		false,
		"⭐ favorite all histories",
	)
	cmd.Flags().BoolVarP(
		&favoriteOps.NoConfirm,
		"no-confirm",
		"",
		false,
		"🚫 do not confirm before favoriting all the histories",
	)

	cmd.SetRunE(
		func(cmd *c.Command, args []string) error {
			return runFavorite(
				cmd,
				args,
				output,
			)
		},
	)

	return cmd
}

// runFavorite runs the favorite command.
func runFavorite(
	cmd *c.Command,
	args []string,
	output *string,
) error {
	if len(args) == 0 && !favoriteOps.All {
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
	fuc := jrpApp.NewFavoriteUseCase(historyRepo)

	if favoriteOps.All && !favoriteOps.NoConfirm {
		if answer, err := presenter.RunPrompt(
			"Proceed with favoriting all the histories? [y/N]",
		); err != nil {
			return err
		} else if answer != "y" && answer != "Y" {
			o := formatter.Yellow("🚫 Cancelled favoriting all the histories.")
			*output = o
			return nil
		}
	}

	if err := fuc.Run(
		cmd.Context(),
		ids,
		favoriteOps.All,
	); err != nil && err.Error() == "no histories to favorite" {
		o := formatter.Yellow("⚡ No histories to favorite...")
		*output = o
		return nil
	} else if err != nil {
		return err
	}

	o := formatter.Green("✅ Favorited successfully!")
	*output = o

	return nil
}

const (
	// favoriteHelpTemplate is the help template of the favorite command.
	favoriteHelpTemplate = `⭐ Favorite the histories of the "generate" command.

You can specify the histories to favorite with ID arguments.
You have to get ID from the "history" command.
Multiple ID's can be specified separated by spaces.

This command can make the histories easier to find.
And you will not be able to remove the histories with executing "history remove" and "history clear".

Also, you can favorite all the histories with the "-a" or "--all" flag.

` + favoriteUsageTemplate
	// favoriteUsageTemplate is the usage template of the favorite command.
	favoriteUsageTemplate = `Usage:
  jrp favorite [flag] [arguments]
  jrp fav      [flag] [arguments]
  jrp f        [flag] [arguments]

Flags:
  -a, --all    ⭐ favorite all the histories
  -no-confirm  🚫 do not confirm before favoriting all the histories
  -h, --help   🤝 help for favorite

Arguments:
  ID  🆔 favorite with the ID of the history (e.g. : 1 2 3)
`
)
