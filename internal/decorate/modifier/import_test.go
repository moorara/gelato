package modifier

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/log"
)

func TestImportModifier(t *testing.T) {
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

	importGenDecl := &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Name: &ast.Ident{
					Name: "u",
				},
				Path: &ast.BasicLit{
					Value: "net/url",
				},
			},
		},
	}

	tests := []struct {
		name         string
		depth        int
		node         ast.Node
		expectedNode ast.Node
	}{
		{
			name:  "InvalidGenDecl",
			depth: 2,
			node: &ast.GenDecl{
				Tok: token.TYPE,
			},
			expectedNode: &ast.GenDecl{
				Tok: token.TYPE,
			},
		},
		{
			name:         "ImportGenDecl",
			depth:        2,
			node:         importGenDecl,
			expectedNode: importGenDecl,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := &importModifier{
				modifier: modifier{
					depth:  tc.depth,
					logger: clogger,
				},
			}

			node := m.Modify(tc.node)

			assert.Equal(t, tc.expectedNode, node)
		})
	}
}
