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
		Name: &ast.Ident{Name: "NewController"},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{Name: "ug"},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "gateway"},
								Sel: &ast.Ident{Name: "UserGateway"},
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.Ident{Name: "Controller"},
					},
					{
						Type: &ast.Ident{Name: "error"},
					},
				},
			},
		},
	}

	unexportedFunc := &ast.FuncDecl{
		Recv: nil,
		Name: &ast.Ident{Name: "newController"},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{Name: "ug"},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "gateway"},
								Sel: &ast.Ident{Name: "UserGateway"},
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.Ident{Name: "controller"},
						},
					},
					{
						Type: &ast.Ident{Name: "error"},
					},
				},
			},
		},
	}

	exportedMethod := &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{Name: "c"},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{Name: "controller"},
					},
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
						Names: []*ast.Ident{
							{Name: "a"},
							{Name: "b"},
						},
						Type: &ast.StarExpr{
							X: &ast.Ident{Name: "int"},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "entity"},
								Sel: &ast.Ident{Name: "CalculateResponse"},
							},
						},
					},
					{
						Type: &ast.Ident{Name: "error"},
					},
				},
			},
		},
	}

	unexportedMethod := &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{Name: "c"},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{Name: "controller"},
					},
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
						Names: []*ast.Ident{
							{Name: "a"},
							{Name: "b"},
						},
						Type: &ast.Ident{Name: "int"},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{Name: "resp"},
						},
						Type: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "entity"},
							Sel: &ast.Ident{Name: "CalculateResponse"},
						},
					},
					{
						Names: []*ast.Ident{
							{Name: "err"},
						},
						Type: &ast.Ident{Name: "error"},
					},
				},
			},
		},
	}

	tests := []struct {
		name         string
		depth        int
		node         ast.Node
		expectedNode ast.Node
		expectedFunc funcType
	}{
		{
			name:         "ExportedFunc",
			depth:        2,
			node:         exportedFunc,
			expectedNode: exportedFunc,
			expectedFunc: funcType{
				Exported: true,
				Name:     "NewController",
				Receiver: nil,
				Inputs: fields{
					{
						Names:   []string{"ug"},
						Star:    true,
						Package: "gateway",
						Type:    "UserGateway",
					},
				},
				Outputs: fields{
					{
						Names:   nil,
						Star:    false,
						Package: "",
						Type:    "Controller",
					},
					{
						Names:   nil,
						Star:    false,
						Package: "",
						Type:    "error",
					},
				},
			},
		},
		{
			name:         "UnexportedFunc",
			depth:        2,
			node:         unexportedFunc,
			expectedNode: unexportedFunc,
			expectedFunc: funcType{
				Exported: false,
				Name:     "newController",
				Receiver: nil,
				Inputs: fields{
					{
						Names:   []string{"ug"},
						Star:    true,
						Package: "gateway",
						Type:    "UserGateway",
					},
				},
				Outputs: fields{
					{
						Names:   nil,
						Star:    true,
						Package: "",
						Type:    "controller",
					},
					{
						Names:   nil,
						Star:    false,
						Package: "",
						Type:    "error",
					},
				},
			},
		},
		{
			name:         "ExportedMethod",
			depth:        2,
			node:         exportedMethod,
			expectedNode: exportedMethod,
			expectedFunc: funcType{
				Exported: true,
				Name:     "Calculate",
				Receiver: &receiver{
					Name: "c",
					Star: true,
					Type: "controller",
				},
				Inputs: fields{
					{
						Names:   []string{"a", "b"},
						Star:    true,
						Package: "",
						Type:    "int",
					},
				},
				Outputs: fields{
					{
						Names:   nil,
						Star:    true,
						Package: "entity",
						Type:    "CalculateResponse",
					},
					{
						Names:   nil,
						Star:    false,
						Package: "",
						Type:    "error",
					},
				},
			},
		},
		{
			name:         "UnexportedMethod",
			depth:        2,
			node:         unexportedMethod,
			expectedNode: unexportedMethod,
			expectedFunc: funcType{
				Exported: false,
				Name:     "calculate",
				Receiver: &receiver{
					Name: "c",
					Star: true,
					Type: "controller",
				},
				Inputs: fields{
					{
						Names:   []string{"a", "b"},
						Star:    false,
						Package: "",
						Type:    "int",
					},
				},
				Outputs: fields{
					{
						Names:   []string{"resp"},
						Star:    false,
						Package: "entity",
						Type:    "CalculateResponse",
					},
					{
						Names:   []string{"err"},
						Star:    false,
						Package: "",
						Type:    "error",
					},
				},
			},
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

			node := m.Modify(tc.node)

			assert.Equal(t, tc.expectedNode, node)
			assert.Equal(t, tc.expectedFunc, m.outputs.Func)
		})
	}
}
