package build

import (
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
	DecorateMock struct {
		InPath   string
		OutError error
	}

	MockDecorateService struct {
		DecorateIndex int
		DecorateMocks []DecorateMock
	}
)

func (m *MockDecorateService) Decorate(path string) error {
	i := m.DecorateIndex
	m.DecorateIndex++
	m.DecorateMocks[i].InPath = path
	return m.DecorateMocks[i].OutError
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
