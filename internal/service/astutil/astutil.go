package astutil

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/tools/imports"
)

const debugFile = "gelato-debug.log"

// WriteFile formats and writes an ast.File to disk.
func WriteFile(path string, file *ast.File, fset *token.FileSet) error {
	buf := new(bytes.Buffer)
	if err := format.Node(buf, fset, file); err != nil {
		return fmt.Errorf("gofmt error: %s", err)
	}

	// Format the modified Go file
	b, err := imports.Process(path, buf.Bytes(), &imports.Options{
		TabWidth:  8,
		TabIndent: true,
		Comments:  true,
		Fragment:  true,
	})

	if err != nil {
		// Try writing a log file for debugging purposes
		_ = ioutil.WriteFile(debugFile, buf.Bytes(), 0644)

		return fmt.Errorf("goimports error: %s", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	if _, err := f.Write(b); err != nil {
		return err
	}

	return nil
}
