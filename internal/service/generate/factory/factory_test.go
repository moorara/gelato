package factory

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/log"
)

func TestGenerator(t *testing.T) {
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

	tests := []struct {
		name         string
		pkgPath      string
		pkgName      string
		file         *ast.File
		expectedFile *ast.File
	}{
		{
			name:    "OK",
			pkgPath: "github.com/octocat/app/example",
			pkgName: "example",
			file:    &ast.File{},
			expectedFile: &ast.File{
				Name: &ast.Ident{
					Name: "examplefactory",
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
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGenerator(clogger)

			file := g.Generate(tc.pkgPath, tc.pkgName, tc.file)

			assert.Equal(t, tc.expectedFile, file)
		})
	}
}
