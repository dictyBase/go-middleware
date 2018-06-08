// Package logrus is a logrus(https://godoc.org/github.com/Sirupsen/logrus)
// net/http middleware to log information about http handler. It can output
// both in text and JSON format. Both output format can be completely
// customized.
package logrus

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// LogResponseWriter is a custom type that extends http.ResponseWriter interface
// to capture and provide an easy access to http status code
type LogResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

// Status is a easy way to retrieve the status code
func (w *LogResponseWriter) Status() int {
	return w.status
}

// Size provides the size of response object
func (w *LogResponseWriter) Size() int {
	return w.size
}

// Header returns the header to satisfy the http.ResponseWriter interface
func (w *LogResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// Write capture the size of the data written and satisfy the http.ResponseWriter interface
func (w *LogResponseWriter) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	written, err := w.ResponseWriter.Write(data)
	w.size += written
	return written, err
}

// WriteHeader capture the status code and satisfies the http.ResponseWriter interface
func (w *LogResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	w.status = code
}

type timer interface {
	Now() time.Time
	Since(time.Time) time.Duration
}

type realClock struct{}

// Now return the current time
func (rc *realClock) Now() time.Time {
	return time.Now()
}

// Since returns the time passed
func (rc *realClock) Since(t time.Time) time.Duration {
	return time.Since(t)
}

// Logger is a middleware handler that logs the request as it goes in and the response as it goes out.
type Logger struct {
	// Logger is the log.Logger instance used to log messages with the Logger middleware
	Logrus *logrus.Logger
	// Name is the name of the application as recorded in latency metrics
	Name string

	clock timer
}

// NewJSONLogger returns a new *Logger that logs in JSON format
func NewJSONLogger() *Logger {
	log := logrus.New()
	log.Level = logrus.InfoLevel
	log.Formatter = &logrus.JSONFormatter{
		TimestampFormat: "02/Jan/2006:15:04:05",
	}
	log.Out = os.Stderr
	return &Logger{
		Logrus: log,
		Name:   "web",
		clock:  &realClock{},
	}
}

// NewJSONFileLogger writes to a file in JSON format
func NewJSONFileLogger(w io.Writer) *Logger {
	logger := NewJSONLogger()
	logger.Logrus.Out = w
	return logger
}

// NewLogger returns a new *Logger
func NewLogger() *Logger {
	log := logrus.New()
	log.Level = logrus.InfoLevel
	log.Formatter = &logrus.TextFormatter{
		TimestampFormat: "02/Jan/2006:15:04:05",
		FullTimestamp:   true,
	}
	log.Out = os.Stderr
	return &Logger{
		Logrus: log,
		Name:   "web",
		clock:  &realClock{},
	}
}

// NewFileLogger writes to a file
func NewFileLogger(w io.Writer) *Logger {
	logger := NewLogger()
	logger.Logrus.Out = w
	return logger
}

// NewCustomMiddleware builds a *Logger with the given level and formatter
func NewCustomMiddleware(level logrus.Level, formatter logrus.Formatter, name string) *Logger {
	log := logrus.New()
	log.Level = level
	log.Formatter = formatter
	log.Out = os.Stderr
	return &Logger{
		Logrus: log,
		Name:   name,
		clock:  &realClock{},
	}
}

// NewMiddlewareFromLogger returns a new *Logger which writes to a given logrus logger.
func NewMiddlewareFromLogger(logger *logrus.Logger, name string) *Logger {
	return &Logger{Logrus: logger, Name: name, clock: &realClock{}}
}

// MiddlewareFn works with http.HandlerFunc type
func (l *Logger) MiddlewareFn(fn http.HandlerFunc) http.HandlerFunc {
	newfn := func(w http.ResponseWriter, r *http.Request) {
		start := l.clock.Now()

		// Try to get the real IP
		remoteAddr := r.RemoteAddr
		if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			remoteAddr = realIP
		}
		entry := l.Logrus.WithFields(logrus.Fields{
			"request":    r.RequestURI,
			"method":     r.Method,
			"remote":     remoteAddr,
			"user-agent": r.UserAgent(),
			"referer":    r.Referer(),
		})
		res := &LogResponseWriter{ResponseWriter: w}
		fn(res, r)

		latency := l.clock.Since(start)
		entry.WithFields(logrus.Fields{
			"status": res.Status(),
			"took":   latency,
			"size":   res.Size(),
		}).Info("completed handling request")
	}
	return newfn
}

// Middleware works with http.Handler type
func (l *Logger) Middleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := l.clock.Now()

		// Try to get the real IP
		remoteAddr := r.RemoteAddr
		if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			remoteAddr = realIP
		}
		entry := l.Logrus.WithFields(logrus.Fields{
			"request":    r.RequestURI,
			"method":     r.Method,
			"remote":     remoteAddr,
			"user-agent": r.UserAgent(),
			"referer":    r.Referer(),
		})
		res := &LogResponseWriter{ResponseWriter: w}
		h.ServeHTTP(res, r)

		latency := l.clock.Since(start)
		entry.WithFields(logrus.Fields{
			"status": res.Status(),
			"took":   latency,
			"size":   res.Size(),
		}).Info("completed handling request")
	}
	return http.HandlerFunc(fn)
}
