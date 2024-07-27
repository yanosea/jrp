package cmd

import (
	"database/sql"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"

	"github.com/yanosea/jrp/constant"
	"github.com/yanosea/jrp/util"
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
		RunE: func(cmd *cobra.Command, args []string) error {

			o.Out = globalOption.Out
			o.ErrOut = globalOption.ErrOut

			return o.generate()
		},
	}

	cmd.PersistentFlags().IntVarP(&o.Number, constant.GENERATE_FLAG_NUMBER, constant.GENERATE_FLAG_NUMBER_SHORTHAND, 1, constant.GENERATE_FLAG_NUMBER_DESCRIPTION)

	o.Out = globalOption.Out
	o.ErrOut = globalOption.ErrOut
	cmd.SetOut(o.Out)
	cmd.SetErr(o.ErrOut)

	cmd.SetHelpTemplate(constant.GENARETE_HELP_TEMPLATE)

	return cmd
}

func (o *generateOption) generate() error {
	// get the directory of wnjpn.db from environment
	dbFileDirPath, err := util.GetDBFileDirPath()
	if err != nil {
		return err
	}

	// end the program if the database file doesn't exist
	dbFilePath := filepath.Join(dbFileDirPath, constant.WNJPN_DB_FILE_NAME)
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		fmt.Println(color.YellowString(constant.GENERATE_MESSAGE_NOTIFY_DOWNLOAD_REQUIRED))
		return nil
	}

	// connect to the database
	db, err := sql.Open("sqlite", "file:"+dbFilePath)
	if err != nil {
		return err
	}
	defer db.Close()

	// get all rows from the word table where the lang is Japanese and the pos is adjective, verb, or noun
	rows, err := db.Query(constant.GENERATE_SQL_GET_ALL_JAPANESE_AVN_WORDS)
	if err != nil {
		return err
	}
	defer rows.Close()

	allAVNWords := make([]Word, 0)
	for rows.Next() {
		var word Word
		err = rows.Scan(&word.Lemma, &word.Pos)
		if err != nil {
			return err
		}
		allAVNWords = append(allAVNWords, word)
	}

	// separate the words into adjectives and verbs, and nouns
	var allAVWords []Word
	var allNWords []Word

	for _, word := range allAVNWords {
		if word.Pos.Valid && word.Pos.String == "n" {
			allNWords = append(allNWords, word)
		} else {
			allAVWords = append(allAVWords, word)
		}
	}

	// generate random number
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// generate the words
	for i := 0; i < o.Number; i++ {
		randomIndexA := rand.Intn(len(allAVWords))
		randomIndexB := rand.Intn(len(allNWords))
		randomWord := allAVWords[randomIndexA]
		randomWord2 := allNWords[randomIndexB]
		fmt.Println(randomWord.Lemma.String + randomWord2.Lemma.String)
	}

	return nil
}
