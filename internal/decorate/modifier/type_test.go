package modifier

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/log"
)

func TestTypeModifier(t *testing.T) {
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

	interfaceGenDecl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{
					Name: "DomainController",
				},
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{},
								Type:  &ast.FuncType{},
							},
						},
					},
				},
			},
		},
	}

	structGenDecl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{
					Name: "domainController",
				},
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{},
								Type:  &ast.SelectorExpr{},
							},
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name                string
		depth               int
		node                ast.Node
		expectedNode        ast.Node
		expectedExported    bool
		expectedIsInterface bool
		expectedIsStruct    bool
		expectedTypeName    string
	}{
		{
			name:  "InvalidGenDecl",
			depth: 2,
			node: &ast.GenDecl{
				Tok: token.IMPORT,
			},
			expectedNode: &ast.GenDecl{
				Tok: token.IMPORT,
			},
			expectedIsInterface: false,
			expectedIsStruct:    false,
			expectedExported:    false,
			expectedTypeName:    "",
		},
		{
			name:                "InterfaceGenDecl",
			depth:               2,
			node:                interfaceGenDecl,
			expectedNode:        interfaceGenDecl,
			expectedIsInterface: true,
			expectedIsStruct:    false,
			expectedExported:    true,
			expectedTypeName:    "DomainController",
		},
		{
			name:                "StructGenDecl",
			depth:               2,
			node:                structGenDecl,
			expectedNode:        structGenDecl,
			expectedIsInterface: false,
			expectedIsStruct:    true,
			expectedExported:    false,
			expectedTypeName:    "domainController",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := &typeModifier{
				modifier: modifier{
					depth:  tc.depth,
					logger: clogger,
				},
			}

			node := m.Apply(tc.node)

			assert.Equal(t, tc.expectedNode, node)
			assert.Equal(t, tc.expectedIsInterface, m.outputs.IsInterface)
			assert.Equal(t, tc.expectedIsStruct, m.outputs.IsStruct)
			assert.Equal(t, tc.expectedExported, m.outputs.Exported)
			assert.Equal(t, tc.expectedTypeName, m.outputs.TypeName)
		})
	}
}
