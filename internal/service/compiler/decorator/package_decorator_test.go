package decorator

import (
	"go/ast"
	"go/token"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/service/compiler"
)

func TestPkgDecorator_Package(t *testing.T) {
	tests := []struct {
		name             string
		info             *compiler.PackageInfo
		pkg              *ast.Package
		expectedContinue bool
	}{
		{
			name: "HandlerPackage",
			info: &compiler.PackageInfo{
				ImportPath: "github.com/octocat/service/internal/handler",
			},
			pkg: &ast.Package{
				Name: "handler",
			},
			expectedContinue: true,
		},
		{
			name: "HandlerSubpackage",
			info: &compiler.PackageInfo{
				ImportPath: "github.com/octocat/service/internal/handler/lookup",
			},
			pkg: &ast.Package{
				Name: "lookup",
			},
			expectedContinue: true,
		},
		{
			name: "ControllerPackage",
			info: &compiler.PackageInfo{
				ImportPath: "github.com/octocat/service/internal/controller",
			},
			pkg: &ast.Package{
				Name: "controller",
			},
			expectedContinue: true,
		},
		{
			name: "ControllerSubpackage",
			info: &compiler.PackageInfo{
				ImportPath: "github.com/octocat/service/internal/controller/lookup",
			},
			pkg: &ast.Package{
				Name: "lookup",
			},
			expectedContinue: true,
		},
		{
			name: "GatewayPackage",
			info: &compiler.PackageInfo{
				ImportPath: "github.com/octocat/service/internal/gateway",
			},
			pkg: &ast.Package{
				Name: "gateway",
			},
			expectedContinue: true,
		},
		{
			name: "GatewaySubpackage",
			info: &compiler.PackageInfo{
				ImportPath: "github.com/octocat/service/internal/gateway/lookup",
			},
			pkg: &ast.Package{
				Name: "lookup",
			},
			expectedContinue: true,
		},
		{
			name: "RepositoryPackage",
			info: &compiler.PackageInfo{
				ImportPath: "github.com/octocat/service/internal/repository",
			},
			pkg: &ast.Package{
				Name: "repository",
			},
			expectedContinue: true,
		},
		{
			name: "RepositorySubpackage",
			info: &compiler.PackageInfo{
				ImportPath: "github.com/octocat/service/internal/repository/lookup",
			},
			pkg: &ast.Package{
				Name: "lookup",
			},
			expectedContinue: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := &pkgDecorator{}

			cont := d.Package(tc.info, tc.pkg)

			assert.Equal(t, tc.expectedContinue, cont)
		})
	}
}

func TestPkgDecorator_FilePre(t *testing.T) {
	tests := []struct {
		name             string
		info             *compiler.FileInfo
		file             *ast.File
		expectedContinue bool
	}{
		{
			name:             "OK",
			info:             &compiler.FileInfo{},
			file:             &ast.File{},
			expectedContinue: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := &pkgDecorator{}

			cont := d.FilePre(tc.info, tc.file)

			assert.Equal(t, tc.expectedContinue, cont)
		})
	}
}

