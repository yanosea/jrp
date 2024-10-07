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

// favoriteAddOption is the struct for favorite add command.
type favoriteAddOption struct {
	Out                   ioproxy.WriterInstanceInterface
	ErrOut                ioproxy.WriterInstanceInterface
	Args                  []string
	DBFileDirPathProvider dbfiledirpathprovider.DBFileDirPathProvidable
	JrpRepository         repository.JrpRepositoryInterface
	Utility               utility.UtilityInterface
}

// NewFavoriteAddCommand creates a new favorite add command.
func NewFavoriteAddCommand(g *GlobalOption) *cobraproxy.CommandInstance {
	o := &favoriteAddOption{
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

	cmd.FieldCommand.Use = constant.FAVORITE_ADD_USE
	cmd.FieldCommand.Aliases = constant.GetFavoriteAddAliases()
	cmd.FieldCommand.RunE = o.favoriteAddRunE

	cmd.SetOut(g.Out)
	cmd.SetErr(g.ErrOut)
	cmd.SetHelpTemplate(constant.FAVORITE_ADD_HELP_TEMPLATE)

	return cmd
}

// favoriteAddRunE is the function that is called when the favorite add command is executed.
func (o *favoriteAddOption) favoriteAddRunE(_ *cobra.Command, _ []string) error {
	if len(o.Args) <= 2 {
		// if no arguments is given, set default value to args
		o.Args = []string{constant.FAVORITE_USE, constant.FAVORITE_ADD_USE, ""}
	}

	// set ID
	strconvProxy := strconvproxy.New()
	var IDs []int
	for _, arg := range o.Args[2:] {
		if id, err := strconvProxy.Atoi(arg); err != nil {
			continue
		} else {
			IDs = append(IDs, id)
		}
	}
	if len(IDs) == 0 {
		// if no ID is specified, print write and return
		colorProxy := colorproxy.New()
		o.Utility.PrintlnWithWriter(o.Out, colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_NO_ID_SPECIFIED))
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
	return o.favoriteAdd(filepathProxy.Join(jrpDBFileDirPath, repository.JRP_DB_FILE_NAME), IDs)
}

// favoriteAdd adds the specified ID to favorite.
func (o *favoriteAddOption) favoriteAdd(jrpDBFilePath string, IDs []int) error {
	// if ID is specified, add to favorite
	res, err := o.JrpRepository.AddFavoriteByIDs(jrpDBFilePath, IDs)
	o.writeFavoriteAddResult(res)

	return err
}

// writeFavoriteAddResult writes the result of favorite add.
func (o *favoriteAddOption) writeFavoriteAddResult(result repository.AddStatus) {
	var out = o.Out
	var message string
	colorProxy := colorproxy.New()
	if result == repository.AddedFailed {
		out = o.ErrOut
		message = colorProxy.RedString(constant.FAVORITE_ADD_MESSAGE_ADDED_FAILURE)
	} else if result == repository.AddedNone {
		message = colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NONE)
	} else if result == repository.AddedNotAll {
		message = colorProxy.YellowString(constant.FAVORITE_ADD_MESSAGE_ADDED_NOT_ALL)
	} else {
		message = colorProxy.GreenString(constant.FAVORITE_ADD_MESSAGE_ADDED_SUCCESSFULLY)
	}
	o.Utility.PrintlnWithWriter(out, message)
}
