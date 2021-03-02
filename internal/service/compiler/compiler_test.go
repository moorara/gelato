package compiler

import (
	"errors"
	goast "go/ast"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/log"
)

func TestNew(t *testing.T) {
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
		name      string
		logger    *log.ColorfulLogger
		consumers []*Consumer
	}{
		{
			name:      "OK",
			logger:    clogger,
			consumers: []*Consumer{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := New(tc.logger, tc.consumers...)

			assert.NotNil(t, c)
			assert.NotNil(t, c.parser)
			assert.Equal(t, tc.logger, c.parser.logger)
			assert.Equal(t, tc.consumers, c.parser.consumers)
		})
	}
}

func TestCompiler_Compile(t *testing.T) {
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
		consumers     []*Consumer
		path          string
		opts          ParseOptions
		expectedError string
	}{
		{
			name: "Success_SkipPackages",
			consumers: []*Consumer{
				{
					Name:    "tester",
					Package: func(*PackageInfo, *goast.Package) bool { return false },
				},
			},
			path: "./test/valid",
			opts: ParseOptions{
				SkipTestFiles: true,
			},
			expectedError: "",
		},
		{
			name: "Success_SkipFiles",
			consumers: []*Consumer{
				{
					Name:    "tester",
					Package: func(*PackageInfo, *goast.Package) bool { return true },
					FilePre: func(*FileInfo, *goast.File) bool { return false },
				},
			},
			path: "./test/valid",
			opts: ParseOptions{
				SkipTestFiles: true,
			},
			expectedError: "",
		},
		{
			name: "Success",
			consumers: []*Consumer{
				{
					Name:      "tester",
					Package:   func(*PackageInfo, *goast.Package) bool { return true },
					FilePre:   func(*FileInfo, *goast.File) bool { return true },
					Import:    func(*FileInfo, *goast.ImportSpec) {},
					Struct:    func(*TypeInfo, *goast.StructType) {},
					Interface: func(*TypeInfo, *goast.InterfaceType) {},
					FuncType:  func(*TypeInfo, *goast.FuncType) {},
					FuncDecl:  func(*FuncInfo, *goast.FuncType, *goast.BlockStmt) {},
					FilePost:  func(*FileInfo, *goast.File) error { return nil },
				},
			},
			path: "./test/valid",
			opts: ParseOptions{
				SkipTestFiles: true,
			},
			expectedError: "",
		},
		{
			name: "FilePostFails",
			consumers: []*Consumer{
				{
					Name:      "tester",
					Package:   func(*PackageInfo, *goast.Package) bool { return true },
					FilePre:   func(*FileInfo, *goast.File) bool { return true },
					Import:    func(*FileInfo, *goast.ImportSpec) {},
					Struct:    func(*TypeInfo, *goast.StructType) {},
					Interface: func(*TypeInfo, *goast.InterfaceType) {},
					FuncType:  func(*TypeInfo, *goast.FuncType) {},
					FuncDecl:  func(*FuncInfo, *goast.FuncType, *goast.BlockStmt) {},
					FilePost:  func(*FileInfo, *goast.File) error { return errors.New("file error") },
				},
			},
			path: "./test/valid",
			opts: ParseOptions{
				SkipTestFiles: true,
			},
			expectedError: "file error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := New(clogger, tc.consumers...)

			err := c.Compile(tc.path, tc.opts)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
