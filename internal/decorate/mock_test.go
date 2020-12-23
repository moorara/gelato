package decorate

import "go/ast"

type (
	ModifyMock struct {
		InModule  string
		InDecDir  string
		InRelPath string
		InNode    ast.Node
		OutNode   ast.Node
	}

	MockModifier struct {
		ModifyIndex int
		ModifyMocks []ModifyMock
	}
)

func (m *MockModifier) Modify(module, decDir, relPath string, node ast.Node) ast.Node {
	i := m.ModifyIndex
	m.ModifyIndex++
	m.ModifyMocks[i].InModule = module
	m.ModifyMocks[i].InDecDir = decDir
	m.ModifyMocks[i].InRelPath = relPath
	m.ModifyMocks[i].InNode = node
	return m.ModifyMocks[i].OutNode
}
