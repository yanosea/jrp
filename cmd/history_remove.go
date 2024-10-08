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
	"github.com/yanosea/jrp/app/proxy/strconv"
	"github.com/yanosea/jrp/app/proxy/strings"
	"github.com/yanosea/jrp/app/proxy/user"
	"github.com/yanosea/jrp/cmd/constant"
)

// historyRemoveOption is the struct for history remove command.
type historyRemoveOption struct {
	Out                   ioproxy.WriterInstanceInterface
	ErrOut                ioproxy.WriterInstanceInterface
	Args                  []string
	All                   bool
	Force                 bool
	NoConfirm             bool
	DBFileDirPathProvider dbfiledirpathprovider.DBFileDirPathProvidable
	JrpRepository         repository.JrpRepositoryInterface
	PromptuiProxy         promptuiproxy.Promptui
	Utility               utility.UtilityInterface
}

// NewHistoryRemoveCommand creates a new history remove command.
func NewHistoryRemoveCommand(g *GlobalOption, promptuiProxy promptuiproxy.Promptui) *cobraproxy.CommandInstance {
	o := &historyRemoveOption{
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

	cmd.FieldCommand.Use = constant.HISTORY_REMOVE_USE
	cmd.FieldCommand.Aliases = constant.GetHistoryRemoveAliases()
	cmd.FieldCommand.RunE = o.historyRemoveRunE

	cmd.PersistentFlags().BoolVarP(
		&o.All,
		constant.HISTORY_REMOVE_FLAG_ALL,
		constant.HISTORY_REMOVE_FLAG_ALL_SHORTHAND,
		constant.HISTORY_REMOVE_FLAG_ALL_DEFAULT,
		constant.HISTORY_REMOVE_FLAG_ALL_DESCRIPTION,
	)
	cmd.PersistentFlags().BoolVarP(
		&o.Force,
		constant.HISTORY_REMOVE_FLAG_FORCE,
		constant.HISTORY_REMOVE_FLAG_FORCE_SHORTHAND,
		constant.HISTORY_REMOVE_FLAG_FORCE_DEFAULT,
		constant.HISTORY_REMOVE_FLAG_FORCE_DESCRIPTION,
	)
	cmd.PersistentFlags().BoolVarP(
		&o.NoConfirm,
		constant.HISTORY_REMOVE_FLAG_NO_CONFIRM,
		constant.HISTORY_REMOVE_FLAG_NO_CONFIRM_SHORTHAND,
		constant.HISTORY_REMOVE_FLAG_NO_CONFIRM_DEFAULT,
		constant.HISTORY_REMOVE_FLAG_NO_CONFIRM_DESCRIPTION,
	)

	cmd.SetOut(g.Out)
	cmd.SetErr(g.ErrOut)
	cmd.SetHelpTemplate(constant.HISTORY_REMOVE_HELP_TEMPLATE)

	return cmd
}

// historyRemoveRunE is the function that is called when the history remove command is executed.
func (o *historyRemoveOption) historyRemoveRunE(_ *cobra.Command, _ []string) error {
	if len(o.Args) <= 2 {
		// if no arguments is given, set default value to args
		o.Args = []string{constant.HISTORY_USE, constant.HISTORY_REMOVE_USE, ""}
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
		o.Utility.PrintlnWithWriter(o.Out, colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_NO_ID_SPECIFIED))
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

	// write prompt if all flag is set and no-confirm option is not set
	if o.All && !o.NoConfirm {
		colorProxy := colorproxy.New()
		prompt := o.PromptuiProxy.NewPrompt()
		prompt.SetLabel(colorProxy.YellowString(constant.HISTORY_REMOVE_PROMPT_REMOVE_ALL_LABEL))
		result, err := prompt.Run()
		if err != nil {
			return err
		}
		if result != "y" && result != "Y" {
			// write canceled message and do nothing
			o.Utility.PrintlnWithWriter(o.Out, constant.HISTORY_REMOVE_MESSAGE_REMOVE_ALL_CANCELED)
			return nil
		}
	}

	filepathProxy := filepathproxy.New()
	return o.historyRemove(filepathProxy.Join(jrpDBFileDirPath, repository.JRP_DB_FILE_NAME), IDs)
}

// historyRemove is the function that removes history by IDs.
func (o *historyRemoveOption) historyRemove(jrpDBFilePath string, IDs []int) error {
	var res repository.RemoveStatus
	var err error
	if o.All {
		// if all flag is set, remove all history
		res, err = o.JrpRepository.RemoveHistoryAll(jrpDBFilePath, o.Force)
	} else {
		// if IDs are specified, remove history by IDs
		res, err = o.JrpRepository.RemoveHistoryByIDs(jrpDBFilePath, IDs, o.Force)
	}
	o.writeHistoryRemoveResult(res)

	return err
}

// writeHistoryRemoveResult writes the result of history remove.
func (o *historyRemoveOption) writeHistoryRemoveResult(result repository.RemoveStatus) {
	var out = o.Out
	var message string
	colorProxy := colorproxy.New()
	if result == repository.RemovedFailed {
		out = o.ErrOut
		message = colorProxy.RedString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_FAILURE)
	} else if result == repository.RemovedNone {
		message = colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_NONE)
	} else if result == repository.RemovedNotAll {
		message = colorProxy.YellowString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_NOT_ALL)
	} else {
		message = colorProxy.GreenString(constant.HISTORY_REMOVE_MESSAGE_REMOVED_SUCCESSFULLY)
	}
	o.Utility.PrintlnWithWriter(out, message)
}
