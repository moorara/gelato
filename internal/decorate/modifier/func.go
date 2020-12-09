package modifier

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

type funcModifier struct {
	modifier
	outputs struct {
		Exported bool
		FuncName string
	}
}

func (m *funcModifier) Apply(n ast.Node) ast.Node {
	m.outputs.Exported = false
	m.outputs.FuncName = ""

	return astutil.Apply(n, m.Pre, m.Post)
}

func (m *funcModifier) Pre(c *astutil.Cursor) bool {
	m.depth++

	switch n := c.Node().(type) {
	case *ast.FuncDecl:
		return true
	case *ast.Ident:
		if _, ok := c.Parent().(*ast.FuncDecl); ok {
			m.outputs.FuncName = n.Name
			m.outputs.Exported = (n.Name == strings.Title(n.Name))
		}
	case *ast.FieldList:
		return true
	case *ast.Field:
		return true
	}

	return false
}

func (m *funcModifier) Post(c *astutil.Cursor) bool {
	m.depth--
	return true
}
