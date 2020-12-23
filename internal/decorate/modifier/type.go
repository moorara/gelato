package modifier

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

type typeModifier struct {
	modifier
	inputs struct {
		origPkgName   string
		interfaceName string
	}
	outputs struct {
		TypeName  string
		Exported  bool
		Interface *interfaceType
		Struct    *structType
	}
}

func (m *typeModifier) Modify(origPkgName, interfaceName string, n ast.Node) ast.Node {
	m.inputs.origPkgName = origPkgName
	m.inputs.interfaceName = interfaceName
	m.outputs.TypeName = ""
	m.outputs.Exported = false
	m.outputs.Interface = nil
	m.outputs.Struct = nil

	return astutil.Apply(n, m.pre, m.post)
}

func (m *typeModifier) createStructFieldList() []*ast.Field {
	return []*ast.Field{
		{
			Names: []*ast.Ident{
				{
					// TODO: Resolve NamePos
					Name: implementationID,
				},
			},
			Type: &ast.SelectorExpr{
				X:   &ast.Ident{Name: m.inputs.origPkgName},
				Sel: &ast.Ident{Name: m.inputs.interfaceName},
			},
		},
	}
}

func (m *typeModifier) pre(c *astutil.Cursor) bool {
	m.depth++

	switch n := c.Node().(type) {
	case *ast.GenDecl:
		return n.Tok == token.TYPE

	case *ast.TypeSpec:
		name := n.Name.Name
		m.outputs.TypeName = name
		m.outputs.Exported = name == strings.Title(name)
		return true

	case *ast.InterfaceType:
		m.outputs.Interface = &interfaceType{
			Exported: m.outputs.Exported,
			Name:     m.outputs.TypeName,
		}
		return true

	case *ast.StructType:
		m.outputs.Struct = &structType{
			Exported: m.outputs.Exported,
			Name:     m.outputs.TypeName,
		}
		return true

	case *ast.FuncType:
		return true

	case *ast.FieldList:
		switch c.Name() {
		case "Fields":
			// Modify the struct field list
			// TODO: verify this is the right FieldList to modify (as opposed to a FieldList in InterfaceType or a different StructType type)
			n.List = m.createStructFieldList()
			return false
		case "Methods":
		case "Params":
		case "Results":
		}
		return true

	case *ast.Field:
		return true

	case *ast.StarExpr:
		return true

	case *ast.SelectorExpr:
		return true

	case *ast.Ident:
		return true
	}

	return false
}

func (m *typeModifier) post(c *astutil.Cursor) bool {
	m.depth--

	switch c.Node().(type) {
	case *ast.FieldList:
		switch c.Name() {
		case "Fields":
		case "Methods":
		case "Params":
		case "Results":
		}
	}

	return true
}
