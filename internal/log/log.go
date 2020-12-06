package log

import (
	"log"
	"os"
	"sync"

	"github.com/moorara/color"
)

// Level is the logging verbosity level.
type Level int

const (
	// Trace shows logs in all levels.
	Trace Level = iota
	// Debug shows logs in Debug, Info, Warn, Error, and Fatal levels.
	Debug
	// Info shows logs in Info, Warn, Error, and Fatal levels.
	Info
	// Warn shows logs in Warn, Error, and Fatal levels.
	Warn
	// Error shows logs in Error and Fatal levels.
	Error
	// Fatal shows logs in Fatal level.
	Fatal
	// None does not show any logs.
	None
)

// Logger is the interface for a simple logger.
type Logger interface {
	GetLevel() Level
	SetLevel(l Level)
	Tracef(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}

// logger implements the Logger interface for logging to standard output.
type logger struct {
	sync.Mutex
	level  Level
	logger *log.Logger
}

// New creates a new logger.
func New(level Level) Logger {
	return &logger{
		level:  level,
		logger: log.New(os.Stdout, "", 0),
	}
}

func (l *logger) GetLevel() Level {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	return l.level
}

func (l *logger) SetLevel(level Level) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	l.level = level
}

func (l *logger) Tracef(format string, v ...interface{}) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	if l.level <= Trace {
		l.logger.Printf(format, v...)
	}
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	if l.level <= Debug {
		l.logger.Printf(format, v...)
	}
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	if l.level <= Info {
		l.logger.Printf(format, v...)
	}
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	if l.level <= Warn {
		l.logger.Printf(format, v...)
	}
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	if l.level <= Error {
		l.logger.Printf(format, v...)
	}
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	if l.level <= Fatal {
		l.logger.Printf(format, v...)
	}
}

// coloredLogger implements the Logger interface for logging to standard output with a color.
type coloredLogger struct {
	logger *logger
	color  *color.Color
}

// NewColored creates a new colored logger.
func NewColored(level Level, color *color.Color) Logger {
	return &coloredLogger{
		logger: &logger{
			level:  level,
			logger: log.New(os.Stdout, "", 0),
		},
		color: color,
	}
}

func (l *coloredLogger) GetLevel() Level {
	return l.logger.GetLevel()
}

func (l *coloredLogger) SetLevel(level Level) {
	l.logger.SetLevel(level)
}

func (l *coloredLogger) Tracef(format string, v ...interface{}) {
	msg := l.color.Sprintf(format, v...)
	l.logger.Tracef(msg)
}

func (l *coloredLogger) Debugf(format string, v ...interface{}) {
	msg := l.color.Sprintf(format, v...)
	l.logger.Debugf(msg)
}

func (l *coloredLogger) Infof(format string, v ...interface{}) {
	msg := l.color.Sprintf(format, v...)
	l.logger.Infof(msg)
}

func (l *coloredLogger) Warnf(format string, v ...interface{}) {
	msg := l.color.Sprintf(format, v...)
	l.logger.Warnf(msg)
}

func (l *coloredLogger) Errorf(format string, v ...interface{}) {
	msg := l.color.Sprintf(format, v...)
	l.logger.Errorf(msg)
}

func (l *coloredLogger) Fatalf(format string, v ...interface{}) {
	msg := l.color.Sprintf(format, v...)
	l.logger.Fatalf(msg)
}
