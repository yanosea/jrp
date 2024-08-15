package logic

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"strconv"

	"github.com/fatih/color"

	"github.com/yanosea/jrp/constant"
	"github.com/yanosea/jrp/internal/database"
	"github.com/yanosea/jrp/internal/fs"
	"github.com/yanosea/jrp/internal/usermanager"
	"github.com/yanosea/jrp/model"
)

type Generator interface {
	Generate(num int) ([]string, error)
}

type JapaneseRandomPhraseGenerator struct {
	User       usermanager.UserProvider
	DbProvider database.DatabaseProvider
	FileSystem fs.FileManager
}

func NewJapaneseRandomPhraseGenerator(u usermanager.UserProvider, d database.DatabaseProvider, f fs.FileManager) *JapaneseRandomPhraseGenerator {
	return &JapaneseRandomPhraseGenerator{
		User:       u,
		DbProvider: d,
		FileSystem: f,
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
		return make([]string, 0), nil
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
		if err := rows.Scan(&word.Lemma, &word.Pos); err != nil {
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
		randomIndexA := rand.Intn(len(allAVWords))
		randomIndexB := rand.Intn(len(allNWords))
		randomWordA := allAVWords[randomIndexA]
		randomWordB := allNWords[randomIndexB]
		jrp = append(jrp, randomWordA.Lemma.String+randomWordB.Lemma.String)
	}

	return jrp, nil
}
