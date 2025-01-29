// logger/logger.go
package logger

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	debug   *log.Logger
	info    *log.Logger
	error   *log.Logger
	isDebug bool
}

// NewLogger creates a new Logger instance
func NewLogger(debug bool) *Logger {
	flags := log.Ldate | log.Ltime | log.LUTC | log.Lshortfile

	return &Logger{
		debug:   log.New(os.Stdout, "DEBUG: ", flags),
		info:    log.New(os.Stdout, "INFO: ", flags),
		error:   log.New(os.Stderr, "ERROR: ", flags),
		isDebug: debug,
	}
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.isDebug {
		l.debug.Printf(format, v...)
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.info.Printf(format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.error.Printf(format, v...)
}

// SetOutput allows changing the output writer for all loggers
func (l *Logger) SetOutput(w io.Writer) {
	l.debug.SetOutput(w)
	l.info.SetOutput(w)
	l.error.SetOutput(w)
}
