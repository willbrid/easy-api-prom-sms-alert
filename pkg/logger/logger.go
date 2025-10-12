package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	INFO LogLevel = iota
	WARNING
	ERROR
	FATAL
)

const DefaultCallerDepth = 2

type LogLevel int

type ILogger interface {
	Info(message string, args ...any)
	Warning(message string, args ...any)
	Error(message string, args ...any)
	Fatal(message string, args ...any)
}

type Logger struct {
	logger *log.Logger
}

func NewLogger() *Logger {
	logger := log.New(os.Stdout, "", 0)

	return &Logger{logger}
}

var logLevels = map[LogLevel]string{
	INFO:    "INFO",
	WARNING: "WARNING",
	ERROR:   "ERROR",
	FATAL:   "FATAL",
}

func (l *Logger) msg(level LogLevel, message string, args ...any) {
	logLevel, ok := logLevels[level]
	if !ok {
		logLevel = "UNKNOWN"
	}

	datetime := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	filename := ""
	if ok {
		filename = filepath.Base(file)
	}

	l.logger.SetOutput(os.Stdout)
	l.logger.SetFlags(0)
	l.logger.Printf("[%s][%s][%s:%d] : %s", datetime, logLevel, filename, line, fmt.Sprintf(message, args...))
}

func (l *Logger) Info(message string, args ...any) {
	l.msg(INFO, message, args...)
}

func (l *Logger) Warning(message string, args ...any) {
	l.msg(WARNING, message, args...)
}

func (l *Logger) Error(message string, args ...any) {
	l.msg(ERROR, message, args...)
}

func (l *Logger) Fatal(message string, args ...any) {
	l.msg(FATAL, message, args...)

	os.Exit(1)
}
