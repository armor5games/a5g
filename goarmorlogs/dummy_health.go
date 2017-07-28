package goarmorlogs

import (
	"fmt"

	"github.com/armor5games/goarmor/goarmorapi"
	"github.com/sirupsen/logrus"
)

type DummyHealth struct {
	Logger *logrus.Logger `json:"logger"`
}

func NewDummyHealth(l *logrus.Logger) *DummyHealth {
	return &DummyHealth{Logger: l}
}

func (l *DummyHealth) Event(eventName string) {
	l.Logger.Debug(eventName)
}

func (l *DummyHealth) EventKv(eventName string, kvs map[string]string) {
	l.Logger.WithFields(dummyHealthKVToLogrusFields(kvs)).Debug(eventName)
}

func (l *DummyHealth) EventErr(eventName string, err error) error {
	err = fmt.Errorf("%s %s", eventName, err.Error())
	l.Logger.Error(err.Error())
	return err
}

func (l *DummyHealth) EventErrKv(eventName string, err error, kvs map[string]string) error {
	logrusKV := dummyHealthKVToLogrusFields(kvs)
	goarmorapiKV := goarmorapi.KV(logrusKV)

	err = fmt.Errorf("%s %s", eventName, err.Error())
	l.Logger.WithFields(logrusKV).Error(err.Error())

	return fmt.Errorf("%s %s", err.Error(), goarmorapiKV.String())
}

func (l *DummyHealth) Timing(eventName string, nanoseconds int64) {
	l.Logger.
		WithFields(logrus.Fields{"elapsedNanoseconds": nanoseconds}).Debug(eventName)
}

func (l *DummyHealth) TimingKv(eventName string, nanoseconds int64, kvs map[string]string) {
	f := dummyHealthKVToLogrusFields(kvs)
	f["elapsedNanoseconds"] = nanoseconds
	l.Logger.WithFields(f).Debug(eventName)
}

func dummyHealthKVToLogrusFields(keyValues map[string]string) logrus.Fields {
	if len(keyValues) == 0 {
		return nil
	}

	f := make(logrus.Fields)

	for k, v := range keyValues {
		f[k] = v
	}

	return f
}
