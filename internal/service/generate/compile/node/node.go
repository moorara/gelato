package node

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
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

// Comment creates an ast.Comment node.
func (f *Factory) Comment(text string) *ast.Comment {
	pos := token.Pos(f.offset)
	f.offset += len(text) + 1 // new line

	return &ast.Comment{
		Slash: pos,
		Text:  text,
	}
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

// AnnotateFuncDecl adds positional information to an ast.FuncDecl node.
func (f *Factory) AnnotateFuncDecl(node *ast.FuncDecl) {
	f.offset++ // newline

	astutil.Apply(node,
		// Pre-order traversal
		func(c *astutil.Cursor) bool {
			switch n := c.Node().(type) {
			case *ast.FuncDecl:
				return true

			case *ast.FuncType:
				n.Func = token.Pos(f.offset)
				f.offset += len("func") + 1 // space

				funcName := c.Parent().(*ast.FuncDecl).Name
				funcName.NamePos = token.Pos(f.offset)
				f.offset += len(funcName.Name)

				return true

			case *ast.FieldList:
				f.AnnotateFieldList(n)

			case *ast.BlockStmt:
				f.AnnotateBlockStmt(n)
			}

			return false
		},
		// Post-order traversal
		func(c *astutil.Cursor) bool {
			return true
		},
	)
}

// AnnotateFieldList adds positional information to an ast.FieldList node.
func (f *Factory) AnnotateFieldList(node *ast.FieldList) {
	astutil.Apply(node,
		// Pre-order traversal
		func(c *astutil.Cursor) bool {
			switch n := c.Node().(type) {
			case *ast.FieldList:
				n.Opening = token.Pos(f.offset)
				f.offset++ // (
				return true

			case *ast.Field:
				return true

			case *ast.StarExpr:
				return true

			case *ast.SelectorExpr:
				return true

			case *ast.Ident:
				n.NamePos = token.Pos(f.offset)
				f.offset += len(n.Name)

				switch c.Name() {
				case "Names":
					f.offset++ // space or comma
				case "Type":
				case "X":
					f.offset++ // dot
				case "Sel":
				}

				return true
			}

			return false
		},
		// Post-order traversal
		func(c *astutil.Cursor) bool {
			switch n := c.Node().(type) {
			case *ast.FieldList:
				n.Closing = token.Pos(f.offset)
				f.offset += 2 // ), space

			case *ast.Field:
				f.offset++ // comma
			}

			return true
		},
	)
}

// AnnotateBlockStmt adds positional information to an ast.BlockStmt node.
func (f *Factory) AnnotateBlockStmt(node *ast.BlockStmt) {
	astutil.Apply(node,
		// Pre-order traversal
		func(c *astutil.Cursor) bool {
			switch n := c.Node().(type) {
			case *ast.BlockStmt:
				n.Lbrace = token.Pos(f.offset)
				f.offset += 3 // {, newline, tab
				return true

			case *ast.ReturnStmt:
				n.Return = token.Pos(f.offset)
				f.offset += len("return") + 1 // space
				return true

			case *ast.CallExpr:
				return true

			case *ast.SelectorExpr:
				return true

			case *ast.Ident:
				n.NamePos = token.Pos(f.offset)
				f.offset += len(n.Name)
				return true
			}

			return false
		},
		// Post-order traversal
		func(c *astutil.Cursor) bool {
			switch n := c.Node().(type) {
			case *ast.BlockStmt:
				n.Rbrace = token.Pos(f.offset)
				f.offset += 3 // newline, }, newline

			case *ast.CallExpr:
				n.Rparen = token.Pos(f.offset)
				f.offset++ // )
			}

			switch c.Name() {
			case "X":
				if _, ok := c.Parent().(*ast.SelectorExpr); ok {
					f.offset++ // dot
				}

			case "Fun":
				if n, ok := c.Parent().(*ast.CallExpr); ok {
					n.Lparen = token.Pos(f.offset)
					f.offset++ // (
				}

			case "Results":
				if _, ok := c.Parent().(*ast.ReturnStmt); ok {
					f.offset++ // comma or space
				}
			}

			return true
		},
	)
}
