package decorate

import (
	"go/ast"
	"strings"
)

type visitor struct {
	depth   int
	loggers *loggers
}

func (v *visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	indent := strings.Repeat("  ", v.depth)
	v.loggers.yellow.Tracef("%s%T", indent, n)

	switch n := n.(type) {
	case *ast.GenDecl:
	case *ast.FuncDecl:
	case *ast.Ident:
		v.loggers.red.Tracef("%s  %s", indent, n.Name)
	case *ast.ImportSpec:
		v.loggers.red.Tracef("%s  %s", indent, n.Path.Value)
	}

	return &visitor{
		depth:   v.depth + 1,
		loggers: v.loggers,
	}
}
