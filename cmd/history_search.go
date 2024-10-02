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
	"github.com/yanosea/jrp/app/proxy/sort"
	"github.com/yanosea/jrp/app/proxy/sql"
	"github.com/yanosea/jrp/app/proxy/strconv"
	"github.com/yanosea/jrp/app/proxy/strings"
	"github.com/yanosea/jrp/app/proxy/tablewriter"
	"github.com/yanosea/jrp/app/proxy/user"
	"github.com/yanosea/jrp/cmd/constant"
)

// historySearchOption is the struct for history search command.
type historySearchOption struct {
	Out                   ioproxy.WriterInstanceInterface
	ErrOut                ioproxy.WriterInstanceInterface
	Args                  []string
	And                   bool
	Number                int
	All                   bool
	Plain                 bool
	DBFileDirPathProvider dbfiledirpathprovider.DBFileDirPathProvidable
	JrpRepository         repository.JrpRepositoryInterface
	JrpWriter             jrpwriter.JrpWritable
	Utility               utility.UtilityInterface
}

// NewHistorySearchCommand creates a new history search command.
func NewHistorySearchCommand(g *GlobalOption) *cobraproxy.CommandInstance {
	o := &historySearchOption{
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

	cmd.FieldCommand.Use = constant.HISTORY_SEARCH_USE
	cmd.FieldCommand.Aliases = constant.GetHistorySearchAliases()
	cmd.FieldCommand.Short = constant.HISTORY_SEARCH_SHORT
	cmd.FieldCommand.Long = constant.HISTORY_SEARCH_LONG
	cmd.FieldCommand.RunE = o.historySearchRunE

	cmd.PersistentFlags().BoolVarP(
		&o.And,
		constant.HISTORY_SEARCH_FLAG_AND,
		constant.HISTORY_SEARCH_FLAG_AND_SHORTHAND,
		constant.HISTORY_SEARCH_FLAG_AND_DEFAULT,
		constant.HISTORY_SEARCH_FLAG_AND_DESCRIPTION,
	)
	cmd.PersistentFlags().IntVarP(&o.Number,
		constant.HISTORY_SEARCH_FLAG_NUMBER,
		constant.HISTORY_SEARCH_FLAG_NUMBER_SHORTHAND,
		constant.HISTORY_SEARCH_FLAG_NUMBER_DEFAULT,
		constant.HISTORY_SEARCH_FLAG_NUMBER_DESCRIPTION,
	)
	cmd.PersistentFlags().BoolVarP(
		&o.All,
		constant.HISTORY_SEARCH_FLAG_ALL,
		constant.HISTORY_SEARCH_FLAG_ALL_SHORTHAND,
		constant.HISTORY_SEARCH_FLAG_ALL_DEFAULT,
		constant.HISTORY_SEARCH_FLAG_ALL_DESCRIPTION,
	)
	cmd.PersistentFlags().BoolVarP(&o.Plain,
		constant.HISTORY_SEARCH_FLAG_PLAIN,
		constant.HISTORY_SEARCH_FLAG_PLAIN_SHORTHAND,
		constant.HISTORY_SEARCH_FLAG_PLAIN_DEFAULT,
		constant.HISTORY_SEARCH_FLAG_PLAIN_DESCRIPTION,
	)

	cmd.SetOut(g.Out)
	cmd.SetErr(g.ErrOut)
	cmd.SetHelpTemplate(constant.HISTORY_SEARCH_HELP_TEMPLATE)

	return cmd
}

// historySearchRunE is the function that is called when the history search command is executed.
func (o *historySearchOption) historySearchRunE(_ *cobra.Command, _ []string) error {
	if len(o.Args) <= 2 {
		// if no arguments is given, set default value to args
		o.Args = []string{constant.HISTORY_USE, constant.HISTORY_SEARCH_USE, ""}
	}

	// set keywords
	var keywords []string
	for _, arg := range o.Args[2:] {
		if arg != "" {
			keywords = append(keywords, arg)
		}
	}
	if keywords == nil || len(keywords) == 0 {
		// if no keywords are provided, write message and return
		colorProxy := colorproxy.New()
		o.Utility.PrintlnWithWriter(o.Out, colorProxy.YellowString(constant.HISTORY_SEARCH_MESSAGE_NO_KEYWORDS_PROVIDED))
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
	return o.historySearch(filepathProxy.Join(jrpDBFileDirPath, repository.JRP_DB_FILE_NAME), keywords)
}

// historySearch searches the history with the given keywords.
func (o *historySearchOption) historySearch(jrpDBFilePath string, keywords []string) error {
	var histories []model.Jrp
	var err error
	if o.All {
		// if all flag is set, search all history
		histories, err = o.JrpRepository.SearchAllHistory(jrpDBFilePath, keywords, o.And)
	} else {
		// search history with the given number
		histories, err = o.JrpRepository.SearchHistoryWithNumber(
			jrpDBFilePath,
			o.Number,
			keywords,
			o.And,
		)
	}
	o.writeHistorySearchResult(histories)

	return err
}

// writeHistorySearchResult writes the history search result.
func (o *historySearchOption) writeHistorySearchResult(histories []model.Jrp) {
	if len(histories) != 0 {
		if o.Plain {
			for _, hist := range histories {
				// if plain flag is set, write only the phrase
				o.Utility.PrintlnWithWriter(o.Out, hist.Phrase)
			}
		} else {
			// if plain flag is not set, write the history as a table
			o.JrpWriter.WriteAsTable(o.Out, histories)
		}
	} else {
		// if no history is found, write the message
		colorProxy := colorproxy.New()
		o.Utility.PrintlnWithWriter(o.Out, colorProxy.YellowString(constant.HISTORY_SEARCH_MESSAGE_NO_RESULT_FOUND))
	}
}
