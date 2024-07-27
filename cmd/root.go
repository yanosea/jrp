package cmd

import (
	"io"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/constant"
	"github.com/yanosea/jrp/logic"
	"github.com/yanosea/jrp/util"
)

var version = "develop"

type GlobalOption struct {
	Out    io.Writer
	ErrOut io.Writer
}

type rootOption struct {
	Args   []string
	Number int

	Out    io.Writer
	ErrOut io.Writer
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

func newRootCommand(outWriter, errWriter io.Writer) (*cobra.Command, error) {
	glbo := &GlobalOption{
		Out:    outWriter,
		ErrOut: errWriter,
	}
	ro := &rootOption{}

	cmd := &cobra.Command{
		Use:           constant.ROOT_USE,
		Short:         constant.ROOT_SHORT,
		Long:          constant.ROOT_LONG,
		Version:       version,
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ro.Out = glbo.Out
			ro.ErrOut = glbo.ErrOut
			ro.Args = args

			return ro.rootGenerate()
		},
	}

	cmd.PersistentFlags().IntVarP(&ro.Number, constant.ROOT_FLAG_NUMBER, constant.ROOT_FLAG_NUMBER_SHORTHAND, 1, constant.ROOT_FLAG_NUMBER_DESCRIPTION)

	cmd.SetOut(outWriter)
	cmd.SetErr(errWriter)
	cmd.SetHelpTemplate(constant.ROOT_HELP_TEMPLATE)

	cmd.AddCommand(
		newCompletionCommand(glbo),
		newDownloadCommand(glbo),
		newGenerateCommand(glbo),
		newVersionCommand(glbo),
	)

	return cmd, nil
}

func (o *rootOption) rootGenerate() error {
	if len(o.Args) == 0 {
		if err := logic.Generate(o.Number); err != nil {
			return err
		}
		return nil
	}

	argNum, err := strconv.Atoi(o.Args[0])
	if err != nil || argNum <= 0 {
		if err := logic.Generate(o.Number); err != nil {
			return err
		}
		return nil
	}

	if o.Number == 1 {
		if err := logic.Generate(argNum); err != nil {
			return err
		}
	} else {
		if err := logic.Generate(o.Number); err != nil {
			return err
		}
	}

	return nil
}
