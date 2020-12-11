package modifier

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

type importModifier struct {
	modifier
}

func (m *importModifier) Modify(n ast.Node) ast.Node {
	return astutil.Apply(n, m.pre, m.post)
}

func (m *importModifier) pre(c *astutil.Cursor) bool {
	m.depth++

	switch n := c.Node().(type) {
	case *ast.GenDecl:
		return n.Tok == token.IMPORT
	case *ast.ImportSpec:
		return true
	case *ast.Ident:
		return true
	case *ast.BasicLit:
		return true
	}

	return false
}

func (m *importModifier) post(c *astutil.Cursor) bool {
	m.depth--
	return true
}
