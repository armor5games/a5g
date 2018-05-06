package a5glogs

import (
	"github.com/armor5games/a5g/a5gfields"
	"github.com/sirupsen/logrus"
)

type logrusEntryWrapper struct{ logger *logrus.Entry }

func (l *logrusEntryWrapper) With(a ...a5gfields.Field) Logger {
	if len(a) == 0 {
		return &logrusEntryWrapper{logger: l.logger}
	}
	m := make(logrus.Fields)
	for _, v := range a {
		m[v.Key()] = v.Value()
	}
	return &logrusEntryWrapper{logger: l.logger.WithFields(m)}
}

func (l *logrusEntryWrapper) Debug(s string, a ...a5gfields.Field) {
	if len(a) == 0 {
		l.logger.Debug(s)
		return
	}
	m := make(logrus.Fields)
	for _, v := range a {
		m[v.Key()] = v.Value()
	}
	l.logger.WithFields(m).Debug(s)
}

func (l *logrusEntryWrapper) Info(s string, a ...a5gfields.Field) {
	if len(a) == 0 {
		l.logger.Info(s)
		return
	}
	m := make(logrus.Fields)
	for _, v := range a {
		m[v.Key()] = v.Value()
	}
	l.logger.WithFields(m).Info(s)
}

func (l *logrusEntryWrapper) Warn(s string, a ...a5gfields.Field) {
	if len(a) == 0 {
		l.logger.Warn(s)
		return
	}
	m := make(logrus.Fields)
	for _, v := range a {
		m[v.Key()] = v.Value()
	}
	l.logger.WithFields(m).Warn(s)
}

func (l *logrusEntryWrapper) Error(s string, a ...a5gfields.Field) {
	if len(a) == 0 {
		l.logger.Error(s)
		return
	}
	m := make(logrus.Fields)
	for _, v := range a {
		m[v.Key()] = v.Value()
	}
	l.logger.WithFields(m).Error(s)
}

func (l *logrusEntryWrapper) Panic(s string, a ...a5gfields.Field) {
	if len(a) == 0 {
		l.logger.Panic(s)
		return
	}
	m := make(logrus.Fields)
	for _, v := range a {
		m[v.Key()] = v.Value()
	}
	l.logger.WithFields(m).Panic(s)
}

func (l *logrusEntryWrapper) Fatal(s string, a ...a5gfields.Field) {
	if len(a) == 0 {
		l.logger.Fatal(s)
		return
	}
	m := make(logrus.Fields)
	for _, v := range a {
		m[v.Key()] = v.Value()
	}
	l.logger.WithFields(m).Fatal(s)
}
