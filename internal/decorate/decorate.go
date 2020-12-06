package decorate

import (
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/moorara/color"

	"github.com/moorara/gelato/internal/log"
)

const (
	decoratedDir = ".build"
)

// Decorator decorates a Go application.
type Decorator struct {
	loggers struct {
		red     log.Logger
		green   log.Logger
		yellow  log.Logger
		blue    log.Logger
		magenta log.Logger
		cyan    log.Logger
		white   log.Logger
	}
}

// New creates a new decorator.
func New() *Decorator {
	d := &Decorator{}
	d.loggers.red = log.NewColored(log.None, color.New(color.FgRed))
	d.loggers.green = log.NewColored(log.None, color.New(color.FgGreen))
	d.loggers.yellow = log.NewColored(log.None, color.New(color.FgYellow))
	d.loggers.blue = log.NewColored(log.None, color.New(color.FgBlue))
	d.loggers.magenta = log.NewColored(log.None, color.New(color.FgMagenta))
	d.loggers.cyan = log.NewColored(log.None, color.New(color.FgCyan))
	d.loggers.white = log.NewColored(log.None, color.New(color.FgWhite))

	return d
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
	d.loggers.red.SetLevel(level)
	d.loggers.green.SetLevel(level)
	d.loggers.yellow.SetLevel(level)
	d.loggers.blue.SetLevel(level)
	d.loggers.magenta.SetLevel(level)
	d.loggers.cyan.SetLevel(level)
	d.loggers.white.SetLevel(level)

	// Sanitize the path
	if _, err := os.Stat(path); err != nil {
		return err
	}

	d.loggers.white.Infof("Decorating ...")

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
			d.loggers.yellow.Debugf("     Package: %s", pkg.Name)
			for name, _ := range pkg.Files {
				d.loggers.green.Debugf("      File: %s", name)
			}
		}

		return nil
	})
}
