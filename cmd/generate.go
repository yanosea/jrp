package cmd

import (
	"io"

	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"

	"github.com/yanosea/jrp/constant"
	"github.com/yanosea/jrp/logic"
)

type generateOption struct {
	Out    io.Writer
	ErrOut io.Writer
	Args   []string
	Number int
}

func newGenerateCommand(g *GlobalOption) *cobra.Command {
	o := &generateOption{
		Out:    g.Out,
		ErrOut: g.ErrOut,
	}

	cmd := &cobra.Command{
		Use:     constant.GENERATE_USE,
		Aliases: constant.GetGenerateAliases(),
		Short:   constant.GENERATE_SHORT,
		Long:    constant.GENERATE_LONG,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
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
	env := logic.OsEnv{}
	user := logic.OsUser{}

	japaneseRandomPhraseGenaretaer := logic.NewJapaneseRandomPhraseGenerator(o.Number, o.Args, env, user)
	num := japaneseRandomPhraseGenaretaer.DefineNumber()

	if err := japaneseRandomPhraseGenaretaer.Generate(num); err != nil {
		return err
	}
	return nil
}
