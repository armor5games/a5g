package a5glogs

import (
	"github.com/armor5games/a5g/a5gfields"
	"github.com/pkg/errors"
)

type MGOOutputer interface {
	Output(int, string) error
}

type mgoOutputer struct{ Logger Logger }

func NewMGOOutputer(l Logger) (*mgoOutputer, error) {
	if l == nil {
		return nil, errors.New("logger missing")
	}
	return &mgoOutputer{Logger: l}, nil
}

func (l *mgoOutputer) Output(callDepth int, s string) error {
	l.Logger.With(a5gfields.Int("callDepth", callDepth)).Debug(s)
	return nil
}
