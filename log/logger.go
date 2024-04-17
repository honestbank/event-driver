package log

import (
	"log"
)

type Level string

const (
	DEBUG Level = "DEBUG"
	INFO  Level = "INFO"
	WARN  Level = "WARN"
	ERROR Level = "ERROR"
)

type Logger struct {
	domain  string
	verbose bool
}

func New(domain string) *Logger {
	return &Logger{
		domain: domain,
	}
}

func (l *Logger) Verbose() {
	l.verbose = true
}

func (l *Logger) Debug(message string, args ...any) {
	l.Printf(DEBUG, message, args...)
}

func (l *Logger) Info(message string, args ...any) {
	l.Printf(INFO, message, args...)
}

func (l *Logger) Warn(message string, args ...any) {
	l.Printf(WARN, message, args...)
}

func (l *Logger) Error(message string, args ...any) {
	l.Printf(ERROR, message, args...)
}

func (l *Logger) Printf(logLevel Level, message string, args ...any) {
	if !l.verbose && logLevel == DEBUG {
		return
	}
	log.Printf("[%s] %s: "+message, append([]any{logLevel, l.domain}, args...)...)
}
