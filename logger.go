package golang_todoist_api

import (
	"fmt"
)

type logger interface {
	Output(int, string) error
}

type ilogger interface {
	logger
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})
}

type Debug interface {
	Debug() bool

	Debugf(format string, v ...interface{})
	Debugln(v ...interface{})
}

type internalLog struct {
	logger
}

func (t internalLog) Println(v ...interface{}) {
	err := t.Output(2, fmt.Sprintln(v...))
	if err != nil {
		return
	}
}

func (t internalLog) Printf(format string, v ...interface{}) {
	err := t.Output(2, fmt.Sprintf(format, v...))
	if err != nil {
		return
	}
}

func (t internalLog) Print(v ...interface{}) {
	err := t.Output(2, fmt.Sprint(v...))
	if err != nil {
		return
	}
}

type discard struct{}

func (t discard) Debug() bool {
	return false
}

func (t discard) Debugf(format string, v ...interface{}) {}

func (t discard) Debugln(v ...interface{}) {}
