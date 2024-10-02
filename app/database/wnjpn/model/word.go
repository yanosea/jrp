package model

import (
	"github.com/yanosea/jrp/app/proxy/sql"
)

// Word is a struct that represents word in wnjpn db file.
type Word struct {
	WordID int
	Lang   *sqlproxy.NullStringInstance
	Lemma  *sqlproxy.NullStringInstance
	Pron   *sqlproxy.NullStringInstance
	Pos    *sqlproxy.NullStringInstance
}
