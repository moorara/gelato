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

	"github.com/moorara/color"

	"github.com/moorara/gelato/internal/log"
)

const (
	decoratedDir  = ".build"
	gatewayPkg    = "gateway"
	repositoryPkg = "repository"
	controllerPkg = "controller"
	handlerPkg    = "handler"
	serverPkg     = "server"
)

type loggers struct {
	red     log.Logger
	green   log.Logger
	yellow  log.Logger
	blue    log.Logger
	magenta log.Logger
	cyan    log.Logger
	white   log.Logger
}

func (l *loggers) SetLevel(level log.Level) {
	l.red.SetLevel(level)
	l.green.SetLevel(level)
	l.yellow.SetLevel(level)
	l.blue.SetLevel(level)
	l.magenta.SetLevel(level)
	l.cyan.SetLevel(level)
	l.white.SetLevel(level)
}

// Decorator decorates a Go application.
type Decorator struct {
	loggers *loggers
}

// New creates a new decorator.
func New() *Decorator {
	return &Decorator{
		loggers: &loggers{
			red:     log.NewColored(log.None, color.New(color.FgRed)),
			green:   log.NewColored(log.None, color.New(color.FgGreen)),
			yellow:  log.NewColored(log.None, color.New(color.FgYellow)),
			blue:    log.NewColored(log.None, color.New(color.FgBlue)),
			magenta: log.NewColored(log.None, color.New(color.FgMagenta)),
			cyan:    log.NewColored(log.None, color.New(color.FgCyan)),
			white:   log.NewColored(log.None, color.New(color.FgWhite)),
		},
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
		if file.IsDir() && file.Name() != decoratedDir {
			subdir := filepath.Join(relPath, file.Name())
			if err := d.directories(basePath, subdir, visit); err != nil {
				return err
			}
		}
	}

	return nil
}

// Decorate decorates a Go application.
func (d *Decorator) Decorate(level log.Level, path string) error {
	// Update loggers
	d.loggers.SetLevel(level)

	// Sanitize the path
	if _, err := os.Stat(path); err != nil {
		return err
	}

	d.loggers.white.Infof("Decorating ...")

	visitor := &visitor{
		depth:   4,
		loggers: d.loggers,
	}

	return d.directories(path, ".", func(basePath, relPath string) error {
		// Creating a new directory for the decorated package
		newDir := filepath.Join(basePath, decoratedDir, relPath)
		if err := os.MkdirAll(newDir, os.ModePerm); err != nil {
			return err
		}
		d.loggers.blue.Tracef("  Directory created: %s", newDir)

		// Parse all Go packages and files in the currecnt directory
		fset := token.NewFileSet()
		pkgDir := filepath.Join(basePath, relPath)
		d.loggers.cyan.Debugf("  Parsing directory: %s", pkgDir)
		pkgs, err := parser.ParseDir(fset, pkgDir, nil, parser.AllErrors)
		if err != nil {
			return err
		}
		d.loggers.cyan.Tracef("  Directory parsed: %s", pkgDir)

		// Visit all parsed Go files in the currecnt directory
		for _, pkg := range pkgs {
			d.loggers.magenta.Debugf("     Package: %s", pkg.Name)
			for name, file := range pkg.Files {
				if !strings.HasSuffix(name, "_test.go") {
					d.loggers.green.Debugf("      File: %s", name)

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

					d.loggers.green.Debugf("      File written: %s", newName)
				}
			}
		}

		return nil
	})
}
