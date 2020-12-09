package visitor

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/log"
)

func TestDebugVisitor(t *testing.T) {
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
		name            string
		depth           int
		node            ast.Node
		expectedVisitor ast.Visitor
	}{
		{
			name:            "NilNode",
			depth:           2,
			node:            nil,
			expectedVisitor: nil,
		},
		{
			name:  "FileNode",
			depth: 2,
			node:  &ast.File{},
			expectedVisitor: &DebugVisitor{
				visitor: visitor{
					depth:  3,
					logger: clogger,
				},
			},
		},
		{
			name:  "BasicLitNode",
			depth: 2,
			node: &ast.BasicLit{
				Value: "context",
			},
			expectedVisitor: &DebugVisitor{
				visitor: visitor{
					depth:  3,
					logger: clogger,
				},
			},
		},
		{
			name:  "GenDeclNode",
			depth: 2,
			node: &ast.GenDecl{
				Tok: token.IMPORT,
			},
			expectedVisitor: &DebugVisitor{
				visitor: visitor{
					depth:  3,
					logger: clogger,
				},
			},
		},
		{
			name:  "IdentNode",
			depth: 2,
			node: &ast.Ident{
				Name: "id",
			},
			expectedVisitor: &DebugVisitor{
				visitor: visitor{
					depth:  3,
					logger: clogger,
				},
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
			expectedVisitor: &DebugVisitor{
				visitor: visitor{
					depth:  3,
					logger: clogger,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := NewDebug(tc.depth, clogger)
			visitor := v.Visit(tc.node)

			assert.Equal(t, tc.expectedVisitor, visitor)
		})
	}
}
