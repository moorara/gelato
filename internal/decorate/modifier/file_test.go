package modifier

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/log"
)

func TestFileModifier(t *testing.T) {
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

	fileNode := &ast.File{
		Name: &ast.Ident{
			Name: "controller",
		},
		Decls: []ast.Decl{
			// Imports
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Name: nil,
						Path: &ast.BasicLit{
							Value: "context",
						},
					},
					&ast.ImportSpec{
						Name: nil,
						Path: &ast.BasicLit{
							Value: "github.com/octokit/service/internal/entity",
						},
					},
					&ast.ImportSpec{
						Name: nil,
						Path: &ast.BasicLit{
							Value: "github.com/octokit/service/internal/gateway",
						},
					},
					&ast.ImportSpec{
						Name: nil,
						Path: &ast.BasicLit{
							Value: "github.com/octokit/service/internal/repository",
						},
					},
				},
			},
			// Interface
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "Controller",
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
			},
			// Struct
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "controller",
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
			},
			// Exported Function
			&ast.FuncDecl{
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
			},
			// Unexported Function
			&ast.FuncDecl{
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
			},
			// Exported Method
			&ast.FuncDecl{
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
			},
			// UnexportedMethod
			&ast.FuncDecl{
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
			},
		},
	}

	tests := []struct {
		name         string
		depth        int
		module       string
		dir          string
		node         ast.Node
		expectedNode ast.Node
	}{
		{
			name:         "FileNode",
			depth:        2,
			module:       "github.com/octocat/Hello-World",
			dir:          ".build",
			node:         fileNode,
			expectedNode: fileNode,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := NewFile(tc.depth, clogger)

			node := m.Modify(tc.module, tc.dir, tc.node)

			assert.Equal(t, tc.expectedNode, node)
		})
	}
}
