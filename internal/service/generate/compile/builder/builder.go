package builder

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/moorara/gelato/internal/service/generate/compile/namer"
	"github.com/moorara/gelato/internal/service/generate/compile/node"
)

// TODO: provide an option for handling unexported fields?

// Builder is used for creating declarations for a struct builder.
type Builder struct {
	factory *node.Factory
}

// New creates a new builder.
func New(factory *node.Factory) *Builder {
	return &Builder{
		factory: factory,
	}
}

// CreateDecls creates all declarations for a struct builder.
func (b *Builder) CreateDecls(pkgName, typeName string, node *ast.StructType) []ast.Decl {
	decls := []ast.Decl{}
	decls = append(decls, createFuncDecl(pkgName, typeName))
	decls = append(decls, createBuilderStructDecl(pkgName, typeName))
	decls = append(decls, createBuildFuncDecl(typeName))

	for _, field := range node.Fields.List {
		if len(field.Names) > 0 {
			for _, id := range field.Names {
				// Only consider exported fields
				if namer.IsExported(id.Name) {
					decls = append(decls, createBuilderMethodDecl(typeName, field.Type, id))
				}
			}
		} else {
			// Embedded field
			id := &ast.Ident{
				Name: namer.InferName(field.Type),
			}

			// Only consider exported fields
			if namer.IsExported(id.Name) {
				decls = append(decls, createBuilderMethodDecl(typeName, field.Type, id))
			}
		}
	}

	decls = append(decls, createBuilderValueDecl(pkgName, typeName))
	decls = append(decls, createBuilderPointerDecl(pkgName, typeName))

	return decls
}

func createFuncDecl(pkgName, typeName string) ast.Decl {
	return &ast.FuncDecl{
		Name: &ast.Ident{Name: typeName},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.SelectorExpr{
							X:   &ast.Ident{Name: pkgName},
							Sel: &ast.Ident{Name: typeName},
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
										Name: "Build" + typeName,
									},
								},
								Sel: &ast.Ident{Name: "Value"},
							},
						},
					},
				},
			},
		},
	}
}

func createBuilderStructDecl(pkgName, typeName string) ast.Decl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{
					Name: typeName + "Builder",
				},
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "v"},
								},
								Type: &ast.SelectorExpr{
									X:   &ast.Ident{Name: pkgName},
									Sel: &ast.Ident{Name: typeName},
								},
							},
						},
					},
				},
			},
		},
	}
}

func createBuildFuncDecl(typeName string) ast.Decl {
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: "Build" + typeName,
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.Ident{Name: typeName + "Builder"},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						// TODO:
						&ast.CompositeLit{
							Type: &ast.Ident{Name: typeName + "Builder"},
						},
					},
				},
			},
		},
	}
}

func createBuilderMethodDecl(typeName string, typ ast.Expr, id *ast.Ident) ast.Decl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{Name: "b"},
					},
					Type: &ast.Ident{Name: typeName + "Builder"},
				},
			},
		},
		Name: &ast.Ident{
			Name: fmt.Sprintf("With%s", id.Name),
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							&ast.Ident{Name: namer.ConvertToUnexported(id.Name)},
						},
						Type: typ,
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.Ident{Name: typeName + "Builder"},
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
							Sel: &ast.Ident{Name: id.Name},
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.Ident{Name: id.Name},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.Ident{Name: "b"},
					},
				},
			},
		},
	}
}

func createBuilderValueDecl(pkgName, typeName string) ast.Decl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{Name: "b"},
					},
					Type: &ast.Ident{Name: typeName + "Builder"},
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
							X:   &ast.Ident{Name: pkgName},
							Sel: &ast.Ident{Name: typeName},
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
	}
}

func createBuilderPointerDecl(pkgName, typeName string) ast.Decl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{Name: "b"},
					},
					Type: &ast.Ident{Name: typeName + "Builder"},
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
								X:   &ast.Ident{Name: pkgName},
								Sel: &ast.Ident{Name: typeName},
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
	}
}
