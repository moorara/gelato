package visitor

import (
	"go/ast"
	"strings"

	"github.com/moorara/gelato/internal/log"
)

// DebugVisitor is used for debugging an ast.Node.
// It implements the ast.Visitor interface.
type DebugVisitor struct {
	visitor
}

// NewDebug creates a new debug visitor.
func NewDebug(depth int, logger *log.ColorfulLogger) *DebugVisitor {
	return &DebugVisitor{
		visitor: visitor{
			depth:  depth,
			logger: logger,
		},
	}
}

// Visit is called for an AST node and returns a new visitor for visiting the child nodes.
func (v *DebugVisitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	indent := strings.Repeat("  ", v.depth)
	v.logger.Yellow.Tracef("%s%T", indent, n)

	switch n := n.(type) {
	case *ast.BasicLit:
		v.logger.Red.Tracef("%s  %s", indent, n.Value)
	case *ast.GenDecl:
		v.logger.Red.Tracef("%s  %s", indent, n.Tok)
	case *ast.Ident:
		v.logger.Red.Tracef("%s  %s", indent, n.Name)
	case *ast.ImportSpec:
		v.logger.Red.Tracef("%s  %s", indent, n.Path.Value)
	}

	return &DebugVisitor{
		visitor: v.visitor.Next(),
	}
}
