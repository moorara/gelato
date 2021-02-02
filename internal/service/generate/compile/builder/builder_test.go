package builder

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/service/generate/compile/node"
)

func TestCreateBuilderDecls(t *testing.T) {
	tests := []struct {
		name          string
		pkgName       string
		typeName      string
		node          *ast.StructType
		expectedDecls []ast.Decl
	}{
		{
			name:     "Request",
			pkgName:  "lookup",
			typeName: "Request",
			node: &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								&ast.Ident{Name: "ID"},
							},
							Type: &ast.Ident{Name: "string"},
						},
					},
				},
			},
			expectedDecls: []ast.Decl{
				// Type func
				&ast.FuncDecl{
					Name: &ast.Ident{
						NamePos: 7,
						Name:    "Request",
					},
					Type: &ast.FuncType{
						Func: 2,
						Params: &ast.FieldList{
							Opening: 14,
							Closing: 15,
						},
						Results: &ast.FieldList{
							Opening: 17,
							List: []*ast.Field{
								{
									Type: &ast.SelectorExpr{
										X: &ast.Ident{
											NamePos: 18,
											Name:    "lookup",
										},
										Sel: &ast.Ident{
											NamePos: 25,
											Name:    "Request",
										},
									},
								},
							},
							Closing: 34,
						},
					},
					Body: &ast.BlockStmt{
						Lbrace: 36,
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Return: 39,
								Results: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.Ident{
													NamePos: 46,
													Name:    "BuildRequest",
												},
												Lparen: 58,
												Rparen: 59,
											},
											Sel: &ast.Ident{
												NamePos: 61,
												Name:    "Value",
											},
										},
										Lparen: 66,
										Rparen: 67,
									},
								},
							},
						},
						Rbrace: 70,
					},
				},
				// Builder struct
				&ast.GenDecl{
					TokPos: 73,
					Tok:    token.TYPE,
					Specs: []ast.Spec{
						&ast.TypeSpec{
							Name: &ast.Ident{
								NamePos: 78,
								Name:    "RequestBuilder",
							},
							Type: &ast.StructType{
								Struct: 93,
								Fields: &ast.FieldList{
									Opening: 100,
									List: []*ast.Field{
										{
											Names: []*ast.Ident{
												{
													NamePos: 101,
													Name:    "v",
												},
											},
											Type: &ast.SelectorExpr{
												X: &ast.Ident{
													NamePos: 103,
													Name:    "lookup",
												},
												Sel: &ast.Ident{
													NamePos: 110,
													Name:    "Request",
												},
											},
										},
									},
									Closing: 119,
								},
							},
						},
					},
				},
				// Build func
				&ast.FuncDecl{
					Name: &ast.Ident{
						Name: "BuildRequest",
					},
					Type: &ast.FuncType{
						Params: &ast.FieldList{},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.Ident{Name: "RequestBuilder"},
								},
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.CompositeLit{
										Type: &ast.Ident{Name: "RequestBuilder"},
									},
								},
							},
						},
					},
				},
				// Builder method
				&ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "b"},
								},
								Type: &ast.Ident{Name: "RequestBuilder"},
							},
						},
					},
					Name: &ast.Ident{
						Name: "WithID",
					},
					Type: &ast.FuncType{
						Params: &ast.FieldList{
							List: []*ast.Field{
								{
									Names: []*ast.Ident{
										&ast.Ident{Name: "id"},
									},
									Type: &ast.Ident{Name: "string"},
								},
							},
						},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.Ident{Name: "RequestBuilder"},
								},
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X:   &ast.Ident{Name: "b"},
											Sel: &ast.Ident{Name: "v"},
										},
										Sel: &ast.Ident{Name: "ID"},
									},
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.Ident{Name: "ID"},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.Ident{Name: "b"},
								},
							},
						},
					},
				},
				// Value method
				&ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "b"},
								},
								Type: &ast.Ident{Name: "RequestBuilder"},
							},
						},
					},
					Name: &ast.Ident{Name: "Value"},
					Type: &ast.FuncType{
						Params: &ast.FieldList{},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "lookup"},
										Sel: &ast.Ident{Name: "Request"},
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
										X:   &ast.Ident{Name: "b"},
										Sel: &ast.Ident{Name: "v"},
									},
								},
							},
						},
					},
				},
				// Pointer method
				&ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "b"},
								},
								Type: &ast.Ident{Name: "RequestBuilder"},
							},
						},
					},
					Name: &ast.Ident{Name: "Pointer"},
					Type: &ast.FuncType{
						Params: &ast.FieldList{},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X:   &ast.Ident{Name: "lookup"},
											Sel: &ast.Ident{Name: "Request"},
										},
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
										X: &ast.SelectorExpr{
											X:   &ast.Ident{Name: "b"},
											Sel: &ast.Ident{Name: "v"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:     "Response",
			pkgName:  "lookup",
			typeName: "Response",
			node: &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								&ast.Ident{Name: "Name"},
							},
							Type: &ast.Ident{Name: "string"},
						},
					},
				},
			},
			expectedDecls: []ast.Decl{
				// Type func
				&ast.FuncDecl{
					Name: &ast.Ident{
						NamePos: 7,
						Name:    "Response",
					},
					Type: &ast.FuncType{
						Func: 2,
						Params: &ast.FieldList{
							Opening: 15,
							Closing: 16,
						},
						Results: &ast.FieldList{
							Opening: 18,
							List: []*ast.Field{
								{
									Type: &ast.SelectorExpr{
										X: &ast.Ident{
											NamePos: 19,
											Name:    "lookup",
										},
										Sel: &ast.Ident{
											NamePos: 26,
											Name:    "Response",
										},
									},
								},
							},
							Closing: 36,
						},
					},
					Body: &ast.BlockStmt{
						Lbrace: 38,
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Return: 41,
								Results: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.Ident{
													NamePos: 48,
													Name:    "BuildResponse",
												},
												Lparen: 61,
												Rparen: 62,
											},
											Sel: &ast.Ident{
												NamePos: 64,
												Name:    "Value",
											},
										},
										Lparen: 69,
										Rparen: 70,
									},
								},
							},
						},
						Rbrace: 73,
					},
				},
				// Builder struct
				&ast.GenDecl{
					TokPos: 76,
					Tok:    token.TYPE,
					Specs: []ast.Spec{
						&ast.TypeSpec{
							Name: &ast.Ident{
								NamePos: 81,
								Name:    "ResponseBuilder",
							},
							Type: &ast.StructType{
								Struct: 97,
								Fields: &ast.FieldList{
									Opening: 104,
									List: []*ast.Field{
										{
											Names: []*ast.Ident{
												{
													NamePos: 105,
													Name:    "v",
												},
											},
											Type: &ast.SelectorExpr{
												X: &ast.Ident{
													NamePos: 107,
													Name:    "lookup",
												},
												Sel: &ast.Ident{
													NamePos: 114,
													Name:    "Response",
												},
											},
										},
									},
									Closing: 124,
								},
							},
						},
					},
				},
				// Build func
				&ast.FuncDecl{
					Name: &ast.Ident{
						Name: "BuildResponse",
					},
					Type: &ast.FuncType{
						Params: &ast.FieldList{},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.Ident{Name: "ResponseBuilder"},
								},
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.CompositeLit{
										Type: &ast.Ident{Name: "ResponseBuilder"},
									},
								},
							},
						},
					},
				},
				// Builder method
				&ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "b"},
								},
								Type: &ast.Ident{Name: "ResponseBuilder"},
							},
						},
					},
					Name: &ast.Ident{
						Name: "WithName",
					},
					Type: &ast.FuncType{
						Params: &ast.FieldList{
							List: []*ast.Field{
								{
									Names: []*ast.Ident{
										&ast.Ident{Name: "name"},
									},
									Type: &ast.Ident{Name: "string"},
								},
							},
						},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.Ident{Name: "ResponseBuilder"},
								},
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X:   &ast.Ident{Name: "b"},
											Sel: &ast.Ident{Name: "v"},
										},
										Sel: &ast.Ident{Name: "Name"},
									},
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.Ident{Name: "Name"},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.Ident{Name: "b"},
								},
							},
						},
					},
				},
				// Value method
				&ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "b"},
								},
								Type: &ast.Ident{Name: "ResponseBuilder"},
							},
						},
					},
					Name: &ast.Ident{Name: "Value"},
					Type: &ast.FuncType{
						Params: &ast.FieldList{},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "lookup"},
										Sel: &ast.Ident{Name: "Response"},
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
										X:   &ast.Ident{Name: "b"},
										Sel: &ast.Ident{Name: "v"},
									},
								},
							},
						},
					},
				},
				// Pointer method
				&ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "b"},
								},
								Type: &ast.Ident{Name: "ResponseBuilder"},
							},
						},
					},
					Name: &ast.Ident{Name: "Pointer"},
					Type: &ast.FuncType{
						Params: &ast.FieldList{},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X:   &ast.Ident{Name: "lookup"},
											Sel: &ast.Ident{Name: "Response"},
										},
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
										X: &ast.SelectorExpr{
											X:   &ast.Ident{Name: "b"},
											Sel: &ast.Ident{Name: "v"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:     "EmbeddedStruct",
			pkgName:  "account",
			typeName: "Account",
			node: &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "common"},
								Sel: &ast.Ident{Name: "Address"},
							},
						},
					},
				},
			},
			expectedDecls: []ast.Decl{
				// Type func
				&ast.FuncDecl{
					Name: &ast.Ident{
						NamePos: 7,
						Name:    "Account",
					},
					Type: &ast.FuncType{
						Func: 2,
						Params: &ast.FieldList{
							Opening: 14,
							Closing: 15,
						},
						Results: &ast.FieldList{
							Opening: 17,
							List: []*ast.Field{
								{
									Type: &ast.SelectorExpr{
										X: &ast.Ident{
											NamePos: 18,
											Name:    "account",
										},
										Sel: &ast.Ident{
											NamePos: 26,
											Name:    "Account",
										},
									},
								},
							},
							Closing: 35,
						},
					},
					Body: &ast.BlockStmt{
						Lbrace: 37,
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Return: 40,
								Results: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.Ident{
													NamePos: 47,
													Name:    "BuildAccount",
												},
												Lparen: 59,
												Rparen: 60,
											},
											Sel: &ast.Ident{
												NamePos: 62,
												Name:    "Value",
											},
										},
										Lparen: 67,
										Rparen: 68,
									},
								},
							},
						},
						Rbrace: 71,
					},
				},
				// Builder struct
				&ast.GenDecl{
					TokPos: 74,
					Tok:    token.TYPE,
					Specs: []ast.Spec{
						&ast.TypeSpec{
							Name: &ast.Ident{
								NamePos: 79,
								Name:    "AccountBuilder",
							},
							Type: &ast.StructType{
								Struct: 94,
								Fields: &ast.FieldList{
									Opening: 101,
									List: []*ast.Field{
										{
											Names: []*ast.Ident{
												{
													NamePos: 102,
													Name:    "v",
												},
											},
											Type: &ast.SelectorExpr{
												X: &ast.Ident{
													NamePos: 104,
													Name:    "account",
												},
												Sel: &ast.Ident{
													NamePos: 112,
													Name:    "Account",
												},
											},
										},
									},
									Closing: 121,
								},
							},
						},
					},
				},
				// Build func
				&ast.FuncDecl{
					Name: &ast.Ident{
						Name: "BuildAccount",
					},
					Type: &ast.FuncType{
						Params: &ast.FieldList{},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.Ident{Name: "AccountBuilder"},
								},
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.CompositeLit{
										Type: &ast.Ident{Name: "AccountBuilder"},
									},
								},
							},
						},
					},
				},
				// Builder method
				&ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "b"},
								},
								Type: &ast.Ident{Name: "AccountBuilder"},
							},
						},
					},
					Name: &ast.Ident{
						Name: "WithAddress",
					},
					Type: &ast.FuncType{
						Params: &ast.FieldList{
							List: []*ast.Field{
								{
									Names: []*ast.Ident{
										&ast.Ident{Name: "address"},
									},
									Type: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "common"},
										Sel: &ast.Ident{Name: "Address"},
									},
								},
							},
						},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.Ident{Name: "AccountBuilder"},
								},
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X:   &ast.Ident{Name: "b"},
											Sel: &ast.Ident{Name: "v"},
										},
										Sel: &ast.Ident{Name: "Address"},
									},
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.Ident{Name: "Address"},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.Ident{Name: "b"},
								},
							},
						},
					},
				},
				// Value method
				&ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "b"},
								},
								Type: &ast.Ident{Name: "AccountBuilder"},
							},
						},
					},
					Name: &ast.Ident{Name: "Value"},
					Type: &ast.FuncType{
						Params: &ast.FieldList{},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "account"},
										Sel: &ast.Ident{Name: "Account"},
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
										X:   &ast.Ident{Name: "b"},
										Sel: &ast.Ident{Name: "v"},
									},
								},
							},
						},
					},
				},
				// Pointer method
				&ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "b"},
								},
								Type: &ast.Ident{Name: "AccountBuilder"},
							},
						},
					},
					Name: &ast.Ident{Name: "Pointer"},
					Type: &ast.FuncType{
						Params: &ast.FieldList{},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X:   &ast.Ident{Name: "account"},
											Sel: &ast.Ident{Name: "Account"},
										},
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
										X: &ast.SelectorExpr{
											X:   &ast.Ident{Name: "b"},
											Sel: &ast.Ident{Name: "v"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:     "UnexportedField",
			pkgName:  "example",
			typeName: "Example",
			node: &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								&ast.Ident{Name: "internal"},
							},
							Type: &ast.Ident{Name: "string"},
						},
					},
				},
			},
			expectedDecls: []ast.Decl{
				// Type func
				&ast.FuncDecl{
					Name: &ast.Ident{
						NamePos: 7,
						Name:    "Example",
					},
					Type: &ast.FuncType{
						Func: 2,
						Params: &ast.FieldList{
							Opening: 14,
							Closing: 15,
						},
						Results: &ast.FieldList{
							Opening: 17,
							List: []*ast.Field{
								{
									Type: &ast.SelectorExpr{
										X: &ast.Ident{
											NamePos: 18,
											Name:    "example",
										},
										Sel: &ast.Ident{
											NamePos: 26,
											Name:    "Example",
										},
									},
								},
							},
							Closing: 35,
						},
					},
					Body: &ast.BlockStmt{
						Lbrace: 37,
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Return: 40,
								Results: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.Ident{
													NamePos: 47,
													Name:    "BuildExample",
												},
												Lparen: 59,
												Rparen: 60,
											},
											Sel: &ast.Ident{
												NamePos: 62,
												Name:    "Value",
											},
										},
										Lparen: 67,
										Rparen: 68,
									},
								},
							},
						},
						Rbrace: 71,
					},
				},
				// Builder struct
				&ast.GenDecl{
					TokPos: 74,
					Tok:    token.TYPE,
					Specs: []ast.Spec{
						&ast.TypeSpec{
							Name: &ast.Ident{
								NamePos: 79,
								Name:    "ExampleBuilder",
							},
							Type: &ast.StructType{
								Struct: 94,
								Fields: &ast.FieldList{
									Opening: 101,
									List: []*ast.Field{
										{
											Names: []*ast.Ident{
												{
													NamePos: 102,
													Name:    "v",
												},
											},
											Type: &ast.SelectorExpr{
												X: &ast.Ident{
													NamePos: 104,
													Name:    "example",
												},
												Sel: &ast.Ident{
													NamePos: 112,
													Name:    "Example",
												},
											},
										},
									},
									Closing: 121,
								},
							},
						},
					},
				},
				// Build func
				&ast.FuncDecl{
					Name: &ast.Ident{
						Name: "BuildExample",
					},
					Type: &ast.FuncType{
						Params: &ast.FieldList{},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.Ident{Name: "ExampleBuilder"},
								},
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.CompositeLit{
										Type: &ast.Ident{Name: "ExampleBuilder"},
									},
								},
							},
						},
					},
				},
				// Builder method
				// Value method
				&ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "b"},
								},
								Type: &ast.Ident{Name: "ExampleBuilder"},
							},
						},
					},
					Name: &ast.Ident{Name: "Value"},
					Type: &ast.FuncType{
						Params: &ast.FieldList{},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "example"},
										Sel: &ast.Ident{Name: "Example"},
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
										X:   &ast.Ident{Name: "b"},
										Sel: &ast.Ident{Name: "v"},
									},
								},
							},
						},
					},
				},
				// Pointer method
				&ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "b"},
								},
								Type: &ast.Ident{Name: "ExampleBuilder"},
							},
						},
					},
					Name: &ast.Ident{Name: "Pointer"},
					Type: &ast.FuncType{
						Params: &ast.FieldList{},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X:   &ast.Ident{Name: "example"},
											Sel: &ast.Ident{Name: "Example"},
										},
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
										X: &ast.SelectorExpr{
											X:   &ast.Ident{Name: "b"},
											Sel: &ast.Ident{Name: "v"},
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

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			factory := node.NewFactory()
			builder := New(factory)
			decls := builder.CreateDecls(tc.pkgName, tc.typeName, tc.node)

			assert.Equal(t, tc.expectedDecls, decls)
		})
	}
}
