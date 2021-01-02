package archive

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/moorara/gelato/internal/log"
	"github.com/stretchr/testify/assert"
)

func TestNewTarArchive(t *testing.T) {
	tests := []struct {
		name  string
		level log.Level
	}{
		{
			name:  "OK",
			level: log.None,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			arch := NewTarArchive(tc.level)

			assert.NotNil(t, arch)
		})
	}
}

func TestTarArchive_Extract(t *testing.T) {
	logger := log.New(log.None)
	clogger := &log.ColorfulLogger{
		Red:     logger,
		Green:   logger,
		Yellow:  logger,
		Blue:    logger,
		Magenta: logger,
		Cyan:    logger,
		White:   logger,
	}

	tests := []struct {
		name          string
		archFile      string
		f             Selector
		expectedError string
	}{
		{
			name:          "InvalidArchive",
			archFile:      "test/invalid.tar.gz",
			f:             nil,
			expectedError: "error on creating gzip reader: EOF",
		},
		{
			name:     "Success",
			archFile: "test/github.tar.gz",
			f: func(path string) (string, bool) {
				return path, true
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dest, err := ioutil.TempDir("", "gelato-test-*")
			assert.NoError(t, err)
			defer os.RemoveAll(dest)

			f, err := os.Open(tc.archFile)
			assert.NoError(t, err)

			arch := &TarArchive{
				logger: clogger,
			}

			err = arch.Extract(dest, f, tc.f)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
