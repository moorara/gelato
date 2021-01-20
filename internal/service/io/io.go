package io

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// GoModule returns the name of go module in a give path.
func GoModule(path string) (string, error) {
	f, err := os.Open(filepath.Join(path, "go.mod"))
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

type visitFunc func(basePath, relPath string) error

func isPackageDir(name string) bool {
	startsWithDot := strings.HasPrefix(name, ".")
	return !startsWithDot && name != "bin"
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
