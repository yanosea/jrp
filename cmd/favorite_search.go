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

// favoriteSearchOption is the struct for favorite search command.
type favoriteSearchOption struct {
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

// NewFavoriteSearchCommand creates a new favorite search command.
func NewFavoriteSearchCommand(g *GlobalOption) *cobraproxy.CommandInstance {
	o := &favoriteSearchOption{
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

	cmd.FieldCommand.Use = constant.FAVORITE_SEARCH_USE
	cmd.FieldCommand.Aliases = constant.GetFavoriteSearchAliases()
	cmd.FieldCommand.Short = constant.FAVORITE_SEARCH_SHORT
	cmd.FieldCommand.Long = constant.FAVORITE_SEARCH_LONG
	cmd.FieldCommand.RunE = o.favoriteSearchRunE

	cmd.PersistentFlags().BoolVarP(
		&o.And,
		constant.FAVORITE_SEARCH_FLAG_AND,
		constant.FAVORITE_SEARCH_FLAG_AND_SHORTHAND,
		constant.FAVORITE_SEARCH_FLAG_AND_DEFAULT,
		constant.FAVORITE_SEARCH_FLAG_AND_DESCRIPTION,
	)
	cmd.PersistentFlags().IntVarP(&o.Number,
		constant.FAVORITE_SEARCH_FLAG_NUMBER,
		constant.FAVORITE_SEARCH_FLAG_NUMBER_SHORTHAND,
		constant.FAVORITE_SEARCH_FLAG_NUMBER_DEFAULT,
		constant.FAVORITE_SEARCH_FLAG_NUMBER_DESCRIPTION,
	)
	cmd.PersistentFlags().BoolVarP(
		&o.All,
		constant.FAVORITE_SEARCH_FLAG_ALL,
		constant.FAVORITE_SEARCH_FLAG_ALL_SHORTHAND,
		constant.FAVORITE_SEARCH_FLAG_ALL_DEFAULT,
		constant.FAVORITE_SEARCH_FLAG_ALL_DESCRIPTION,
	)
	cmd.PersistentFlags().BoolVarP(&o.Plain,
		constant.FAVORITE_SEARCH_FLAG_PLAIN,
		constant.FAVORITE_SEARCH_FLAG_PLAIN_SHORTHAND,
		constant.FAVORITE_SEARCH_FLAG_PLAIN_DEFAULT,
		constant.FAVORITE_SEARCH_FLAG_PLAIN_DESCRIPTION,
	)

	cmd.SetOut(g.Out)
	cmd.SetErr(g.ErrOut)
	cmd.SetHelpTemplate(constant.FAVORITE_SEARCH_HELP_TEMPLATE)

	return cmd
}

// favoriteSearchRunE is the function that is called when the favorite search command is executed.
func (o *favoriteSearchOption) favoriteSearchRunE(_ *cobra.Command, _ []string) error {
	if len(o.Args) <= 2 {
		// if no arguments is given, set default value to args
		o.Args = []string{constant.FAVORITE_USE, constant.FAVORITE_SEARCH_USE, ""}
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
		o.Utility.PrintlnWithWriter(o.Out, colorProxy.YellowString(constant.FAVORITE_SEARCH_MESSAGE_NO_KEYWORDS_PROVIDED))
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
	return o.favoriteSearch(filepathProxy.Join(jrpDBFileDirPath, repository.JRP_DB_FILE_NAME), keywords)
}

// favoriteSearch searches favorite.
func (o *favoriteSearchOption) favoriteSearch(jrpDBFilePath string, keywords []string) error {
	var favorites []model.Jrp
	var err error
	if o.All {
		// if all flag is set, search all favorite
		favorites, err = o.JrpRepository.SearchAllFavorite(jrpDBFilePath, keywords, o.And)
	} else {
		// search favorite with the given number
		favorites, err = o.JrpRepository.SearchFavoriteWithNumber(
			jrpDBFilePath,
			o.Number,
			keywords,
			o.And,
		)
	}
	o.writeFavoriteSearchResult(favorites)

	return err
}

// writeFavoriteSearchResult writes the favorite search result.
func (o *favoriteSearchOption) writeFavoriteSearchResult(favorites []model.Jrp) {
	if len(favorites) != 0 {
		if o.Plain {
			for _, favorite := range favorites {
				// if plain flag is set, write only the phrase
				o.Utility.PrintlnWithWriter(o.Out, favorite.Phrase)
			}
		} else {
			// if plain flag is not set, write the favorite as a table
			o.JrpWriter.WriteAsTable(o.Out, favorites)
		}
	} else {
		// if no favorite is found, write the message
		colorProxy := colorproxy.New()
		o.Utility.PrintlnWithWriter(o.Out, colorProxy.YellowString(constant.FAVORITE_SEARCH_MESSAGE_NO_RESULT_FOUND))
	}
}
