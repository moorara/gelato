package log

import (
	"testing"

	"github.com/moorara/color"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		name   string
		level  Level
		format string
		args   []interface{}
	}{
		{
			name:   "Trace",
			level:  Trace,
			format: "foo: %s",
			args:   []interface{}{"bar"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l := New(None)

			l.SetLevel(tc.level)
			level := l.GetLevel()
			assert.Equal(t, tc.level, level)

			l.Tracef(tc.format, tc.args...)
			l.Debugf(tc.format, tc.args...)
			l.Infof(tc.format, tc.args...)
			l.Warnf(tc.format, tc.args...)
			l.Errorf(tc.format, tc.args...)
		})
	}
}

func TestColoredLogger(t *testing.T) {
	tests := []struct {
		name   string
		level  Level
		color  *color.Color
		format string
		args   []interface{}
	}{
		{
			name:   "Trace",
			level:  Trace,
			color:  color.New(color.FgWhite),
			format: "foo: %s",
			args:   []interface{}{"bar"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l := NewColored(None, tc.color)

			l.SetLevel(tc.level)
			level := l.GetLevel()
			assert.Equal(t, tc.level, level)

			l.Tracef(tc.format, tc.args...)
			l.Debugf(tc.format, tc.args...)
			l.Infof(tc.format, tc.args...)
			l.Warnf(tc.format, tc.args...)
			l.Errorf(tc.format, tc.args...)
		})
	}
}

func TestColorfulLogger(t *testing.T) {
	tests := []struct {
		name   string
		level  Level
		format string
		args   []interface{}
	}{
		{
			name:   "Trace",
			level:  Trace,
			format: "foo: %s",
			args:   []interface{}{"bar"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l := NewColorful(None)

			l.SetLevel(tc.level)

			l.Red.Tracef(tc.format, tc.args...)
			l.Red.Debugf(tc.format, tc.args...)
			l.Red.Infof(tc.format, tc.args...)
			l.Red.Warnf(tc.format, tc.args...)
			l.Red.Errorf(tc.format, tc.args...)

			l.Green.Tracef(tc.format, tc.args...)
			l.Green.Debugf(tc.format, tc.args...)
			l.Green.Infof(tc.format, tc.args...)
			l.Green.Warnf(tc.format, tc.args...)
			l.Green.Errorf(tc.format, tc.args...)

			l.Yellow.Tracef(tc.format, tc.args...)
			l.Yellow.Debugf(tc.format, tc.args...)
			l.Yellow.Infof(tc.format, tc.args...)
			l.Yellow.Warnf(tc.format, tc.args...)
			l.Yellow.Errorf(tc.format, tc.args...)

			l.Blue.Tracef(tc.format, tc.args...)
			l.Blue.Debugf(tc.format, tc.args...)
			l.Blue.Infof(tc.format, tc.args...)
			l.Blue.Warnf(tc.format, tc.args...)
			l.Blue.Errorf(tc.format, tc.args...)

			l.Magenta.Tracef(tc.format, tc.args...)
			l.Magenta.Debugf(tc.format, tc.args...)
			l.Magenta.Infof(tc.format, tc.args...)
			l.Magenta.Warnf(tc.format, tc.args...)
			l.Magenta.Errorf(tc.format, tc.args...)

			l.Cyan.Tracef(tc.format, tc.args...)
			l.Cyan.Debugf(tc.format, tc.args...)
			l.Cyan.Infof(tc.format, tc.args...)
			l.Cyan.Warnf(tc.format, tc.args...)
			l.Cyan.Errorf(tc.format, tc.args...)

			l.White.Tracef(tc.format, tc.args...)
			l.White.Debugf(tc.format, tc.args...)
			l.White.Infof(tc.format, tc.args...)
			l.White.Warnf(tc.format, tc.args...)
			l.White.Errorf(tc.format, tc.args...)
		})
	}
}
