package io

import (
	"errors"
	"go/ast"
	"go/token"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			name:           "Success",
			path:           "./test/valid",
			expectedModule: "test",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			module, err := GoModule(tc.path)

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

func TestPackageDirectories(t *testing.T) {
	tests := []struct {
		name          string
		basePath      string
		relPath       string
		visit         visitFunc
		expectedError string
	}{
		{
			name:     "InvalidPath",
			basePath: "./invalid",
			relPath:  ".",
			visit: func(_, _ string) error {
				return nil
			},
			expectedError: "open invalid: no such file or directory",
		},
		{
			name:     "Success",
			basePath: "./test/valid",
			relPath:  ".",
			visit: func(_, _ string) error {
				return nil
			},
			expectedError: "",
		},
		{
			name:     "VisitFails_FirstTime",
			basePath: "./test/valid",
			relPath:  ".",
			visit: func(_, _ string) error {
				return errors.New("generic error")
			},
			expectedError: "generic error",
		},
		{
			name:     "VisitFails_SecondTime",
			basePath: "./test/valid",
			relPath:  ".",
			visit: func(_, relPath string) error {
				if relPath == "." {
					return nil
				}
				return errors.New("generic error")
			},
			expectedError: "generic error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := PackageDirectories(tc.basePath, tc.relPath, tc.visit)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestWriteASTFile(t *testing.T) {
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
			path: "./test/foo.go",
			file: &ast.File{
				Name: &ast.Ident{},
			},
			expectedError: "expected 'IDENT', found 'EOF'",
		},
		{
			name:          "InvalidPath",
			path:          "./test",
			file:          mainFile,
			expectedError: "open ./test: is a directory",
		},
		{
			name:          "Success",
			path:          "./test/foo.go",
			file:          mainFile,
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fset := token.NewFileSet()

			err := WriteASTFile(tc.path, tc.file, fset)

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
