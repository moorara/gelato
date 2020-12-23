package modifier

import (
	"go/ast"
	"go/token"
	"path/filepath"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/moorara/gelato/internal/log"
)

// FileModifier is used for modifying an ast.File node.
// It implements the Pre and Post astutil.ApplyFunc functions.
type FileModifier struct {
	modifier
	importModifier *importModifier
	typeModifier   *typeModifier
	funcModifier   *funcModifier
	inputs         struct {
		module  string
		decDir  string
		relPath string
	}
	outputs struct {
		pkgName        string
		ExportedType   string
		UnexportedType string
	}
}

// NewFile creates a new file modifier.
func NewFile(depth int, logger *log.ColorfulLogger) *FileModifier {
	m := modifier{
		depth:  depth,
		logger: logger,
	}

	return &FileModifier{
		modifier:       m,
		importModifier: &importModifier{modifier: m},
		typeModifier:   &typeModifier{modifier: m},
		funcModifier:   &funcModifier{modifier: m},
	}
}

// Modify modifies a ast.File node.
func (m *FileModifier) Modify(module, decDir, relPath string, n ast.Node) ast.Node {
	m.inputs.module = module
	m.inputs.decDir = decDir
	m.inputs.relPath = relPath
	m.outputs.pkgName = ""
	m.outputs.ExportedType = ""
	m.outputs.UnexportedType = ""

	return astutil.Apply(n, m.pre, m.post)
}

// Pre is called for each node before the node's children are traversed (pre-order).
func (m *FileModifier) pre(c *astutil.Cursor) bool {
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
func (m *FileModifier) post(c *astutil.Cursor) bool {
	m.depth--
	return true
}
