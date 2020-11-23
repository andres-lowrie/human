package io

import (
	"bytes"
	"log"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

type level int

const (
	OFF   level = iota
	INFO        // v
	WARN        // vv
	DEBUG       // vvv
)

type Ourlog interface {
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

	// Be smart about formatting output and just regular output. The user should
	// be able to just call anyone of these methods with an arbitrary number of
	// parameters and have it do the right thing
	//
	// That being said, here's what is expected:
	// 	```
	// 	format_string, stuff...  <- like typical calls to the *f methods of log, and fmt
	// 	string...                <- without leading format string
	// 	stuff...                 <- in this case, spew.Dump it out
	// 	```
	//
	first, ok := v[0].(string)
	if !ok {
		l.actualLog.Print(spew.Sdump(v...))
		return
	}

	if strings.ContainsRune(first, '%') {
		l.actualLog.Printf(first, v[1:]...)
	} else {
		l.actualLog.Print(v...)
	}
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
func NewLogger(level level, writeOut bool) Ourlog {
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
