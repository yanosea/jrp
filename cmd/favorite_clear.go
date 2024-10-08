package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/app/database/jrp/repository"
	"github.com/yanosea/jrp/app/library/dbfiledirpathprovider"
	"github.com/yanosea/jrp/app/library/utility"
	"github.com/yanosea/jrp/app/proxy/cobra"
	"github.com/yanosea/jrp/app/proxy/color"
	"github.com/yanosea/jrp/app/proxy/filepath"
	"github.com/yanosea/jrp/app/proxy/fmt"
	"github.com/yanosea/jrp/app/proxy/io"
	"github.com/yanosea/jrp/app/proxy/os"
	"github.com/yanosea/jrp/app/proxy/promptui"
	"github.com/yanosea/jrp/app/proxy/sort"
	"github.com/yanosea/jrp/app/proxy/sql"
	"github.com/yanosea/jrp/app/proxy/strings"
	"github.com/yanosea/jrp/app/proxy/user"
	"github.com/yanosea/jrp/cmd/constant"
)

// favoriteClearOption is the struct for favorite clear command.
type favoriteClearOption struct {
	Out                   ioproxy.WriterInstanceInterface
	ErrOut                ioproxy.WriterInstanceInterface
	Args                  []string
	NoConfirm             bool
	DBFileDirPathProvider dbfiledirpathprovider.DBFileDirPathProvidable
	JrpRepository         repository.JrpRepositoryInterface
	PromptuiProxy         promptuiproxy.Promptui
	Utility               utility.UtilityInterface
}

// NewFavoriteClearCommand creates a new favorite clear command.
func NewFavoriteClearCommand(g *GlobalOption, promptuiProxy promptuiproxy.Promptui) *cobraproxy.CommandInstance {
	o := &favoriteClearOption{
		Out:     g.Out,
		ErrOut:  g.ErrOut,
		Args:    g.Args,
		Utility: g.Utility,
	}
	o.DBFileDirPathProvider = dbfiledirpathprovider.New(
		filepathproxy.New(),
		osproxy.New(),
		userproxy.New(),
	)
	o.JrpRepository = repository.New(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		stringsproxy.New(),
	)
	o.PromptuiProxy = promptuiProxy

	cobraProxy := cobraproxy.New()
	cmd := cobraProxy.NewCommand()

	cmd.FieldCommand.Use = constant.FAVORITE_CLEAR_USE
	cmd.FieldCommand.Aliases = constant.GetFavoriteClearAliases()
	cmd.FieldCommand.RunE = o.favoriteClearRunE

	cmd.PersistentFlags().BoolVarP(
		&o.NoConfirm,
		constant.FAVORITE_CLEAR_FLAG_NO_CONFIRM,
		constant.FAVORITE_CLEAR_FLAG_NO_CONFIRM_SHORTHAND,
		constant.FAVORITE_CLEAR_FLAG_NO_CONFIRM_DEFAULT,
		constant.FAVORITE_CLEAR_FLAG_NO_CONFIRM_DESCRIPTION,
	)

	cmd.SetOut(g.Out)
	cmd.SetErr(g.ErrOut)
	cmd.SetHelpTemplate(constant.FAVORITE_CLEAR_HELP_TEMPLATE)

	return cmd
}

// favoriteClearRunE is the function that is called when the favorite clear command is executed.
func (o *favoriteClearOption) favoriteClearRunE(_ *cobra.Command, _ []string) error {
	// get jrp db file dir path
	jrpDBFileDirPath, err := o.DBFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		return err
	}

	// create the directory if it does not exist
	if err := o.Utility.CreateDirIfNotExist(jrpDBFileDirPath); err != nil {
		return err
	}

	// write prompt if no-confirm option is not set
	if !o.NoConfirm {
		colorProxy := colorproxy.New()
		prompt := o.PromptuiProxy.NewPrompt()
		prompt.SetLabel(colorProxy.YellowString(constant.FAVORITE_CLEAR_PROMPT_LABEL))
		result, err := prompt.Run()
		if err != nil {
			return err
		}
		if result != "y" && result != "Y" {
			// write canceled message and do nothing
			o.Utility.PrintlnWithWriter(o.Out, constant.FAVORITE_CLEAR_MESSAGE_CLEAR_CANCELED)
			return nil
		}
	}

	filepathProxy := filepathproxy.New()
	return o.favoriteClear(filepathProxy.Join(jrpDBFileDirPath, repository.JRP_DB_FILE_NAME))
}

// favoriteClear clears all favorite.
func (o *favoriteClearOption) favoriteClear(jrpDBFilePath string) error {
	// remove all favorite
	res, err := o.JrpRepository.RemoveFavoriteAll(jrpDBFilePath)
	o.writeFavoriteClearResult(res)

	return err
}

// writeFavoriteClearResult writes the result of favorite clear.
func (o *favoriteClearOption) writeFavoriteClearResult(result repository.RemoveStatus) {
	var out = o.Out
	var message string
	colorProxy := colorproxy.New()
	if result == repository.RemovedFailed {
		out = o.ErrOut
		message = colorProxy.RedString(constant.FAVORITE_CLEAR_MESSAGE_CLEARED_FAIRULE)
	} else if result == repository.RemovedNone {
		message = colorProxy.YellowString(constant.FAVORITE_CLEAR_MESSAGE_CLEARED_NONE)
	} else {
		message = colorProxy.GreenString(constant.FAVORITE_CLEAR_MESSAGE_CLEARED_SUCCESSFULLY)
	}
	o.Utility.PrintlnWithWriter(out, message)
}
