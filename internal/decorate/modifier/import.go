package modifier

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

type importModifier struct {
	modifier
}

func (m *importModifier) Apply(n ast.Node) ast.Node {
	return astutil.Apply(n, m.Pre, m.Post)
}

func (m *importModifier) Pre(c *astutil.Cursor) bool {
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

func (m *importModifier) Post(c *astutil.Cursor) bool {
	m.depth--
	return true
}
