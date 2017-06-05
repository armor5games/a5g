// Package gameserverlogs by chi example
// <https://github.com/pressly/chi/blob/master/_examples/logging/main.go>
package gameserverlogs

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pressly/chi/middleware"
	"github.com/sirupsen/logrus"
)

func NewStructuredLogger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{logger})
}

type StructuredLogger struct {
	Logger *logrus.Logger
}

func (l *StructuredLogger) NewLogEntry(req *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{Logger: logrus.NewEntry(l.Logger)}
	logFields := logrus.Fields{}

	logFields["ts"] = time.Now().UTC().Format(time.RFC1123)

	if reqID := middleware.GetReqID(req.Context()); reqID != "" {
		logFields["reqID"] = reqID
	}

	scheme := "http"
	if req.TLS != nil {
		scheme = "https"
	}
	logFields["httpScheme"] = scheme
	logFields["httpProto"] = req.Proto
	logFields["httpMethod"] = req.Method

	logFields["remoteAddr"] = req.RemoteAddr
	logFields["userAgent"] = req.UserAgent()

	logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, req.Host, req.RequestURI)

	entry.Logger = entry.Logger.WithFields(logFields)

	entry.Logger.Debugln("request started")

	return entry
}

type StructuredLoggerEntry struct {
	Logger logrus.FieldLogger
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

func GetLogEntry(req *http.Request) logrus.FieldLogger {
	entry := middleware.GetLogEntry(req).(*StructuredLoggerEntry)
	return entry.Logger
}

func LogEntrySetField(req *http.Request, key string, value interface{}) {
	if entry, ok := req.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		entry.Logger = entry.Logger.WithField(key, value)
	}
}

func LogEntrySetFields(req *http.Request, fields map[string]interface{}) {
	if entry, ok := req.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		entry.Logger = entry.Logger.WithFields(fields)
	}
}
