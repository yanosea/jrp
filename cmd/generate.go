package cmd

import (
	"database/sql"
	"io"
	"strconv"

	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"

	"github.com/yanosea/jrp/constant"
	"github.com/yanosea/jrp/logic"
)

// WordNet Japanese word table structure
type Word struct {
	WordID int
	Lang   sql.NullString
	Lemma  sql.NullString
	Pron   sql.NullString
	Pos    sql.NullString
}

type generateOption struct {
	Args   []string
	Number int

	Out    io.Writer
	ErrOut io.Writer
}

func newGenerateCommand(globalOption *GlobalOption) *cobra.Command {
	o := &generateOption{}
	cmd := &cobra.Command{
		Use:   constant.GENERATE_USE,
		Short: constant.GENERATE_SHORT,
		Long:  constant.GENERATE_LONG,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Out = globalOption.Out
			o.ErrOut = globalOption.ErrOut
			o.Args = args

			return o.generate()
		},
	}

	cmd.PersistentFlags().IntVarP(&o.Number, constant.GENERATE_FLAG_NUMBER, constant.GENERATE_FLAG_NUMBER_SHORTHAND, 1, constant.GENERATE_FLAG_NUMBER_DESCRIPTION)

	cmd.SetOut(o.Out)
	cmd.SetErr(o.ErrOut)
	cmd.SetHelpTemplate(constant.GENARETE_HELP_TEMPLATE)

	return cmd
}

func (o *generateOption) generate() error {
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
