package cmd

import (
	"io"
	"strconv"

	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"

	"github.com/yanosea/jrp/constant"
	"github.com/yanosea/jrp/internal/db"
	"github.com/yanosea/jrp/internal/fs"
	"github.com/yanosea/jrp/internal/rand"
	"github.com/yanosea/jrp/internal/usermanager"
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
	u := usermanager.OSUserProvider{}
	d := db.SQLiteProvider{}
	f := fs.OsFileManager{}
	r := rand.NewDefaultRandomGenerator()

	japaneseRandomPhraseGenaretaer := logic.NewJapaneseRandomPhraseGenerator(u, d, f, r)
	if err := japaneseRandomPhraseGenaretaer.Generate(defineNumber(o.Number, o.Args)); err != nil {
		return err
	}
	return nil
}

func defineNumber(num int, args []string) int {
	if len(args) == 0 {
		return num
	}

	argNum, _ := strconv.Atoi(args[0])

	if argNum > num {
		return argNum
	} else {
		return num
	}
}
