package decorate

import (
	"go/ast"

	"golang.org/x/tools/go/ast/astutil"
)

type (
	VisitMock struct {
		InNode     ast.Node
		OutVisitor ast.Visitor
	}

	MockVisitor struct {
		VisitMock
	}
)

func (m *MockVisitor) Visit(node ast.Node) ast.Visitor {
	m.VisitMock.InNode = node
	return m.VisitMock.OutVisitor
}

type (
	PreMock struct {
		InCursor *astutil.Cursor
		OutBool  bool
	}

	PostMock struct {
		InCursor *astutil.Cursor
		OutBool  bool
	}

	MockModifier struct {
		PreMock
		PostMock
	}
)

func (m *MockModifier) Pre(cursor *astutil.Cursor) bool {
	m.PreMock.InCursor = cursor
	return m.PreMock.OutBool
}

func (m *MockModifier) Post(cursor *astutil.Cursor) bool {
	m.PostMock.InCursor = cursor
	return m.PostMock.OutBool
}
