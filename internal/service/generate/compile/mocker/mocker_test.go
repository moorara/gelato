package mocker

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMockerDecls(t *testing.T) {
	tests := []struct {
		name          string
		pkgName       string
		typeName      string
		node          *ast.InterfaceType
		expectedDecls []ast.Decl
	}{
		{
			name:     "Service",
			pkgName:  "lookup",
			typeName: "Service",
			node: &ast.InterfaceType{
				Methods: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								&ast.Ident{Name: "Lookup"},
							},
							Type: &ast.FuncType{
								Params: &ast.FieldList{
									List: []*ast.Field{
										{
											Type: &ast.StarExpr{
												X: &ast.Ident{Name: "Request"},
											},
										},
									},
								},
								Results: &ast.FieldList{
									List: []*ast.Field{
										{
											Type: &ast.StarExpr{
												X: &ast.Ident{Name: "Response"},
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
			expectedDecls: []ast.Decl{
				// Mocker struct
				&ast.GenDecl{
					Tok: token.TYPE,
					Specs: []ast.Spec{
						&ast.TypeSpec{
							Name: &ast.Ident{
								Name: "ServiceMocker",
							},
							Type: &ast.StructType{
								Fields: &ast.FieldList{
									List: []*ast.Field{
										{
											Names: []*ast.Ident{
												{Name: "t"},
											},
											Type: &ast.StarExpr{
												X: &ast.SelectorExpr{
													X:   &ast.Ident{Name: "testing"},
													Sel: &ast.Ident{Name: "T"},
												},
											},
										},
										{
											Names: []*ast.Ident{
												{Name: "exps"},
											},
											Type: &ast.StarExpr{
												X: &ast.Ident{Name: "ServiceExpectations"},
											},
										},
									},
								},
							},
						},
					},
				},
				// Mock func
				&ast.FuncDecl{
					Name: &ast.Ident{
						Name: "MockService",
					},
					Type: &ast.FuncType{
						Params: &ast.FieldList{
							List: []*ast.Field{
								{
									Names: []*ast.Ident{
										{Name: "t"},
									},
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X:   &ast.Ident{Name: "testing"},
											Sel: &ast.Ident{Name: "T"},
										},
									},
								},
							},
						},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.StarExpr{
										X: &ast.Ident{Name: "ServiceMocker"},
									},
								},
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.UnaryExpr{
										Op: token.AND,
										X: &ast.CompositeLit{
											Type: &ast.Ident{Name: "ServiceMocker"},
											Elts: []ast.Expr{
												&ast.KeyValueExpr{
													Key:   &ast.Ident{Name: "t"},
													Value: &ast.Ident{Name: "t"},
												},
												&ast.KeyValueExpr{
													Key: &ast.Ident{Name: "exps"},
													Value: &ast.CallExpr{
														Fun: &ast.Ident{Name: "new"},
														Args: []ast.Expr{
															&ast.Ident{Name: "ServiceExpectations"},
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
				},
				// Mocker Expect method
				&ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "m"},
								},
								Type: &ast.StarExpr{
									X: &ast.Ident{Name: "ServiceMocker"},
								},
							},
						},
					},
					Name: &ast.Ident{
						Name: "Expect",
					},
					Type: &ast.FuncType{
						Params: &ast.FieldList{},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.StarExpr{
										X: &ast.Ident{Name: "ServiceExpectations"},
									},
								},
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.SelectorExpr{
										X:   &ast.Ident{Name: "m"},
										Sel: &ast.Ident{Name: "exps"},
									},
								},
							},
						},
					},
				},
				// Mocker Assert method
				&ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "m"},
								},
								Type: &ast.StarExpr{
									X: &ast.Ident{Name: "ServiceMocker"},
								},
							},
						},
					},
					Name: &ast.Ident{
						Name: "Assert",
					},
					Type: &ast.FuncType{
						Params: &ast.FieldList{},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{Name: "errs"},
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.Ident{Name: "new"},
										Args: []ast.Expr{
											&ast.SelectorExpr{
												X:   &ast.Ident{Name: "bytes"},
												Sel: &ast.Ident{Name: "Buffer"},
											},
										},
									},
								},
							},
							// &ast.RangeStmt{ TODO: },
							&ast.IfStmt{
								Init: &ast.AssignStmt{
									Lhs: []ast.Expr{
										&ast.Ident{Name: "s"},
									},
									Tok: token.DEFINE,
									Rhs: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   &ast.Ident{Name: "errs"},
												Sel: &ast.Ident{Name: "String"},
											},
										},
									},
								},
								Cond: &ast.BinaryExpr{
									X:  &ast.Ident{Name: "s"},
									Op: token.NEQ,
									Y: &ast.BasicLit{
										Kind:  token.STRING,
										Value: `""`,
									},
								},
								Body: &ast.BlockStmt{
									List: []ast.Stmt{
										&ast.ExprStmt{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.SelectorExpr{
														X:   &ast.Ident{Name: "m"},
														Sel: &ast.Ident{Name: "t"},
													},
													Sel: &ast.Ident{Name: "Fatal"},
												},
												Args: []ast.Expr{
													&ast.Ident{Name: "s"},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				// Mocker Impl method
				&ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "m"},
								},
								Type: &ast.StarExpr{
									X: &ast.Ident{Name: "ServiceMocker"},
								},
							},
						},
					},
					Name: &ast.Ident{
						Name: "Impl",
					},
					Type: &ast.FuncType{
						Params: &ast.FieldList{},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "lookup"},
										Sel: &ast.Ident{Name: "Service"},
									},
								},
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.UnaryExpr{
										Op: token.AND,
										X: &ast.CompositeLit{
											Type: &ast.Ident{Name: "ServiceImpl"},
											Elts: []ast.Expr{
												&ast.KeyValueExpr{
													Key: &ast.Ident{Name: "t"},
													Value: &ast.SelectorExpr{
														X:   &ast.Ident{Name: "m"},
														Sel: &ast.Ident{Name: "t"},
													},
												},
												&ast.KeyValueExpr{
													Key: &ast.Ident{Name: "exps"},
													Value: &ast.SelectorExpr{
														X:   &ast.Ident{Name: "m"},
														Sel: &ast.Ident{Name: "exps"},
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
				// Expectations struct
				&ast.GenDecl{
					Tok: token.TYPE,
					Specs: []ast.Spec{
						&ast.TypeSpec{
							Name: &ast.Ident{
								Name: "ServiceExpectations",
							},
							Type: &ast.StructType{
								Fields: &ast.FieldList{
									List: []*ast.Field{
										{
											Names: []*ast.Ident{
												{Name: "lookupExpectations"},
											},
											Type: &ast.ArrayType{
												Elt: &ast.StarExpr{
													X: &ast.Ident{Name: "LookupExpectation"},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				// Expectations methods
				&ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "e"},
								},
								Type: &ast.StarExpr{
									X: &ast.Ident{Name: "ServiceExpectations"},
								},
							},
						},
					},
					Name: &ast.Ident{Name: "Lookup"},
					Type: &ast.FuncType{
						Params: &ast.FieldList{},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.StarExpr{
										X: &ast.Ident{Name: "LookupExpectation"},
									},
								},
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{Name: "expectation"},
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.Ident{Name: "new"},
										Args: []ast.Expr{
											&ast.Ident{Name: "LookupExpectation"},
										},
									},
								},
							},
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.SelectorExpr{
										X:   &ast.Ident{Name: "e"},
										Sel: &ast.Ident{Name: "lookupExpectations"},
									},
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.Ident{Name: "append"},
										Args: []ast.Expr{
											&ast.SelectorExpr{
												X:   &ast.Ident{Name: "e"},
												Sel: &ast.Ident{Name: "lookupExpectations"},
											},
											&ast.Ident{Name: "expectation"},
										},
									},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.Ident{Name: "expectation"},
								},
							},
						},
					},
				},
				// Expectation structs
				&ast.GenDecl{
					Tok: token.TYPE,
					Specs: []ast.Spec{
						&ast.TypeSpec{
							Name: &ast.Ident{
								Name: "LookupExpectation",
							},
							Type: &ast.StructType{
								Fields: &ast.FieldList{
									List: []*ast.Field{
										{
											Names: []*ast.Ident{
												{Name: "isCalled"},
											},
											Type: &ast.Ident{Name: "bool"},
										},
										{
											Names: []*ast.Ident{
												{Name: "inputs"},
											},
											Type: &ast.StarExpr{
												X: &ast.Ident{Name: "lookupInputs"},
											},
										},
										{
											Names: []*ast.Ident{
												{Name: "outputs"},
											},
											Type: &ast.StarExpr{
												X: &ast.Ident{Name: "lookupOutputs"},
											},
										},
										{
											Names: []*ast.Ident{
												{Name: "callback"},
											},
											Type: &ast.FuncType{
												Params: &ast.FieldList{
													List: []*ast.Field{
														{
															Type: &ast.StarExpr{
																X: &ast.Ident{Name: "Request"},
															},
														},
													},
												},
												Results: &ast.FieldList{
													List: []*ast.Field{
														{
															Type: &ast.StarExpr{
																X: &ast.Ident{Name: "Response"},
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
				},
				&ast.GenDecl{
					Tok: token.TYPE,
					Specs: []ast.Spec{
						&ast.TypeSpec{
							Name: &ast.Ident{
								Name: "lookupInputs",
							},
							Type: &ast.StructType{
								Fields: &ast.FieldList{
									List: []*ast.Field{
										{
											Names: []*ast.Ident{
												{Name: "request"},
											},
											Type: &ast.StarExpr{
												X: &ast.Ident{Name: "Request"},
											},
										},
									},
								},
							},
						},
					},
				},
				&ast.GenDecl{
					Tok: token.TYPE,
					Specs: []ast.Spec{
						&ast.TypeSpec{
							Name: &ast.Ident{
								Name: "lookupOutputs",
							},
							Type: &ast.StructType{
								Fields: &ast.FieldList{
									List: []*ast.Field{
										{
											Names: []*ast.Ident{
												{Name: "response"},
											},
											Type: &ast.StarExpr{
												X: &ast.Ident{Name: "Response"},
											},
										},
										{
											Names: []*ast.Ident{
												{Name: "error"},
											},
											Type: &ast.Ident{Name: "error"},
										},
									},
								},
							},
						},
					},
				},
				// Expectation WithArgs method
				&ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "e"},
								},
								Type: &ast.StarExpr{
									X: &ast.Ident{Name: "LookupExpectation"},
								},
							},
						},
					},
					Name: &ast.Ident{Name: "WithArgs"},
					Type: &ast.FuncType{
						Params: &ast.FieldList{
							List: []*ast.Field{
								{
									Names: []*ast.Ident{
										{Name: "request"},
									},
									Type: &ast.StarExpr{
										X: &ast.Ident{Name: "Request"},
									},
								},
							},
						},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.StarExpr{
										X: &ast.Ident{Name: "LookupExpectation"},
									},
								},
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.SelectorExpr{
										X:   &ast.Ident{Name: "e"},
										Sel: &ast.Ident{Name: "inputs"},
									},
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.UnaryExpr{
										Op: token.AND,
										X: &ast.CompositeLit{
											Type: &ast.Ident{Name: "lookupInputs"},
											Elts: []ast.Expr{
												&ast.KeyValueExpr{
													Key:   &ast.Ident{Name: "request"},
													Value: &ast.Ident{Name: "request"},
												},
											},
										},
									},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.Ident{Name: "e"},
								},
							},
						},
					},
				},
				// Expectation Return method
				&ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "e"},
								},
								Type: &ast.StarExpr{
									X: &ast.Ident{Name: "LookupExpectation"},
								},
							},
						},
					},
					Name: &ast.Ident{Name: "Return"},
					Type: &ast.FuncType{
						Params: &ast.FieldList{
							List: []*ast.Field{
								{
									Names: []*ast.Ident{
										{Name: "response"},
									},
									Type: &ast.StarExpr{
										X: &ast.Ident{Name: "Response"},
									},
								},
								{
									Names: []*ast.Ident{
										{Name: "error"},
									},
									Type: &ast.Ident{Name: "error"},
								},
							},
						},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.StarExpr{
										X: &ast.Ident{Name: "LookupExpectation"},
									},
								},
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.SelectorExpr{
										X:   &ast.Ident{Name: "e"},
										Sel: &ast.Ident{Name: "outputs"},
									},
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.UnaryExpr{
										Op: token.AND,
										X: &ast.CompositeLit{
											Type: &ast.Ident{Name: "lookupOutputs"},
											Elts: []ast.Expr{
												&ast.KeyValueExpr{
													Key:   &ast.Ident{Name: "response"},
													Value: &ast.Ident{Name: "response"},
												},
												&ast.KeyValueExpr{
													Key:   &ast.Ident{Name: "error"},
													Value: &ast.Ident{Name: "error"},
												},
											},
										},
									},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.Ident{Name: "e"},
								},
							},
						},
					},
				},
				// Expectation Call method
				&ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "e"},
								},
								Type: &ast.StarExpr{
									X: &ast.Ident{Name: "LookupExpectation"},
								},
							},
						},
					},
					Name: &ast.Ident{Name: "Call"},
					Type: &ast.FuncType{
						Params: &ast.FieldList{
							List: []*ast.Field{
								{
									Names: []*ast.Ident{
										{Name: "callback"},
									},
									Type: &ast.FuncType{
										Params: &ast.FieldList{
											List: []*ast.Field{
												{
													Type: &ast.StarExpr{
														X: &ast.Ident{Name: "Request"},
													},
												},
											},
										},
										Results: &ast.FieldList{
											List: []*ast.Field{
												{
													Type: &ast.StarExpr{
														X: &ast.Ident{Name: "Response"},
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
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.StarExpr{
										X: &ast.Ident{Name: "LookupExpectation"},
									},
								},
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.SelectorExpr{
										X:   &ast.Ident{Name: "e"},
										Sel: &ast.Ident{Name: "callback"},
									},
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.Ident{Name: "callback"},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.Ident{Name: "e"},
								},
							},
						},
					},
				},
				// Implementation struct
				&ast.GenDecl{
					Tok: token.TYPE,
					Specs: []ast.Spec{
						&ast.TypeSpec{
							Name: &ast.Ident{
								Name: "ServiceImpl",
							},
							Type: &ast.StructType{
								Fields: &ast.FieldList{
									List: []*ast.Field{
										{
											Names: []*ast.Ident{
												{Name: "t"},
											},
											Type: &ast.StarExpr{
												X: &ast.SelectorExpr{
													X:   &ast.Ident{Name: "testing"},
													Sel: &ast.Ident{Name: "T"},
												},
											},
										},
										{
											Names: []*ast.Ident{
												{Name: "exps"},
											},
											Type: &ast.StarExpr{
												X: &ast.Ident{Name: "ServiceExpectations"},
											},
										},
									},
								},
							},
						},
					},
				},
				// Implementation methods
				&ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "i"},
								},
								Type: &ast.StarExpr{
									X: &ast.Ident{Name: "ServiceImpl"},
								},
							},
						},
					},
					Name: &ast.Ident{Name: "Lookup"},
					Type: &ast.FuncType{
						Params: &ast.FieldList{
							List: []*ast.Field{
								{
									Names: []*ast.Ident{
										{Name: "request"},
									},
									Type: &ast.StarExpr{
										X: &ast.Ident{Name: "Request"},
									},
								},
							},
						},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.StarExpr{
										X: &ast.Ident{Name: "Response"},
									},
								},
								{
									Type: &ast.Ident{Name: "error"},
								},
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.Ident{Name: "nil"},
									&ast.Ident{Name: "nil"},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			decls := CreateMockerDecls(tc.pkgName, tc.typeName, tc.node)

			assert.Equal(t, tc.expectedDecls, decls)
		})
	}
}
