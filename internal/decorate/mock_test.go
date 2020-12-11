package decorate

import "go/ast"

type (
	ModifyMock struct {
		InModule string
		InDir    string
		InNode   ast.Node
		OutNode  ast.Node
	}

	MockModifier struct {
		ModifyIndex int
		ModifyMocks []ModifyMock
	}
)

func (m *MockModifier) Modify(module, dir string, node ast.Node) ast.Node {
	i := m.ModifyIndex
	m.ModifyIndex++
	m.ModifyMocks[i].InModule = module
	m.ModifyMocks[i].InDir = dir
	m.ModifyMocks[i].InNode = node
	return m.ModifyMocks[i].OutNode
}
