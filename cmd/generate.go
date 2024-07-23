package cmd

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/yanosea/jrp/util"
)

const (
	genarete_help_template = `✨ Generate Japanese random phrases.

You can generate Japanese random phrase.
You can specify the number of phrases to generate by the flag "-n" or "--number".

Usage:
  jrp generate [flags]

Flags:
	-n, --number    🔢 number of phrases to generate (default 1). You can abbreviate "generate" sub command.
  -h, --help      🤝 help for generate
`
	generate_use   = "generate"
	generate_short = "✨ Generate Japanese random phrases."
	generate_long  = `✨ Generate Japanese random phrases.

You can generate Japanese random phrase.
You can specify the number of phrases to generate by the flag "-n" or "--number".
`
	generate_flag_number             = "number"
	generate_flag_number_shorthand   = "n"
	generate_flag_number_description = "number of phrases to generate"
)

// WordNet Japanese word table structure
type Word struct {
	WordID int
	Lang   string
	Lemma  string
	Pron   string
	Pos    string
}

// TableName returns the true table name
func (w *Word) TableName() string {
	return "word"
}

type generateOption struct {
	Number int

	Out    io.Writer
	ErrOut io.Writer
}

func newGenerateCommand(globalOption *GlobalOption) *cobra.Command {
	o := &generateOption{}
	cmd := &cobra.Command{
		Use:   generate_use,
		Short: generate_short,
		Long:  generate_long,
		RunE: func(cmd *cobra.Command, args []string) error {

			o.Out = globalOption.Out
			o.ErrOut = globalOption.ErrOut

			return o.generate()
		},
	}

	cmd.PersistentFlags().IntVarP(&o.Number, generate_flag_number, generate_flag_number_shorthand, 1, generate_flag_number_description)

	o.Out = globalOption.Out
	o.ErrOut = globalOption.ErrOut
	cmd.SetOut(o.Out)
	cmd.SetErr(o.ErrOut)

	cmd.SetHelpTemplate(genarete_help_template)

	return cmd
}

func (o *generateOption) generate() error {
	// get the directory of wnjpn.db from environment
	var dbFileDirPath = util.GetDBFileDirPath()

	// connect to the database
	db, err := gorm.Open(sqlite.Open(dbFileDirPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&Word{})
	if err != nil {
		panic(err)
	}

	// get all Japanese words
	var jpns []Word
	db.Where(&Word{Lang: "jpn"}).Find(&jpns)

	// filter the words by part of speech
	var wordsA []Word
	var wordsB []Word
	for _, jpn := range jpns {
		if jpn.Pos == "a" || jpn.Pos == "v" {
			wordsA = append(wordsA, jpn)
		}
		if jpn.Pos == "n" {
			wordsB = append(wordsB, jpn)
		}
	}

	// generate random number
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// generate the words
	for i := 0; i < o.Number; i++ {
		randomIndexA := rand.Intn(len(wordsA))
		randomIndexB := rand.Intn(len(wordsB))
		randomWord := wordsA[randomIndexA].Lemma
		randomWord2 := wordsB[randomIndexB].Lemma
		fmt.Println(randomWord + randomWord2)
	}

	return nil
}
