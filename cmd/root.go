package cmd

import (
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/yanosea/jrp/constant"
	"github.com/yanosea/jrp/internal/cmdwrapper"
	"github.com/yanosea/jrp/internal/db"
	"github.com/yanosea/jrp/internal/fs"
	"github.com/yanosea/jrp/internal/usermanager"
	"github.com/yanosea/jrp/logic"
	"github.com/yanosea/jrp/util"
)

var version = "develop"

type GlobalOption struct {
	Out            io.Writer
	ErrOut         io.Writer
	NewRootCommand func(ow, ew io.Writer) cmdwrapper.ICommand
}

type rootOption struct {
	Out    io.Writer
	ErrOut io.Writer
	Args   []string
	Number int
}

func NewGlobalOption(out io.Writer, errOut io.Writer) *GlobalOption {
	return &GlobalOption{
		Out:            out,
		ErrOut:         errOut,
		NewRootCommand: newRootCommand,
	}
}

func (g *GlobalOption) Execute() int {
	rootCmd := g.NewRootCommand(g.Out, g.ErrOut)
	if err := rootCmd.Execute(); err != nil {
		util.PrintlnWithWriter(g.ErrOut, color.RedString(err.Error()))
		return 1
	}
	return 0
}

func newRootCommand(ow, ew io.Writer) cmdwrapper.ICommand {
	g := &GlobalOption{
		Out:    ow,
		ErrOut: ew,
	}
	o := &rootOption{
		Out:    g.Out,
		ErrOut: g.ErrOut,
	}

	cmd := &cobra.Command{
		Use:           constant.ROOT_USE,
		Short:         constant.ROOT_SHORT,
		Long:          constant.ROOT_LONG,
		Version:       version,
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				o.Args = []string{"1"}
			} else {
				o.Args = args
			}
			return o.rootGenerate()
		},
	}

	cmd.PersistentFlags().IntVarP(&o.Number, constant.ROOT_FLAG_NUMBER, constant.ROOT_FLAG_NUMBER_SHORTHAND, 1, constant.ROOT_FLAG_NUMBER_DESCRIPTION)

	cmd.SetOut(ow)
	cmd.SetErr(ew)
	cmd.SetHelpTemplate(constant.ROOT_HELP_TEMPLATE)

	cmd.AddCommand(
		newCompletionCommand(g),
		newDownloadCommand(g),
		newGenerateCommand(g),
		newVersionCommand(g),
	)

	return cmdwrapper.NewCommandWrapper(cmd)
}

func (o *rootOption) rootGenerate() error {
	u := usermanager.OSUserProvider{}
	d := db.SQLiteProvider{}
	f := fs.OsFileManager{}

	japaneseRandomPhraseGenaretaer := logic.NewJapaneseRandomPhraseGenerator(u, d, f)
	jrp, err := japaneseRandomPhraseGenaretaer.Generate(logic.DefineNumber(o.Number, o.Args[0]))
	if err != nil {
		return err
	}

	if len(jrp) != 0 {
		fmt.Println(jrp)
	}

	return nil
}
