package generate

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/moorara/gelato/internal/log"
	"github.com/moorara/gelato/internal/service/io"
)

const (
	genDir  = ".gen"
	mainPkg = "main"
)

// Generator generates mocks and factories for a Go application.
type Generator struct {
	logger *log.ColorfulLogger
}

// New creates a new generator.
func New(level log.Level) *Generator {
	logger := log.NewColorful(level)

	return &Generator{
		logger: logger,
	}
}

// Generate generates mocks and factories for a Go application.
func (g *Generator) Generate(path string, mock, factory bool) error {
	// Sanitize the path
	if _, err := os.Stat(path); err != nil {
		return err
	}

	g.logger.White.Infof("Generating ...")

	return io.PackageDirectories(path, ".", func(basePath, relPath string) error {
		pkgDir := filepath.Join(basePath, relPath)
		mockPkgDir := filepath.Join(basePath, genDir, "mock", relPath)
		factoryPkgDir := filepath.Join(basePath, genDir, "factory", relPath)

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

		// Creating a new directory for the mock package
		if mock {
			if err := os.MkdirAll(mockPkgDir, os.ModePerm); err != nil {
				return err
			}
			g.logger.Blue.Tracef("  Directory created: %s", mockPkgDir)
		}

		// Creating a new directory for the factory package
		if factory {
			if err := os.MkdirAll(factoryPkgDir, os.ModePerm); err != nil {
				return err
			}
			g.logger.Blue.Tracef("  Directory created: %s", factoryPkgDir)
		}

		// Visit all parsed Go files in the currecnt directory
		for _, pkg := range pkgs {
			g.logger.Magenta.Debugf("     Package: %s", pkg.Name)
			for name := range pkg.Files {
				if !strings.HasSuffix(name, "_test.go") {
					g.logger.Green.Debugf("      File: %s", name)
				}
			}
		}

		return nil
	})
}
