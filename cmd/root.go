package cmd

import (
	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/app/database/jrp/model"
	jrprepository "github.com/yanosea/jrp/app/database/jrp/repository"
	wnjpnrepository "github.com/yanosea/jrp/app/database/wnjpn/repository"
	"github.com/yanosea/jrp/app/library/dbfiledirpathprovider"
	"github.com/yanosea/jrp/app/library/generator"
	"github.com/yanosea/jrp/app/library/jrpwriter"
	"github.com/yanosea/jrp/app/library/utility"
	"github.com/yanosea/jrp/app/library/versionprovider"
	"github.com/yanosea/jrp/app/proxy/cobra"
	"github.com/yanosea/jrp/app/proxy/color"
	"github.com/yanosea/jrp/app/proxy/debug"
	"github.com/yanosea/jrp/app/proxy/filepath"
	"github.com/yanosea/jrp/app/proxy/fmt"
	"github.com/yanosea/jrp/app/proxy/io"
	"github.com/yanosea/jrp/app/proxy/keyboard"
	"github.com/yanosea/jrp/app/proxy/os"
	"github.com/yanosea/jrp/app/proxy/rand"
	"github.com/yanosea/jrp/app/proxy/sort"
	"github.com/yanosea/jrp/app/proxy/sql"
	"github.com/yanosea/jrp/app/proxy/strconv"
	"github.com/yanosea/jrp/app/proxy/strings"
	"github.com/yanosea/jrp/app/proxy/tablewriter"
	"github.com/yanosea/jrp/app/proxy/time"
	"github.com/yanosea/jrp/app/proxy/user"
	"github.com/yanosea/jrp/cmd/constant"
)

// ver is the version of the jrp.
var ver = ""

// GlobalOption is the struct for global option.
type GlobalOption struct {
	Out            ioproxy.WriterInstanceInterface
	ErrOut         ioproxy.WriterInstanceInterface
	Args           []string
	Utility        utility.UtilityInterface
	NewRootCommand func(ow, ew ioproxy.WriterInstanceInterface, args []string) cobraproxy.CommandInstanceInterface
}

// rootOption is the struct for root command.
type rootOption struct {
	Out                   ioproxy.WriterInstanceInterface
	ErrOut                ioproxy.WriterInstanceInterface
	Args                  []string
	Number                int
	Prefix                string
	Suffix                string
	DryRun                bool
	Plain                 bool
	Interactive           bool
	Timeout               int
	DBFileDirPathProvider dbfiledirpathprovider.DBFileDirPathProvidable
	Generator             generator.Generatable
	JrpRepository         jrprepository.JrpRepositoryInterface
	JrpWriter             jrpwriter.JrpWritable
	WNJpnRepository       wnjpnrepository.WNJpnRepositoryInterface
	Utility               utility.UtilityInterface
}

// NewGlobalOption creates a new global option.
func NewGlobalOption(fmtProxy fmtproxy.Fmt, osProxy osproxy.Os, strconvProxy strconvproxy.Strconv) *GlobalOption {
	return &GlobalOption{
		Out:    osproxy.Stdout,
		ErrOut: osproxy.Stderr,
		Args:   osproxy.Args[1:],
		Utility: utility.New(
			fmtProxy,
			osProxy,
			strconvProxy,
		),
		NewRootCommand: NewRootCommand,
	}
}

// Execute executes the jrp.
func (g *GlobalOption) Execute() int {
	rootCmd := g.NewRootCommand(g.Out, g.ErrOut, g.Args)
	if err := rootCmd.Execute(); err != nil {
		colorProxy := colorproxy.New()
		g.Utility.PrintlnWithWriter(g.ErrOut, colorProxy.RedString(err.Error()))
		return 1
	}
	return 0
}

