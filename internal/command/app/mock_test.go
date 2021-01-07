package app

import (
	"context"
	"io"

	"github.com/moorara/go-github"

	"github.com/moorara/gelato/internal/service/archive"
	"github.com/moorara/gelato/internal/service/edit"
	"github.com/moorara/gelato/internal/service/git"
)

type (
	DownloadTarArchiveMock struct {
		InContext   context.Context
		InRef       string
		InWriter    io.Writer
		OutResponse *github.Response
		OutError    error
	}

	MockRepoService struct {
		DownloadTarArchiveIndex int
		DownloadTarArchiveMocks []DownloadTarArchiveMock
	}
)

func (m *MockRepoService) DownloadTarArchive(ctx context.Context, ref string, writer io.Writer) (*github.Response, error) {
	i := m.DownloadTarArchiveIndex
	m.DownloadTarArchiveIndex++
	m.DownloadTarArchiveMocks[i].InContext = ctx
	m.DownloadTarArchiveMocks[i].InRef = ref
	m.DownloadTarArchiveMocks[i].InWriter = writer
	return m.DownloadTarArchiveMocks[i].OutResponse, m.DownloadTarArchiveMocks[i].OutError
}

type (
	ExtractMock struct {
		InDest     string
		InReader   io.Reader
		InSelector archive.Selector
		OutError   error
	}

	MockArchiveService struct {
		ExtractIndex int
		ExtractMocks []ExtractMock
	}
)

func (m *MockArchiveService) Extract(dest string, reader io.Reader, selector archive.Selector) error {
	i := m.ExtractIndex
	m.ExtractIndex++
	m.ExtractMocks[i].InDest = dest
	m.ExtractMocks[i].InReader = reader
	m.ExtractMocks[i].InSelector = selector
	return m.ExtractMocks[i].OutError
}

type (
	ReplaceInDirMock struct {
		InRoot   string
		InSpecs  []edit.ReplaceSpec
		OutError error
	}

	MockEditService struct {
		ReplaceInDirIndex int
		ReplaceInDirMocks []ReplaceInDirMock
	}
)

func (m *MockEditService) ReplaceInDir(root string, specs []edit.ReplaceSpec) error {
	i := m.ReplaceInDirIndex
	m.ReplaceInDirIndex++
	m.ReplaceInDirMocks[i].InRoot = root
	m.ReplaceInDirMocks[i].InSpecs = specs
	return m.ReplaceInDirMocks[i].OutError
}

type (
	PathMock struct {
		OutPath  string
		OutError error
	}

	AddRemoteMock struct {
		InName   string
		InURL    string
		OutError error
	}

	SubmoduleMock struct {
		InName       string
		OutSubmodule git.Submodule
		OutError     error
	}

	CreateCommitMock struct {
		InMessage string
		InPaths   []string
		OutHash   string
		OutError  error
	}

	MockGitService struct {
		PathIndex int
		PathMocks []PathMock

		AddRemoteIndex int
		AddRemoteMocks []AddRemoteMock

		SubmoduleIndex int
		SubmoduleMocks []SubmoduleMock

		CreateCommitIndex int
		CreateCommitMocks []CreateCommitMock
	}
)

func (m *MockGitService) Path() (string, error) {
	i := m.PathIndex
	m.PathIndex++
	return m.PathMocks[i].OutPath, m.PathMocks[i].OutError
}

func (m *MockGitService) AddRemote(name string, url string) error {
	i := m.AddRemoteIndex
	m.AddRemoteIndex++
	m.AddRemoteMocks[i].InName = name
	m.AddRemoteMocks[i].InURL = url
	return m.AddRemoteMocks[i].OutError
}

func (m *MockGitService) Submodule(name string) (git.Submodule, error) {
	i := m.SubmoduleIndex
	m.SubmoduleIndex++
	m.SubmoduleMocks[i].InName = name
	return m.SubmoduleMocks[i].OutSubmodule, m.SubmoduleMocks[i].OutError
}

func (m *MockGitService) CreateCommit(message string, paths ...string) (string, error) {
	i := m.CreateCommitIndex
	m.CreateCommitIndex++
	m.CreateCommitMocks[i].InMessage = message
	m.CreateCommitMocks[i].InPaths = paths
	return m.CreateCommitMocks[i].OutHash, m.CreateCommitMocks[i].OutError
}
