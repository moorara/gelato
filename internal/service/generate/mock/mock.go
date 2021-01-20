package mock

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/moorara/gelato/internal/log"
)

// Generator is used for generating mocks for interfaces.
type Generator struct {
	logger  *log.ColorfulLogger
	outputs struct {
		file *ast.File
	}
}

// NewGenerator creates a new mock generator.
func NewGenerator(logger *log.ColorfulLogger) *Generator {
	return &Generator{
		logger: logger,
	}
}

// Generate takes an ast.File node and generates a new ast.File node with mocks.
func (g *Generator) Generate(pkgPath, pkgName string, file *ast.File) *ast.File {
	g.outputs.file = &ast.File{
		Name: &ast.Ident{
			Name: pkgName + "mock",
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

	astutil.Apply(file, g.pre, g.post)

	return g.outputs.file
}

// Pre is called for each node before the node's children are traversed (pre-order).
func (g *Generator) pre(c *astutil.Cursor) bool {
	return true
}

// Post is called for each node after its children are traversed (post-order).
func (g *Generator) post(c *astutil.Cursor) bool {
	return true
}