// NewRootCommand creates a new root command.
func NewRootCommand(ow, ew ioproxy.WriterInstanceInterface, cmdArgs []string) cobraproxy.CommandInstanceInterface {
	util := utility.New(
		fmtproxy.New(),
		osproxy.New(),
		strconvproxy.New(),
	)
	g := &GlobalOption{
		Out:     ow,
		ErrOut:  ew,
		Args:    cmdArgs,
		Utility: util,
	}
	o := &rootOption{
		Out:     g.Out,
		ErrOut:  g.ErrOut,
		Args:    cmdArgs,
		Utility: util,
	}
	o.DBFileDirPathProvider = dbfiledirpathprovider.New(
		filepathproxy.New(),
		osproxy.New(),
		userproxy.New(),
	)
	o.JrpRepository = jrprepository.New(
		fmtproxy.New(),
		sortproxy.New(),
		sqlproxy.New(),
		stringsproxy.New(),
	)
	o.JrpWriter = jrpwriter.New(
		strconvproxy.New(),
		tablewriterproxy.New(),
	)
	o.WNJpnRepository = wnjpnrepository.New(
		sqlproxy.New(),
	)
	o.Generator = generator.New(
		osproxy.New(),
		randproxy.New(),
		sqlproxy.New(),
		timeproxy.New(),
		o.WNJpnRepository,
	)

	v := versionprovider.New(debugproxy.New())

	cobraProxy := cobraproxy.New()
	cmd := cobraProxy.NewCommand()

	cmd.FieldCommand.Use = constant.ROOT_USE
	cmd.FieldCommand.Version = v.GetVersion(ver)
	cmd.FieldCommand.SilenceErrors = true
	cmd.FieldCommand.SilenceUsage = true
	cmd.FieldCommand.Args = cobraProxy.MaximumNArgs(1).FieldPositionalArgs
	cmd.FieldCommand.RunE = o.rootRunE

	cmd.PersistentFlags().IntVarP(
		&o.Number,
		constant.ROOT_FLAG_NUMBER,
		constant.ROOT_FLAG_NUMBER_SHORTHAND,
		constant.ROOT_FLAG_NUMBER_DEFAULT,
		constant.ROOT_FLAG_NUMBER_DESCRIPTION,
	)
	cmd.PersistentFlags().StringVarP(
		&o.Prefix,
		constant.ROOT_FLAG_PREFIX,
		constant.ROOT_FLAG_PREFIX_SHORTHAND,
		constant.ROOT_FLAG_PREFIX_DEFAULT,
		constant.ROOT_FLAG_PREFIX_DESCRIPTION,
	)
	cmd.PersistentFlags().StringVarP(
		&o.Suffix,
		constant.ROOT_FLAG_SUFFIX,
		constant.ROOT_FLAG_SUFFIX_SHORTHAND,
		constant.ROOT_FLAG_SUFFIX_DEFAULT,
		constant.ROOT_FLAG_SUFFIX_DESCRIPTION,
	)
	cmd.PersistentFlags().BoolVarP(
		&o.DryRun,
		constant.ROOT_FLAG_DRY_RUN,
		constant.ROOT_FLAG_DRY_RUN_SHORTHAND,
		constant.ROOT_FLAG_DRY_RUN_DEFAULT,
		constant.ROOT_FLAG_DRY_RUN_DESCRIPTION,
	)
	cmd.PersistentFlags().BoolVarP(
		&o.Plain,
		constant.ROOT_FLAG_PLAIN,
		constant.ROOT_FLAG_PLAIN_SHORTHAND,
		constant.ROOT_FLAG_PLAIN_DEFAULT,
		constant.ROOT_FLAG_PLAIN_DESCRIPTION,
	)
	cmd.PersistentFlags().BoolVarP(
		&o.Interactive,
		constant.ROOT_FLAG_INTERACTIVE,
		constant.ROOT_FLAG_INTERACTIVE_SHORTHAND,
		constant.ROOT_FLAG_INTERACTIVE_DEFAULT,
		constant.ROOT_FLAG_INTERACTIVE_DESCRIPTION,
	)
	cmd.PersistentFlags().IntVarP(
		&o.Timeout,
		constant.ROOT_FLAG_TIMEOUT,
		constant.ROOT_FLAG_TIMEOUT_SHORTHAND,
		constant.ROOT_FLAG_TIMEOUT_DEFAULT,
		constant.ROOT_FLAG_TIMEOUT_DESCRIPTION,
	)

	cmd.SetOut(ow)
	cmd.SetErr(ew)
	cmd.SetHelpTemplate(constant.ROOT_HELP_TEMPLATE)

	cmd.AddCommand(
		NewDownloadCommand(g),
		NewFavoriteCommand(g),
		NewGenerateCommand(g),
		NewHistoryCommand(g),
		NewInteractiveCommand(g, keyboardproxy.New()),
		NewVersionCommand(g),
		NewCompletionCommand(g),
	)

	cmd.SetArgs(cmdArgs)

	return cmd
}

