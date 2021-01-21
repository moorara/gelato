package compile

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/log"
)

func TestCompiler(t *testing.T) {
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
			Name: "exampletest",
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
			// Structs
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{Name: "RequestBuilder"},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Type: &ast.SelectorExpr{
											X:   &ast.Ident{Name: "example"},
											Sel: &ast.Ident{Name: "Request"},
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
						Name: &ast.Ident{Name: "ResponseBuilder"},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Type: &ast.SelectorExpr{
											X:   &ast.Ident{Name: "example"},
											Sel: &ast.Ident{Name: "Response"},
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

	tests := []struct {
		name         string
		pkgPath      string
		pkgName      string
		file         *ast.File
		expectedFile *ast.File
	}{
		{
			name:         "OK",
			pkgPath:      "github.com/octocat/app/example",
			pkgName:      "example",
			file:         inputFile,
			expectedFile: expectedFile,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := New(clogger)

			file := c.Compile(tc.pkgPath, tc.pkgName, tc.file)

			assert.Equal(t, tc.expectedFile, file)
		})
	}
}
