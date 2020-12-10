package modifier

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/ast/astutil"

	"github.com/moorara/gelato/internal/log"
)

func TestDebugModifier(t *testing.T) {
	logger := log.New(log.None)
	clogger := &log.ColorfulLogger{
		Red:     logger,
		Green:   logger,
		Yellow:  logger,
		Blue:    logger,
		Magenta: logger,
		Cyan:    logger,
		White:   logger,
	}

	tests := []struct {
		name         string
		depth        int
		node         ast.Node
		expectedNode ast.Node
	}{
		{
			name:         "FileNode",
			depth:        2,
			node:         &ast.File{},
			expectedNode: &ast.File{},
		},
		{
			name:  "BasicLitNode",
			depth: 2,
			node: &ast.BasicLit{
				Value: "context",
			},
			expectedNode: &ast.BasicLit{
				Value: "context",
			},
		},
		{
			name:  "GenDeclNode",
			depth: 2,
			node: &ast.GenDecl{
				Tok: token.IMPORT,
			},
			expectedNode: &ast.GenDecl{
				Tok: token.IMPORT,
			},
		},
		{
			name:  "IdentNode",
			depth: 2,
			node: &ast.Ident{
				Name: "id",
			},
			expectedNode: &ast.Ident{
				Name: "id",
			},
		},
		{
			name:  "ImportSpecNode",
			depth: 2,
			node: &ast.ImportSpec{
				Path: &ast.BasicLit{
					Value: "context",
				},
			},
			expectedNode: &ast.ImportSpec{
				Path: &ast.BasicLit{
					Value: "context",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := NewDebug(tc.depth, clogger)
			node := astutil.Apply(tc.node, m.Pre, m.Post)

			assert.Equal(t, tc.expectedNode, node)
		})
	}
}
