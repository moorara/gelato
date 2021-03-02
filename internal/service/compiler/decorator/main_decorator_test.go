package decorator

import (
	"go/ast"
	"go/token"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/service/compiler"
)

func TestMainDecorator_Package(t *testing.T) {
	tests := []struct {
		name             string
		info             *compiler.PackageInfo
		pkg              *ast.Package
		expectedContinue bool
	}{
		{
			name: "MainPackage",
			info: &compiler.PackageInfo{
				PackageName: "main",
				ImportPath:  "github.com/octocat/service",
			},
			pkg: &ast.Package{
				Name: "main",
			},
			expectedContinue: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := &mainDecorator{}

			cont := d.Package(tc.info, tc.pkg)

			assert.Equal(t, tc.expectedContinue, cont)
		})
	}
}

func TestMainDecorator_FilePre(t *testing.T) {
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
			d := &mainDecorator{}

			cont := d.FilePre(tc.info, tc.file)

			assert.Equal(t, tc.expectedContinue, cont)
		})
	}
}

func TestMainDecorator_FilePost(t *testing.T) {
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
			d := &mainDecorator{
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

func TestMainDecorator_Import(t *testing.T) {
	tests := []struct {
		name            string
		info            *compiler.FileInfo
		spec            *ast.ImportSpec
		expectedImports []ast.Spec
	}{
		{
			name: "NotDecoratablePackage",
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
		{
			name: "DecoratablePackage",
			info: &compiler.FileInfo{
				PackageInfo: compiler.PackageInfo{
					ModuleName: "github.com/octocat/service",
				},
			},
			spec: &ast.ImportSpec{
				Path: &ast.BasicLit{Value: `"github.com/octocat/service/internal/controller"`},
			},
			expectedImports: []ast.Spec{
				&ast.ImportSpec{
					Path: &ast.BasicLit{Value: `"github.com/octocat/service/internal/controller"`},
				},
				&ast.ImportSpec{
					Name: &ast.Ident{Name: "_controller"},
					Path: &ast.BasicLit{Value: `"github.com/octocat/service/.build/internal/controller"`},
				},
			},
		},
		{
			name: "DecoratableSubpackage",
			info: &compiler.FileInfo{
				PackageInfo: compiler.PackageInfo{
					ModuleName: "github.com/octocat/service",
				},
			},
			spec: &ast.ImportSpec{
				Name: &ast.Ident{Name: "lookupcontroller"},
				Path: &ast.BasicLit{Value: `"github.com/octocat/service/internal/controller/lookup"`},
			},
			expectedImports: []ast.Spec{
				&ast.ImportSpec{
					Name: &ast.Ident{Name: "lookupcontroller"},
					Path: &ast.BasicLit{Value: `"github.com/octocat/service/internal/controller/lookup"`},
				},
				&ast.ImportSpec{
					Name: &ast.Ident{Name: "_lookupcontroller"},
					Path: &ast.BasicLit{Value: `"github.com/octocat/service/.build/internal/controller/lookup"`},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := &mainDecorator{}

			d.Import(tc.info, tc.spec)

			assert.Equal(t, tc.expectedImports, d.imports)
		})
	}
}

func TestMainDecorator_FuncDecl(t *testing.T) {
	tests := []struct {
		name          string
		info          *compiler.FuncInfo
		funcType      *ast.FuncType
		body          *ast.BlockStmt
		expectedDecls []ast.Decl
	}{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := &mainDecorator{}

			d.FuncDecl(tc.info, tc.funcType, tc.body)

			assert.Equal(t, tc.expectedDecls, d.decls)
		})
	}
}
