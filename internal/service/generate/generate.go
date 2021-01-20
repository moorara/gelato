package generate

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/moorara/gelato/internal/log"
	"github.com/moorara/gelato/internal/service/generate/factory"
	"github.com/moorara/gelato/internal/service/generate/mock"
	"github.com/moorara/gelato/internal/service/io"
)

const (
	genDir  = ".gen"
	mainPkg = "main"
)

type (
	generator interface {
		Generate(string, string, *ast.File) *ast.File
	}

	// Generator generates mocks and factories for a Go application.
	Generator struct {
		logger  *log.ColorfulLogger
		factory generator
		mock    generator
	}
)

// New creates a new generator.
func New(level log.Level) *Generator {
	logger := log.NewColorful(level)

	return &Generator{
		logger:  logger,
		factory: factory.NewGenerator(logger),
		mock:    mock.NewGenerator(logger),
	}
}

// Generate generates mocks and factories for a Go application.
func (g *Generator) Generate(path string, mock, factory bool) error {
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
		mockPkgDir := filepath.Join(basePath, genDir, relPath+"mock")
		factoryPkgDir := filepath.Join(basePath, genDir, relPath+"factory")

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

					// Generate a new file for factories
					if factory {
						factoryFile := g.factory.Generate(pkgPath, pkg.Name, file)
						factoryFilePath := filepath.Join(mockPkgDir, filepath.Base(name))
						if err := io.WriteASTFile(factoryFilePath, factoryFile, fset); err != nil {
							return err
						}
						g.logger.Green.Debugf("        Factories written: %s", factoryFilePath)
					}

					// Generate a new file for mocks
					if mock {
						mockFile := g.mock.Generate(pkgPath, pkg.Name, file)
						mockFilePath := filepath.Join(factoryPkgDir, filepath.Base(name))
						if err := io.WriteASTFile(mockFilePath, mockFile, fset); err != nil {
							return err
						}
						g.logger.Green.Debugf("        Mocks written: %s", mockFilePath)
					}
				}
			}
		}

		return nil
	})
}
