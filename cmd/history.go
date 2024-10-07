package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/app/database/jrp/model"
	"github.com/yanosea/jrp/app/database/jrp/repository"
	"github.com/yanosea/jrp/app/library/dbfiledirpathprovider"
	"github.com/yanosea/jrp/app/library/jrpwriter"
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
	"github.com/yanosea/jrp/app/proxy/tablewriter"
	"github.com/yanosea/jrp/app/proxy/user"
	"github.com/yanosea/jrp/cmd/constant"
)

// historyOption is the struct for history command.
type historyOption struct {
	Out                   ioproxy.WriterInstanceInterface
	ErrOut                ioproxy.WriterInstanceInterface
	Args                  []string
	Number                int
	All                   bool
	Plain                 bool
	DBFileDirPathProvider dbfiledirpathprovider.DBFileDirPathProvidable
	JrpRepository         repository.JrpRepositoryInterface
	JrpWriter             jrpwriter.JrpWritable
	Utility               utility.UtilityInterface
}

// NewHistoryCommand creates a new history command.
func NewHistoryCommand(g *GlobalOption) *cobraproxy.CommandInstance {
	o := &historyOption{
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
	o.JrpWriter = jrpwriter.New(
		strconvproxy.New(),
		tablewriterproxy.New(),
	)

	cobraProxy := cobraproxy.New()
	cmd := cobraProxy.NewCommand()
	cmd.FieldCommand.Use = constant.HISTORY_USE
	cmd.FieldCommand.Aliases = constant.GetHistoryAliases()
	cmd.FieldCommand.Args = cobra.MaximumNArgs(1)
	cmd.FieldCommand.RunE = o.historyRunE

	cmd.PersistentFlags().IntVarP(&o.Number,
		constant.HISTORY_FLAG_NUMBER,
		constant.HISTORY_FLAG_NUMBER_SHORTHAND,
		constant.HISTORY_FLAG_NUMBER_DEFAULT,
		constant.HISTORY_FLAG_NUMBER_DESCRIPTION,
	)
	cmd.PersistentFlags().BoolVarP(&o.All,
		constant.HISTORY_FLAG_ALL,
		constant.HISTORY_FLAG_ALL_SHORTHAND,
		constant.HISTORY_FLAG_ALL_DEFAULT,
		constant.HISTORY_FLAG_ALL_DESCRIPTION,
	)
	cmd.PersistentFlags().BoolVarP(&o.Plain,
		constant.HISTORY_FLAG_PLAIN,
		constant.HISTORY_FLAG_PLAIN_SHORTHAND,
		constant.HISTORY_FLAG_PLAIN_DEFAULT,
		constant.HISTORY_FLAG_PLAIN_DESCRIPTION,
	)

	cmd.SetOut(g.Out)
	cmd.SetErr(g.ErrOut)
	cmd.SetHelpTemplate(constant.HISTORY_HELP_TEMPLATE)

	cmd.AddCommand(
		NewHistoryShowCommand(g),
		NewHistoryRemoveCommand(g, promptuiproxy.New()),
		NewHistorySearchCommand(g),
		NewHistoryClearCommand(g, promptuiproxy.New()),
	)

	return cmd
}

// historyRunE is the function that is called when the history command is executed.
func (o *historyOption) historyRunE(_ *cobra.Command, _ []string) error {
	strconvProxy := strconvproxy.New()
	if len(o.Args) <= 1 {
		// if no argument is given, set the default value to args
		o.Args = []string{constant.HISTORY_USE, strconvProxy.Itoa(constant.HISTORY_FLAG_NUMBER_DEFAULT)}
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
	return o.history(filepathProxy.Join(jrpDBFileDirPath, repository.JRP_DB_FILE_NAME))
}

// history searches the history.
func (o *historyOption) history(jrpDBFilePath string) error {
	var histories []model.Jrp
	var err error
	if o.All {
		// if all flag is set, get all history
		histories, err = o.JrpRepository.GetAllHistory(jrpDBFilePath)
	} else {
		if o.Number != constant.HISTORY_FLAG_NUMBER_DEFAULT && o.Number >= 1 {
			// if number flag is set, get history with the given number
			histories, err = o.JrpRepository.GetHistoryWithNumber(jrpDBFilePath, o.Number)
		} else {
			strconvProxy := strconvproxy.New()
			// get history with the given number
			histories, err = o.JrpRepository.GetHistoryWithNumber(
				jrpDBFilePath,
				// get the larger number between the given number flag and the largest number that can be converted from the args
				o.Utility.GetLargerNumber(
					o.Number,
					o.Utility.GetMaxConvertibleString(
						o.Args,
						strconvProxy.Itoa(constant.HISTORY_FLAG_NUMBER_DEFAULT),
					),
				),
			)
		}
	}
	o.writeHistoryResult(histories)

	return err
}

// writeHistoryResult writes the history result.
func (o *historyOption) writeHistoryResult(histories []model.Jrp) {
	if len(histories) != 0 {
		if o.Plain {
			for _, history := range histories {
				// if plain flag is set, write only the phrase
				o.Utility.PrintlnWithWriter(o.Out, history.Phrase)
			}
		} else {
			// if plain flag is not set, write the history as a table
			o.JrpWriter.WriteAsTable(o.Out, histories)
		}
	} else {
		// if no history is found, write the message
		colorProxy := colorproxy.New()
		o.Utility.PrintlnWithWriter(o.Out, colorProxy.YellowString(constant.HISTORY_MESSAGE_NO_HISTORY_FOUND))
	}
}
