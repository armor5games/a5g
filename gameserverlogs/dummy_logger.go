package gameserverlogs

type DummyLogger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Panic(...interface{})
	Panicf(string, ...interface{})
}

type Log struct {
	DummyLogger
}

func NewDummyLogger(l DummyLogger) DummyLogger {
	return &Log{DummyLogger: l}
}

func (l *Log) Print(v ...interface{}) {
	l.DummyLogger.Print(v...)
}

func (l *Log) Printf(s string, v ...interface{}) {
	l.DummyLogger.Printf(s, v...)
}

func (l *Log) Panic(v ...interface{}) {
	l.DummyLogger.Panic(v...)
}

func (l *Log) Panicf(s string, v ...interface{}) {
	l.DummyLogger.Panicf(s, v...)
}
