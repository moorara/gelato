package compile

import "go/ast"

type (
	BuilderCreateDeclsMock struct {
		InPkgName  string
		InTypeName string
		InNode     *ast.StructType
		OutDecls   []ast.Decl
	}

	MockBuilder struct {
		CreateDeclsIndex int
		CreateDeclsMocks []BuilderCreateDeclsMock
	}
)

func (m *MockBuilder) CreateDecls(pkgName, typeName string, node *ast.StructType) []ast.Decl {
	i := m.CreateDeclsIndex
	m.CreateDeclsIndex++
	m.CreateDeclsMocks[i].InPkgName = pkgName
	m.CreateDeclsMocks[i].InTypeName = typeName
	m.CreateDeclsMocks[i].InNode = node
	return m.CreateDeclsMocks[i].OutDecls
}

type (
	MockerCreateDeclsMock struct {
		InPkgName  string
		InTypeName string
		InNode     *ast.InterfaceType
		OutDecls   []ast.Decl
	}

	MockMocker struct {
		CreateDeclsIndex int
		CreateDeclsMocks []MockerCreateDeclsMock
	}
)

func (m *MockMocker) CreateDecls(pkgName, typeName string, node *ast.InterfaceType) []ast.Decl {
	i := m.CreateDeclsIndex
	m.CreateDeclsIndex++
	m.CreateDeclsMocks[i].InPkgName = pkgName
	m.CreateDeclsMocks[i].InTypeName = typeName
	m.CreateDeclsMocks[i].InNode = node
	return m.CreateDeclsMocks[i].OutDecls
}
