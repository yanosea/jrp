package logic

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/fatih/color"
	_ "modernc.org/sqlite"

	"github.com/yanosea/jrp/constant"
	"github.com/yanosea/jrp/model"
)

type Genarater interface {
	DetermineNumber(num int, args []string) int
	Generate() error
}

type JapaneseRandomPhraseGenaretaer struct {
	Num  int
	Args []string
	Env  Env
	User User
}

func NewJapaneseRandomPhraseGenerator(num int, args []string, env Env, user User) *JapaneseRandomPhraseGenaretaer {
	return &JapaneseRandomPhraseGenaretaer{
		Num:  num,
		Args: args,
		Env:  env,
		User: user,
	}
}

// WordNet Japanese word table structure
func (j JapaneseRandomPhraseGenaretaer) DefineNumber() int {
	if len(j.Args) == 0 {
		return j.Num
	}

	argNum, _ := strconv.Atoi(j.Args[0])

	if argNum > j.Num {
		return argNum
	} else {
		return j.Num
	}
}

func (j JapaneseRandomPhraseGenaretaer) Generate(num int) error {
	// create DBFileDirPathGetter instance
	dbFileDirPathGetter := NewDBFileDirPathGetter(j.Env, j.User)
	// get the directory of wnjpn.db from environment
	dbFileDirPath, err := dbFileDirPathGetter.GetFileDirPath()
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

	allAVNWords := make([]model.Word, 0)
	for rows.Next() {
		var word model.Word
		err = rows.Scan(&word.Lemma, &word.Pos)
		if err != nil {
			return err
		}
		allAVNWords = append(allAVNWords, word)
	}

	// separate the words into adjectives and verbs, and nouns
	var allAVWords []model.Word
	var allNWords []model.Word

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
	for i := 0; i < num; i++ {
		randomIndexA := rand.Intn(len(allAVWords))
		randomIndexB := rand.Intn(len(allNWords))
		randomWord := allAVWords[randomIndexA]
		randomWord2 := allNWords[randomIndexB]
		fmt.Println(randomWord.Lemma.String + randomWord2.Lemma.String)
	}

	return nil
}
