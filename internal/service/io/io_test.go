package io

import (
	"errors"
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
