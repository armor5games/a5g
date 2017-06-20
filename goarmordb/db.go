package goarmordb

import "github.com/gocraft/dbr"

type Pooler interface {
	ReadPool() Connector
	WritePool() Connector
}

type Connector interface {
	NewSession() *dbr.Session
	Close() error
}
