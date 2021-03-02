package build

import (
	"github.com/moorara/gelato/internal/service/compiler"
	"github.com/moorara/gelato/pkg/semver"
)

type (
	HEADMock struct {
		OutHash   string
		OutBranch string
		OutError  error
	}

	MockGitService struct {
		HEADIndex int
		HEADMocks []HEADMock
	}
)

func (m *MockGitService) HEAD() (string, string, error) {
	i := m.HEADIndex
	m.HEADIndex++
	return m.HEADMocks[i].OutHash, m.HEADMocks[i].OutBranch, m.HEADMocks[i].OutError
}

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

type (
	RunMock struct {
		InArgs  []string
		OutCode int
	}

	SemVerMock struct {
		OutSemVer semver.SemVer
	}

	MockSemverCommand struct {
		RunIndex int
		RunMocks []RunMock

		SemVerIndex int
		SemVerMocks []SemVerMock
	}
)

func (m *MockSemverCommand) Run(args []string) int {
	i := m.RunIndex
	m.RunIndex++
	return m.RunMocks[i].OutCode
}

func (m *MockSemverCommand) SemVer() semver.SemVer {
	i := m.SemVerIndex
	m.SemVerIndex++
	return m.SemVerMocks[i].OutSemVer
}
