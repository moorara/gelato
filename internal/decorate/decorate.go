package decorate

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	decoratedDir = ".build"
)

// Decorator decorates a Go application.
type Decorator struct {
}

// New creates a new decorator.
func New() *Decorator {
	return &Decorator{}
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
func (d *Decorator) Decorate(path string) error {
	// Sanitize the path
	if _, err := os.Stat(path); err != nil {
		return err
	}

	return d.directories(path, ".", func(basePath, relPath string) error {
		// Creating a new directory for the decorated package
		newDir := filepath.Join(basePath, decoratedDir, relPath)
		if err := os.MkdirAll(newDir, os.ModePerm); err != nil {
			return err
		}

		return nil
	})
}
