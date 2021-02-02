package node

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFactory(t *testing.T) {
	f := NewFactory()

	assert.Equal(t, 1, f.offset)
}

func TestFactory_IncOffset(t *testing.T) {
	f := &Factory{offset: 1}
	f.IncOffset(2)

	assert.Equal(t, 3, f.offset)
}

func TestFactory_PackagePos(t *testing.T) {
	f := &Factory{offset: 1}
	pos := f.PackagePos()

	assert.Equal(t, token.Pos(1), pos)
}

func TestFactory_Comment(t *testing.T) {
	f := &Factory{offset: 1}
	comment := f.Comment("// Comment")

	expected := &ast.Comment{
		Slash: 1,
		Text:  "// Comment",
	}

	assert.Equal(t, expected, comment)
}

func TestFactory_Ident(t *testing.T) {
	f := &Factory{offset: 1}
	ident := f.Ident("foo")

	expected := &ast.Ident{
		NamePos: 1,
		Name:    "foo",
	}

	assert.Equal(t, expected, ident)
}

func TestFactory_ImportDecl(t *testing.T) {
	tests := []struct {
		name         string
		pkgs         []string
		expectedDecl *ast.GenDecl
	}{
		{
			name: "OK",
			pkgs: []string{"fmt", "github.com/octocat/example"},
			expectedDecl: &ast.GenDecl{
				TokPos: 1,
				Tok:    token.IMPORT,
				Lparen: 8,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							ValuePos: 11,
							Value:    `"fmt"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							ValuePos: 18,
							Value:    `"github.com/octocat/example"`,
						},
					},
				},
				Rparen: 48,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := &Factory{offset: 1}
			decl := f.ImportDecl("fmt", "github.com/octocat/example")

			assert.Equal(t, tc.expectedDecl, decl)
		})
	}
}

func TestFactory_AnnotateStructDecl(t *testing.T) {
	tests := []struct {
		name         string
		node         *ast.GenDecl
		expectedNode *ast.GenDecl
	}{
		{
			name: "Request",
			node: &ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{Name: "RequestBuilder"},
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
			expectedNode: &ast.GenDecl{
				TokPos: 2,
				Tok:    token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							NamePos: 7,
							Name:    "RequestBuilder",
						},
						Type: &ast.StructType{
							Struct: 22,
							Fields: &ast.FieldList{
								Opening: 29,
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											{
												NamePos: 30,
												Name:    "v",
											},
										},
										Type: &ast.SelectorExpr{
											X: &ast.Ident{
												NamePos: 32,
												Name:    "lookup",
											},
											Sel: &ast.Ident{
												NamePos: 39,
												Name:    "Request",
											},
										},
									},
								},
								Closing: 48,
							},
						},
					},
				},
			},
		},
		{
			name: "Respose",
			node: &ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{Name: "ResponseBuilder"},
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
			expectedNode: &ast.GenDecl{
				TokPos: 2,
				Tok:    token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							NamePos: 7,
							Name:    "ResponseBuilder",
						},
						Type: &ast.StructType{
							Struct: 23,
							Fields: &ast.FieldList{
								Opening: 30,
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											{
												NamePos: 31,
												Name:    "v",
											},
										},
										Type: &ast.SelectorExpr{
											X: &ast.Ident{
												NamePos: 33,
												Name:    "lookup",
											},
											Sel: &ast.Ident{
												NamePos: 40,
												Name:    "Response",
											},
										},
									},
								},
								Closing: 50,
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := &Factory{offset: 1}
			f.AnnotateStructDecl(tc.node)

			assert.Equal(t, tc.expectedNode, tc.node)
		})
	}
}

func TestFactory_AnnotateFuncDecl(t *testing.T) {
	tests := []struct {
		name         string
		node         *ast.FuncDecl
		expectedNode *ast.FuncDecl
	}{
		{
			name: "OK",
			node: &ast.FuncDecl{
				Name: &ast.Ident{Name: "Get"},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "ctx"},
								},
								Type: &ast.SelectorExpr{
									X:   &ast.Ident{Name: "context"},
									Sel: &ast.Ident{Name: "Context"},
								},
							},
							{
								Names: []*ast.Ident{
									{Name: "req"},
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "entity"},
										Sel: &ast.Ident{Name: "Request"},
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
										Sel: &ast.Ident{Name: "Response"},
									},
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
			expectedNode: &ast.FuncDecl{
				Name: &ast.Ident{
					NamePos: 7,
					Name:    "Get",
				},
				Type: &ast.FuncType{
					Func: 2,
					Params: &ast.FieldList{
						Opening: 10,
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										NamePos: 11,
										Name:    "ctx",
									},
								},
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										NamePos: 15,
										Name:    "context",
									},
									Sel: &ast.Ident{
										NamePos: 23,
										Name:    "Context",
									},
								},
							},
							{
								Names: []*ast.Ident{
									{
										NamePos: 32,
										Name:    "req",
									},
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											NamePos: 37,
											Name:    "entity",
										},
										Sel: &ast.Ident{
											NamePos: 44,
											Name:    "Request",
										},
									},
								},
							},
						},
						Closing: 53,
					},
					Results: &ast.FieldList{
						Opening: 55,
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											NamePos: 57,
											Name:    "entity",
										},
										Sel: &ast.Ident{
											NamePos: 64,
											Name:    "Response",
										},
									},
								},
							},
							{
								Type: &ast.Ident{
									NamePos: 74,
									Name:    "error",
								},
							},
						},
						Closing: 81,
					},
				},
				Body: &ast.BlockStmt{
					Lbrace: 83,
					List: []ast.Stmt{
						&ast.ReturnStmt{
							Return: 86,
							Results: []ast.Expr{
								&ast.Ident{
									NamePos: 93,
									Name:    "nil",
								},
								&ast.Ident{
									NamePos: 97,
									Name:    "nil",
								},
							},
						},
					},
					Rbrace: 102,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := &Factory{offset: 1}
			f.AnnotateFuncDecl(tc.node)

			assert.Equal(t, tc.expectedNode, tc.node)
		})
	}
}

func TestFactory_AnnotateFieldList(t *testing.T) {
	tests := []struct {
		name         string
		node         *ast.FieldList
		expectedNode *ast.FieldList
	}{
		{
			name: "Params",
			node: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{Name: "ctx"},
						},
						Type: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "context"},
							Sel: &ast.Ident{Name: "Context"},
						},
					},
					{
						Names: []*ast.Ident{
							{Name: "req"},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "entity"},
								Sel: &ast.Ident{Name: "Request"},
							},
						},
					},
				},
			},
			expectedNode: &ast.FieldList{
				Opening: 1,
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								NamePos: 2,
								Name:    "ctx",
							},
						},
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								NamePos: 6,
								Name:    "context",
							},
							Sel: &ast.Ident{
								NamePos: 14,
								Name:    "Context",
							},
						},
					},
					{
						Names: []*ast.Ident{
							{
								NamePos: 23,
								Name:    "req",
							},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									NamePos: 28,
									Name:    "entity",
								},
								Sel: &ast.Ident{
									NamePos: 35,
									Name:    "Request",
								},
							},
						},
					},
				},
				Closing: 44,
			},
		},
		{
			name: "Results",
			node: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "entity"},
								Sel: &ast.Ident{Name: "Response"},
							},
						},
					},
					{
						Type: &ast.Ident{Name: "error"},
					},
				},
			},
			expectedNode: &ast.FieldList{
				Opening: 1,
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									NamePos: 3,
									Name:    "entity",
								},
								Sel: &ast.Ident{
									NamePos: 10,
									Name:    "Response",
								},
							},
						},
					},
					{
						Type: &ast.Ident{
							NamePos: 20,
							Name:    "error",
						},
					},
				},
				Closing: 27,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := &Factory{offset: 1}
			f.AnnotateFieldList(tc.node)

			assert.Equal(t, tc.expectedNode, tc.node)
		})
	}
}

func TestFactory_AnnotateBlockStmt(t *testing.T) {
	tests := []struct {
		name         string
		node         *ast.BlockStmt
		expectedNode *ast.BlockStmt
	}{
		{
			name: "OK",
			node: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.CallExpr{
										Fun: &ast.Ident{Name: "BuildFoo"},
									},
									Sel: &ast.Ident{Name: "Value"},
								},
							},
						},
					},
				},
			},
			expectedNode: &ast.BlockStmt{
				Lbrace: 1,
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Return: 4,
						Results: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.CallExpr{
										Fun: &ast.Ident{
											NamePos: 11,
											Name:    "BuildFoo",
										},
										Lparen: 19,
										Rparen: 20,
									},
									Sel: &ast.Ident{
										NamePos: 22,
										Name:    "Value",
									},
								},
								Lparen: 27,
								Rparen: 28,
							},
						},
					},
				},
				Rbrace: 31,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := &Factory{offset: 1}
			f.AnnotateBlockStmt(tc.node)

			assert.Equal(t, tc.expectedNode, tc.node)
		})
	}
}
