package a5glogs

import "github.com/armor5games/a5g/a5gfields"

type Logger interface {
	With(...a5gfields.Field) Logger
	Debug(string, ...a5gfields.Field)
	Info(string, ...a5gfields.Field)
	Warn(string, ...a5gfields.Field)
	Error(string, ...a5gfields.Field)
	Panic(string, ...a5gfields.Field)
	Fatal(string, ...a5gfields.Field)
}
