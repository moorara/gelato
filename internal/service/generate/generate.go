package generate

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/moorara/gelato/internal/log"
	"github.com/moorara/gelato/internal/service/generate/compile"
	"github.com/moorara/gelato/internal/service/io"
)

const (
	genDir  = ".gen"
	mainPkg = "main"
)

type (
	compiler interface {
		Compile(string, *ast.File) *ast.File
	}

	// Generator generates test helpers (mocks, factories, builders, etc.) for a Go application.
	Generator struct {
		logger   *log.ColorfulLogger
		compiler compiler
	}
)

// New creates a new generator.
func New(level log.Level) *Generator {
	logger := log.NewColorful(level)

	return &Generator{
		logger:   logger,
		compiler: compile.New(logger),
	}
}

// Generate generates test helpers (mocks, factories, builders, etc.) for a Go application.
func (g *Generator) Generate(path string) error {
	// Sanitize the path
	if _, err := os.Stat(path); err != nil {
		return err
	}

	g.logger.White.Infof("Generating ...")

	module, err := io.GoModule(path)
	if err != nil {
		return err
	}

	return io.PackageDirectories(path, ".", func(basePath, relPath string) error {
		pkgPath := filepath.Join(module, relPath)
		pkgDir := filepath.Join(basePath, relPath)
		testPkgDir := filepath.Join(basePath, genDir, relPath+"test")

		// Parse all Go packages and files in the currecnt directory
		g.logger.Cyan.Debugf("  Parsing directory: %s", pkgDir)
		fset := token.NewFileSet()
		pkgs, err := parser.ParseDir(fset, pkgDir, nil, parser.AllErrors)
		if err != nil {
			return err
		}
		g.logger.Cyan.Tracef("  Directory parsed: %s", pkgDir)

		// Skip the main package
		if _, ok := pkgs[mainPkg]; ok {
			return nil
		}

		// Visit all parsed Go files in the currecnt directory
		for _, pkg := range pkgs {
			g.logger.Magenta.Debugf("     Package: %s", pkg.Name)
			for name, file := range pkg.Files {
				if !strings.HasSuffix(name, "_test.go") {
					g.logger.Green.Debugf("      File: %s", name)

					// Generate a new file for test helpers
					newFile := g.compiler.Compile(pkgPath, file)
					newFilePath := filepath.Join(testPkgDir, filepath.Base(name))
					if err := io.WriteASTFile(newFilePath, newFile, fset); err != nil {
						return err
					}

					g.logger.Green.Debugf("        File written: %s", newFilePath)
				}
			}
		}

		return nil
	})
}
