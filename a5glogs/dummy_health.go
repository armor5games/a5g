package a5glogs

import (
	"fmt"

	"github.com/armor5games/a5g/a5gfields"
	"github.com/armor5games/a5g/a5gkv"
	"github.com/pkg/errors"
)

type DummyHealth struct{ Logger Logger }

func NewDummyHealth(l Logger) (*DummyHealth, error) {
	if l == nil {
		return nil, errors.New("nil pointer")
	}
	return &DummyHealth{Logger: l}, nil
}

func (l *DummyHealth) Event(eventName string) {
	l.Logger.Debug(eventName)
}

func (l *DummyHealth) EventKv(eventName string, m map[string]string) {
	l.Logger.With(dummyHealthKVToFields(m)...).Debug(eventName)
}

func (l *DummyHealth) EventErr(eventName string, err error) error {
	err = fmt.Errorf("%s %s", eventName, err.Error())
	l.Logger.Error(err.Error())
	return err
}

func (l *DummyHealth) EventErrKv(
	eventName string, err error, m map[string]string) error {
	var (
		a  = dummyHealthKVToFields(m)
		m2 = a5gkv.NewByMapString(m)
	)
	err = fmt.Errorf("%s %s", eventName, err.Error())
	l.Logger.With(a...).Error(err.Error())
	return fmt.Errorf("%s %s", err.Error(), m2.String())
}

func (l *DummyHealth) Timing(eventName string, nanoSeconds int64) {
	l.Logger.
		With(a5gfields.Int64("elapsedNanoseconds", nanoSeconds)).Debug(eventName)
}

func (l *DummyHealth) TimingKv(
	eventName string, nanoSeconds int64, m map[string]string) {
	a := dummyHealthKVToFields(m)
	a = append(a, a5gfields.Int64("elapsedNanoseconds", nanoSeconds))
	l.Logger.With(a...).Debug(eventName)
}

func dummyHealthKVToFields(m map[string]string) []a5gfields.Field {
	if len(m) == 0 {
		return nil
	}
	var a []a5gfields.Field
	for k, v := range m {
		a = append(a, a5gfields.String(k, v))
	}
	return a
}
