package compile

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/moorara/gelato/internal/log"
)

// Compiler creates test helpers (mocks, factories, builders, etc.) for a Go file.
type Compiler struct {
	logger *log.ColorfulLogger
	state  struct {
		TypeName string
	}
	inputs struct {
		pkgPath string
		pkgName string
	}
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
	c.inputs.pkgPath = pkgPath
	c.inputs.pkgName = pkgName

	astutil.Apply(file, c.pre, c.post)
	return c.outputs.file
}

func createFile(pkgPath, pkgName string) *ast.File {
	return &ast.File{
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
}

func createBuilderStruct(pkgName, typeName string) *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{
					Name: typeName + "Builder",
				},
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.SelectorExpr{
									X:   &ast.Ident{Name: pkgName},
									Sel: &ast.Ident{Name: typeName},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Pre is called for each node before the node's children are traversed (pre-order).
func (c *Compiler) pre(cr *astutil.Cursor) bool {
	switch n := cr.Node().(type) {
	case *ast.File:
		c.outputs.file = createFile(c.inputs.pkgPath, c.inputs.pkgName)
		return true

	case *ast.GenDecl:
		switch n.Tok {
		case token.IMPORT:
		case token.TYPE:
			return true
		}

	case *ast.TypeSpec:
		// Continue only if the type is exported
		if name := n.Name.Name; name == strings.Title(name) {
			c.state.TypeName = name
			return true
		}

	case *ast.StructType:
		builderStruct := createBuilderStruct(c.inputs.pkgName, c.state.TypeName)
		c.outputs.file.Decls = append(c.outputs.file.Decls, builderStruct)

	case *ast.InterfaceType:
	}

	return false
}

// Post is called for each node after its children are traversed (post-order).
func (c *Compiler) post(cr *astutil.Cursor) bool {
	switch cr.Node().(type) {
	case *ast.TypeSpec:
		c.state.TypeName = ""
	}

	return true
}
