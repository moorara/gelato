package decorate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/log"
)

func TestNew(t *testing.T) {
	d := New()

	assert.NotNil(t, d)
}

func TestDecorator_Decorate(t *testing.T) {
	tests := []struct {
		name          string
		level         log.Level
		path          string
		expectedError string
	}{
		{
			name:          "PathNotExist",
			level:         log.None,
			path:          "/invalid/path",
			expectedError: "stat /invalid/path: no such file or directory",
		},
		{
			name:          "Success",
			level:         log.None,
			path:          "./test",
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := New()

			// Clean-up
			defer os.RemoveAll(filepath.Join(tc.path, decoratedDir))

			err := d.Decorate(tc.level, tc.path)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
