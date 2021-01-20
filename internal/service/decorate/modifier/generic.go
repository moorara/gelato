package modifier

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/moorara/gelato/internal/log"
)

// GenericModifier is used for decorating a generic package.
// It implements the Pre and Post astutil.ApplyFunc functions.
type GenericModifier struct {
	modifier
	importModifier *genericImportModifier
	typeModifier   *genericTypeModifier
	funcModifier   *genericFuncModifier
	inputs         struct {
		module  string
		relPath string
	}
	outputs struct {
		pkgName        string
		ExportedType   string
		UnexportedType string
	}
}

// NewGeneric creates a new generic modifier.
func NewGeneric(logger *log.ColorfulLogger) *GenericModifier {
	m := modifier{
		depth:  0,
		logger: logger,
	}

	return &GenericModifier{
		modifier:       m,
		importModifier: &genericImportModifier{modifier: m},
		typeModifier:   &genericTypeModifier{modifier: m},
		funcModifier:   &genericFuncModifier{modifier: m},
	}
}

// Modify modifies an ast.File node.
func (m *GenericModifier) Modify(module, relPath string, n ast.Node) ast.Node {
	m.inputs.module = module
	m.inputs.relPath = relPath
	m.outputs.pkgName = ""
	m.outputs.ExportedType = ""
	m.outputs.UnexportedType = ""

	return astutil.Apply(n, m.pre, m.post)
}

// Pre is called for each node before the node's children are traversed (pre-order).
func (m *GenericModifier) pre(c *astutil.Cursor) bool {
	m.depth++

	switch n := c.Node().(type) {
	case *ast.File:
		return true

	case *ast.Ident:
		// Keep the node in the AST if it is the package identifier
		if _, ok := c.Parent().(*ast.File); ok {
			m.outputs.pkgName = n.Name
			return false
		}

	case *ast.GenDecl:
		switch n.Tok {
		case token.IMPORT:
			// If GenDecl is an import, keep it in the AST
			origPkgName := "_" + m.outputs.pkgName
			origPkgPath := filepath.Join(m.inputs.module, m.inputs.relPath)
			m.importModifier.Modify(origPkgName, origPkgPath, n)
			return false
		case token.TYPE:
			// If GenDecl is a type, visit its children using another modifier to determine whether it is an interface, struct, etc.
			origPkgName := "_" + m.outputs.pkgName
			m.typeModifier.Modify(origPkgName, m.outputs.ExportedType, n)
			out := m.typeModifier.outputs

			if out.Exported {
				m.outputs.ExportedType = out.TypeName
			} else {
				m.outputs.UnexportedType = out.TypeName
			}

			if out.Interface != nil && out.Interface.Exported {
				// TODO: save a reference to the interface type
			} else if out.Struct != nil && !out.Struct.Exported {
				// Keep the modified GenDecl in the AST if it is a struct declaration
				// TODO: determine if the struct is implementing the interface
				return false
			}
		}

	case *ast.FuncDecl:
		// Visit the function node children using another modifier to determine wheher or not we should keep it in the AST
		origPkgName := "_" + m.outputs.pkgName
		m.funcModifier.Modify(origPkgName, m.outputs.ExportedType, m.outputs.UnexportedType, n)
		out := m.funcModifier.outputs

		if out.Func.Exported {
			// Keep the modified FuncDecl in the AST if it implements an interface method
			// TODO: determine if the current method has a counterpart in the interface
			return false
		}
	}

	// Remove the node from the AST if it is part of its parent slice
	if c.Index() >= 0 {
		c.Delete()
	}

	return false
}

// Post is called for each node after its children are traversed (post-order).
func (m *GenericModifier) post(c *astutil.Cursor) bool {
	m.depth--
	return true
}

type genericImportModifier struct {
	modifier
	inputs struct {
		origPkgName string
		origPkgPath string
	}
}

func (m *genericImportModifier) Modify(origPkgName, origPkgPath string, n ast.Node) ast.Node {
	m.inputs.origPkgName = origPkgName
	m.inputs.origPkgPath = origPkgPath

	return astutil.Apply(n, m.pre, m.post)
}

func (m *genericImportModifier) pre(c *astutil.Cursor) bool {
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

func (m *genericImportModifier) post(c *astutil.Cursor) bool {
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

type genericTypeModifier struct {
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

func (m *genericTypeModifier) Modify(origPkgName, interfaceName string, n ast.Node) ast.Node {
	m.inputs.origPkgName = origPkgName
	m.inputs.interfaceName = interfaceName
	m.outputs.TypeName = ""
	m.outputs.Exported = false
	m.outputs.Interface = nil
	m.outputs.Struct = nil

	return astutil.Apply(n, m.pre, m.post)
}

func (m *genericTypeModifier) createStructFieldList() []*ast.Field {
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

func (m *genericTypeModifier) pre(c *astutil.Cursor) bool {
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

func (m *genericTypeModifier) post(c *astutil.Cursor) bool {
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

const (
	addToReceiver = 1 + iota
	addToInputs
	addToOutputs
)

type genericFuncModifier struct {
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

func (m *genericFuncModifier) Modify(origPkgName, interfaceName, structName string, n ast.Node) ast.Node {
	m.addTo = 0
	m.inputs.origPkgName = origPkgName
	m.inputs.interfaceName = interfaceName
	m.inputs.structName = structName
	m.outputs.Func = funcType{}

	return astutil.Apply(n, m.pre, m.post)
}

func (m *genericFuncModifier) createNewFuncBody() *ast.BlockStmt {
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

func (m *genericFuncModifier) createDecoratedMethodBody() *ast.BlockStmt {
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

func (m *genericFuncModifier) pre(c *astutil.Cursor) bool {
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

func (m *genericFuncModifier) post(c *astutil.Cursor) bool {
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
