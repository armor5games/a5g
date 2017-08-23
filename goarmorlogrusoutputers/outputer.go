// Package goarmorlogrusoutputers aim to use logrus
// as mgo (http://gopkg.in/mgo.v2)
// logger (http://godoc.org/gopkg.in/mgo.v2#SetLogger).
package goarmorlogrusoutputers

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Outputer interface {
	Output(int, string) error
}

type DummyOutputer struct {
	Logger *logrus.Logger
}

func New(l *logrus.Logger) (*DummyOutputer, error) {
	if l == nil {
		return nil, errors.New("nil pointer")
	}

	return &DummyOutputer{Logger: l}, nil
}

func (l *DummyOutputer) Output(callDepth int, s string) error {
	if l == nil {
		return errors.New("nil pointer")
	}

	l.Logger.WithFields(map[string]interface{}{"callDepth": callDepth}).Debug(s)

	return nil
}
