package edit

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/log"
)

func TestNewEditor(t *testing.T) {
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
			editor := NewEditor(tc.level)

			assert.NotNil(t, editor)
		})
	}
}

func TestEditor_Remove(t *testing.T) {
	assert.NoError(t, os.Mkdir("./temp", 0755))
	defer os.RemoveAll("./temp")

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
		globs         []string
		expectedError string
	}{
		{
			name:          "NoGlob",
			globs:         []string{},
			expectedError: "",
		},
		{
			name:          "NoMatch",
			globs:         []string{""},
			expectedError: "",
		},
		{
			name:          "InvalidPattern",
			globs:         []string{"\\"},
			expectedError: "syntax error in pattern",
		},
		{
			name:          "Success",
			globs:         []string{"./temp"},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			editor := &Editor{
				logger: clogger,
			}

			err := editor.Remove(tc.globs...)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestEditor_Move(t *testing.T) {
	err := os.Mkdir("./temp", 0755)
	assert.NoError(t, err)
	_, err = os.Create("./temp/foo")
	assert.NoError(t, err)
	defer os.RemoveAll("./temp")

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
		specs         []MoveSpec
		expectedError string
	}{
		{
			name:          "NoSpec",
			specs:         []MoveSpec{},
			expectedError: "",
		},
		{
			name: "InvalidSpec",
			specs: []MoveSpec{
				{
					Src:  "./foo",
					Dest: "./bar",
				},
			},
			expectedError: "rename ./foo ./bar: no such file or directory",
		},
		{
			name: "Success",
			specs: []MoveSpec{
				{
					Src:  "./temp/foo",
					Dest: "./temp/bar",
				},
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			editor := &Editor{
				logger: clogger,
			}

			err := editor.Move(tc.specs...)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestEditor_ReplaceInDir(t *testing.T) {
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
		root          string
		specs         []ReplaceSpec
		expectedError string
	}{
		{
			name:          "DirectoryNotExist",
			root:          "./foo",
			specs:         []ReplaceSpec{},
			expectedError: "lstat ./foo: no such file or directory",
		},
		{
			name: "Success",
			root: "./test",
			specs: []ReplaceSpec{
				{
					PathRE: regexp.MustCompile(`\.txt$`),
					OldRE:  regexp.MustCompile(`foo`),
					New:    "bar",
				},
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			editor := &Editor{
				logger: clogger,
			}

			err := editor.ReplaceInDir(tc.root, tc.specs...)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
