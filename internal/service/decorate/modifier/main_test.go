package modifier

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/log"
)

func TestMainModifier(t *testing.T) {
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
			Name: "main",
		},
		Decls: []ast.Decl{
			// Imports
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Name: nil,
						Path: &ast.BasicLit{
							Value: "flag",
						},
					},
					&ast.ImportSpec{
						Name: &ast.Ident{
							Name: "ctrl",
						},
						Path: &ast.BasicLit{
							Value: "github.com/octocat/service/internal/controller",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Value: "github.com/octocat/service/internal/gateway",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Value: "github.com/octocat/service/internal/handler",
						},
					},
					&ast.ImportSpec{
						Name: &ast.Ident{
							Name: "repo",
						},
						Path: &ast.BasicLit{
							Value: "github.com/octocat/service/internal/repository",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Value: "github.com/octocat/service/internal/server",
						},
					},
				},
			},
			// main Function
			&ast.FuncDecl{
				Name: &ast.Ident{Name: "main"},
				Type: &ast.FuncType{
					Params:  &ast.FieldList{},
					Results: nil,
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						// Gateway
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{Name: "g"},
								&ast.Ident{Name: "err"},
							},
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "gateway"},
										Sel: &ast.Ident{Name: "NewGateway"},
									},
								},
							},
						},
						// Repository
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{Name: "r"},
								&ast.Ident{Name: "err"},
							},
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "repository"},
										Sel: &ast.Ident{Name: "NewRepository"},
									},
								},
							},
						},
						// Controller
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{Name: "c"},
								&ast.Ident{Name: "err"},
							},
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "controller"},
										Sel: &ast.Ident{Name: "NewController"},
									},
								},
							},
						},
						// Handler
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{Name: "h"},
								&ast.Ident{Name: "err"},
							},
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "handler"},
										Sel: &ast.Ident{Name: "NewHandler"},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name         string
		module       string
		decDir       string
		node         ast.Node
		expectedNode ast.Node
	}{
		{
			name:         "OK",
			module:       "github.com/octocat/service",
			decDir:       "./build",
			node:         fileNode,
			expectedNode: fileNode,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := NewMain(clogger)

			node := m.Modify(tc.module, tc.decDir, tc.node)

			assert.Equal(t, tc.expectedNode, node)
		})
	}
}
