package todoist

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
