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
			assert.NotNil(t, g.compiler)
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

	testFile := &ast.File{
		Name: &ast.Ident{
			Name: "exampletest",
		},
	}

	tests := []struct {
		name          string
		compiler      *MockCompiler
		path          string
		expectedError string
	}{
		{
			name:          "PathNotExist",
			path:          "/invalid/path",
			expectedError: "stat /invalid/path: no such file or directory",
		},
		{
			name: "Success",
			compiler: &MockCompiler{
				CompileMocks: []CompileMock{
					{OutFile: testFile},
				},
			},
			path:          "./test",
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := &Generator{
				logger:   clogger,
				compiler: tc.compiler,
			}

			// Clean-up
			defer os.RemoveAll(filepath.Join(tc.path, genDir))

			err := g.Generate(tc.path)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
