package generate

import (
	"go/ast"
)

type (
	GenerateMock struct {
		InPkgPath string
		InPkgName string
		InFile    *ast.File
		OutFile   *ast.File
	}

	MockGenerator struct {
		GenerateIndex int
		GenerateMocks []GenerateMock
	}
)

func (m *MockGenerator) Generate(pkgPath, pkgName string, file *ast.File) *ast.File {
	i := m.GenerateIndex
	m.GenerateIndex++
	m.GenerateMocks[i].InPkgPath = pkgPath
	m.GenerateMocks[i].InPkgName = pkgName
	m.GenerateMocks[i].InFile = file
	return m.GenerateMocks[i].OutFile
}
