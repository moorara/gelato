package modifier

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

type funcModifier struct {
	modifier

	addToReceiver bool
	addToInputs   bool
	addToOutputs  bool

	outputs struct {
		Func funcType
	}
}

func (m *funcModifier) Modify(n ast.Node) ast.Node {
	// Reset the state
	m.addToReceiver = false
	m.addToInputs = false
	m.addToOutputs = false
	m.outputs.Func = funcType{}

	return astutil.Apply(n, m.pre, m.post)
}

func (m *funcModifier) pre(c *astutil.Cursor) bool {
	m.depth++

	switch n := c.Node().(type) {
	case *ast.FuncDecl:
		name := n.Name.Name
		m.outputs.Func.Name = name
		m.outputs.Func.Exported = name == strings.Title(name)
		return true

	case *ast.FuncType:
		return true

	case *ast.FieldList:
		switch c.Name() {
		case "Recv":
			m.addToReceiver = true
			m.outputs.Func.Receiver = &receiver{}
		case "Params":
			m.addToInputs = true
		case "Results":
			m.addToOutputs = true
		}
		return true

	case *ast.Field:
		if m.addToReceiver {
			m.outputs.Func.Receiver.Name = n.Names[0].Name
		} else if m.addToInputs {
			m.outputs.Func.Inputs.Append(n)
		} else if m.addToOutputs {
			m.outputs.Func.Outputs.Append(n)
		}
		return true

	case *ast.StarExpr:
		if m.addToReceiver {
			m.outputs.Func.Receiver.Star = true
		} else if m.addToInputs {
			m.outputs.Func.Inputs.SetStar()
		} else if m.addToOutputs {
			m.outputs.Func.Outputs.SetStar()
		}
		return true

	case *ast.SelectorExpr:
		return true

	case *ast.Ident:
		switch c.Parent().(type) {
		case *ast.Field, *ast.StarExpr:
			if m.addToReceiver {
				m.outputs.Func.Receiver.Type = n.Name
			} else if m.addToInputs {
				m.outputs.Func.Inputs.SetType(n)
			} else if m.addToOutputs {
				m.outputs.Func.Outputs.SetType(n)
			}
		case *ast.SelectorExpr:
			// SelectorExpr can only appear for a method input or output
			switch c.Name() {
			case "X":
				if m.addToInputs {
					m.outputs.Func.Inputs.SetPackage(n)
				} else if m.addToOutputs {
					m.outputs.Func.Outputs.SetPackage(n)
				}
			case "Sel":
				if m.addToInputs {
					m.outputs.Func.Inputs.SetType(n)
				} else if m.addToOutputs {
					m.outputs.Func.Outputs.SetType(n)
				}
			}
		}
	}

	return false
}

func (m *funcModifier) post(c *astutil.Cursor) bool {
	m.depth--

	switch c.Node().(type) {
	case *ast.FieldList:
		switch c.Name() {
		case "Recv":
			m.addToReceiver = false
		case "Params":
			m.addToInputs = false
		case "Results":
			m.addToOutputs = false
		}
		return true
	}

	return true
}
