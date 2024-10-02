package model

import (
	"github.com/yanosea/jrp/app/proxy/sql"
	"github.com/yanosea/jrp/app/proxy/time"
)

// Jrp is a struct that represents jrp.
type Jrp struct {
	ID          int
	Phrase      string
	Prefix      *sqlproxy.NullStringInstance
	Suffix      *sqlproxy.NullStringInstance
	IsFavorited int
	CreatedAt   *timeproxy.TimeInstance
	UpdatedAt   *timeproxy.TimeInstance
}
