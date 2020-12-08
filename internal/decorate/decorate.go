package decorate

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/moorara/gelato/internal/log"
)

const (
	decoratedDir = ".build"

	mainPkg       = "main"
	serverPkg     = "server"
	handlerPkg    = "handler"
	controllerPkg = "controller"
	gatewayPkg    = "gateway"
	repositoryPkg = "repository"
)

func isPackageDecoratable(pkgPath string) bool {
	return strings.HasSuffix(pkgPath, "/"+serverPkg) || strings.Contains(pkgPath, "/"+serverPkg+"/") ||
		strings.HasSuffix(pkgPath, "/"+handlerPkg) || strings.Contains(pkgPath, "/"+handlerPkg+"/") ||
		strings.HasSuffix(pkgPath, "/"+controllerPkg) || strings.Contains(pkgPath, "/"+controllerPkg+"/") ||
		strings.HasSuffix(pkgPath, "/"+gatewayPkg) || strings.Contains(pkgPath, "/"+gatewayPkg+"/") ||
		strings.HasSuffix(pkgPath, "/"+repositoryPkg) || strings.Contains(pkgPath, "/"+repositoryPkg+"/")
}

func directories(basePath, relPath string, visit func(string, string) error) error {
	if err := visit(basePath, relPath); err != nil {
		return err
	}

	dir := filepath.Join(basePath, relPath)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() && file.Name() != decoratedDir {
			subdir := filepath.Join(relPath, file.Name())
			if err := directories(basePath, subdir, visit); err != nil {
				return err
			}
		}
	}

	return nil
}

// Decorator decorates a Go application.
type Decorator struct {
	logger     *log.ColorfulLogger
	gateway    ast.Visitor
	repository ast.Visitor
	controller ast.Visitor
	handler    ast.Visitor
	server     ast.Visitor
}

// New creates a new decorator.
func New() *Decorator {
	return &Decorator{
		logger:     log.NewColorful(log.None),
		gateway:    newStructDecorator(gatewayPkg),
		repository: newStructDecorator(repositoryPkg),
		controller: newStructDecorator(controllerPkg),
		handler:    newStructDecorator(handlerPkg),
		server:     newStructDecorator(serverPkg),
	}
}

// Decorate decorates a Go application.
func (d *Decorator) Decorate(level log.Level, path string) error {
	// Update logging level
	d.logger.SetLevel(level)

	// Sanitize the path
	if _, err := os.Stat(path); err != nil {
		return err
	}

	d.logger.White.Infof("Decorating ...")

	visitor := &visitor{
		depth:  4,
		logger: d.logger,
	}

	return directories(path, ".", func(basePath, relPath string) error {
		newDir := filepath.Join(basePath, decoratedDir, relPath)
		pkgDir := filepath.Join(basePath, relPath)

		// Parse all Go packages and files in the currecnt directory
		d.logger.Cyan.Debugf("  Parsing directory: %s", pkgDir)
		fset := token.NewFileSet()
		pkgs, err := parser.ParseDir(fset, pkgDir, nil, parser.AllErrors)
		if err != nil {
			return err
		}
		d.logger.Cyan.Tracef("  Directory parsed: %s", pkgDir)

		// Skip the directory if it is not the main package or it does not need decoration
		if _, exist := pkgs[mainPkg]; !exist && !isPackageDecoratable(pkgDir) {
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
					ast.Walk(visitor, file)

					// Write the modified Go file to disk
					newName := filepath.Join(newDir, filepath.Base(name))
					f, err := os.Create(newName)
					if err != nil {
						return err
					}

					if err := format.Node(f, fset, file); err != nil {
						return err
					}

					d.logger.Green.Debugf("      File written: %s", newName)
				}
			}
		}

		return nil
	})
}
