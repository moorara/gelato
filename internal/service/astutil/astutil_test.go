package astutil

import (
	"go/ast"
	"go/token"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteFile(t *testing.T) {
	mainFile := &ast.File{
		Name: &ast.Ident{
			Name: "main",
		},
	}

	tests := []struct {
		name          string
		path          string
		file          *ast.File
		expectedError string
	}{
		{
			name: "InvalidFile",
			path: "./foo.go",
			file: &ast.File{
				Name: &ast.Ident{},
			},
			expectedError: "expected 'IDENT', found 'EOF'",
		},
		{
			name:          "InvalidPath",
			path:          ".",
			file:          mainFile,
			expectedError: "open .: is a directory",
		},
		{
			name:          "Success",
			path:          "./foo.go",
			file:          mainFile,
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fset := token.NewFileSet()

			err := WriteFile(tc.path, tc.file, fset)

			// Cleanup
			defer os.Remove(tc.path)
			defer os.Remove(debugFile)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Contains(t, err.Error(), tc.expectedError)
			}
		})
	}
}
