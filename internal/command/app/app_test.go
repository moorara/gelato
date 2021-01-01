package app

import (
	"bufio"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/archive"
	"github.com/moorara/gelato/internal/command"
	"github.com/moorara/gelato/internal/spec"
	"github.com/moorara/go-github"
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

func TestNewCommand(t *testing.T) {
	ui := cli.NewMockUi()
	spec := spec.Spec{}
	c, err := NewCommand(ui, spec)

	assert.NoError(t, err)
	assert.NotNil(t, c)
}

func TestCommand_Synopsis(t *testing.T) {
	c := &Command{}
	synopsis := c.Synopsis()

	assert.NotEmpty(t, synopsis)
}

func TestCommand_Help(t *testing.T) {
	c := &Command{}
	help := c.Help()

	assert.NotEmpty(t, help)
}

func TestCommand_Run(t *testing.T) {
	c := &Command{ui: cli.NewMockUi()}
	c.Run([]string{"--undefined"})

	assert.NotNil(t, c.services.repo)
	assert.NotNil(t, c.services.arch)
}

func TestCommand_run(t *testing.T) {
	tests := []struct {
		name             string
		repo             *MockRepoService
		arch             *MockArchiveService
		args             []string
		inputs           string
		expectedExitCode int
	}{
		{
			name:             "UndefinedFlag",
			repo:             &MockRepoService{},
			arch:             &MockArchiveService{},
			args:             []string{"--undefined"},
			expectedExitCode: command.FlagError,
		},
		{
			name:             "InvalidModuleName",
			repo:             &MockRepoService{},
			arch:             &MockArchiveService{},
			args:             []string{},
			inputs:           "",
			expectedExitCode: command.InputError,
		},
		{
			name:             "EmptyModuleName",
			repo:             &MockRepoService{},
			arch:             &MockArchiveService{},
			args:             []string{},
			inputs:           "\n",
			expectedExitCode: command.UnsupportedError,
		},
		{
			name:             "InvalidAppLang",
			repo:             &MockRepoService{},
			arch:             &MockArchiveService{},
			args:             []string{},
			inputs:           "github.com/octocat/service\n",
			expectedExitCode: command.InputError,
		},
		{
			name:             "EmptyAppLang",
			repo:             &MockRepoService{},
			arch:             &MockArchiveService{},
			args:             []string{},
			inputs:           "github.com/octocat/service\n\n",
			expectedExitCode: command.UnsupportedError,
		},
		{
			name:             "UnsupportedAppLang",
			repo:             &MockRepoService{},
			arch:             &MockArchiveService{},
			args:             []string{},
			inputs:           "github.com/octocat/service\njavascript\n",
			expectedExitCode: command.UnsupportedError,
		},
		{
			name:             "InvalidAppType",
			repo:             &MockRepoService{},
			arch:             &MockArchiveService{},
			args:             []string{},
			inputs:           "github.com/octocat/service\ngo\n",
			expectedExitCode: command.InputError,
		},
		{
			name:             "EmptyAppType",
			repo:             &MockRepoService{},
			arch:             &MockArchiveService{},
			args:             []string{},
			inputs:           "github.com/octocat/service\ngo\n\n",
			expectedExitCode: command.UnsupportedError,
		},
		{
			name:             "UnsupportedAppType",
			repo:             &MockRepoService{},
			arch:             &MockArchiveService{},
			args:             []string{},
			inputs:           "github.com/octocat/service\ngo\ncli\n",
			expectedExitCode: command.UnsupportedError,
		},
		{
			name:             "InvalidAppType",
			repo:             &MockRepoService{},
			arch:             &MockArchiveService{},
			args:             []string{},
			inputs:           "github.com/octocat/service\ngo\nhttp-service\n",
			expectedExitCode: command.InputError,
		},
		{
			name:             "EmptyAppType",
			repo:             &MockRepoService{},
			arch:             &MockArchiveService{},
			args:             []string{},
			inputs:           "github.com/octocat/service\ngo\nhttp-service\n\n",
			expectedExitCode: command.UnsupportedError,
		},
		{
			name:             "UnsupportedAppType",
			repo:             &MockRepoService{},
			arch:             &MockArchiveService{},
			args:             []string{},
			inputs:           "github.com/octocat/service\ngo\nhttp-service\ndiagonal\n",
			expectedExitCode: command.UnsupportedError,
		},
		{
			name: "DownloadTarArchiveFails",
			repo: &MockRepoService{
				DownloadTarArchiveMocks: []DownloadTarArchiveMock{
					{OutError: errors.New("error on downloading repository archive")},
				},
			},
			arch: &MockArchiveService{},
			args: []string{
				"-module=github.com/octocat/service",
				"-language=go",
				"-type=http-service",
				"-layout=vertical",
			},
			inputs:           "",
			expectedExitCode: command.GitHubError,
		},
		{
			name: "ArchiveWalkFails",
			repo: &MockRepoService{
				DownloadTarArchiveMocks: []DownloadTarArchiveMock{
					{OutResponse: &github.Response{}},
				},
			},
			arch: &MockArchiveService{
				ExtractMocks: []ExtractMock{
					{OutError: errors.New("error on extracting archive")},
				},
			},
			args: []string{
				"-module=github.com/octocat/service",
				"-language=go",
				"-type=http-service",
				"-layout=vertical",
			},
			inputs:           "",
			expectedExitCode: command.ExtractionError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// cli.Ui.Ask() method creates a new bufio.Reader every time.
			// Simply assigning an strings.Reader to mockUI.InputReader causes the bufio.Reader.ReadString() to error the second time the cli.Ui.Ask() method is called.
			// We need to assign a bufio.Reader to mockUI.InputReader, so bufio.NewReader() (called in cli.Ui.Ask()) will reuse it instead of creating a new one.
			var inputReader io.Reader
			inputReader = strings.NewReader(tc.inputs)
			inputReader = bufio.NewReader(inputReader)

			mockUI := cli.NewMockUi()
			mockUI.InputReader = inputReader
			c := &Command{ui: mockUI}
			c.services.repo = tc.repo
			c.services.arch = tc.arch

			exitCode := c.run(tc.args)

			assert.Equal(t, tc.expectedExitCode, exitCode)
		})
	}
}
