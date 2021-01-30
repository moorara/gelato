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
	f := &Factory{offset: 1}
	decl := f.ImportDecl("fmt", "github.com/octocat/example")

	expected := &ast.GenDecl{
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
	}

	assert.Equal(t, expected, decl)
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
