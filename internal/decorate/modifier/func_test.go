package modifier

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/log"
)

func TestFuncModifier(t *testing.T) {
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

	exportedFunc := &ast.FuncDecl{
		Recv: nil,
		Name: &ast.Ident{
			Name: "NewDomain",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{},
						Type:  &ast.FuncType{},
					},
				},
			},
		},
	}

	unexportedFunc := &ast.FuncDecl{
		Recv: nil,
		Name: &ast.Ident{
			Name: "newDomain",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{},
						Type:  &ast.FuncType{},
					},
				},
			},
		},
	}

	exportedMethod := &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{},
					Type:  &ast.FuncType{},
				},
			},
		},
		Name: &ast.Ident{
			Name: "Calculate",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{},
						Type:  &ast.FuncType{},
					},
				},
			},
		},
	}

	unexportedMethod := &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{},
					Type:  &ast.FuncType{},
				},
			},
		},
		Name: &ast.Ident{
			Name: "calculate",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{},
						Type:  &ast.FuncType{},
					},
				},
			},
		},
	}

	tests := []struct {
		name             string
		depth            int
		node             ast.Node
		expectedNode     ast.Node
		expectedExported bool
		expectedFuncName string
	}{
		{
			name:             "ExportedFunc",
			depth:            2,
			node:             exportedFunc,
			expectedNode:     exportedFunc,
			expectedExported: true,
			expectedFuncName: "NewDomain",
		},
		{
			name:             "UnexportedFunc",
			depth:            2,
			node:             unexportedFunc,
			expectedNode:     unexportedFunc,
			expectedExported: false,
			expectedFuncName: "newDomain",
		},
		{
			name:             "ExportedMethod",
			depth:            2,
			node:             exportedMethod,
			expectedNode:     exportedMethod,
			expectedExported: true,
			expectedFuncName: "Calculate",
		},
		{
			name:             "UnexportedMethod",
			depth:            2,
			node:             unexportedMethod,
			expectedNode:     unexportedMethod,
			expectedExported: false,
			expectedFuncName: "calculate",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := &funcModifier{
				modifier: modifier{
					depth:  tc.depth,
					logger: clogger,
				},
			}

			node := m.Apply(tc.node)

			assert.Equal(t, tc.expectedNode, node)
			assert.Equal(t, tc.expectedExported, m.outputs.Exported)
			assert.Equal(t, tc.expectedFuncName, m.outputs.FuncName)
		})
	}
}
