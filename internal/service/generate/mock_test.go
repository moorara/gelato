package generate

import "go/ast"

type (
	CompileMock struct {
		InPkgPath string
		InFile    *ast.File
		OutFile   *ast.File
	}

	MockCompiler struct {
		CompileIndex int
		CompileMocks []CompileMock
	}
)

func (m *MockCompiler) Compile(pkgPath string, file *ast.File) *ast.File {
	i := m.CompileIndex
	m.CompileIndex++
	m.CompileMocks[i].InPkgPath = pkgPath
	m.CompileMocks[i].InFile = file
	return m.CompileMocks[i].OutFile
}
