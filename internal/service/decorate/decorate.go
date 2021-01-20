package decorate

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/imports"

	"github.com/moorara/gelato/internal/log"
	"github.com/moorara/gelato/internal/service/decorate/modifier"
	"github.com/moorara/gelato/internal/service/io"
)

const (
	decoratedDir  = ".build"
	mainPkg       = "main"
	handlerPkg    = "handler"
	controllerPkg = "controller"
	gatewayPkg    = "gateway"
	repositoryPkg = "repository"
)

func isMainPackage(pkgs map[string]*ast.Package) bool {
	_, exist := pkgs[mainPkg]
	return exist
}

func isGenericPackage(pkgPath string) bool {
	return strings.HasSuffix(pkgPath, "/"+handlerPkg) || strings.Contains(pkgPath, "/"+handlerPkg+"/") ||
		strings.HasSuffix(pkgPath, "/"+controllerPkg) || strings.Contains(pkgPath, "/"+controllerPkg+"/") ||
		strings.HasSuffix(pkgPath, "/"+gatewayPkg) || strings.Contains(pkgPath, "/"+gatewayPkg+"/") ||
		strings.HasSuffix(pkgPath, "/"+repositoryPkg) || strings.Contains(pkgPath, "/"+repositoryPkg+"/")
}

type (
	mainModifier interface {
		Modify(string, string, ast.Node) ast.Node
	}

	genericModifier interface {
		Modify(string, string, ast.Node) ast.Node
	}

	// Decorator decorates a Go application.
	Decorator struct {
		logger          *log.ColorfulLogger
		mainModifier    mainModifier
		genericModifier genericModifier
	}
)

// New creates a new decorator.
func New(level log.Level) *Decorator {
	logger := log.NewColorful(level)

	return &Decorator{
		logger:          logger,
		mainModifier:    modifier.NewMain(logger),
		genericModifier: modifier.NewGeneric(logger),
	}
}

// Decorate decorates a Go application.
func (d *Decorator) Decorate(path string) error {
	// Sanitize the path
	if _, err := os.Stat(path); err != nil {
		return err
	}

	d.logger.White.Infof("Decorating ...")

	module, err := io.GoModule(path)
	if err != nil {
		return err
	}

	return io.PackageDirectories(path, ".", func(basePath, relPath string) error {
		pkgDir := filepath.Join(basePath, relPath)
		newDir := filepath.Join(basePath, decoratedDir, relPath)

		// Parse all Go packages and files in the currecnt directory
		d.logger.Cyan.Debugf("  Parsing directory: %s", pkgDir)
		fset := token.NewFileSet()
		pkgs, err := parser.ParseDir(fset, pkgDir, nil, parser.AllErrors)
		if err != nil {
			return err
		}
		d.logger.Cyan.Tracef("  Directory parsed: %s", pkgDir)

		// Skip the directory if it does not need decoration
		if !isMainPackage(pkgs) && !isGenericPackage(pkgDir) {
			return nil
		}

		// Creating a new directory for the decorated package
		if err := os.MkdirAll(newDir, os.ModePerm); err != nil {
			return err
		}
		d.logger.Blue.Tracef("  Directory created: %s", newDir)

		// Visit all parsed Go files in the currecnt directory
		for _, pkg := range pkgs {
			d.logger.Magenta.Debugf("     Package: %s", pkg.Name)
			for name, file := range pkg.Files {
				if !strings.HasSuffix(name, "_test.go") {
					d.logger.Green.Debugf("      File: %s", name)

					// Visit all nodes in the current file AST
					switch {
					case isMainPackage(pkgs):
						d.mainModifier.Modify(module, decoratedDir, file)
					case isGenericPackage(pkgDir):
						d.genericModifier.Modify(module, relPath, file)
					}

					buf := new(bytes.Buffer)
					if err := format.Node(buf, fset, file); err != nil {
						return err
					}

					// Format the modified Go file
					newName := filepath.Join(newDir, filepath.Base(name))
					b, err := imports.Process(newName, buf.Bytes(), &imports.Options{
						TabWidth:  8,
						TabIndent: true,
						Comments:  true,
						Fragment:  true,
					})

					if err != nil {
						return err
					}

					// Write the Go file to disk
					f, err := os.Create(newName)
					if err != nil {
						return err
					}

					if _, err := f.Write(b); err != nil {
						return err
					}

					d.logger.Green.Debugf("      File written: %s", newName)
				}
			}
			d.logger.White.Infof("  Decorated %s package", pkg.Name)
		}

		return nil
	})
}
