package generator

import (
	jrp "github.com/yanosea/jrp/app/database/jrp/model"
	wnjpn "github.com/yanosea/jrp/app/database/wnjpn/model"
	"github.com/yanosea/jrp/app/database/wnjpn/repository"
	"github.com/yanosea/jrp/app/proxy/os"
	"github.com/yanosea/jrp/app/proxy/rand"
	"github.com/yanosea/jrp/app/proxy/sql"
	"github.com/yanosea/jrp/app/proxy/time"
)

// Generatable is an interface for Generator.
type Generatable interface {
	GenerateJrp(wnJpnDBFilePath string, num int, word string, mode GenerateMode) (GenerateResult, []jrp.Jrp, error)
}

// Generatoor is a struct that implements Generatable interface.
type Generator struct {
	OsProxy         osproxy.Os
	RandProxy       randproxy.Rand
	SqlProxy        sqlproxy.Sql
	TimeProxy       timeproxy.Time
	WNJpnRepository repository.WNJpnRepositoryInterface
}

// New is a constructor of Generator.
func New(
	osProxy osproxy.Os,
	randProxy randproxy.Rand,
	sqlProxy sqlproxy.Sql,
	timeProxy timeproxy.Time,
	wnJpnRepository repository.WNJpnRepositoryInterface,
) *Generator {
	return &Generator{
		OsProxy:         osProxy,
		RandProxy:       randProxy,
		SqlProxy:        sqlProxy,
		TimeProxy:       timeProxy,
		WNJpnRepository: wnJpnRepository,
	}
}

// GenerateJrp generates jrps.
func (g *Generator) GenerateJrp(wnJpnDBFilePath string, num int, word string, mode GenerateMode) (GenerateResult, []jrp.Jrp, error) {
	if _, err := g.OsProxy.Stat(wnJpnDBFilePath); g.OsProxy.IsNotExist(err) {
		// if WordNet Japan sqlite database file does not exist, return warning
		return DBFileNotFound, nil, nil
	}

	// define prefix, suffix and query
	prefix, suffix := g.getPrefixAndSuffix(word, mode)

	// execute query and get all words
	allWords, err := g.getAllWords(wnJpnDBFilePath, mode)
	if err != nil {
		return GeneratedFailed, nil, err
	}

	// separate all words into AV and N words
	allAVWords, allNWords := g.separateWords(allWords)

	// get jrps
	jrps := g.getJrps(num, allAVWords, allNWords, prefix, suffix, mode)

	return GeneratedSuccessfully, jrps, nil
}

// getAllWords gets all words based on mode.
func (g *Generator) getAllWords(wnJpnDBFilePath string, mode GenerateMode) ([]wnjpn.Word, error) {
	var allWords []wnjpn.Word
	var err error
	switch mode {
	case WithNoPrefixOrSuffix:
		allWords, err = g.WNJpnRepository.GetAllAVNWords(wnJpnDBFilePath)
	case WithPrefix:
		allWords, err = g.WNJpnRepository.GetAllNWords(wnJpnDBFilePath)
	case WithSuffix:
		allWords, err = g.WNJpnRepository.GetAllAVWords(wnJpnDBFilePath)
	}
	return allWords, err
}

// getJrps gets jrps based on mode.
func (g *Generator) getJrps(num int,
	allAVWords []wnjpn.Word,
	allNWords []wnjpn.Word,
	argPrefix string,
	argSuffix string,
	mode GenerateMode,
) []jrp.Jrp {
	jrps := make([]jrp.Jrp, 0)
	createdAt := g.TimeProxy.Now()

	for i := 0; i < num; i++ {
		var prefixWord string
		var suffixWord string
		switch mode {
		case WithNoPrefixOrSuffix:
			// get random number for prefix
			randomIndexForPrefix := g.RandProxy.Intn(len(allAVWords))
			// get random prefix word
			randomPrefix := allAVWords[randomIndexForPrefix]
			// get random number for suffix
			randomIndexForSuffix := g.RandProxy.Intn(len(allNWords))
			// get random suffix word
			randomSuffix := allNWords[randomIndexForSuffix]
			// set prefix word and suffix word
			prefixWord = randomPrefix.Lemma.FieldNullString.String
			suffixWord = randomSuffix.Lemma.FieldNullString.String
			// set argPrefix and argSuffix to empty string
			argPrefix = ""
			argSuffix = ""
		case WithPrefix:
			// get random number for suffix
			randomIndexSuffix := g.RandProxy.Intn(len(allAVWords))
			// get random prefix word
			randomSuffix := allAVWords[randomIndexSuffix]
			// set prefix word and suffix word
			prefixWord = argPrefix
			suffixWord = randomSuffix.Lemma.FieldNullString.String
			// set argSuffix to empty string
			argSuffix = ""
		case WithSuffix:
			// get random number for argPrefix
			randomIndexPrefix := g.RandProxy.Intn(len(allNWords))
			// get random prefix word
			randomPrefix := allNWords[randomIndexPrefix]
			// set prefix word and suffix word
			prefixWord = randomPrefix.Lemma.FieldNullString.String
			suffixWord = argSuffix
			// set argPrefix to empty string
			argPrefix = ""
		}

		jrp := jrp.Jrp{
			Phrase:    prefixWord + suffixWord,
			Prefix:    g.SqlProxy.StringToNullString(argPrefix),
			Suffix:    g.SqlProxy.StringToNullString(argSuffix),
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
		}

		jrps = append(jrps, jrp)
	}

	return jrps
}

// getPrefixAndSuffix gets prefix word and suffix word based on mode.
func (g *Generator) getPrefixAndSuffix(word string, mode GenerateMode) (string, string) {
	var prefix, suffix string
	// define prefix and suffix
	switch mode {
	case WithNoPrefixOrSuffix:
		prefix = ""
		suffix = ""
	case WithPrefix:
		prefix = word
		suffix = ""
	case WithSuffix:
		prefix = ""
		suffix = word
	}

	return prefix, suffix
}

// separateWords separates all words into AV and N words.
func (g *Generator) separateWords(allWords []wnjpn.Word) ([]wnjpn.Word, []wnjpn.Word) {
	allAVWords := []wnjpn.Word{}
	allNWords := []wnjpn.Word{}
	for _, word := range allWords {
		if word.Pos.FieldNullString.Valid && word.Pos.FieldNullString.String == "n" {
			// if word is noun
			allNWords = append(allNWords, word)
		} else {
			// if word is adjective or verb
			allAVWords = append(allAVWords, word)
		}
	}

	return allAVWords, allNWords
}
