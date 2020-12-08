package decorate

import "go/ast"

// structDecorator decorates a struct implementing an interface.
// It implements the ast.Visitor interface.
type structDecorator struct {
	layer string
}

func newStructDecorator(layer string) *structDecorator {
	return &structDecorator{
		layer: layer,
	}
}

func (d *structDecorator) Visit(n ast.Node) ast.Visitor {
	return d
}
