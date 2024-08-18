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
	Generate(num int, prefix string, suffix string) ([]string, error)
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

func GetFirstConvertibleToString(args []string) string {
	for _, arg := range args {
		if _, err := strconv.Atoi(arg); err == nil {
			return arg
		}
	}
	return "1"
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

func (j JapaneseRandomPhraseGenerator) Generate(num int, prefix string, suffix string) ([]string, error) {
	// get db file path
	dbFileDirPathGetter := NewDBFileDirPathGetter(j.User)
	dbFileDirPath, err := dbFileDirPathGetter.GetFileDirPath()
	if err != nil {
		return nil, err
	}
	dbFilePath := filepath.Join(dbFileDirPath, constant.WNJPN_DB_FILE_NAME)
	if !j.FileSystem.Exists(dbFilePath) {
		// if db file does not exist, notify to download
		fmt.Println(color.YellowString(constant.GENERATE_MESSAGE_NOTIFY_DOWNLOAD_REQUIRED))
		return make([]string, 0), nil
	}

	// define query
	var query string
	if prefix == "" && suffix == "" {
		query = constant.GENERATE_SQL_GET_ALL_JAPANESE_AVN_WORDS
	} else if prefix != "" && suffix == "" {
		query = constant.GENERATE_SQL_GET_ALL_JAPANESE_N_WORDS
	} else if prefix == "" && suffix != "" {
		query = constant.GENERATE_SQL_GET_ALL_JAPANESE_AV_WORDS
	} else {
		// if both prefix and suffix are provided, notify to use only one
		fmt.Println(color.YellowString(constant.GENERATE_MESSAGE_NOTIFY_USE_ONLY_ONE))
		return make([]string, 0), nil
	}

	// connect to db
	db, err := j.DbProvider.Connect(dbFilePath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// get all words from db
	rows, err := j.DbProvider.Query(db, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	allWords := make([]model.Word, 0)
	for rows.Next() {
		var word model.Word
		if err := rows.Scan(&word.Lemma, &word.Pos); err != nil {
			return nil, err
		}
		allWords = append(allWords, word)
	}

	var allAVWords []model.Word
	var allNWords []model.Word

	for _, word := range allWords {
		if word.Pos.Valid && word.Pos.String == "n" {
			allNWords = append(allNWords, word)
		} else {
			allAVWords = append(allAVWords, word)
		}
	}

	jrp := make([]string, 0)
	for i := 0; i < num; i++ {
		var prefixWord string
		if prefix != "" {
			prefixWord = prefix
		} else {
			randomIndexPrefix := rand.Intn(len(allAVWords))
			randomWordPrefix := allAVWords[randomIndexPrefix]
			prefixWord = randomWordPrefix.Lemma.String
		}

		var suffixWord string
		if suffix != "" {
			suffixWord = suffix
		} else {
			randomIndexSuffix := rand.Intn(len(allNWords))
			randomWordSuffix := allNWords[randomIndexSuffix]
			suffixWord = randomWordSuffix.Lemma.String
		}

		jrp = append(jrp, prefixWord+suffixWord)
	}

	return jrp, nil
}
