package modifier

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

type typeModifier struct {
	modifier

	typeName     string
	typeExported bool

	outputs struct {
		Interface *interfaceType
		Struct    *structType
	}
}

func (m *typeModifier) Modify(n ast.Node) ast.Node {
	// Reset the state
	m.typeName = ""
	m.typeExported = false
	m.outputs.Interface = nil
	m.outputs.Struct = nil

	return astutil.Apply(n, m.pre, m.post)
}

func (m *typeModifier) pre(c *astutil.Cursor) bool {
	m.depth++

	switch n := c.Node().(type) {
	case *ast.GenDecl:
		return n.Tok == token.TYPE

	case *ast.TypeSpec:
		return true

	case *ast.Ident:
		if _, ok := c.Parent().(*ast.TypeSpec); ok {
			m.typeName = n.Name
			m.typeExported = n.Name == strings.Title(n.Name)
		}

	case *ast.InterfaceType:
		m.outputs.Interface = &interfaceType{
			Exported: m.typeExported,
			Name:     m.typeName,
		}
		return true

	case *ast.StructType:
		m.outputs.Struct = &structType{
			Exported: m.typeExported,
			Name:     m.typeName,
		}
		return true

	case *ast.FieldList:
		// Modify the struct field list
		// TODO: verify this is the right FieldList to modify (as opposed to a FieldList in InterfaceType or a different StructType type)
		n.List = []*ast.Field{
			{
				Names: []*ast.Ident{
					{
						Name: "impl",
					},
				},
				Type: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "pkg",
					},
					Sel: &ast.Ident{
						Name: m.typeName,
					},
				},
			},
		}
		return false
	}

	return false
}

func (m *typeModifier) post(c *astutil.Cursor) bool {
	m.depth--
	return true
}
