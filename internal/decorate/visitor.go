package decorate

import (
	"go/ast"
	"strings"

	"github.com/moorara/gelato/internal/log"
)

type visitor struct {
	depth  int
	logger *log.ColorfulLogger
}

func (v *visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	indent := strings.Repeat("  ", v.depth)
	v.logger.Yellow.Tracef("%s%T", indent, n)

	switch n := n.(type) {
	case *ast.GenDecl:
	case *ast.FuncDecl:
	case *ast.Ident:
		v.logger.Red.Tracef("%s  %s", indent, n.Name)
	case *ast.ImportSpec:
		v.logger.Red.Tracef("%s  %s", indent, n.Path.Value)
	}

	return &visitor{
		depth:  v.depth + 1,
		logger: v.logger,
	}
}
