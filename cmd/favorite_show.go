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

// favoriteShowOption is the struct for favorite show command.
type favoriteShowOption struct {
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

// NewFavoriteShowCommand creates a new favorite show command.
func NewFavoriteShowCommand(g *GlobalOption) *cobraproxy.CommandInstance {
	o := &favoriteShowOption{
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

	cmd.FieldCommand.Use = constant.FAVORITE_SHOW_USE
	cmd.FieldCommand.Aliases = constant.GetFavoriteShowAliases()
	cmd.FieldCommand.Args = cobra.MaximumNArgs(1)
	cmd.FieldCommand.RunE = o.favoriteShowRunE

	cmd.PersistentFlags().IntVarP(
		&o.Number,
		constant.FAVORITE_SHOW_FLAG_NUMBER,
		constant.FAVORITE_SHOW_FLAG_NUMBER_SHORTHAND,
		constant.FAVORITE_SHOW_FLAG_NUMBER_DEFAULT,
		constant.FAVORITE_SHOW_FLAG_NUMBER_DESCRIPTION,
	)
	cmd.PersistentFlags().BoolVarP(&o.All,
		constant.FAVORITE_SHOW_FLAG_ALL,
		constant.FAVORITE_SHOW_FLAG_ALL_SHORTHAND,
		constant.FAVORITE_SHOW_FLAG_ALL_DEFAULT,
		constant.FAVORITE_SHOW_FLAG_ALL_DESCRIPTION,
	)
	cmd.PersistentFlags().BoolVarP(&o.Plain,
		constant.FAVORITE_SHOW_FLAG_PLAIN,
		constant.FAVORITE_SHOW_FLAG_PLAIN_SHORTHAND,
		constant.FAVORITE_SHOW_FLAG_PLAIN_DEFAULT,
		constant.FAVORITE_SHOW_FLAG_PLAIN_DESCRIPTION,
	)

	cmd.SetOut(g.Out)
	cmd.SetErr(g.ErrOut)
	cmd.SetHelpTemplate(constant.FAVORITE_SHOW_HELP_TEMPLATE)

	cmd.SetArgs(o.Args)
	return cmd
}

// favoriteShowRunE is the function that is called when the favorite show command is executed.
func (o *favoriteShowOption) favoriteShowRunE(_ *cobra.Command, _ []string) error {
	strconvProxy := strconvproxy.New()
	if len(o.Args) <= 2 {
		// ifno arguments are given, set default value to args
		o.Args = []string{constant.FAVORITE_USE, constant.FAVORITE_SHOW_USE, strconvProxy.Itoa(constant.FAVORITE_SHOW_FLAG_NUMBER_DEFAULT)}
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
	return o.favoriteShow(filepathProxy.Join(jrpDBFileDirPath, repository.JRP_DB_FILE_NAME))
}

// favoriteShow is the function that gets the favorite.
func (o *favoriteShowOption) favoriteShow(jrpDBFilePath string) error {
	var favorites []*model.Jrp
	var err error
	if o.All {
		// if all flag is set, get all favorite
		favorites, err = o.JrpRepository.GetAllFavorite(jrpDBFilePath)
	} else {
		strconvProxy := strconvproxy.New()
		// get the larger number between the given number flag and the largest number that can be converted from the args
		num := o.Utility.GetLargerNumber(
			o.Number,
			o.Utility.GetMaxConvertibleString(
				o.Args,
				strconvProxy.Itoa(constant.FAVORITE_SHOW_FLAG_NUMBER_DEFAULT),
			),
		)
		if o.Number != num && o.Number > 0 {
			// if the number flag is littler than the default number, set the number flag value to num
			num = o.Number
		}
		// get favorite with the given number
		favorites, err = o.JrpRepository.GetFavoriteWithNumber(
			jrpDBFilePath,
			num,
		)
	}
	o.writeFavoriteShowResult(favorites)

	return err
}

// writeFavoriteShowResult writes the favorite show result.
func (o *favoriteShowOption) writeFavoriteShowResult(favorites []*model.Jrp) {
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
		o.Utility.PrintlnWithWriter(o.Out, colorProxy.YellowString(constant.FAVORITE_SHOW_MESSAGE_NO_FAVORITE_FOUND))
	}
}