// rootRunE is the function to run root command.
func (o *rootOption) rootRunE(_ *cobra.Command, _ []string) error {
	if o.Interactive {
		// if interactive flag is set, switch to interactive command
		return switchToInteractiveCommand(
			o.Out,
			o.ErrOut,
			o.Args,
			o.Utility,
			o.Prefix,
			o.Suffix,
			o.Plain,
			o.Timeout,
		)
	}

	var word string
	var mode generator.GenerateMode
	if o.Prefix != "" && o.Suffix != "" {
		// if both prefix and suffix are provided, notify to use only one
		colorProxy := colorproxy.New()
		o.Utility.PrintlnWithWriter(o.Out, colorProxy.YellowString(constant.GENERATE_MESSAGE_NOTIFY_USE_ONLY_ONE))
		return nil
	} else if o.Prefix != "" {
		word = o.Prefix
		mode = generator.WithPrefix
	} else if o.Suffix != "" {
		word = o.Suffix
		mode = generator.WithSuffix
	}

	if len(o.Args) == 0 {
		strconvProxy := strconvproxy.New()
		// if no args are given, set the default value to the args
		o.Args = []string{strconvProxy.Itoa(constant.ROOT_FLAG_NUMBER_DEFAULT)}
	}

	// get jrp db file dir path
	wnJpnDBFileDirPath, err := o.DBFileDirPathProvider.GetWNJpnDBFileDirPath()
	if err != nil {
		return err
	}

	// get jrp db file dir path
	jrpDBFileDirPath, err := o.DBFileDirPathProvider.GetJrpDBFileDirPath()
	if err != nil {
		return err
	}

	// execute root command
	filepathProxy := filepathproxy.New()
	return o.root(
		filepathProxy.Join(wnJpnDBFileDirPath, wnjpnrepository.WNJPN_DB_FILE_NAME),
		filepathProxy.Join(jrpDBFileDirPath, jrprepository.JRP_DB_FILE_NAME),
		word,
		mode,
	)
}

// root generates jrpss and saves them.
func (o *rootOption) root(
	wnJpnDBFilePath string,
	jrpDBFilePath string,
	word string,
	mode generator.GenerateMode,
) error {
	var jrps []*model.Jrp
	var err error
	jrps, err = o.rootGenerate(wnJpnDBFilePath, word, mode)
	if err != nil {
		return err
	}
	err = o.rootSave(jrpDBFilePath, jrps)
	if err != nil {
		return err
	}
	o.writeRootResult(jrps)

	return nil
}

// rootGenerate generates jrpss.
func (o *rootOption) rootGenerate(wnJpnDBFilePath string, word string, mode generator.GenerateMode) ([]*model.Jrp, error) {
	strconvProxy := strconvproxy.New()
	res, jrps, err := o.Generator.GenerateJrp(
		wnJpnDBFilePath,
		// get the larger number between the given number flag and the largest number that can be converted from the args
		o.Utility.GetLargerNumber(
			o.Number,
			o.Utility.GetMaxConvertibleString(
				o.Args,
				strconvProxy.Itoa(constant.ROOT_FLAG_NUMBER_DEFAULT),
			),
		),
		word,
		mode,
	)
	o.writeRootGenerateResult(res)

	return jrps, err
}

// writeRootGenerateResult writes the result of the generation.
func (o *rootOption) writeRootGenerateResult(result generator.GenerateResult) {
	var out = o.Out
	var message string
	colorProxy := colorproxy.New()
	if result == generator.GeneratedFailed {
		out = o.ErrOut
		message = colorProxy.RedString(constant.ROOT_MESSAGE_GENERATE_FAILURE)
	} else if result == generator.DBFileNotFound {
		message = colorProxy.YellowString(constant.ROOT_MESSAGE_NOTIFY_DOWNLOAD_REQUIRED)
	}

	if message != "" {
		// if success, do not write any message
		o.Utility.PrintlnWithWriter(out, message)
	}
}

// rootSave saves jrpss.
func (o *rootOption) rootSave(jrpDBFilePath string, jrps []*model.Jrp) error {
	var res jrprepository.SaveStatus
	var err error
	if !o.DryRun && len(jrps) != 0 {
		// if the dry-run flag is not set and the generated phrases are not empty, save the generated phrases
		res, err = o.JrpRepository.SaveHistory(jrpDBFilePath, jrps)
	}
	o.writeRootSaveResult(res)
	return err
}

// writeRootSaveResult writes the result of saving jrps.
func (o *rootOption) writeRootSaveResult(result jrprepository.SaveStatus) {
	var out = o.Out
	var message string
	colorProxy := colorproxy.New()
	if result == jrprepository.SavedFailed {
		out = o.ErrOut
		message = colorProxy.RedString(constant.ROOT_MESSAGE_SAVED_FAILURE)
	} else if result == jrprepository.SavedNone {
		message = colorProxy.YellowString(constant.ROOT_MESSAGE_SAVED_NONE)
	} else if result == jrprepository.SavedNotAll {
		message = colorProxy.YellowString(constant.ROOT_MESSAGE_SAVED_NOT_ALL)
	}

	if message != "" {
		// if success, do not write any message
		o.Utility.PrintlnWithWriter(out, message)
	}
}

// writeRootResult writes the result of the root command.
func (o *rootOption) writeRootResult(jrps []*model.Jrp) {
	if len(jrps) != 0 {
		if o.Plain {
			for _, jrp := range jrps {
				// if plain flag is set, write only the phrase
				o.Utility.PrintlnWithWriter(o.Out, jrp.Phrase)
			}
		} else {
			// if plain flag is not set, write the favorite as a table
			o.JrpWriter.WriteGenerateResultAsTable(o.Out, jrps, !o.DryRun)
		}
	}
}
