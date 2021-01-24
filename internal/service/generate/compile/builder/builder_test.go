package builder

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
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
					Name: &ast.Ident{Name: "Request"},
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
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.Ident{
													Name: "BuildRequest",
												},
											},
											Sel: &ast.Ident{Name: "Value"},
										},
									},
								},
							},
						},
					},
				},
				// Builder struct
				&ast.GenDecl{
					Tok: token.TYPE,
					Specs: []ast.Spec{
						&ast.TypeSpec{
							Name: &ast.Ident{
								Name: "RequestBuilder",
							},
							Type: &ast.StructType{
								Fields: &ast.FieldList{
									List: []*ast.Field{
										{
											Names: []*ast.Ident{
												{Name: "v"},
											},
											Type: &ast.SelectorExpr{
												X:   &ast.Ident{Name: "lookup"},
												Sel: &ast.Ident{Name: "Request"},
											},
										},
									},
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
					Name: &ast.Ident{Name: "Response"},
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
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.Ident{
													Name: "BuildResponse",
												},
											},
											Sel: &ast.Ident{Name: "Value"},
										},
									},
								},
							},
						},
					},
				},
				// Builder struct
				&ast.GenDecl{
					Tok: token.TYPE,
					Specs: []ast.Spec{
						&ast.TypeSpec{
							Name: &ast.Ident{
								Name: "ResponseBuilder",
							},
							Type: &ast.StructType{
								Fields: &ast.FieldList{
									List: []*ast.Field{
										{
											Names: []*ast.Ident{
												{Name: "v"},
											},
											Type: &ast.SelectorExpr{
												X:   &ast.Ident{Name: "lookup"},
												Sel: &ast.Ident{Name: "Response"},
											},
										},
									},
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
					Name: &ast.Ident{Name: "Account"},
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
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.Ident{
													Name: "BuildAccount",
												},
											},
											Sel: &ast.Ident{Name: "Value"},
										},
									},
								},
							},
						},
					},
				},
				// Builder struct
				&ast.GenDecl{
					Tok: token.TYPE,
					Specs: []ast.Spec{
						&ast.TypeSpec{
							Name: &ast.Ident{
								Name: "AccountBuilder",
							},
							Type: &ast.StructType{
								Fields: &ast.FieldList{
									List: []*ast.Field{
										{
											Names: []*ast.Ident{
												{Name: "v"},
											},
											Type: &ast.SelectorExpr{
												X:   &ast.Ident{Name: "account"},
												Sel: &ast.Ident{Name: "Account"},
											},
										},
									},
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
					Name: &ast.Ident{Name: "Example"},
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
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.CallExpr{
												Fun: &ast.Ident{
													Name: "BuildExample",
												},
											},
											Sel: &ast.Ident{Name: "Value"},
										},
									},
								},
							},
						},
					},
				},
				// Builder struct
				&ast.GenDecl{
					Tok: token.TYPE,
					Specs: []ast.Spec{
						&ast.TypeSpec{
							Name: &ast.Ident{
								Name: "ExampleBuilder",
							},
							Type: &ast.StructType{
								Fields: &ast.FieldList{
									List: []*ast.Field{
										{
											Names: []*ast.Ident{
												{Name: "v"},
											},
											Type: &ast.SelectorExpr{
												X:   &ast.Ident{Name: "example"},
												Sel: &ast.Ident{Name: "Example"},
											},
										},
									},
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
			decls := CreateBuilderDecls(tc.pkgName, tc.typeName, tc.node)

			assert.Equal(t, tc.expectedDecls, decls)
		})
	}
}
