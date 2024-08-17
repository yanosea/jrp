package cmd

import (
	"io"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/constant"
	"github.com/yanosea/jrp/internal/buildinfo"
	"github.com/yanosea/jrp/internal/cmdwrapper"
	"github.com/yanosea/jrp/internal/database"
	"github.com/yanosea/jrp/internal/fs"
	"github.com/yanosea/jrp/internal/usermanager"
	"github.com/yanosea/jrp/logic"
	"github.com/yanosea/jrp/util"
)

var version = ""

type GlobalOption struct {
	Out            io.Writer
	ErrOut         io.Writer
	Args           []string
	NewRootCommand func(ow, ew io.Writer, args []string) cmdwrapper.ICommand
}

type RootOption struct {
	Out       io.Writer
	ErrOut    io.Writer
	Args      []string
	Number    int
	Generator logic.Generator
}

func NewGlobalOption(out io.Writer, errOut io.Writer, args []string) *GlobalOption {
	return &GlobalOption{
		Out:    out,
		ErrOut: errOut,
		Args:   args,
		NewRootCommand: func(ow, ew io.Writer, _ []string) cmdwrapper.ICommand {
			return NewRootCommand(ow, ew, args)
		},
	}
}

func (g *GlobalOption) Execute() int {
	rootCmd := g.NewRootCommand(g.Out, g.ErrOut, g.Args)
	if err := rootCmd.Execute(); err != nil {
		util.PrintlnWithWriter(g.ErrOut, color.RedString(err.Error()))
		return 1
	}
	return 0
}

func NewRootCommand(ow, ew io.Writer, cmdArgs []string) cmdwrapper.ICommand {
	g := &GlobalOption{
		Out:    ow,
		ErrOut: ew,
		Args:   cmdArgs,
	}
	o := &RootOption{
		Out:    g.Out,
		ErrOut: g.ErrOut,
		Args:   cmdArgs,
	}
	v := logic.NewJrpVersionGetter(buildinfo.RealBuildInfoProvider{})

	cmd := &cobra.Command{
		Use:           constant.ROOT_USE,
		Short:         constant.ROOT_SHORT,
		Long:          constant.ROOT_LONG,
		Version:       v.GetVersion(version),
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.MaximumNArgs(1),
		RunE:          o.RootRunE,
	}

	cmd.PersistentFlags().IntVarP(&o.Number, constant.ROOT_FLAG_NUMBER, constant.ROOT_FLAG_NUMBER_SHORTHAND, 1, constant.ROOT_FLAG_NUMBER_DESCRIPTION)

	cmd.SetOut(ow)
	cmd.SetErr(ew)
	cmd.SetHelpTemplate(constant.ROOT_HELP_TEMPLATE)

	cmd.AddCommand(
		newCompletionCommand(g),
		NewDownloadCommand(g),
		NewGenerateCommand(g),
		newVersionCommand(g),
	)

	cmd.SetArgs(cmdArgs)
	return cmdwrapper.NewCommandWrapper(cmd)
}

func (o *RootOption) RootRunE(_ *cobra.Command, _ []string) error {
	if len(o.Args) == 0 {
		o.Args = []string{"1"}
	}

	o.Generator = logic.NewJapaneseRandomPhraseGenerator(usermanager.OSUserProvider{}, database.SQLiteProvider{}, fs.OsFileManager{})

	return o.RootGenerate()
}

func (o *RootOption) RootGenerate() error {
	jrps, err := o.Generator.Generate(logic.DefineNumber(o.Number, o.Args[0]))
	if err != nil {
		return err
	}

	for _, jrp := range jrps {
		util.PrintlnWithWriter(o.Out, jrp)
	}

	return nil
}
