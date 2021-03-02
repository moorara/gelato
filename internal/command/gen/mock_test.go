package gen

import (
	"github.com/moorara/gelato/internal/service/compiler"
)

type (
	CompileMock struct {
		InPath    string
		InOptions compiler.ParseOptions
		OutError  error
	}

	MockCompilerService struct {
		CompileIndex int
		CompileMocks []CompileMock
	}
)

func (m *MockCompilerService) Compile(path string, opts compiler.ParseOptions) error {
	i := m.CompileIndex
	m.CompileIndex++
	m.CompileMocks[i].InPath = path
	m.CompileMocks[i].InOptions = opts
	return m.CompileMocks[i].OutError
}
