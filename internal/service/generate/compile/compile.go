package compile

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/moorara/gelato/internal/log"
)

// Compiler creates test helpers (mocks, factories, builders, etc.) for a Go file.
type Compiler struct {
	logger  *log.ColorfulLogger
	outputs struct {
		file *ast.File
	}
}

// New creates a new compiler.
func New(logger *log.ColorfulLogger) *Compiler {
	return &Compiler{
		logger: logger,
	}
}

// Compile takes an ast.File node and generates a new ast.File node with test helpers (mocks, factories, builders, etc.).
func (c *Compiler) Compile(pkgPath, pkgName string, file *ast.File) *ast.File {
	c.outputs.file = &ast.File{
		Name: &ast.Ident{
			Name: pkgName + "test",
		},
		Decls: []ast.Decl{
			// Imports
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Value: fmt.Sprintf("%q", pkgPath),
						},
					},
				},
			},
		},
	}

	astutil.Apply(file, c.pre, c.post)

	return c.outputs.file
}

// Pre is called for each node before the node's children are traversed (pre-order).
func (c *Compiler) pre(cursor *astutil.Cursor) bool {
	return true
}

// Post is called for each node after its children are traversed (post-order).
func (c *Compiler) post(cursor *astutil.Cursor) bool {
	return true
}
