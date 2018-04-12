package a5gdb

import "github.com/gocraft/dbr"

type Pooler interface {
	ReadPool() Connector
	WritePool() Connector
	Validate() error
}

type Connector interface {
	NewSession() *dbr.Session
	Close() error
}
