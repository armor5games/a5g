package goarmormongodb

import mgo "gopkg.in/mgo.v2"

type Pooler interface {
	Pool() Connector
	Enabled() bool
}

type Connector interface {
	NewSession() (*mgo.Session, error)
	Close() error
}
