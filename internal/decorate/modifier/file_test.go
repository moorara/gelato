package modifier

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/ast/astutil"

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
			},
			// Struct
			&ast.GenDecl{
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
			},
			// Exported Function
			&ast.FuncDecl{
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
			},
			// Unexported Function
			&ast.FuncDecl{
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
			},
			// Exported Method
			&ast.FuncDecl{
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
			},
			// UnexportedMethod
			&ast.FuncDecl{
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
			},
		},
	}

	tests := []struct {
		name         string
		depth        int
		node         ast.Node
		expectedNode ast.Node
	}{
		{
			name:         "FileNode",
			depth:        2,
			node:         fileNode,
			expectedNode: fileNode,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := NewFile(tc.depth, clogger)
			node := astutil.Apply(tc.node, m.Pre, m.Post)

			assert.Equal(t, tc.expectedNode, node)
		})
	}
}
