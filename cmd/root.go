package cmd

import (
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/constant"
	"github.com/yanosea/jrp/util"
)

var version = "develop"

type GlobalOption struct {
	Out    io.Writer
	ErrOut io.Writer
}

type rootOption struct {
	Number int

	Out    io.Writer
	ErrOut io.Writer
}

func Execute() int {
	o := os.Stdout
	e := os.Stderr

	rootCmd, err := NewRootCommand(o, e)
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

func NewRootCommand(outWriter, errWriter io.Writer) (*cobra.Command, error) {
	glbo := &GlobalOption{}
	ro := &rootOption{}

	cmd := &cobra.Command{
		Use:           constant.ROOT_USE,
		Short:         constant.ROOT_SHORT,
		Long:          constant.ROOT_LONG,
		Version:       version,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.PersistentFlags().IntVarP(&ro.Number, constant.ROOT_FLAG_NUMBER, constant.ROOT_FLAG_NUMBER_SHORTHAND, 1, constant.ROOT_FLAG_NUMBER_DESCRIPTION)

	glbo.Out = outWriter
	glbo.ErrOut = errWriter
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
