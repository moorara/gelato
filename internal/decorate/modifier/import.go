package modifier

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

type importModifier struct {
	modifier
	inputs struct {
		origPkgName string
		origPkgPath string
	}
}

func (m *importModifier) Modify(origPkgName, origPkgPath string, n ast.Node) ast.Node {
	m.inputs.origPkgName = origPkgName
	m.inputs.origPkgPath = origPkgPath

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

	switch n := c.Node().(type) {
	case *ast.GenDecl:
		n.Specs = append(n.Specs, &ast.ImportSpec{
			Name: &ast.Ident{
				// TODO: Resolve NamePos
				Name: m.inputs.origPkgName,
			},
			Path: &ast.BasicLit{
				// TODO: Resolve ValuePos
				Value: fmt.Sprintf("%q", m.inputs.origPkgPath),
			},
		})
	}

	return true
}
