// Package a5glogs by chi example
// <https://github.com/pressly/chi/blob/master/_examples/logging/main.go>
package a5glogs

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

func NewStructuredLogger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{logger})
}

type StructuredLogger struct {
	Logger *logrus.Logger `json:"logger"`
}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{Logger: logrus.NewEntry(l.Logger)}
	logFields := logrus.Fields{}
	logFields["ts"] = time.Now().UTC().Format(time.RFC1123)
	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields["reqID"] = reqID
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	logFields["httpScheme"] = scheme
	logFields["httpProto"] = r.Proto
	logFields["httpMethod"] = r.Method
	logFields["remoteAddr"] = r.RemoteAddr
	logFields["userAgent"] = r.UserAgent()
	logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)
	entry.Logger = entry.Logger.WithFields(logFields)
	entry.Logger.Debugln("request started")
	return entry
}

type StructuredLoggerEntry struct {
	Logger logrus.FieldLogger `json:"logger"`
}

func (l *StructuredLoggerEntry) Write(status, bytes int, elapsed time.Duration) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"respStatus":      status,
		"respBytesLength": bytes,
		"respElaspedMs":   float64(elapsed.Nanoseconds()) / 1000000.0,
	})
	l.Logger.Debugln("request complete")
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

// Helper methods used by the application to get the request-scoped
// logger entry and set additional fields between handlers.
//
// This is a useful pattern to use to set state on the entry as it
// passes through the handler chain, which at any point can be logged
// with a call to .Print(), .Info(), etc.

func GetLogEntry(r *http.Request) logrus.FieldLogger {
	entry, _ := middleware.GetLogEntry(r).(*StructuredLoggerEntry)
	return entry.Logger
}

func LogEntrySetField(r *http.Request, key string, value interface{}) {
	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		entry.Logger = entry.Logger.WithField(key, value)
	}
}

func LogEntrySetFields(r *http.Request, fields map[string]interface{}) {
	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		entry.Logger = entry.Logger.WithFields(fields)
	}
}
