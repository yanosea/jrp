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

// favoriteOption is the struct for favorite command.
type favoriteOption struct {
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

// NewFavoriteCommand creates a new favorite command.
func NewFavoriteCommand(g *GlobalOption) *cobraproxy.CommandInstance {
	o := &favoriteOption{
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

	cmd.FieldCommand.Use = constant.FAVORITE_USE
	cmd.FieldCommand.Aliases = constant.GetFavoriteAliases()
	cmd.FieldCommand.Args = cobra.MaximumNArgs(1)
	cmd.FieldCommand.RunE = o.favoriteRunE

	cmd.PersistentFlags().IntVarP(
		&o.Number,
		constant.FAVORITE_FLAG_NUMBER,
		constant.FAVORITE_FLAG_NUMBER_SHORTHAND,
		constant.FAVORITE_FLAG_NUMBER_DEFAULT,
		constant.FAVORITE_FLAG_NUMBER_DESCRIPTION,
	)
	cmd.PersistentFlags().BoolVarP(&o.All,
		constant.FAVORITE_FLAG_ALL,
		constant.FAVORITE_FLAG_ALL_SHORTHAND,
		constant.FAVORITE_FLAG_ALL_DEFAULT,
		constant.FAVORITE_FLAG_ALL_DESCRIPTION,
	)
	cmd.PersistentFlags().BoolVarP(&o.Plain,
		constant.FAVORITE_FLAG_PLAIN,
		constant.FAVORITE_FLAG_PLAIN_SHORTHAND,
		constant.FAVORITE_FLAG_PLAIN_DEFAULT,
		constant.FAVORITE_FLAG_PLAIN_DESCRIPTION,
	)

	cmd.SetOut(g.Out)
	cmd.SetErr(g.ErrOut)
	cmd.SetHelpTemplate(constant.FAVORITE_HELP_TEMPLATE)

	cmd.AddCommand(
		NewFavoriteShowCommand(g),
		NewFavoriteAddCommand(g),
		NewFavoriteRemoveCommand(g, promptuiproxy.New()),
		NewFavoriteSearchCommand(g),
		NewFavoriteClearCommand(g, promptuiproxy.New()),
	)

	return cmd
}

// favoriteRunE is the function that is called when the favorite command is executed.
func (o *favoriteOption) favoriteRunE(_ *cobra.Command, _ []string) error {
	strconvProxy := strconvproxy.New()
	if len(o.Args) <= 1 {
		// if no argument is given, set the default value to args
		o.Args = []string{constant.FAVORITE_USE, strconvProxy.Itoa(constant.FAVORITE_FLAG_NUMBER_DEFAULT)}
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
	return o.favorite(filepathProxy.Join(jrpDBFileDirPath, repository.JRP_DB_FILE_NAME))
}

// favorite shows the favorite.
func (o *favoriteOption) favorite(jrpDBFilePath string) error {
	var favorites []model.Jrp
	var err error
	if o.All {
		// if all flag is set, get all favorite
		favorites, err = o.JrpRepository.GetAllFavorite(jrpDBFilePath)
	} else {
		if o.Number != constant.FAVORITE_FLAG_NUMBER_DEFAULT && o.Number >= 1 {
			// if number flag is set, get favorites with the given number
			favorites, err = o.JrpRepository.GetFavoriteWithNumber(jrpDBFilePath, o.Number)
		} else {
			strconvProxy := strconvproxy.New()
			// get favorite with the given number
			favorites, err = o.JrpRepository.GetFavoriteWithNumber(
				jrpDBFilePath,
				// get the larger number between the given number flag and the largest number that can be converted from the args
				o.Utility.GetLargerNumber(
					o.Number,
					o.Utility.GetMaxConvertibleString(
						o.Args,
						strconvProxy.Itoa(constant.FAVORITE_FLAG_NUMBER_DEFAULT),
					),
				),
			)
		}
	}
	o.writeFavoriteResult(favorites)

	return err
}

// writeFavoriteResult writes the favorite result.
func (o *favoriteOption) writeFavoriteResult(favorites []model.Jrp) {
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
		o.Utility.PrintlnWithWriter(o.Out, colorProxy.YellowString(constant.FAVORITE_MESSAGE_NO_FAVORITE_FOUND))
	}
}
