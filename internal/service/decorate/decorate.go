package decorate

import (
	"bufio"
	"bytes"
	"errors"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/imports"

	"github.com/moorara/gelato/internal/service/decorate/modifier"
	"github.com/moorara/gelato/internal/log"
)

const (
	goModFile     = "go.mod"
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
	mainModifier interface {
		Modify(string, string, ast.Node) ast.Node
	}

	genericModifier interface {
		Modify(string, string, ast.Node) ast.Node
	}

	// Decorator decorates a Go application.
	Decorator struct {
		decoratedDir    string
		logger          *log.ColorfulLogger
		mainModifier    mainModifier
		genericModifier genericModifier
	}
)

// New creates a new decorator.
func New(decoratedDir string, level log.Level) *Decorator {
	logger := log.NewColorful(level)

	return &Decorator{
		decoratedDir:    decoratedDir,
		logger:          logger,
		mainModifier:    modifier.NewMain(4, logger),
		genericModifier: modifier.NewGeneric(4, logger),
	}
}

func (d *Decorator) directories(basePath, relPath string, visit func(string, string) error) error {
	if err := visit(basePath, relPath); err != nil {
		return err
	}

	dir := filepath.Join(basePath, relPath)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() && file.Name() != d.decoratedDir {
			subdir := filepath.Join(relPath, file.Name())
			if err := d.directories(basePath, subdir, visit); err != nil {
				return err
			}
		}
	}

	return nil
}

// Decorate decorates a Go application.
func (d *Decorator) Decorate(path string) error {
	// Sanitize the path
	if _, err := os.Stat(path); err != nil {
		return err
	}

	d.logger.White.Infof("Decorating ...")

	module, err := getGoModule(path)
	if err != nil {
		return err
	}

	return d.directories(path, ".", func(basePath, relPath string) error {
		pkgDir := filepath.Join(basePath, relPath)
		newDir := filepath.Join(basePath, d.decoratedDir, relPath)

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
						d.mainModifier.Modify(module, d.decoratedDir, file)
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
