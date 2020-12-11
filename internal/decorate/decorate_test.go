package decorate

import (
	"errors"
	"go/ast"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/log"
)

func TestIsPackageDecoratable(t *testing.T) {
	tests := []struct {
		name       string
		pkgPath    string
		expectedOK bool
	}{
		{
			name:       "main",
			pkgPath:    "github.com/octocat/Hello-World",
			expectedOK: false,
		},
		{
			name:       "internal",
			pkgPath:    "github.com/octocat/Hello-World/internal",
			expectedOK: false,
		},
		{
			name:       "controller",
			pkgPath:    "github.com/octocat/Hello-World/internal/controller",
			expectedOK: true,
		},
		{
			name:       "controller/sub",
			pkgPath:    "github.com/octocat/Hello-World/internal/controller/sub",
			expectedOK: true,
		},
		{
			name:       "gateway",
			pkgPath:    "github.com/octocat/Hello-World/internal/gateway",
			expectedOK: true,
		},
		{
			name:       "gateway/sub",
			pkgPath:    "github.com/octocat/Hello-World/internal/gateway/sub",
			expectedOK: true,
		},
		{
			name:       "handler",
			pkgPath:    "github.com/octocat/Hello-World/internal/handler",
			expectedOK: true,
		},
		{
			name:       "handler/sub",
			pkgPath:    "github.com/octocat/Hello-World/internal/handler/sub",
			expectedOK: true,
		},
		{
			name:       "repository",
			pkgPath:    "github.com/octocat/Hello-World/internal/repository",
			expectedOK: true,
		},
		{
			name:       "repository/sub",
			pkgPath:    "github.com/octocat/Hello-World/internal/repository/sub",
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ok := isPackageDecoratable(tc.pkgPath)

			assert.Equal(t, tc.expectedOK, ok)
		})
	}
}

func TestDirectories(t *testing.T) {
	tests := []struct {
		name          string
		basePath      string
		relPath       string
		visit         func(string, string) error
		expectedError string
	}{
		{
			name:     "Success",
			basePath: "./test",
			relPath:  ".",
			visit: func(_, _ string) error {
				return nil
			},
			expectedError: "",
		},
		{
			name:     "VisitFails_FirstTime",
			basePath: "./test",
			relPath:  ".",
			visit: func(_, _ string) error {
				return errors.New("generic error")
			},
			expectedError: "generic error",
		},
		{
			name:     "VisitFails_SecondTime",
			basePath: "./test",
			relPath:  ".",
			visit: func(_, relPath string) error {
				if relPath == "." {
					return nil
				}
				return errors.New("generic error")
			},
			expectedError: "generic error",
		},
		{
			name:     "InvalidPath",
			basePath: "./invalid",
			relPath:  ".",
			visit: func(_, _ string) error {
				return nil
			},
			expectedError: "open invalid: no such file or directory",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := directories(tc.basePath, tc.relPath, tc.visit)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestNew(t *testing.T) {
	d := New()

	assert.NotNil(t, d)
}

func TestDecorator_Decorate(t *testing.T) {
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
		visitor       ast.Visitor
		modifier      *MockModifier
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
			name: "Success",
			modifier: &MockModifier{
				ModifyMocks: []ModifyMock{
					{
						OutNode: &ast.File{},
					},
					{
						OutNode: &ast.File{},
					},
				},
			},
			level:         log.None,
			path:          "./test",
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := &Decorator{
				logger:   clogger,
				modifier: tc.modifier,
			}

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
