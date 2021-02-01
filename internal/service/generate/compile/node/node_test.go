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
										NamePos: 31,
										Name:    "req",
									},
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											NamePos: 35,
											Name:    "entity",
										},
										Sel: &ast.Ident{
											NamePos: 42,
											Name:    "Request",
										},
									},
								},
							},
						},
						Closing: 50,
					},
					Results: &ast.FieldList{
						Opening: 52,
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											NamePos: 53,
											Name:    "entity",
										},
										Sel: &ast.Ident{
											NamePos: 60,
											Name:    "Response",
										},
									},
								},
							},
							{
								Type: &ast.Ident{
									NamePos: 69,
									Name:    "error",
								},
							},
						},
						Closing: 75,
					},
				},
				Body: &ast.BlockStmt{
					Lbrace: 77,
					List: []ast.Stmt{
						&ast.ReturnStmt{
							Return: 80,
							Results: []ast.Expr{
								&ast.Ident{
									NamePos: 87,
									Name:    "nil",
								},
								&ast.Ident{
									NamePos: 91,
									Name:    "nil",
								},
							},
						},
					},
					Rbrace: 95,
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
								NamePos: 22,
								Name:    "req",
							},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									NamePos: 26,
									Name:    "entity",
								},
								Sel: &ast.Ident{
									NamePos: 33,
									Name:    "Request",
								},
							},
						},
					},
				},
				Closing: 41,
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
									NamePos: 2,
									Name:    "entity",
								},
								Sel: &ast.Ident{
									NamePos: 9,
									Name:    "Response",
								},
							},
						},
					},
					{
						Type: &ast.Ident{
							NamePos: 18,
							Name:    "error",
						},
					},
				},
				Closing: 24,
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
				Rbrace: 30,
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
