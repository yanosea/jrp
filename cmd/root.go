package cmd

import (
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/util"
)

var version = "develop"

const (
	root_help_template = `🎲 jrp is the CLI tool to generate random japanese phrases.

You can specify how many phrases to generate.

Usage:
  jrp [flags]
  jrp [command]

Available Commands:
	download    📥 Download Japanese Wordnet and English WordNet in an sqlite3 database from the official site.
	generate    ✨ Generate Japanese random phrases.
  completion  🔧 Generate the autocompletion script for the specified shell.
  version     🔖 Show the version of jrp.

Flags:
	-n, --number    🔢 number of phrases to generate (default 1). You can abbreviate "generate" sub command.
  -h, --help      🤝 help for jrp
  -v, --version   🔖 version for jrp

Use "jrp [command] --help" for more information about a command.
`
	root_use   = "jrp"
	root_short = "🎲 jrp is the CLI tool to generate random japanese phrases."
	root_long  = `🎲 jrp is the CLI tool to generate random japanese phrases.

You can specify how many phrases to generate.`
	root_flag_number             = "number"
	root_flag_number_shorthand   = "n"
	root_flag_number_description = "number of phrases to generate"
)

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
		Use:           root_use,
		Short:         root_short,
		Long:          root_long,
		Version:       version,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.PersistentFlags().IntVarP(&ro.Number, root_flag_number, root_flag_number_shorthand, 1, root_flag_number_description)

	glbo.Out = outWriter
	glbo.ErrOut = errWriter
	cmd.SetOut(outWriter)
	cmd.SetErr(errWriter)

	cmd.SetHelpTemplate(root_help_template)

	cmd.AddCommand(
		newCompletionCommand(glbo),
		newDownloadCommand(glbo),
		newGenerateCommand(glbo),
		newVersionCommand(glbo),
	)

	return cmd, nil
}
