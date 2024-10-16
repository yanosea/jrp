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

// interactiveOption is the struct for interactive command.
type interactiveOption struct {
	Out                   ioproxy.WriterInstanceInterface
	ErrOut                ioproxy.WriterInstanceInterface
	Prefix                string
	Suffix                string
	Plain                 bool
	Timeout               int
	DBFileDirPathProvider dbfiledirpathprovider.DBFileDirPathProvidable
	Generator             generator.Generatable
	JrpRepository         jrprepository.JrpRepositoryInterface
	JrpWriter             jrpwriter.JrpWritable
	WNJpnRepository       wnjpnrepository.WNJpnRepositoryInterface
	KeyboardProxy         keyboardproxy.Keyboard
	Utility               utility.UtilityInterface
}

// NewInteractiveCommand creates a new interactive command.
func NewInteractiveCommand(g *GlobalOption, keyboardProxy keyboardproxy.Keyboard) *cobraproxy.CommandInstance {
	o := &interactiveOption{
		Out:     g.Out,
		ErrOut:  g.ErrOut,
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
	o.KeyboardProxy = keyboardProxy

	cobraProxy := cobraproxy.New()
	cmd := cobraProxy.NewCommand()

	cmd.FieldCommand.Use = constant.INTERACTIVE_USE
	cmd.FieldCommand.Aliases = constant.GetInteractiveAliases()
	cmd.FieldCommand.RunE = o.interactiveRunE

	cmd.PersistentFlags().StringVarP(
		&o.Prefix,
		constant.INTERACTIVE_FLAG_PREFIX,
		constant.INTERACTIVE_FLAG_PREFIX_SHORTHAND,
		constant.INTERACTIVE_FLAG_PREFIX_DEFAULT,
		constant.INTERACTIVE_FLAG_PREFIX_DESCRIPTION,
	)
	cmd.PersistentFlags().StringVarP(
		&o.Suffix,
		constant.INTERACTIVE_FLAG_SUFFIX,
		constant.INTERACTIVE_FLAG_SUFFIX_SHORTHAND,
		constant.INTERACTIVE_FLAG_SUFFIX_DEFAULT,
		constant.INTERACTIVE_FLAG_SUFFIX_DESCRIPTION,
	)
	cmd.PersistentFlags().BoolVarP(
		&o.Plain,
		constant.INTERACTIVE_FLAG_PLAIN,
		constant.INTERACTIVE_FLAG_PLAIN_SHORTHAND,
		constant.INTERACTIVE_FLAG_PLAIN_DEFAULT,
		constant.INTERACTIVE_FLAG_PLAIN_DESCRIPTION,
	)
	cmd.PersistentFlags().IntVarP(
		&o.Timeout,
		constant.INTERACTIVE_FLAG_TIMEOUT,
		constant.INTERACTIVE_FLAG_TIMEOUT_SHORTHAND,
		constant.INTERACTIVE_FLAG_TIMEOUT_DEFAULT,
		constant.INTERACTIVE_FLAG_TIMEOUT_DESCRIPTION,
	)

	cmd.SetOut(o.Out)
	cmd.SetErr(o.ErrOut)
	cmd.SetHelpTemplate(constant.INTERACTIVE_HELP_TEMPLATE)

	return cmd
}

// interactiveRunE is the function that is called when the interactive command is executed.
func (o *interactiveOption) interactiveRunE(_ *cobra.Command, _ []string) error {
	var word string
	var mode generator.GenerateMode
	if o.Prefix != "" && o.Suffix != "" {
		// if both prefix and suffix are provided, notify to use only one
		colorProxy := colorproxy.New()
		o.Utility.PrintlnWithWriter(o.Out, colorProxy.YellowString(constant.INTERACTIVE_MESSAGE_NOTIFY_USE_ONLY_ONE))
		return nil
	} else if o.Prefix != "" {
		word = o.Prefix
		mode = generator.WithPrefix
	} else if o.Suffix != "" {
		word = o.Suffix
		mode = generator.WithSuffix
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
	return o.interactive(
		filepathProxy.Join(wnJpnDBFileDirPath, wnjpnrepository.WNJPN_DB_FILE_NAME),
		filepathProxy.Join(jrpDBFileDirPath, jrprepository.JRP_DB_FILE_NAME),
		word,
		mode,
	)
}

// interactive starts generating jrp interactively.
func (o *interactiveOption) interactive(
	wnJpnDBFilePath string,
	jrpDBFilePath string,
	word string,
	mode generator.GenerateMode,
) error {
	var jrp []*model.Jrp
	var res generator.GenerateResult
	var err error
	var interactiveAnswer constant.InteractiveAnswer
	phase := 1
	// loop until the user wants to exit
	for {
		// generate jrp
		jrp, res, err = o.interactiveGenerate(wnJpnDBFilePath, word, mode)
		if err != nil || res != generator.GeneratedSuccessfully {
			// if failed to generate, exit
			return err
		}
		// leave a blank line
		o.Utility.PrintlnWithWriter(o.Out, "")
		// write phase
		o.writePhase(phase)
		// write generated jrp
		o.writeInteractiveGeneratedJrp(jrp)
		// get interactive status
		interactiveAnswer, err = o.getInteractiveInteractiveAnswer(o.Timeout)
		if err != nil {
			// if failed to get answer, exit
			return err
		}
		if interactiveAnswer == constant.InteractiveAnswerSaveAndFavoriteAndContinue ||
			interactiveAnswer == constant.InteractiveAnswerSaveAndContinue ||
			interactiveAnswer == constant.InteractiveAnswerSaveAndFavoriteAndExit ||
			interactiveAnswer == constant.InteractiveAnswerSaveAndExit {
			// save jrp
			err = o.interactiveSave(jrpDBFilePath, jrp, interactiveAnswer)
			if err != nil {
				// if failed to save, exit
				return err
			}
		}
		if interactiveAnswer == constant.InteractiveAnswerSaveAndFavoriteAndContinue ||
			interactiveAnswer == constant.InteractiveAnswerSaveAndFavoriteAndExit {
			// favorite jrp
			err = o.interactiveFavorite(jrpDBFilePath, jrp)
			if err != nil {
				// if failed to favorite, exit
				return err
			}
		}
		if interactiveAnswer == constant.InteractiveAnswerSkipAndContinue {
			// write skip message
			o.Utility.PrintlnWithWriter(o.Out, constant.INTERACTIVE_MESSAGE_SKIP)
		}
		if interactiveAnswer == constant.InteractiveAnswerSaveAndFavoriteAndExit ||
			interactiveAnswer == constant.InteractiveAnswerSaveAndExit ||
			interactiveAnswer == constant.InteractiveAnswerSkipAndExit {
			// if the user wants to exit, break the loop
			o.Utility.PrintlnWithWriter(o.Out, constant.INTERACTIVE_MESSAGE_EXIT)
			break
		}
		// leave a blank line
		o.Utility.PrintlnWithWriter(o.Out, "")
		// increment phase
		phase++
	}

	return nil
}

// interactiveGenerate generates jrp
func (o *interactiveOption) interactiveGenerate(wnJpnDBFilePath string, word string, mode generator.GenerateMode) ([]*model.Jrp, generator.GenerateResult, error) {
	res, jrps, err := o.Generator.GenerateJrp(wnJpnDBFilePath, 1, word, mode)
	o.writeInteractiveGenerateResult(res)

	return jrps, res, err
}

// writeInteractiveGenerateResult writes the result of generating jrp.
func (o *interactiveOption) writeInteractiveGenerateResult(result generator.GenerateResult) {
	var out = o.Out
	var message string
	colorProxy := colorproxy.New()
	if result == generator.GeneratedFailed {
		out = o.ErrOut
		message = colorProxy.RedString(constant.INTERACTIVE_MESSAGE_GENERATE_FAILURE)
	} else if result == generator.DBFileNotFound {
		message = colorProxy.YellowString(constant.INTERACTIVE_MESSAGE_NOTIFY_DOWNLOAD_REQUIRED)
	}

	if message != "" {
		// if success, do not write any message
		o.Utility.PrintlnWithWriter(out, message)
	}
}

// writeInteractiveGeneratedJrp writes generated jrp.
func (o *interactiveOption) writeInteractiveGeneratedJrp(jrp []*model.Jrp) {
	if len(jrp) != 0 {
		if o.Plain {
			for _, jrp := range jrp {
				// if plain flag is set, write only the phrase
				o.Utility.PrintlnWithWriter(o.Out, jrp.Phrase)
				// leave a blank line
				o.Utility.PrintlnWithWriter(o.Out, "")
			}
		} else {
			// if plain flag is not set, write the result as table
			o.JrpWriter.WriteInteractiveResultAsTable(o.Out, jrp)
		}
	}
}

// interactiveSave saves jrp.
func (o *interactiveOption) interactiveSave(jrpDBFilePath string, jrps []*model.Jrp, interactiveAnswer constant.InteractiveAnswer) error {
	res, err := o.JrpRepository.SaveHistory(jrpDBFilePath, jrps)
	o.writeInteractiveSaveResult(res, interactiveAnswer)

	return err
}

// writeInteractiveSaveResult writes the result of saving jrp.
func (o *interactiveOption) writeInteractiveSaveResult(result jrprepository.SaveStatus, interactiveAnswer constant.InteractiveAnswer) {
	var out = o.Out
	var message string
	colorProxy := colorproxy.New()
	if result == jrprepository.SavedFailed {
		out = o.ErrOut
		message = colorProxy.RedString(constant.INTERACTIVE_MESSAGE_SAVED_FAILURE)
	} else if result == jrprepository.SavedNone {
		message = colorProxy.YellowString(constant.INTERACTIVE_MESSAGE_SAVED_NONE)
	} else if result == jrprepository.SavedNotAll {
		message = colorProxy.YellowString(constant.INTERACTIVE_MESSAGE_SAVED_NOT_ALL)
	} else if interactiveAnswer == constant.InteractiveAnswerSaveAndFavoriteAndContinue {
		message = ""
	} else if interactiveAnswer == constant.InteractiveAnswerSaveAndFavoriteAndExit {
		message = ""
	} else {
		message = colorProxy.GreenString(constant.INTERACTIVE_MESSAGE_SAVED_SUCCESSFULLY)
	}

	if message != "" {
		// if success and the answer watns to favorite, do not write any message
		o.Utility.PrintlnWithWriter(out, message)
	}
}

// interactiveFavorite favorites jrp.
func (o *interactiveOption) interactiveFavorite(jrpDBFilePath string, jrps []*model.Jrp) error {
	res, err := o.JrpRepository.AddFavoriteByIDs(jrpDBFilePath, []int{jrps[0].ID})
	o.writeInteractiveFavoriteResult(res)

	return err
}

// writeInteractiveFavoriteResult writes the result of favoriting jrp.
func (o *interactiveOption) writeInteractiveFavoriteResult(result jrprepository.AddStatus) {
	var out = o.Out
	var message string
	colorProxy := colorproxy.New()
	if result == jrprepository.AddedFailed {
		out = o.ErrOut
		message = colorProxy.RedString(constant.INTERACTIVE_MESSAGE_FAVORITED_FAILURE)
	} else if result == jrprepository.AddedNone {
		message = colorProxy.YellowString(constant.INTERACTIVE_MESSAGE_FAVORITED_NONE)
	} else if result == jrprepository.AddedNotAll {
		message = colorProxy.YellowString(constant.INTERACTIVE_MESSAGE_FAVORITED_NOT_ALL)
	} else {
		message = colorProxy.GreenString(constant.INTERACTIVE_MESSAGE_FAVORITED_SUCCESSFULLY)
	}
	o.Utility.PrintlnWithWriter(out, message)
}

// getInteractiveInteractiveAnswer gets the interactive answer.
func (o *interactiveOption) getInteractiveInteractiveAnswer(timeoutSec int) (constant.InteractiveAnswer, error) {
	// write prompt
	o.Utility.PrintlnWithWriter(o.Out, constant.INTERACTIVE_PROMPT_LABEL)
	// open keyboard
	if err := o.KeyboardProxy.Open(); err != nil {
		return constant.InteractiveAnswerSkipAndExit, err
	}
	defer o.KeyboardProxy.Close()
	// get answer
	answer, _, err := o.KeyboardProxy.GetKey(timeoutSec)
	if err != nil {
		return constant.InteractiveAnswerSkipAndExit, err
	}
	var interactiveAnswer constant.InteractiveAnswer
	if string(answer) == "u" || string(answer) == "U" {
		interactiveAnswer = constant.InteractiveAnswerSaveAndFavoriteAndContinue
	} else if string(answer) == "i" || string(answer) == "I" {
		interactiveAnswer = constant.InteractiveAnswerSaveAndFavoriteAndExit
	} else if string(answer) == "j" || string(answer) == "J" {
		interactiveAnswer = constant.InteractiveAnswerSaveAndContinue
	} else if string(answer) == "k" || string(answer) == "K" {
		interactiveAnswer = constant.InteractiveAnswerSaveAndExit
	} else if string(answer) == "m" || string(answer) == "M" {
		interactiveAnswer = constant.InteractiveAnswerSkipAndContinue
	} else {
		interactiveAnswer = constant.InteractiveAnswerSkipAndExit
	}

	return interactiveAnswer, nil
}

// writePhase writes the phase.
func (o *interactiveOption) writePhase(phase int) {
	colorProxy := colorproxy.New()
	strconvProxy := strconvproxy.New()
	p := constant.INTERACTIVE_MESSAGE_PHASE + strconvProxy.Itoa(phase)
	o.Utility.PrintlnWithWriter(o.Out, colorProxy.BlueString(p))
	// leave a blank line
	o.Utility.PrintlnWithWriter(o.Out, "")
}

// switchToInteractiveCommand switches to interactive command.
func switchToInteractiveCommand(
	out ioproxy.WriterInstanceInterface,
	errOut ioproxy.WriterInstanceInterface,
	args []string,
	utility utility.UtilityInterface,
	prefix string,
	suffix string,
	plain bool,
	timeout int,
	keyboardProxy keyboardproxy.Keyboard,
) error {
	// regenerate GlobalOption
	g := &GlobalOption{
		Out:     out,
		ErrOut:  errOut,
		Args:    args,
		Utility: utility,
	}
	interactiveCmd := NewInteractiveCommand(g, keyboardProxy)
	strconvProxy := strconvproxy.New()
	// consolidate flags
	flags := map[string]string{
		constant.INTERACTIVE_FLAG_PREFIX:  prefix,
		constant.INTERACTIVE_FLAG_SUFFIX:  suffix,
		constant.INTERACTIVE_FLAG_PLAIN:   strconvProxy.FormatBool(plain),
		constant.INTERACTIVE_FLAG_TIMEOUT: strconvProxy.Itoa(timeout),
	}
	// set each flag
	if err := setFlagsFunc(interactiveCmd, flags); err != nil {
		return err
	}
	// run interactive command
	if err := runSwitchedInteractiveCommandFunc(interactiveCmd); err != nil {
		return err
	}

	return nil
}

// setFlagsFunc is the function that sets flags to the command.
var setFlagsFunc = setFlags

// setFlags sets flags to the command.
func setFlags(cmd *cobraproxy.CommandInstance, flags map[string]string) error {
	for name, value := range flags {
		if err := cmd.PersistentFlags().Set(name, value); err != nil {
			return err
		}
	}

	return nil
}

// runSwitchedInteractiveCommandFunc is the function that runs switched interactive command.
var runSwitchedInteractiveCommandFunc = runSwitchedInteractiveCommand

// runSwitchedInteractiveCommand runs switched interactive command.
func runSwitchedInteractiveCommand(switchedInteractiveCommand *cobraproxy.CommandInstance) error {
	return switchedInteractiveCommand.RunE(nil, nil)
}
