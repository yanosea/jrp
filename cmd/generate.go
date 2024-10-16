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
	"github.com/yanosea/jrp/app/proxy/cobra"
	"github.com/yanosea/jrp/app/proxy/color"
	"github.com/yanosea/jrp/app/proxy/filepath"
	"github.com/yanosea/jrp/app/proxy/fmt"
	"github.com/yanosea/jrp/app/proxy/io"
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

// generateOption is the struct for generate command.
type generateOption struct {
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

// NewGenerateCommand creates a new generate command.
func NewGenerateCommand(g *GlobalOption) *cobraproxy.CommandInstance {
	o := &generateOption{
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

	cobraProxy := cobraproxy.New()
	cmd := cobraProxy.NewCommand()

	cmd.FieldCommand.Use = constant.GENERATE_USE
	cmd.FieldCommand.Aliases = constant.GetGenerateAliases()
	cmd.FieldCommand.Args = cobra.MaximumNArgs(1)
	cmd.FieldCommand.RunE = o.generateRunE

	cmd.PersistentFlags().IntVarP(
		&o.Number,
		constant.GENERATE_FLAG_NUMBER,
		constant.GENERATE_FLAG_NUMBER_SHORTHAND,
		constant.GENERATE_FLAG_NUMBER_DEFAULT,
		constant.GENERATE_FLAG_NUMBER_DESCRIPTION,
	)
	cmd.PersistentFlags().StringVarP(
		&o.Prefix,
		constant.GENERATE_FLAG_PREFIX,
		constant.GENERATE_FLAG_PREFIX_SHORTHAND,
		constant.GENERATE_FLAG_PREFIX_DEFAULT,
		constant.GENERATE_FLAG_PREFIX_DESCRIPTION,
	)
	cmd.PersistentFlags().StringVarP(
		&o.Suffix,
		constant.GENERATE_FLAG_SUFFIX,
		constant.GENERATE_FLAG_SUFFIX_SHORTHAND,
		constant.GENERATE_FLAG_SUFFIX_DEFAULT,
		constant.GENERATE_FLAG_SUFFIX_DESCRIPTION,
	)
	cmd.PersistentFlags().BoolVarP(
		&o.DryRun,
		constant.GENERATE_FLAG_DRY_RUN,
		constant.GENERATE_FLAG_DRY_RUN_SHORTHAND,
		constant.GENERATE_FLAG_DRY_RUN_DEFAULT,
		constant.GENERATE_FLAG_DRY_RUN_DESCRIPTION,
	)
	cmd.PersistentFlags().BoolVarP(
		&o.Plain,
		constant.GENERATE_FLAG_PLAIN,
		constant.GENERATE_FLAG_PLAIN_SHORTHAND,
		constant.GENERATE_FLAG_PLAIN_DEFAULT,
		constant.GENERATE_FLAG_PLAIN_DESCRIPTION,
	)
	cmd.PersistentFlags().BoolVarP(
		&o.Interactive,
		constant.GENERATE_FLAG_INTERACTIVE,
		constant.GENERATE_FLAG_INTERACTIVE_SHORTHAND,
		constant.GENERATE_FLAG_INTERACTIVE_DEFAULT,
		constant.GENERATE_FLAG_INTERACTIVE_DESCRIPTION,
	)
	cmd.PersistentFlags().IntVarP(
		&o.Timeout,
		constant.GENERATE_FLAG_TIMEOUT,
		constant.GENERATE_FLAG_TIMEOUT_SHORTHAND,
		constant.GENERATE_FLAG_TIMEOUT_DEFAULT,
		constant.GENERATE_FLAG_TIMEOUT_DESCRIPTION,
	)

	cmd.SetOut(o.Out)
	cmd.SetErr(o.ErrOut)
	cmd.SetHelpTemplate(constant.GENARETE_HELP_TEMPLATE)

	cmd.SetArgs(o.Args)

	return cmd
}

// generateRunE is the function that is called when the generate command is executed.
func (o *generateOption) generateRunE(_ *cobra.Command, _ []string) error {
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

	if len(o.Args) <= 1 {
		strconvProxy := strconvproxy.New()
		// if no args are given, set the default value to the args
		o.Args = []string{constant.GENERATE_USE, strconvProxy.Itoa(constant.GENERATE_FLAG_NUMBER_DEFAULT)}
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

	filepathProxy := filepathproxy.New()
	return o.generate(
		filepathProxy.Join(wnJpnDBFileDirPath, wnjpnrepository.WNJPN_DB_FILE_NAME),
		filepathProxy.Join(jrpDBFileDirPath, jrprepository.JRP_DB_FILE_NAME),
		word,
		mode,
	)
}

// generate generates jrpss and saves them.
func (o *generateOption) generate(
	wnJpnDBFilePath string,
	jrpDBFilePath string,
	word string,
	mode generator.GenerateMode,
) error {
	var jrps []*model.Jrp
	var err error
	jrps, err = o.generateGenerate(wnJpnDBFilePath, word, mode)
	if err != nil {
		return err
	}
	err = o.generateSave(jrpDBFilePath, jrps)
	if err != nil {
		return err
	}
	o.writeGenerateResult(jrps)

	return nil
}

// generateGenerate generates jrps.
func (o *generateOption) generateGenerate(wnJpnDBFilePath string, word string, mode generator.GenerateMode) ([]*model.Jrp, error) {
	strconvProxy := strconvproxy.New()
	res, jrps, err := o.Generator.GenerateJrp(
		wnJpnDBFilePath,
		// get the larger number between the given number flag and the largest number that can be converted from the args
		o.Utility.GetLargerNumber(
			o.Number,
			o.Utility.GetMaxConvertibleString(
				o.Args,
				strconvProxy.Itoa(constant.GENERATE_FLAG_NUMBER_DEFAULT),
			),
		),
		word,
		mode,
	)
	o.writeGenerateGenerateResult(res)

	return jrps, err
}

// writeGenerateGenerateResult writes the result of generating jrps.
func (o *generateOption) writeGenerateGenerateResult(result generator.GenerateResult) {
	var out = o.Out
	var message string
	colorProxy := colorproxy.New()
	if result == generator.GeneratedFailed {
		out = o.ErrOut
		message = colorProxy.RedString(constant.GENERATE_MESSAGE_GENERATE_FAILURE)
	} else if result == generator.DBFileNotFound {
		message = colorProxy.YellowString(constant.GENERATE_MESSAGE_NOTIFY_DOWNLOAD_REQUIRED)
	}

	if message != "" {
		// if success, do not write any message
		o.Utility.PrintlnWithWriter(out, message)
	}
}

// generateSave saves jrps.
func (o *generateOption) generateSave(jrpDBFilePath string, jrps []*model.Jrp) error {
	var res jrprepository.SaveStatus
	var err error
	if !o.DryRun && len(jrps) != 0 {
		// if the dry-run flag is not set and the generated phrases are not empty, save the generated phrases
		res, err = o.JrpRepository.SaveHistory(jrpDBFilePath, jrps)
	}
	o.writeGenerateSaveResult(res)

	return err
}

// writeGenerateSaveResult writes the result of saving jrps.
func (o *generateOption) writeGenerateSaveResult(result jrprepository.SaveStatus) {
	var out = o.Out
	var message string
	colorProxy := colorproxy.New()
	if result == jrprepository.SavedFailed {
		out = o.ErrOut
		message = colorProxy.RedString(constant.GENERATE_MESSAGE_SAVED_FAILURE)
	} else if result == jrprepository.SavedNone {
		message = colorProxy.YellowString(constant.GENERATE_MESSAGE_SAVED_NONE)
	} else if result == jrprepository.SavedNotAll {
		message = colorProxy.YellowString(constant.GENERATE_MESSAGE_SAVED_NOT_ALL)
	}

	if message != "" {
		// if success, do not write any message
		o.Utility.PrintlnWithWriter(out, message)
	}
}

// writeGenerateResult writes the result of generate command.
func (o *generateOption) writeGenerateResult(jrps []*model.Jrp) {
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
