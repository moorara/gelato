package modifier

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/moorara/gelato/internal/log"
)

// DebugModifier is used for debugging an ast.Node.
// It implements the Pre and Post astutil.ApplyFunc functions.
type DebugModifier struct {
	modifier
}

// NewDebug creates a new debug modifier.
func NewDebug(logger *log.ColorfulLogger) *DebugModifier {
	return &DebugModifier{
		modifier: modifier{
			depth:  0,
			logger: logger,
		},
	}
}

// Modify prints debugging information for an ast.File node.
func (m *DebugModifier) Modify(n ast.Node) ast.Node {
	return astutil.Apply(n, m.pre, m.post)
}

// Pre is called for each node before the node's children are traversed (pre-order).
func (m *DebugModifier) pre(c *astutil.Cursor) bool {
	m.depth++

	indent := strings.Repeat("  ", m.depth)
	m.logger.Yellow.Tracef("%s[%s] %T", indent, c.Name(), c.Node())

	switch n := c.Node().(type) {
	case *ast.BasicLit:
		m.logger.Red.Tracef("%s  %s", indent, n.Value)
	case *ast.GenDecl:
		m.logger.Red.Tracef("%s  %s", indent, n.Tok)
	case *ast.Ident:
		m.logger.Red.Tracef("%s  %s", indent, n.Name)
	case *ast.ImportSpec:
		m.logger.Red.Tracef("%s  %s", indent, n.Path.Value)
	}

	return true
}

// Post is called for each node after its children are traversed (post-order).
func (m *DebugModifier) post(c *astutil.Cursor) bool {
	m.depth--
	return true
}
