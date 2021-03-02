package decorator

import (
	"github.com/moorara/gelato/internal/service/compiler"
	"github.com/moorara/gelato/internal/log"
)

const decoratedDir = ".build"

// New creates a new compiler for generating decorated applications.
func New(level log.Level) *compiler.Compiler {
	logger := log.NewColorful(level)

	md := &mainDecorator{}
	mainConsumer := &compiler.Consumer{
		Name:     "mainDecorator",
		Package:  md.Package,
		FilePre:  md.FilePre,
		FilePost: md.FilePost,
		Import:   md.Import,
		FuncDecl: md.FuncDecl,
	}

	pd := &pkgDecorator{}
	pkgConsumer := &compiler.Consumer{
		Name:      "packageDecorator",
		Package:   pd.Package,
		FilePre:   pd.FilePre,
		FilePost:  pd.FilePost,
		Import:    pd.Import,
		Struct:    pd.Struct,
		Interface: pd.Interface,
		FuncDecl:  pd.FuncDecl,
	}

	return compiler.New(logger, mainConsumer, pkgConsumer)
}
