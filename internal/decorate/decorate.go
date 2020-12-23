package decorate

import (
	"bufio"
	"errors"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/moorara/gelato/internal/decorate/modifier"
	"github.com/moorara/gelato/internal/log"
)

const (
	goModFile    = "go.mod"
	decoratedDir = ".build"

	mainPkg       = "main"
	handlerPkg    = "handler"
	controllerPkg = "controller"
	gatewayPkg    = "gateway"
	repositoryPkg = "repository"
)

func isPackageDecoratable(pkgPath string) bool {
	return strings.HasSuffix(pkgPath, "/"+handlerPkg) || strings.Contains(pkgPath, "/"+handlerPkg+"/") ||
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

func getGoModule(path string) (string, error) {
	f, err := os.Open(filepath.Join(path, goModFile))
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if line := scanner.Text(); strings.HasPrefix(line, "module ") {
			return strings.TrimPrefix(line, "module "), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", errors.New("invalid go.mod file: no module name found")
}

type (
	fileModifier interface {
		Modify(string, string, string, ast.Node) ast.Node
	}

	// Decorator decorates a Go application.
	Decorator struct {
		logger   *log.ColorfulLogger
		modifier fileModifier
	}
)

// New creates a new decorator.
func New() *Decorator {
	logger := log.NewColorful(log.None)

	return &Decorator{
		logger:   logger,
		modifier: modifier.NewFile(4, logger),
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

	module, err := getGoModule(path)
	if err != nil {
		return err
	}

	return directories(path, ".", func(basePath, relPath string) error {
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
		if !isPackageDecoratable(pkgDir) {
			return nil
		}

		if _, exist := pkgs[mainPkg]; exist {
			// TODO: main package requires a special decoration!
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
					d.modifier.Modify(module, decoratedDir, relPath, file)

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
