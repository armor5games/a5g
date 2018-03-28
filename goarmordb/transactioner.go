package goarmordb

import "github.com/gocraft/dbr"

type Transactioner interface {
	Pooler
	Transaction() *dbr.Tx
}
