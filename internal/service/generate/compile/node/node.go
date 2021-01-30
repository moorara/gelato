package node

import (
	"fmt"
	"go/ast"
	"go/token"
)

// Factory is used for generating ast.Node objects with positional information.
type Factory struct {
	offset int
}

// NewFactory creates a new factory for generating ast.Node objects with positional information.
func NewFactory() *Factory {
	return &Factory{
		offset: 1,
	}
}

// IncOffset increases the offset by a number.
func (f *Factory) IncOffset(inc int) {
	f.offset += inc
}

// PackagePos returns the positions for the package keyword.
func (f *Factory) PackagePos() token.Pos {
	pos := token.Pos(f.offset)
	f.offset += len("package") + 1 // space
	return pos
}

// Ident creates an ast.Ident node.
func (f *Factory) Ident(name string) *ast.Ident {
	pos := token.Pos(f.offset)
	f.offset += len(name) + 1 // space

	return &ast.Ident{
		NamePos: pos,
		Name:    name,
	}
}

// ImportDecl creates an ast.GenDecl node for an import.
func (f *Factory) ImportDecl(pkgs ...string) *ast.GenDecl {
	pos := token.Pos(f.offset)
	f.offset += len("import") + 1 // whitespaces

	decl := &ast.GenDecl{
		TokPos: pos,
		Tok:    token.IMPORT,
	}

	decl.Lparen = token.Pos(f.offset)
	f.offset += 3 // (, newline, tab

	for _, pkg := range pkgs {
		val := fmt.Sprintf("%q", pkg)
		pos := token.Pos(f.offset)
		f.offset += len(val) + 2 // newline, tab

		decl.Specs = append(decl.Specs, &ast.ImportSpec{
			Path: &ast.BasicLit{
				ValuePos: pos,
				Value:    val,
			},
		})
	}

	decl.Rparen = token.Pos(f.offset)
	f.offset += 2 // ), newline

	return decl
}

// Comment creates an ast.Comment node.
func (f *Factory) Comment(text string) *ast.Comment {
	pos := token.Pos(f.offset)
	f.offset += len(text) + 1 // new line

	return &ast.Comment{
		Slash: pos,
		Text:  text,
	}
}
