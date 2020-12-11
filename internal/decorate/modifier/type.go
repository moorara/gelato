package modifier

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

type typeModifier struct {
	modifier
	outputs struct {
		IsInterface bool
		IsStruct    bool
		Exported    bool
		TypeName    string
	}
}

func (m *typeModifier) Modify(n ast.Node) ast.Node {
	m.outputs.IsInterface = false
	m.outputs.IsStruct = false
	m.outputs.Exported = false
	m.outputs.TypeName = ""

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
			m.outputs.TypeName = n.Name
			m.outputs.Exported = n.Name == strings.Title(n.Name)
		}
	case *ast.InterfaceType:
		m.outputs.IsInterface = true
		return true
	case *ast.StructType:
		m.outputs.IsStruct = true
		return true
	case *ast.FieldList:
		return false
	}

	return false
}

func (m *typeModifier) post(c *astutil.Cursor) bool {
	m.depth--
	return true
}
