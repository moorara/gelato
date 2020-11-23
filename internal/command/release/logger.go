package release

import (
	"sync"

	"github.com/mitchellh/cli"
	"github.com/moorara/changelog/log"
	"github.com/moorara/color"
)

const indent = "  "

// logger implements the github.com/moorara/changelog/log.Logger interface.
type logger struct {
	sync.Mutex
	ui     cli.Ui
	colors struct {
		info  *color.Color
		warn  *color.Color
		err   *color.Color
		fatal *color.Color
	}
}

func newLogger(ui cli.Ui) *logger {
	cyan := color.New(color.FgCyan)
	yellow := color.New(color.FgYellow)
	magenta := color.New(color.FgMagenta)
	red := color.New(color.FgRed)

	l := &logger{
		ui: ui,
	}

	l.colors.info = cyan
	l.colors.warn = yellow
	l.colors.err = magenta
	l.colors.fatal = red

	return l
}

func (l *logger) ChangeVerbosity(v log.Verbosity) {
	// Noop
}

func (l *logger) Debug(v ...interface{}) {
	// Noop
}

func (l *logger) Debugf(format string, v ...interface{}) {
	// Noop
}

func (l *logger) Info(v ...interface{}) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	msg := indent + l.colors.info.Sprint(v...)
	l.ui.Output(msg)
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	msg := indent + l.colors.info.Sprintf(format, v...)
	l.ui.Output(msg)
}

func (l *logger) Warn(v ...interface{}) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	msg := indent + l.colors.warn.Sprint(v...)
	l.ui.Output(msg)
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	msg := indent + l.colors.warn.Sprintf(format, v...)
	l.ui.Output(msg)
}

func (l *logger) Error(v ...interface{}) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	msg := indent + l.colors.err.Sprint(v...)
	l.ui.Output(msg)
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	msg := indent + l.colors.err.Sprintf(format, v...)
	l.ui.Output(msg)
}

func (l *logger) Fatal(v ...interface{}) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	msg := indent + l.colors.fatal.Sprint(v...)
	l.ui.Output(msg)
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	msg := indent + l.colors.fatal.Sprintf(format, v...)
	l.ui.Output(msg)
}
