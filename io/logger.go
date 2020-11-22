package io

import (
	"bytes"
	"log"
	"os"
)

type level int

const (
	OFF   level = iota
	INFO        // v
	WARN        // vv
	DEBUG       // vvv
)

type ourlog interface {
	Info(v ...interface{})
	Warn(v ...interface{})
	Debug(v ...interface{})
}

type logger struct {
	level     level
	actualLog *log.Logger
	Buf       bytes.Buffer
}

func (l *logger) doWork(level level, v ...interface{}) {
	if level > l.level {
		return
	}

	l.actualLog.Print(v...)
}

func (l *logger) Info(v ...interface{}) {
	l.doWork(INFO, v...)
}

func (l *logger) Warn(v ...interface{}) {
	l.doWork(WARN, v...)
}

func (l *logger) Debug(v ...interface{}) {
	l.doWork(DEBUG, v...)
}

// Constructor
// writeOut enables the writing to standard out. The reasons we have this as a
// setting is so we can unit test really since when a user calls the program
// with `-v` it's implied they want to write out the logs
func NewLogger(level level, writeOut bool) ourlog {
	if level == 0 {
		return empty{}
	}

	l := &logger{level: level}
	l.actualLog = log.New(&l.Buf, "", 0)
	if writeOut {
		l.actualLog.SetOutput(os.Stdout)
	}

	return l
}

// In addition to being the default logger (ie: no logging), this one can also
// be used for testing by other packages
type empty struct{}

func (l empty) Debug(v ...interface{}) {
	return
}

func (l empty) Info(v ...interface{}) {
	return
}

func (l empty) Warn(v ...interface{}) {
	return
}
