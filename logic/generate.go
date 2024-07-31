package logic

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/fatih/color"

	"github.com/yanosea/jrp/constant"
	"github.com/yanosea/jrp/internal/db"
	"github.com/yanosea/jrp/internal/fs"
	"github.com/yanosea/jrp/internal/rand"
	"github.com/yanosea/jrp/internal/usermanager"
	"github.com/yanosea/jrp/model"
)

type Generator interface {
	Generate(num int) error
}

type JapaneseRandomPhraseGenerator struct {
	User            usermanager.UserProvider
	DbProvider      db.DatabaseProvider
	FileSystem      fs.FileManager
	RandomGenerator rand.RandomGenerator
}

func NewJapaneseRandomPhraseGenerator(u usermanager.UserProvider, d db.DatabaseProvider, f fs.FileManager, r rand.RandomGenerator) *JapaneseRandomPhraseGenerator {
	return &JapaneseRandomPhraseGenerator{
		User:            u,
		DbProvider:      d,
		FileSystem:      f,
		RandomGenerator: r,
	}
}

func DefineNumber(num int, argNum string) int {
	if num <= 0 {
		num = 1
	}

	argNumConv, err := strconv.Atoi(argNum)
	if err != nil {
		argNumConv = 1
	}
	if argNumConv <= 0 {
		argNumConv = 1
	}

	if argNumConv > num {
		return argNumConv
	} else {
		return num
	}
}

func (j JapaneseRandomPhraseGenerator) Generate(num int) ([]string, error) {
	dbFileDirPathGetter := NewDBFileDirPathGetter(j.User)

	dbFileDirPath, err := dbFileDirPathGetter.GetFileDirPath()
	if err != nil {
		return nil, err
	}

	dbFilePath := filepath.Join(dbFileDirPath, constant.WNJPN_DB_FILE_NAME)
	if !j.FileSystem.Exists(dbFilePath) {
		fmt.Println(color.YellowString(constant.GENERATE_MESSAGE_NOTIFY_DOWNLOAD_REQUIRED))
		return nil, nil
	}

	db, err := j.DbProvider.Connect(dbFilePath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := j.DbProvider.Query(db, constant.GENERATE_SQL_GET_ALL_JAPANESE_AVN_WORDS)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	allAVNWords := make([]model.Word, 0)
	for rows.Next() {
		var word model.Word
		err = rows.Scan(&word.Lemma, &word.Pos)
		if err != nil {
			return nil, err
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

	jrp := make([]string, 0)
	for i := 0; i < num; i++ {
		randomIndexA := j.RandomGenerator.Intn(len(allAVWords))
		randomIndexB := j.RandomGenerator.Intn(len(allNWords))
		randomWordA := allAVWords[randomIndexA]
		randomWordB := allNWords[randomIndexB]
		jrp = append(jrp, randomWordA.Lemma.String+randomWordB.Lemma.String)
	}

	return jrp, nil
}
