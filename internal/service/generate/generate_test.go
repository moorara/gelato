package generate

import (
	"go/ast"
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
			assert.NotNil(t, g.factory)
			assert.NotNil(t, g.mock)
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

	factoryFile := &ast.File{
		Name: &ast.Ident{
			Name: "examplefactory",
		},
	}

	mockFile := &ast.File{
		Name: &ast.Ident{
			Name: "examplemock",
		},
	}

	tests := []struct {
		name             string
		factoryGenerator *MockGenerator
		mockGenerator    *MockGenerator
		path             string
		mock             bool
		factory          bool
		expectedError    string
	}{
		{
			name:          "PathNotExist",
			path:          "/invalid/path",
			mock:          true,
			factory:       true,
			expectedError: "stat /invalid/path: no such file or directory",
		},
		{
			name: "Success_Factory",
			factoryGenerator: &MockGenerator{
				GenerateMocks: []GenerateMock{
					{OutFile: factoryFile},
				},
			},
			path:          "./test",
			mock:          false,
			factory:       true,
			expectedError: "",
		},
		{
			name: "Success_Mock",
			mockGenerator: &MockGenerator{
				GenerateMocks: []GenerateMock{
					{OutFile: mockFile},
				},
			},
			path:          "./test",
			mock:          true,
			factory:       false,
			expectedError: "",
		},
		{
			name: "Success_Factory_Mock",
			factoryGenerator: &MockGenerator{
				GenerateMocks: []GenerateMock{
					{OutFile: factoryFile},
				},
			},
			mockGenerator: &MockGenerator{
				GenerateMocks: []GenerateMock{
					{OutFile: mockFile},
				},
			},
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

			g.factory = tc.factoryGenerator
			g.mock = tc.mockGenerator

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
