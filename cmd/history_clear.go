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
	"github.com/yanosea/jrp/app/proxy/strings"
	"github.com/yanosea/jrp/app/proxy/user"
	"github.com/yanosea/jrp/cmd/constant"
)

// historyClearOption is the struct for history clear command.
type historyClearOption struct {
	Out                   ioproxy.WriterInstanceInterface
	ErrOut                ioproxy.WriterInstanceInterface
	Args                  []string
	Force                 bool
	DBFileDirPathProvider dbfiledirpathprovider.DBFileDirPathProvidable
	JrpRepository         repository.JrpRepositoryInterface
	Utility               utility.UtilityInterface
}

// NewHistoryClearCommand creates a new history clear command.
func NewHistoryClearCommand(g *GlobalOption) *cobraproxy.CommandInstance {
	o := &historyClearOption{
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

	cmd.FieldCommand.Use = constant.HISTORY_CLEAR_USE
	cmd.FieldCommand.Aliases = constant.GetHistoryClearAliases()
	cmd.FieldCommand.Short = constant.HISTORY_CLEAR_SHORT
	cmd.FieldCommand.Long = constant.HISTORY_CLEAR_LONG
	cmd.FieldCommand.RunE = o.historyClearRunE

	cmd.PersistentFlags().BoolVarP(
		&o.Force,
		constant.HISTORY_CLEAR_FLAG_FORCE,
		constant.HISTORY_CLEAR_FLAG_FORCE_SHORTHAND,
		constant.HISTORY_CLEAR_FLAG_FORCE_DEFAULT,
		constant.HISTORY_CLEAR_FLAG_FORCE_DESCRIPTION,
	)
	cmd.SetOut(g.Out)
	cmd.SetErr(g.ErrOut)
	cmd.SetHelpTemplate(constant.HISTORY_CLEAR_HELP_TEMPLATE)

	return cmd
}

// historyClearRunE is the function that is called when the history clear command is executed.
func (o *historyClearOption) historyClearRunE(_ *cobra.Command, _ []string) error {
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
	return o.historyClear(filepathProxy.Join(jrpDBFileDirPath, repository.JRP_DB_FILE_NAME))
}

// historyClear clears all history.
func (o *historyClearOption) historyClear(jrpDBFilePath string) error {
	// remove all history
	res, err := o.JrpRepository.RemoveHistoryAll(jrpDBFilePath, o.Force)
	o.writeHistoryClearResult(res)

	return err
}

// writeHistoryClearResult writes the result of history clear.
func (o *historyClearOption) writeHistoryClearResult(result repository.RemoveStatus) {
	var out = o.Out
	var message string
	colorProxy := colorproxy.New()
	if result == repository.RemovedFailed {
		out = o.ErrOut
		message = colorProxy.RedString(constant.HISTORY_CLEAR_MESSAGE_CLEARED_FAIRULE)
	} else if result == repository.RemovedNone && !o.Force {
		message = colorProxy.YellowString(constant.HISTORY_CLEAR_MESSAGE_CLEARED_NONE)
	} else {
		message = colorProxy.GreenString(constant.HISTORY_CLEAR_MESSAGE_CLEARED_SUCCESSFULLY)
	}
	o.Utility.PrintlnWithWriter(out, message)
}
