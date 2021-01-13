package io

import (
	"io/ioutil"
	"path/filepath"
)

type visitFunc func(basePath, relPath string) error

func isPackageDir(name string) bool {
	return name != "bin" && name != ".build" && name != ".gen"
}

// PackageDirectories visits all package directories in a given path.
func PackageDirectories(basePath, relPath string, visit visitFunc) error {
	if err := visit(basePath, relPath); err != nil {
		return err
	}

	dir := filepath.Join(basePath, relPath)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		// Any directory that starts with "." is NOT considered
		if file.IsDir() && isPackageDir(file.Name()) {
			subdir := filepath.Join(relPath, file.Name())
			if err := PackageDirectories(basePath, subdir, visit); err != nil {
				return err
			}
		}
	}

	return nil
}
