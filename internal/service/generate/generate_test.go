package generate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/log"
)

func TestNew(t *testing.T) {
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
			g := New(tc.level)

			assert.NotNil(t, g)
			assert.NotNil(t, g.logger)
		})
	}
}

func TestGenerator_Generate(t *testing.T) {
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
		path          string
		mock          bool
		factory       bool
		expectedError string
	}{
		{
			name:          "PathNotExist",
			path:          "/invalid/path",
			mock:          true,
			factory:       true,
			expectedError: "stat /invalid/path: no such file or directory",
		},
		{
			name:          "Success_Mock",
			path:          "./test",
			mock:          true,
			factory:       false,
			expectedError: "",
		},
		{
			name:          "Success_Factory",
			path:          "./test",
			mock:          false,
			factory:       true,
			expectedError: "",
		},
		{
			name:          "Success_Mock_Factory",
			path:          "./test",
			mock:          true,
			factory:       true,
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := &Generator{
				logger: clogger,
			}

			// Clean-up
			defer os.RemoveAll(filepath.Join(tc.path, genDir))

			err := g.Generate(tc.path, tc.mock, tc.factory)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
