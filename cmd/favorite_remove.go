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
	"github.com/yanosea/jrp/app/proxy/sort"
	"github.com/yanosea/jrp/app/proxy/sql"
	"github.com/yanosea/jrp/app/proxy/strconv"
	"github.com/yanosea/jrp/app/proxy/strings"
	"github.com/yanosea/jrp/app/proxy/user"
	"github.com/yanosea/jrp/cmd/constant"
)

// favoriteRemoveOption is the struct for favorite remove command.
type favoriteRemoveOption struct {
	Out                   ioproxy.WriterInstanceInterface
	ErrOut                ioproxy.WriterInstanceInterface
	Args                  []string
	All                   bool
	DBFileDirPathProvider dbfiledirpathprovider.DBFileDirPathProvidable
	JrpRepository         repository.JrpRepositoryInterface
	Utility               utility.UtilityInterface
}

// NewFavoriteRemoveCommand creates a new favorite remove command.
func NewFavoriteRemoveCommand(g *GlobalOption) *cobraproxy.CommandInstance {
	o := &favoriteRemoveOption{
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

	cobraProxy := cobraproxy.New()
	cmd := cobraProxy.NewCommand()

	cmd.FieldCommand.Use = constant.FAVORITE_REMOVE_USE
	cmd.FieldCommand.Aliases = constant.GetFavoriteRemoveAliases()
	cmd.FieldCommand.RunE = o.favoriteRemoveRunE

	cmd.PersistentFlags().BoolVarP(
		&o.All,
		constant.FAVORITE_REMOVE_FLAG_ALL,
		constant.FAVORITE_REMOVE_FLAG_ALL_SHORTHAND,
		constant.FAVORITE_REMOVE_FLAG_ALL_DEFAULT,
		constant.FAVORITE_REMOVE_FLAG_ALL_DESCRIPTION,
	)

	cmd.SetOut(g.Out)
	cmd.SetErr(g.ErrOut)
	cmd.SetHelpTemplate(constant.FAVORITE_REMOVE_HELP_TEMPLATE)

	return cmd
}

// favoriteRemoveRunE is the function that is called when the favorite remove command is executed.
func (o *favoriteRemoveOption) favoriteRemoveRunE(_ *cobra.Command, _ []string) error {
	if len(o.Args) <= 2 {
		// if no arguments is given, set default value to args
		o.Args = []string{constant.FAVORITE_USE, constant.FAVORITE_REMOVE_USE, ""}
	}

	// set ID
	strconvProxy := strconvproxy.New()
	var IDs []int
	if !o.All {
		for _, arg := range o.Args[2:] {
			if id, err := strconvProxy.Atoi(arg); err != nil {
				continue
			} else {
				IDs = append(IDs, id)
			}
		}
	}
	if len(IDs) == 0 && !o.All {
		// if no ID is specified, print write and return
		colorProxy := colorproxy.New()
		o.Utility.PrintlnWithWriter(o.Out, colorProxy.YellowString(constant.FAVORITE_REMOVE_MESSAGE_NO_ID_SPECIFIED))
		return nil
	}

	// get jrp db file dir path
	jrpDBFileDirPath, err := o.DBFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		return err
	}

	// create the directory if it does not exist
	if err := o.Utility.CreateDirIfNotExist(jrpDBFileDirPath); err != nil {
		return err
	}

	filepathProxy := filepathproxy.New()
	return o.favoriteRemove(filepathProxy.Join(jrpDBFileDirPath, repository.JRP_DB_FILE_NAME), IDs)
}

// favoriteRemove removes favorite by IDs.
func (o *favoriteRemoveOption) favoriteRemove(jrpDBFilePath string, IDs []int) error {
	var res repository.RemoveStatus
	var err error
	if o.All {
		// if all flag is set, remove all favorite
		res, err = o.JrpRepository.RemoveFavoriteAll(jrpDBFilePath)
	} else {
		// if IDs are specified, remove favorite by IDs
		res, err = o.JrpRepository.RemoveFavoriteByIDs(jrpDBFilePath, IDs)
	}
	o.writeFavoriteRemoveResult(res)

	return err
}

// writeFavoriteRemoveResult writes the result of favorite remove.
func (o *favoriteRemoveOption) writeFavoriteRemoveResult(result repository.RemoveStatus) {
	var out = o.Out
	var message string
	colorProxy := colorproxy.New()
	if result == repository.RemovedFailed {
		out = o.ErrOut
		message = colorProxy.RedString(constant.FAVORITE_REMOVE_MESSAGE_REMOVED_FAILURE)
	} else if result == repository.RemovedNone {
		message = colorProxy.YellowString(constant.FAVORITE_REMOVE_MESSAGE_REMOVED_NONE)
	} else if result == repository.RemovedNotAll {
		message = colorProxy.YellowString(constant.FAVORITE_REMOVE_MESSAGE_REMOVED_NOT_ALL)
	} else {
		message = colorProxy.GreenString(constant.FAVORITE_REMOVE_MESSAGE_REMOVED_SUCCESSFULLY)
	}
	o.Utility.PrintlnWithWriter(out, message)
}
