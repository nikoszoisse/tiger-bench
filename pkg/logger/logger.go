package logger

import (
	"fmt"
	"log"
	"os"
)

var DefaultLogger = NewLogger(os.Stdout, InfoLevel)

const (
	DebugLevel = 3
	InfoLevel  = 2
	ErrorLevel = 1
	NoLogging  = 0
)

type Logging interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
	SetLevel(level int)
}

type logger struct {
	log   *log.Logger
	level int
}

func NewLogger(dest *os.File, level int) Logging {
	return &logger{
		log:   log.New(dest, "", log.LstdFlags|log.Lshortfile),
		level: level,
	}
}

func (l *logger) SetLevel(level int) {
	l.level = level
}

func (l *logger) Debug(format string, values ...interface{}) {
	if l.level >= DebugLevel {
		_ = l.log.Output(2, "[DEBUG] "+fmt.Sprintf(format, values...)+"\n")
	}
}

func (l *logger) Info(format string, values ...interface{}) {
	if l.level >= InfoLevel {
		_ = l.log.Output(2, "[INFO] "+fmt.Sprintf(format, values...)+"\n")
	}
}

func (l *logger) Error(format string, values ...interface{}) {
	if l.level >= ErrorLevel {
		_ = l.log.Output(2, "[Error] "+fmt.Sprintf(format, values...)+"\n")
	}
}