func TestPkgDecorator_FilePost(t *testing.T) {
	tests := []struct {
		name          string
		imports       []ast.Spec
		decls         []ast.Decl
		info          *compiler.FileInfo
		file          *ast.File
		expectedError string
	}{
		{
			name:          "NoDeclaration",
			imports:       nil,
			decls:         nil,
			info:          &compiler.FileInfo{},
			file:          &ast.File{},
			expectedError: "",
		},
		{
			name: "WriteFileFails",
			imports: []ast.Spec{
				&ast.ImportSpec{
					Path: &ast.BasicLit{Value: `"fmt"`},
				},
			},
			decls: []ast.Decl{
				&ast.GenDecl{
					Tok: token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{
								&ast.Ident{Name: "dummy"},
							},
							Type: &ast.Ident{Name: "string"},
						},
					},
				},
			},
			info: &compiler.FileInfo{
				PackageInfo: compiler.PackageInfo{
					ModuleName:  "github.com/octocat/service",
					PackageName: "controller",
					ImportPath:  "github.com/octocat/service/internal/controller",
					BaseDir:     "/dev/null",
					RelativeDir: "internal/controller",
				},
				FileName: "controller.go",
				FileSet:  token.NewFileSet(),
			},
			file:          &ast.File{},
			expectedError: "mkdir /dev/null: not a directory",
		},
		{
			name: "Success",
			imports: []ast.Spec{
				&ast.ImportSpec{
					Path: &ast.BasicLit{Value: `"fmt"`},
				},
			},
			decls: []ast.Decl{
				&ast.GenDecl{
					Tok: token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{
								&ast.Ident{Name: "dummy"},
							},
							Type: &ast.Ident{Name: "string"},
						},
					},
				},
			},
			info: &compiler.FileInfo{
				PackageInfo: compiler.PackageInfo{
					ModuleName:  "github.com/octocat/service",
					PackageName: "controller",
					ImportPath:  "github.com/octocat/service/internal/controller",
					BaseDir:     "./service",
					RelativeDir: "internal/controller",
				},
				FileName: "controller.go",
				FileSet:  token.NewFileSet(),
			},
			file:          &ast.File{},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := &pkgDecorator{
				imports: tc.imports,
				decls:   tc.decls,
			}

			err := d.FilePost(tc.info, tc.file)

			// Cleanup
			defer os.RemoveAll("./service")

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestPkgDecorator_Import(t *testing.T) {
	tests := []struct {
		name            string
		info            *compiler.FileInfo
		spec            *ast.ImportSpec
		expectedImports []ast.Spec
	}{
		{
			name: "OK",
			info: &compiler.FileInfo{},
			spec: &ast.ImportSpec{
				Path: &ast.BasicLit{Value: `"fmt"`},
			},
			expectedImports: []ast.Spec{
				&ast.ImportSpec{
					Path: &ast.BasicLit{Value: `"fmt"`},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := &pkgDecorator{}

			d.Import(tc.info, tc.spec)

			assert.Equal(t, tc.expectedImports, d.imports)
		})
	}
}

func TestPkgDecorator_Interface(t *testing.T) {
	tests := []struct {
		name                  string
		info                  *compiler.TypeInfo
		node                  *ast.InterfaceType
		expectedInterfaceName string
	}{
		{
			name: "Exported",
			info: &compiler.TypeInfo{
				TypeName: "Example",
			},
			node:                  &ast.InterfaceType{},
			expectedInterfaceName: "Example",
		},
		{
			name: "Unexported",
			info: &compiler.TypeInfo{
				TypeName: "example",
			},
			node:                  &ast.InterfaceType{},
			expectedInterfaceName: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := &pkgDecorator{}

			d.Interface(tc.info, tc.node)

			assert.Equal(t, tc.expectedInterfaceName, d.interfaceName)
		})
	}
}

