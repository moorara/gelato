package modifier

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

const (
	addToReceiver = 1 + iota
	addToInputs
	addToOutputs
)

type funcModifier struct {
	modifier
	addTo  int
	inputs struct {
		origPkgName   string
		interfaceName string
		structName    string
	}
	outputs struct {
		Func funcType
	}
}

func (m *funcModifier) Modify(origPkgName, interfaceName, structName string, n ast.Node) ast.Node {
	m.addTo = 0
	m.inputs.origPkgName = origPkgName
	m.inputs.interfaceName = interfaceName
	m.inputs.structName = structName
	m.outputs.Func = funcType{}

	return astutil.Apply(n, m.pre, m.post)
}

func (m *funcModifier) createNewFuncBody() *ast.BlockStmt {
	argsExprs := []ast.Expr{}
	for _, field := range m.outputs.Func.Inputs {
		for _, name := range field.Names {
			argsExprs = append(argsExprs, &ast.Ident{Name: name})
		}
	}

	returnsExprs := []ast.Expr{}
	for i := 0; i < len(m.outputs.Func.Outputs)-1; i++ {
		returnsExprs = append(returnsExprs, &ast.Ident{Name: "nil"})
	}
	returnsExprs = append(returnsExprs, &ast.Ident{Name: errorID})

	return &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.AssignStmt{
				// TODO: TokPos
				Lhs: []ast.Expr{
					&ast.Ident{Name: implementationID},
					&ast.Ident{Name: errorID},
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: m.inputs.origPkgName},
							Sel: &ast.Ident{Name: m.outputs.Func.Name},
						},
						Args: argsExprs,
					},
				},
			},
			&ast.IfStmt{
				// TODO: If
				Cond: &ast.BinaryExpr{
					X:  &ast.Ident{Name: errorID},
					Op: token.NEQ,
					Y:  &ast.Ident{Name: "nil"},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ReturnStmt{
							Results: returnsExprs,
						},
					},
				},
			},
			&ast.ReturnStmt{
				// TODO: Return
				Results: []ast.Expr{
					&ast.UnaryExpr{
						Op: token.AND,
						X: &ast.CompositeLit{
							Type: &ast.Ident{Name: m.inputs.structName},
							Elts: []ast.Expr{
								&ast.KeyValueExpr{
									Key:   &ast.Ident{Name: implementationID},
									Value: &ast.Ident{Name: implementationID},
								},
							},
						},
					},
					&ast.Ident{Name: "nil"},
				},
			},
		},
	}
}

func (m *funcModifier) createDecoratedMethodBody() *ast.BlockStmt {
	argsExprs := []ast.Expr{}
	for _, field := range m.outputs.Func.Inputs {
		for _, name := range field.Names {
			argsExprs = append(argsExprs, &ast.Ident{Name: name})
		}
	}

	callExpr := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X: &ast.SelectorExpr{
				X:   &ast.Ident{Name: m.outputs.Func.Receiver.Name},
				Sel: &ast.Ident{Name: implementationID},
			},
			Sel: &ast.Ident{Name: m.outputs.Func.Name},
		},
		Args: argsExprs,
	}

	var stmt ast.Stmt
	if len(m.outputs.Func.Outputs) == 0 {
		stmt = &ast.ExprStmt{
			X: callExpr,
		}
	} else {
		stmt = &ast.ReturnStmt{
			// TODO: Return
			Results: []ast.Expr{
				callExpr,
			},
		}
	}

	return &ast.BlockStmt{
		List: []ast.Stmt{stmt},
	}
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
			m.addTo = addToReceiver
			m.outputs.Func.Receiver = &receiver{}
		case "Params":
			m.addTo = addToInputs
		case "Results":
			m.addTo = addToOutputs
		}
		return true

	case *ast.Field:
		switch m.addTo {
		case addToReceiver:
			m.outputs.Func.Receiver.Name = n.Names[0].Name
		case addToInputs:
			m.outputs.Func.Inputs.Append(n)
		case addToOutputs:
			m.outputs.Func.Outputs.Append(n)

			// Check if this is a New... function for creating an interface implementation
			if m.outputs.Func.Exported && m.outputs.Func.Receiver == nil {
				if id, ok := n.Type.(*ast.Ident); ok && id.Name == m.inputs.interfaceName {
					// Reference the return interface type from the original package
					n.Type = &ast.SelectorExpr{
						X:   &ast.Ident{Name: m.inputs.origPkgName},
						Sel: &ast.Ident{Name: m.inputs.interfaceName},
					}
				}
			}
		}
		return true

	case *ast.StarExpr:
		switch m.addTo {
		case addToReceiver:
			m.outputs.Func.Receiver.Star = true
		case addToInputs:
			m.outputs.Func.Inputs.SetStar()
		case addToOutputs:
			m.outputs.Func.Outputs.SetStar()
		}
		return true

	case *ast.SelectorExpr:
		return true

	case *ast.Ident:
		switch c.Parent().(type) {
		case *ast.Field, *ast.StarExpr:
			switch m.addTo {
			case addToReceiver:
				m.outputs.Func.Receiver.Type = n.Name
			case addToInputs:
				m.outputs.Func.Inputs.SetType(n)
			case addToOutputs:
				m.outputs.Func.Outputs.SetType(n)
			}
		case *ast.SelectorExpr:
			// SelectorExpr can only appear for a method input or output
			switch c.Name() {
			case "X":
				if m.addTo == addToInputs {
					m.outputs.Func.Inputs.SetPackage(n)
				} else if m.addTo == addToOutputs {
					m.outputs.Func.Outputs.SetPackage(n)
				}
			case "Sel":
				if m.addTo == addToInputs {
					m.outputs.Func.Inputs.SetType(n)
				} else if m.addTo == addToOutputs {
					m.outputs.Func.Outputs.SetType(n)
				}
			}
		}
	}

	return false
}

func (m *funcModifier) post(c *astutil.Cursor) bool {
	m.depth--

	switch n := c.Node().(type) {
	case *ast.FuncDecl:
		// Re-write the function body
		if m.outputs.Func.Exported {
			if m.outputs.Func.Receiver == nil { // New... function
				n.Body = m.createNewFuncBody()
			} else { // Struct method
				n.Body = m.createDecoratedMethodBody()
			}
		}

	case *ast.FieldList:
		switch c.Name() {
		case "Recv":
			m.addTo = 0
		case "Params":
			m.addTo = 0
		case "Results":
			m.addTo = 0
		}
	}

	return true
}
