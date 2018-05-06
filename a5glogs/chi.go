// Package a5glogs by chi example
// <https://github.com/pressly/chi/blob/master/_examples/logging/main.go>
package a5glogs

import (
	"fmt"
	"net/http"
	"time"

	"github.com/armor5games/a5g/a5gfields"
	"github.com/armor5games/a5g/a5gkv"
	"github.com/armor5games/a5g/a5gvalues"
	"github.com/go-chi/chi/middleware"
)

func NewChiLogger(l Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&chiLogger{l})
}

type chiLogger struct{ logger Logger }

func (l *chiLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	v := &chiLoggerEntry{logger: l.logger}
	m := a5gkv.New()
	m["ts"] = a5gvalues.String(time.Now().UTC().Format(time.RFC1123))
	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		m["reqID"] = a5gvalues.String(reqID)
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	m["httpScheme"] = a5gvalues.String(scheme)
	m["httpProto"] = a5gvalues.String(r.Proto)
	m["httpMethod"] = a5gvalues.String(r.Method)
	m["remoteAddr"] = a5gvalues.String(r.RemoteAddr)
	m["userAgent"] = a5gvalues.String(r.UserAgent())
	m["uri"] = a5gvalues.String(fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI))
	v.logger = v.logger.With(m.Fields()...)
	v.logger.Debug("request started")
	return v
}

type chiLoggerEntry struct{ logger Logger }

func (l *chiLoggerEntry) Write(
	respStatus, respBytesLen int, elapsedDuration time.Duration) {
	m := a5gkv.New()
	m["respStatus"] = a5gvalues.Int(respStatus)
	m["respBytesLength"] = a5gvalues.Int(respBytesLen)
	m["respElaspedMs"] = a5gvalues.Float64(float64(elapsedDuration.Nanoseconds()) / 1000000.0)
	l.logger = l.logger.With(m.Fields()...)
	l.logger.Debug("request complete")
}

func (l *chiLoggerEntry) Panic(v interface{}, stack []byte) {
	m := a5gkv.New()
	m["stack"] = a5gvalues.Bytes(stack)
	m["panic"] = a5gvalues.EmptyInterface(v)
	l.logger = l.logger.With(m.Fields()...)
}

// Helper methods used by the application to get the request-scoped
// logger entry and set additional fields between handlers.
//
// This is a useful pattern to use to set state on the entry as it
// passes through the handler chain, which at any point can be logged
// with a call to .Print(), .Info(), etc.

func GetLogEntry(r *http.Request) Logger {
	v, _ := middleware.GetLogEntry(r).(*chiLoggerEntry)
	return v.logger
}

func LogEntrySetField(r *http.Request, key string, value fmt.Stringer) {
	v, ok := r.Context().Value(middleware.LogEntryCtxKey).(*chiLoggerEntry)
	if ok {
		v.logger = v.logger.With(a5gfields.New(key, value))
	}
}

func LogEntrySetFields(r *http.Request, a ...a5gfields.Field) {
	v, ok := r.Context().Value(middleware.LogEntryCtxKey).(*chiLoggerEntry)
	if ok {
		v.logger = v.logger.With(a...)
	}
}
