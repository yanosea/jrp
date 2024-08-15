package cmd

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/constant"
	"github.com/yanosea/jrp/internal/database"
	"github.com/yanosea/jrp/internal/fs"
	"github.com/yanosea/jrp/internal/usermanager"
	"github.com/yanosea/jrp/logic"
	"github.com/yanosea/jrp/util"
)

type GenerateOption struct {
	Out       io.Writer
	ErrOut    io.Writer
	Args      []string
	Number    int
	Generator logic.Generator
}

func NewGenerateCommand(g *GlobalOption) *cobra.Command {
	o := &GenerateOption{
		Out:    g.Out,
		ErrOut: g.ErrOut,
		Args:   g.Args,
	}

	cmd := &cobra.Command{
		Use:     constant.GENERATE_USE,
		Aliases: constant.GetGenerateAliases(),
		Short:   constant.GENERATE_SHORT,
		Long:    constant.GENERATE_LONG,
		Args:    cobra.MaximumNArgs(1),
		RunE:    o.GenerateRunE,
	}

	cmd.PersistentFlags().IntVarP(&o.Number, constant.GENERATE_FLAG_NUMBER, constant.GENERATE_FLAG_NUMBER_SHORTHAND, 1, constant.GENERATE_FLAG_NUMBER_DESCRIPTION)

	cmd.SetOut(o.Out)
	cmd.SetErr(o.ErrOut)
	cmd.SetHelpTemplate(constant.GENARETE_HELP_TEMPLATE)

	cmd.SetArgs(o.Args)
	return cmd
}

func (o *GenerateOption) GenerateRunE(_ *cobra.Command, _ []string) error {
	if len(o.Args) == 0 {
		o.Args = []string{"1"}
	}

	o.Generator = logic.NewJapaneseRandomPhraseGenerator(usermanager.OSUserProvider{}, database.SQLiteProvider{}, fs.OsFileManager{})

	return o.Generate()
}

func (o *GenerateOption) Generate() error {
	jrps, err := o.Generator.Generate(logic.DefineNumber(o.Number, o.Args[1]))
	if err != nil {
		return err
	}

	for _, jrp := range jrps {
		util.PrintlnWithWriter(o.Out, jrp)
	}

	return nil
}
