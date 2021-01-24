package compile

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/log"
)

func TestNew(t *testing.T) {
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

	c := New(clogger)

	assert.NotNil(t, c)
	assert.NotNil(t, c.logger)
	assert.NotNil(t, c.funcs.createBuilderDecls)
	assert.NotNil(t, c.funcs.createMockerDecls)
}

func TestCompiler_Compile(t *testing.T) {
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

	inputFile := &ast.File{
		Name: &ast.Ident{
			Name: "lookup",
		},
		Decls: []ast.Decl{
			// Imports
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Value: "fmt",
						},
					},
				},
			},
			// Structs
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{Name: "Request"},
						Type: &ast.StructType{
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
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{Name: "Response"},
						Type: &ast.StructType{
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
					},
				},
			},
			// Interfaces
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{Name: "Service"},
						Type: &ast.InterfaceType{
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
					},
				},
			},
		},
	}

	expectedFile := &ast.File{
		Name: &ast.Ident{
			Name: "lookuptest",
		},
		Decls: []ast.Decl{
			// Imports
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Value: `"github.com/octocat/app/example"`,
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name               string
		createBuilderDecls createBuilderDeclsFunc
		createMockerDecls  createMockerDeclsFunc
		pkgPath            string
		file               *ast.File
		expectedFile       *ast.File
	}{
		{
			name: "OK",
			createBuilderDecls: func(pkgName, typeName string, node *ast.StructType) []ast.Decl {
				return []ast.Decl{}
			},
			createMockerDecls: func(pkgName, typeName string, node *ast.InterfaceType) []ast.Decl {
				return []ast.Decl{}
			},
			pkgPath:      "github.com/octocat/app/example",
			file:         inputFile,
			expectedFile: expectedFile,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &Compiler{logger: clogger}
			c.funcs.createBuilderDecls = tc.createBuilderDecls
			c.funcs.createMockerDecls = tc.createMockerDecls

			file := c.Compile(tc.pkgPath, tc.file)

			assert.Equal(t, tc.expectedFile, file)
		})
	}
}
