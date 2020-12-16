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
				Name: &ast.Ident{Name: "Controller"},
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "Calculate"},
								},
								Type: &ast.FuncType{
									Params: &ast.FieldList{
										List: []*ast.Field{
											{
												Type: &ast.SelectorExpr{
													X:   &ast.Ident{Name: "context"},
													Sel: &ast.Ident{Name: "Context"},
												},
											},
											{
												Type: &ast.StarExpr{
													X: &ast.SelectorExpr{
														X:   &ast.Ident{Name: "entity"},
														Sel: &ast.Ident{Name: "CalculateRequest"},
													},
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
				Name: &ast.Ident{Name: "controller"},
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "userGateway"},
								},
								Type: &ast.SelectorExpr{
									X:   &ast.Ident{Name: "gateway"},
									Sel: &ast.Ident{Name: "UserGateway"},
								},
							},
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name              string
		depth             int
		pkg               string
		node              ast.Node
		expectedNode      ast.Node
		expectedInterface *interfaceType
		expectedStruct    *structType
	}{
		{
			name:  "InvalidGenDecl",
			depth: 2,
			pkg:   "controller",
			node: &ast.GenDecl{
				Tok: token.IMPORT,
			},
			expectedNode: &ast.GenDecl{
				Tok: token.IMPORT,
			},
			expectedInterface: nil,
			expectedStruct:    nil,
		},
		{
			name:         "InterfaceGenDecl",
			depth:        2,
			pkg:          "controller",
			node:         interfaceGenDecl,
			expectedNode: interfaceGenDecl,
			expectedInterface: &interfaceType{
				Exported: true,
				Name:     "Controller",
			},
			expectedStruct: nil,
		},
		{
			name:              "StructGenDecl",
			depth:             2,
			pkg:               "controller",
			node:              structGenDecl,
			expectedNode:      structGenDecl,
			expectedInterface: nil,
			expectedStruct: &structType{
				Exported: false,
				Name:     "controller",
			},
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

			node := m.Modify(tc.pkg, tc.node)

			assert.Equal(t, tc.expectedNode, node)
			assert.Equal(t, tc.expectedInterface, m.outputs.Interface)
			assert.Equal(t, tc.expectedStruct, m.outputs.Struct)
		})
	}
}
