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
			l := NewColored(None, color.New(color.FgWhite))

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
