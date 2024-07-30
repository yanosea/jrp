package logic

import (
	"fmt"
	"path/filepath"

	"github.com/fatih/color"

	"github.com/yanosea/jrp/constant"
	"github.com/yanosea/jrp/internal/db"
	"github.com/yanosea/jrp/internal/fs"
	"github.com/yanosea/jrp/internal/rand"
	"github.com/yanosea/jrp/internal/usermanager"
	"github.com/yanosea/jrp/model"
)

type Genarater interface {
	Generate(num int) error
}

type JapaneseRandomPhraseGenaretaer struct {
	User            usermanager.UserProvider
	DbProvider      db.DatabaseProvider
	FileSystem      fs.FileManager
	RandomGenerator rand.RandomGenerator
}

func NewJapaneseRandomPhraseGenerator(u usermanager.UserProvider, d db.DatabaseProvider, f fs.FileManager, r rand.RandomGenerator) *JapaneseRandomPhraseGenaretaer {
	return &JapaneseRandomPhraseGenaretaer{
		User:            u,
		DbProvider:      d,
		FileSystem:      f,
		RandomGenerator: r,
	}
}

func (j JapaneseRandomPhraseGenaretaer) Generate(num int) error {
	dbFileDirPathGetter := NewDBFileDirPathGetter(j.User)

	dbFileDirPath, err := dbFileDirPathGetter.GetFileDirPath()
	if err != nil {
		return err
	}

	dbFilePath := filepath.Join(dbFileDirPath, constant.WNJPN_DB_FILE_NAME)
	if !j.FileSystem.Exists(dbFilePath) {
		fmt.Println(color.YellowString(constant.GENERATE_MESSAGE_NOTIFY_DOWNLOAD_REQUIRED))
		return nil
	}

	db, err := j.DbProvider.Connect(dbFilePath)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := j.DbProvider.Query(db, constant.GENERATE_SQL_GET_ALL_JAPANESE_AVN_WORDS)
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

	var allAVWords []model.Word
	var allNWords []model.Word

	for _, word := range allAVNWords {
		if word.Pos.Valid && word.Pos.String == "n" {
			allNWords = append(allNWords, word)
		} else {
			allAVWords = append(allAVWords, word)
		}
	}

	for i := 0; i < num; i++ {
		randomIndexA := j.RandomGenerator.Intn(len(allAVWords))
		randomIndexB := j.RandomGenerator.Intn(len(allNWords))
		randomWord := allAVWords[randomIndexA]
		randomWord2 := allNWords[randomIndexB]
		fmt.Println(randomWord.Lemma.String + randomWord2.Lemma.String)
	}

	return nil
}