func TestPkgDecorator_Struct(t *testing.T) {
	tests := []struct {
		name               string
		interfaceName      string
		info               *compiler.TypeInfo
		node               *ast.StructType
		expectedStructName string
		expectedDecls      []ast.Decl
	}{
		{
			name:          "Exported",
			interfaceName: "Example",
			info: &compiler.TypeInfo{
				FileInfo: compiler.FileInfo{
					PackageInfo: compiler.PackageInfo{
						PackageName: "controller",
					},
				},
				TypeName: "Example",
			},
			node:               &ast.StructType{},
			expectedStructName: "",
			expectedDecls:      nil,
		},
		{
			name:          "Unexported",
			interfaceName: "Example",
			info: &compiler.TypeInfo{
				FileInfo: compiler.FileInfo{
					PackageInfo: compiler.PackageInfo{
						PackageName: "controller",
					},
				},
				TypeName: "example",
			},
			node:               &ast.StructType{},
			expectedStructName: "example",
			expectedDecls: []ast.Decl{
				&ast.GenDecl{
					Tok: token.TYPE,
					Specs: []ast.Spec{
						&ast.TypeSpec{
							Name: &ast.Ident{Name: "example"},
							Type: &ast.StructType{
								Fields: &ast.FieldList{
									List: []*ast.Field{
										// Implementation
										{
											Names: []*ast.Ident{
												{Name: "impl"},
											},
											Type: &ast.SelectorExpr{
												X:   &ast.Ident{Name: "_controller"},
												Sel: &ast.Ident{Name: "Example"},
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

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := &pkgDecorator{
				interfaceName: tc.interfaceName,
			}

			d.Struct(tc.info, tc.node)

			assert.Equal(t, tc.expectedStructName, d.structName)
			assert.Equal(t, tc.expectedDecls, d.decls)
		})
	}
}

func TestPkgDecorator_FuncDecl(t *testing.T) {
	tests := []struct {
		name          string
		structName    string
		info          *compiler.FuncInfo
		funcType      *ast.FuncType
		body          *ast.BlockStmt
		expectedDecls []ast.Decl
	}{
		{
			name:       "NewFunction",
			structName: "example",
			info: &compiler.FuncInfo{
				FileInfo: compiler.FileInfo{
					PackageInfo: compiler.PackageInfo{
						PackageName: "controller",
					},
				},
				FuncName: "NewExample",
			},
			funcType: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								{Name: "gw"},
							},
							Type: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X:   &ast.Ident{Name: "gateway"},
									Sel: &ast.Ident{Name: "Gateway"},
								},
							},
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.Ident{Name: "Example"},
						},
						{
							Type: &ast.Ident{Name: "error"},
						},
					},
				},
			},
			body: &ast.BlockStmt{},
			expectedDecls: []ast.Decl{
				&ast.FuncDecl{
					Name: &ast.Ident{Name: "NewExample"},
					Type: &ast.FuncType{
						Params: &ast.FieldList{
							List: []*ast.Field{
								{
									Names: []*ast.Ident{
										{Name: "gw"},
									},
									Type: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X:   &ast.Ident{Name: "gateway"},
											Sel: &ast.Ident{Name: "Gateway"},
										},
									},
								},
							},
						},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.Ident{Name: "Example"},
								},
								{
									Type: &ast.Ident{Name: "error"},
								},
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{Name: "impl"},
									&ast.Ident{Name: "err"},
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   &ast.Ident{Name: "_controller"},
											Sel: &ast.Ident{Name: "NewExample"},
										},
										Args: []ast.Expr{
											&ast.Ident{Name: "gw"},
										},
									},
								},
							},
							&ast.IfStmt{
								Cond: &ast.BinaryExpr{
									X:  &ast.Ident{Name: "err"},
									Op: token.NEQ,
									Y:  &ast.Ident{Name: "nil"},
								},
								Body: &ast.BlockStmt{
									List: []ast.Stmt{
										&ast.ReturnStmt{
											Results: []ast.Expr{
												&ast.Ident{Name: "nil"},
												&ast.Ident{Name: "err"},
											},
										},
									},
								},
							},
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.UnaryExpr{
										Op: token.AND,
										X: &ast.CompositeLit{
											Type: &ast.Ident{Name: "example"},
											Elts: []ast.Expr{
												&ast.KeyValueExpr{
													Key:   &ast.Ident{Name: "impl"},
													Value: &ast.Ident{Name: "impl"},
												},
											},
										},
									},
									&ast.Ident{Name: "nil"},
								},
							},
						},
					},
				},
			},
		},
		{
			name:       "StructMethod_NoReturn",
			structName: "example",
			info: &compiler.FuncInfo{
				FileInfo: compiler.FileInfo{
					PackageInfo: compiler.PackageInfo{
						PackageName: "controller",
					},
				},
				FuncName: "Set",
				RecvName: "c",
				RecvType: &ast.Ident{Name: "example"},
			},
			funcType: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								{Name: "id"},
							},
							Type: &ast.Ident{Name: "int"},
						},
						{
							Names: []*ast.Ident{
								{Name: "name"},
							},
							Type: &ast.Ident{Name: "string"},
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{},
				},
			},
			body: &ast.BlockStmt{},
			expectedDecls: []ast.Decl{
				&ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "c"},
								},
								Type: &ast.StarExpr{
									X: &ast.Ident{Name: "example"},
								},
							},
						},
					},
					Name: &ast.Ident{Name: "Set"},
					Type: &ast.FuncType{
						Params: &ast.FieldList{
							List: []*ast.Field{
								{
									Names: []*ast.Ident{
										{Name: "id"},
									},
									Type: &ast.Ident{Name: "int"},
								},
								{
									Names: []*ast.Ident{
										{Name: "name"},
									},
									Type: &ast.Ident{Name: "string"},
								},
							},
						},
						Results: &ast.FieldList{
							List: []*ast.Field{},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X:   &ast.Ident{Name: "c"},
											Sel: &ast.Ident{Name: "impl"},
										},
										Sel: &ast.Ident{Name: "Set"},
									},
									Args: []ast.Expr{
										&ast.Ident{Name: "id"},
										&ast.Ident{Name: "name"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:       "StructMethod_WithReturn",
			structName: "example",
			info: &compiler.FuncInfo{
				FileInfo: compiler.FileInfo{
					PackageInfo: compiler.PackageInfo{
						PackageName: "controller",
					},
				},
				FuncName: "Get",
				RecvName: "c",
				RecvType: &ast.Ident{Name: "example"},
			},
			funcType: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								{Name: "id"},
							},
							Type: &ast.Ident{Name: "int"},
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.Ident{Name: "string"},
						},
					},
				},
			},
			body: &ast.BlockStmt{},
			expectedDecls: []ast.Decl{
				&ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "c"},
								},
								Type: &ast.StarExpr{
									X: &ast.Ident{Name: "example"},
								},
							},
						},
					},
					Name: &ast.Ident{Name: "Get"},
					Type: &ast.FuncType{
						Params: &ast.FieldList{
							List: []*ast.Field{
								{
									Names: []*ast.Ident{
										{Name: "id"},
									},
									Type: &ast.Ident{Name: "int"},
								},
							},
						},
						Results: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.Ident{Name: "string"},
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
											X: &ast.SelectorExpr{
												X:   &ast.Ident{Name: "c"},
												Sel: &ast.Ident{Name: "impl"},
											},
											Sel: &ast.Ident{Name: "Get"},
										},
										Args: []ast.Expr{
											&ast.Ident{Name: "id"},
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
			d := &pkgDecorator{
				structName: tc.structName,
			}

			d.FuncDecl(tc.info, tc.funcType, tc.body)

			assert.Equal(t, tc.expectedDecls, d.decls)
		})
	}
}
