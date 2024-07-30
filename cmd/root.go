package cmd

import (
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/constant"
	"github.com/yanosea/jrp/internal/env"
	"github.com/yanosea/jrp/internal/usermanager"
	"github.com/yanosea/jrp/logic"
	"github.com/yanosea/jrp/util"
)

var version = "develop"

type GlobalOption struct {
	Out    io.Writer
	ErrOut io.Writer
}

type rootOption struct {
	Out    io.Writer
	ErrOut io.Writer
	Args   []string
	Number int
}

func Execute() int {
	o := os.Stdout
	e := os.Stderr

	rootCmd, err := newRootCommand(o, e)
	if err != nil {
		util.PrintlnWithWriter(e, color.RedString(err.Error()))
		return 1
	}

	if err = rootCmd.Execute(); err != nil {
		util.PrintlnWithWriter(e, color.RedString(err.Error()))
		return 1
	}

	return 0
}

func newRootCommand(ow, ew io.Writer) (*cobra.Command, error) {
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
			o.Args = args
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

	return cmd, nil
}

func (o *rootOption) rootGenerate() error {
	e := env.OsEnvironment{}
	u := usermanager.OSUserProvider{}

	japaneseRandomPhraseGenaretaer := logic.NewJapaneseRandomPhraseGenerator(e, u)
	if err := japaneseRandomPhraseGenaretaer.Generate(defineNumber(o.Number, o.Args)); err != nil {
		return err
	}
	return nil
}
