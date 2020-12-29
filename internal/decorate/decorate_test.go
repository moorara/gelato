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

func TestIsMainPackage(t *testing.T) {
	tests := []struct {
		name       string
		pkgs       map[string]*ast.Package
		expectedOK bool
	}{
		{
			name: "MainPackage",
			pkgs: map[string]*ast.Package{
				"main": {},
			},
			expectedOK: true,
		},
		{
			name: "NoneMainPackage",
			pkgs: map[string]*ast.Package{
				"cmd": {},
			},
			expectedOK: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ok := isMainPackage(tc.pkgs)

			assert.Equal(t, tc.expectedOK, ok)
		})
	}
}

func TestIsGenericPackage(t *testing.T) {
	tests := []struct {
		name       string
		pkgPath    string
		expectedOK bool
	}{
		{
			name:       "main",
			pkgPath:    "github.com/octocat/service",
			expectedOK: false,
		},
		{
			name:       "internal",
			pkgPath:    "github.com/octocat/service/internal",
			expectedOK: false,
		},
		{
			name:       "controller",
			pkgPath:    "github.com/octocat/service/internal/controller",
			expectedOK: true,
		},
		{
			name:       "controller/sub",
			pkgPath:    "github.com/octocat/service/internal/controller/sub",
			expectedOK: true,
		},
		{
			name:       "gateway",
			pkgPath:    "github.com/octocat/service/internal/gateway",
			expectedOK: true,
		},
		{
			name:       "gateway/sub",
			pkgPath:    "github.com/octocat/service/internal/gateway/sub",
			expectedOK: true,
		},
		{
			name:       "handler",
			pkgPath:    "github.com/octocat/service/internal/handler",
			expectedOK: true,
		},
		{
			name:       "handler/sub",
			pkgPath:    "github.com/octocat/service/internal/handler/sub",
			expectedOK: true,
		},
		{
			name:       "repository",
			pkgPath:    "github.com/octocat/service/internal/repository",
			expectedOK: true,
		},
		{
			name:       "repository/sub",
			pkgPath:    "github.com/octocat/service/internal/repository/sub",
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ok := isGenericPackage(tc.pkgPath)

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

func TestGetGoModule(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedModule string
		expectedError  string
	}{
		{
			name:          "NoModFile",
			path:          "./test",
			expectedError: "open test/go.mod: no such file or directory",
		},
		{
			name:          "InvalidModFile",
			path:          "./test/invalid",
			expectedError: "invalid go.mod file: no module name found",
		},
		{
			name:           "Success_Horizontal",
			path:           "./test/horizontal",
			expectedModule: "github.com/moorara/test/horizontal",
		},
		{
			name:           "Success_Vertical",
			path:           "./test/vertical",
			expectedModule: "github.com/moorara/test/vertical",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			module, err := getGoModule(tc.path)

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedModule, module)
			} else {
				assert.Empty(t, module)
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
		name            string
		visitor         ast.Visitor
		mainModifier    *MockMainModifier
		genericModifier *MockGenericModifier
		level           log.Level
		path            string
		expectedError   string
	}{
		{
			name:          "PathNotExist",
			level:         log.None,
			path:          "/invalid/path",
			expectedError: "stat /invalid/path: no such file or directory",
		},
		{
			name: "Success_Horizontal",
			mainModifier: &MockMainModifier{
				ModifyMocks: []MainModifyMock{
					{
						OutNode: &ast.File{},
					},
				},
			},
			genericModifier: &MockGenericModifier{
				ModifyMocks: []GenericModifyMock{},
			},
			level:         log.None,
			path:          "./test/horizontal",
			expectedError: "",
		},
		{
			name: "Success_Vertical",
			mainModifier: &MockMainModifier{
				ModifyMocks: []MainModifyMock{
					{
						OutNode: &ast.File{},
					},
				},
			},
			genericModifier: &MockGenericModifier{
				ModifyMocks: []GenericModifyMock{},
			},
			level:         log.None,
			path:          "./test/vertical",
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := &Decorator{
				logger:          clogger,
				mainModifier:    tc.mainModifier,
				genericModifier: tc.genericModifier,
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
