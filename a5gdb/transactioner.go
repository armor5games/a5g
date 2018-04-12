package a5gdb

import "github.com/gocraft/dbr"

type Transactioner interface {
	Pooler
	Transaction() *dbr.Tx
}
