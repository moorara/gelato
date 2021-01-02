package decorate

import "go/ast"

type (
	MainModifyMock struct {
		InModule string
		InDecDir string
		InNode   ast.Node
		OutNode  ast.Node
	}

	MockMainModifier struct {
		ModifyIndex int
		ModifyMocks []MainModifyMock
	}
)

func (m *MockMainModifier) Modify(module, decDir string, node ast.Node) ast.Node {
	i := m.ModifyIndex
	m.ModifyIndex++
	m.ModifyMocks[i].InModule = module
	m.ModifyMocks[i].InDecDir = decDir
	m.ModifyMocks[i].InNode = node
	return m.ModifyMocks[i].OutNode
}

type (
	GenericModifyMock struct {
		InModule  string
		InRelPath string
		InNode    ast.Node
		OutNode   ast.Node
	}

	MockGenericModifier struct {
		ModifyIndex int
		ModifyMocks []GenericModifyMock
	}
)

func (m *MockGenericModifier) Modify(module, relPath string, node ast.Node) ast.Node {
	i := m.ModifyIndex
	m.ModifyIndex++
	m.ModifyMocks[i].InModule = module
	m.ModifyMocks[i].InRelPath = relPath
	m.ModifyMocks[i].InNode = node
	return m.ModifyMocks[i].OutNode
}
