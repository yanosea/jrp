package testutility

import (
	"github.com/yanosea/jrp/app/database/jrp/model"
	"github.com/yanosea/jrp/app/proxy/fmt"
	"github.com/yanosea/jrp/app/proxy/sort"
	"github.com/yanosea/jrp/app/proxy/sql"
	"github.com/yanosea/jrp/app/proxy/strconv"
	"github.com/yanosea/jrp/app/proxy/strings"
	"github.com/yanosea/jrp/app/proxy/time"
)

// JrpCheckerInterface is an interface for checking Jrp.
type JrpCheckerInterface interface {
	GetJrpSeq(jrpDBFilePath string) (int, error)
	IsExist(jrpDBFilePath string, id int) (bool, error)
	IsFavorited(jrpDBFilePath string, id int) (bool, error)
	IsSameJrps(got, want []model.Jrp) bool
}

// JrpChecker is a struct for checking Jrp.
type JrpChecker struct {
	FmtProxy     fmtproxy.Fmt
	SortProxy    sortproxy.Sort
	SqlProxy     sqlproxy.Sql
	StrconvProxy strconvproxy.Strconv
	StringsProxy stringsproxy.Strings
}

// NewJrpChecker is a constructor for JrpChecker.
func NewJrpChecker(
	fmtProxy fmtproxy.Fmt,
	sortProxy sortproxy.Sort,
	sqlProxy sqlproxy.Sql,
	strconvProxy strconvproxy.Strconv,
	stringsProxy stringsproxy.Strings,
) *JrpChecker {
	return &JrpChecker{
		FmtProxy:     fmtProxy,
		SortProxy:    sortProxy,
		SqlProxy:     sqlProxy,
		StrconvProxy: strconvProxy,
		StringsProxy: stringsProxy,
	}
}

// GetJrpSeq returns the sequence of jrp.
func (j *JrpChecker) GetJrpSeq(jrpDBFilePath string) (int, error) {
	db, err := j.SqlProxy.Open(sqlproxy.Sqlite, jrpDBFilePath)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT seq FROM sqlite_sequence WHERE sqlite_sequence.name = 'jrp';")
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var seq int
	for rows.Next() {
		if err := rows.Scan(&seq); err != nil {
			return 0, err
		}
	}

	return seq, nil
}

// IsExist checks if jrp exists.
func (j *JrpChecker) IsExist(jrpDBFilePath string, id int) (bool, error) {
	db, err := j.SqlProxy.Open(sqlproxy.Sqlite, jrpDBFilePath)
	if err != nil {
		return false, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT COUNT(*) FROM jrp WHERE jrp.ID = (?);", j.StrconvProxy.Itoa(id))
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return false, err
		}
	}

	return count == 1, nil
}

// IsFavorited checks if jrp is favorited.
func (j *JrpChecker) IsFavorited(jrpDBFilePath string, id int) (bool, error) {
	db, err := j.SqlProxy.Open(sqlproxy.Sqlite, jrpDBFilePath)
	if err != nil {
		return false, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT COUNT(*) FROM jrp WHERE jrp.IsFavorite = 1 AND jrp.ID = (?);", j.StrconvProxy.Itoa(id))
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return false, err
		}
	}

	return count == 1, nil
}

// IsSameJrps checks if jrps are the same.
func (j *JrpChecker) IsSameJrps(got, want []model.Jrp) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i].IsFavorited == 1 {
			if got[i].ID != want[i].ID ||
				got[i].Phrase != want[i].Phrase ||
				!isSameNullStringInstance(got[i].Prefix, want[i].Prefix) ||
				!isSameNullStringInstance(got[i].Suffix, want[i].Suffix) ||
				got[i].IsFavorited != want[i].IsFavorited ||
				!isSameTimeInstance(got[i].CreatedAt, want[i].CreatedAt) {
				return false
			}
			return true
		}

		if got[i].ID != want[i].ID ||
			got[i].Phrase != want[i].Phrase ||
			!isSameNullStringInstance(got[i].Prefix, want[i].Prefix) ||
			!isSameNullStringInstance(got[i].Suffix, want[i].Suffix) ||
			got[i].IsFavorited != want[i].IsFavorited ||
			!isSameTimeInstance(got[i].CreatedAt, want[i].CreatedAt) ||
			!isSameTimeInstance(got[i].UpdatedAt, want[i].UpdatedAt) {
			return false
		}
	}
	return true
}

// isSameNullStringInstance checks if NullStringInstance is the same.
func isSameNullStringInstance(a, b *sqlproxy.NullStringInstance) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.FieldNullString.String == b.FieldNullString.String
}

// isSameTimeInstance checks if TimeInstance is the same.
func isSameTimeInstance(a, b *timeproxy.TimeInstance) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.FieldTime.Equal(b.FieldTime)
}
