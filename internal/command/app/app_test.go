package app

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/command"
	"github.com/moorara/gelato/internal/service/git"
	"github.com/moorara/gelato/internal/spec"
	"github.com/moorara/go-github"
)

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
	assert.NotNil(t, c.services.edit)
	assert.NotNil(t, c.funcs.gitInit)
	assert.NotNil(t, c.funcs.gitOpen)
}

func TestCommand_run(t *testing.T) {
	tests := []struct {
		name             string
		repo             *MockRepoService
		arch             *MockArchiveService
		edit             *MockEditService
		gitInit          gitFunc
		gitOpen          gitFunc
		args             []string
		inputs           string
		expectedExitCode int
	}{
		{
			name:             "UndefinedFlag",
			repo:             &MockRepoService{},
			arch:             &MockArchiveService{},
			edit:             &MockEditService{},
			gitInit:          nil,
			gitOpen:          nil,
			args:             []string{"--undefined"},
			expectedExitCode: command.FlagError,
		},
		{
			name:             "InvalidAppLang",
			repo:             &MockRepoService{},
			arch:             &MockArchiveService{},
			edit:             &MockEditService{},
			gitInit:          nil,
			gitOpen:          nil,
			args:             []string{},
			inputs:           "",
			expectedExitCode: command.InputError,
		},
		{
			name:             "EmptyAppLang",
			repo:             &MockRepoService{},
			arch:             &MockArchiveService{},
			edit:             &MockEditService{},
			gitInit:          nil,
			gitOpen:          nil,
			args:             []string{},
			inputs:           "\n",
			expectedExitCode: command.UnsupportedError,
		},
		{
			name:             "UnsupportedAppLang",
			repo:             &MockRepoService{},
			arch:             &MockArchiveService{},
			edit:             &MockEditService{},
			gitInit:          nil,
			gitOpen:          nil,
			args:             []string{},
			inputs:           "javascript\n",
			expectedExitCode: command.UnsupportedError,
		},
		{
			name:             "InvalidAppType",
			repo:             &MockRepoService{},
			arch:             &MockArchiveService{},
			edit:             &MockEditService{},
			gitInit:          nil,
			gitOpen:          nil,
			args:             []string{},
			inputs:           "go\n",
			expectedExitCode: command.InputError,
		},
		{
			name:             "EmptyAppType",
			repo:             &MockRepoService{},
			arch:             &MockArchiveService{},
			edit:             &MockEditService{},
			gitInit:          nil,
			gitOpen:          nil,
			args:             []string{},
			inputs:           "go\n\n",
			expectedExitCode: command.UnsupportedError,
		},
		{
			name:             "UnsupportedAppType",
			repo:             &MockRepoService{},
			arch:             &MockArchiveService{},
			edit:             &MockEditService{},
			gitInit:          nil,
			gitOpen:          nil,
			args:             []string{},
			inputs:           "go\ncli\n",
			expectedExitCode: command.UnsupportedError,
		},
		{
			name:             "InvalidAppType",
			repo:             &MockRepoService{},
			arch:             &MockArchiveService{},
			edit:             &MockEditService{},
			gitInit:          nil,
			gitOpen:          nil,
			args:             []string{},
			inputs:           "go\nhttp-service\n",
			expectedExitCode: command.InputError,
		},
		{
			name:             "EmptyAppType",
			repo:             &MockRepoService{},
			arch:             &MockArchiveService{},
			edit:             &MockEditService{},
			gitInit:          nil,
			gitOpen:          nil,
			args:             []string{},
			inputs:           "go\nhttp-service\n\n",
			expectedExitCode: command.UnsupportedError,
		},
		{
			name:             "UnsupportedAppType",
			repo:             &MockRepoService{},
			arch:             &MockArchiveService{},
			edit:             &MockEditService{},
			gitInit:          nil,
			gitOpen:          nil,
			args:             []string{},
			inputs:           "go\nhttp-service\ndiagonal\n",
			expectedExitCode: command.UnsupportedError,
		},
		{
			name:             "InvalidModuleName",
			repo:             &MockRepoService{},
			arch:             &MockArchiveService{},
			edit:             &MockEditService{},
			gitInit:          nil,
			gitOpen:          nil,
			args:             []string{},
			inputs:           "go\nhttp-service\nvertical\n",
			expectedExitCode: command.InputError,
		},
		{
			name:             "EmptyModuleName",
			repo:             &MockRepoService{},
			arch:             &MockArchiveService{},
			edit:             &MockEditService{},
			gitInit:          nil,
			gitOpen:          nil,
			args:             []string{},
			inputs:           "go\nhttp-service\nvertical\n\n",
			expectedExitCode: command.UnsupportedError,
		},
		{
			name:             "InvalidDockerID",
			repo:             &MockRepoService{},
			arch:             &MockArchiveService{},
			edit:             &MockEditService{},
			gitInit:          nil,
			gitOpen:          nil,
			args:             []string{},
			inputs:           "go\nhttp-service\nvertical\ngithub.com/octocat/service\n",
			expectedExitCode: command.InputError,
		},
		{
			name:             "EmptyDockerID",
			repo:             &MockRepoService{},
			arch:             &MockArchiveService{},
			edit:             &MockEditService{},
			gitInit:          nil,
			gitOpen:          nil,
			args:             []string{},
			inputs:           "go\nhttp-service\nvertical\ngithub.com/octocat/service\n\n",
			expectedExitCode: command.UnsupportedError,
		},
		{
			name: "Microrepo_DownloadTarArchiveFails",
			repo: &MockRepoService{
				DownloadTarArchiveMocks: []DownloadTarArchiveMock{
					{OutError: errors.New("error on downloading repository archive")},
				},
			},
			arch:    &MockArchiveService{},
			edit:    &MockEditService{},
			gitInit: nil,
			gitOpen: nil,
			args: []string{
				"-language=go",
				"-type=http-service",
				"-layout=vertical",
				"-module=github.com/octocat/service",
				"-docker=octocat",
			},
			inputs:           "",
			expectedExitCode: command.GitHubError,
		},
		{
			name: "Microrepo_ExtractFails",
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
			edit:    &MockEditService{},
			gitInit: nil,
			gitOpen: nil,
			args: []string{
				"-language=go",
				"-type=http-service",
				"-layout=vertical",
				"-module=github.com/octocat/service",
				"-docker=octocat",
			},
			inputs:           "",
			expectedExitCode: command.ExtractionError,
		},
		{
			name: "Microrepo_GitInitFails",
			repo: &MockRepoService{
				DownloadTarArchiveMocks: []DownloadTarArchiveMock{
					{OutResponse: &github.Response{}},
				},
			},
			arch: &MockArchiveService{
				ExtractMocks: []ExtractMock{
					{OutError: nil},
				},
			},
			edit: &MockEditService{},
			gitInit: func(string) (gitService, error) {
				return nil, errors.New("git error")
			},
			gitOpen: nil,
			args: []string{
				"-language=go",
				"-type=http-service",
				"-layout=vertical",
				"-module=github.com/octocat/service",
				"-docker=octocat",
			},
			inputs:           "",
			expectedExitCode: command.GitError,
		},
		{
			name: "Monorepo_GitOpenFails",
			repo: &MockRepoService{
				DownloadTarArchiveMocks: []DownloadTarArchiveMock{
					{OutResponse: &github.Response{}},
				},
			},
			arch: &MockArchiveService{
				ExtractMocks: []ExtractMock{
					{OutError: nil},
				},
			},
			edit:    &MockEditService{},
			gitInit: nil,
			gitOpen: func(string) (gitService, error) {
				return nil, errors.New("git error")
			},
			args: []string{
				"-language=go",
				"-type=http-service",
				"-layout=vertical",
				"-module=github.com/octocat/monorepo/services/domain/product/name",
				"-docker=octocat",
				"-monorepo",
			},
			inputs:           "",
			expectedExitCode: command.GitError,
		},
		{
			name: "Monorepo_GitSubmoduleFails",
			repo: &MockRepoService{
				DownloadTarArchiveMocks: []DownloadTarArchiveMock{
					{OutResponse: &github.Response{}},
				},
			},
			arch: &MockArchiveService{
				ExtractMocks: []ExtractMock{
					{OutError: nil},
				},
			},
			edit:    &MockEditService{},
			gitInit: nil,
			gitOpen: func(string) (gitService, error) {
				return &MockGitService{
					SubmoduleMocks: []SubmoduleMock{
						{OutError: errors.New("submodule not found")},
					},
				}, nil
			},
			args: []string{
				"-language=go",
				"-type=http-service",
				"-layout=vertical",
				"-module=github.com/octocat/monorepo/services/domain/product/name",
				"-docker=octocat",
				"-monorepo",
			},
			inputs:           "",
			expectedExitCode: command.GitError,
		},
		{
			name: "Monorepo_GitPathFails",
			repo: &MockRepoService{
				DownloadTarArchiveMocks: []DownloadTarArchiveMock{
					{OutResponse: &github.Response{}},
				},
			},
			arch: &MockArchiveService{
				ExtractMocks: []ExtractMock{
					{OutError: nil},
				},
			},
			edit:    &MockEditService{},
			gitInit: nil,
			gitOpen: func(string) (gitService, error) {
				return &MockGitService{
					SubmoduleMocks: []SubmoduleMock{
						{
							OutSubmodule: git.Submodule{
								Name:   "make",
								Path:   "services/common/make",
								URL:    "git@github.com:moorara/make.git",
								Branch: "main",
							},
						},
					},
					PathMocks: []PathMock{
						{OutError: errors.New("git error")},
					},
				}, nil
			},
			args: []string{
				"-language=go",
				"-type=http-service",
				"-layout=vertical",
				"-module=github.com/octocat/monorepo/services/domain/product/name",
				"-docker=octocat",
				"-monorepo",
			},
			inputs:           "",
			expectedExitCode: command.MiscError,
		},
		{
			name: "Microrepo_ReplaceFails",
			repo: &MockRepoService{
				DownloadTarArchiveMocks: []DownloadTarArchiveMock{
					{OutResponse: &github.Response{}},
				},
			},
			arch: &MockArchiveService{
				ExtractMocks: []ExtractMock{
					{OutError: nil},
				},
			},
			edit: &MockEditService{
				ReplaceInDirMocks: []ReplaceInDirMock{
					{OutError: errors.New("error on replacing")},
				},
			},
			gitInit: func(string) (gitService, error) {
				return &MockGitService{}, nil
			},
			gitOpen: nil,
			args: []string{
				"-language=go",
				"-type=http-service",
				"-layout=vertical",
				"-module=github.com/octocat/service",
				"-docker=octocat",
			},
			inputs:           "",
			expectedExitCode: command.OSError,
		},
		{
			name: "Monorepo_ReplaceFails",
			repo: &MockRepoService{
				DownloadTarArchiveMocks: []DownloadTarArchiveMock{
					{OutResponse: &github.Response{}},
				},
			},
			arch: &MockArchiveService{
				ExtractMocks: []ExtractMock{
					{OutError: nil},
				},
			},
			edit: &MockEditService{
				ReplaceInDirMocks: []ReplaceInDirMock{
					{OutError: errors.New("error on replacing")},
				},
			},
			gitInit: nil,
			gitOpen: func(string) (gitService, error) {
				return &MockGitService{
					SubmoduleMocks: []SubmoduleMock{
						{
							OutSubmodule: git.Submodule{
								Name:   "make",
								Path:   "services/common/make",
								URL:    "git@github.com:moorara/make.git",
								Branch: "main",
							},
						},
					},
					PathMocks: []PathMock{
						{OutPath: "/home/user/code/github.com/octocat/monorepo"},
					},
				}, nil
			},
			args: []string{
				"-language=go",
				"-type=http-service",
				"-layout=vertical",
				"-module=github.com/octocat/monorepo/services/domain/product/name",
				"-docker=octocat",
				"-monorepo",
			},
			inputs:           "",
			expectedExitCode: command.OSError,
		},
		{
			name: "Microrepo_GitCommitFails",
			repo: &MockRepoService{
				DownloadTarArchiveMocks: []DownloadTarArchiveMock{
					{OutResponse: &github.Response{}},
				},
			},
			arch: &MockArchiveService{
				ExtractMocks: []ExtractMock{
					{OutError: nil},
				},
			},
			edit: &MockEditService{
				ReplaceInDirMocks: []ReplaceInDirMock{
					{OutError: nil},
				},
			},
			gitInit: func(string) (gitService, error) {
				return &MockGitService{
					CreateCommitMocks: []CreateCommitMock{
						{OutError: errors.New("git error")},
					},
				}, nil
			},
			gitOpen: nil,
			args: []string{
				"-language=go",
				"-type=http-service",
				"-layout=vertical",
				"-module=github.com/octocat/service",
				"-docker=octocat",
			},
			inputs:           "",
			expectedExitCode: command.GitError,
		},
		{
			name: "Microrepo_GitMoveBranchFails",
			repo: &MockRepoService{
				DownloadTarArchiveMocks: []DownloadTarArchiveMock{
					{OutResponse: &github.Response{}},
				},
			},
			arch: &MockArchiveService{
				ExtractMocks: []ExtractMock{
					{OutError: nil},
				},
			},
			edit: &MockEditService{
				ReplaceInDirMocks: []ReplaceInDirMock{
					{OutError: nil},
				},
			},
			gitInit: func(string) (gitService, error) {
				return &MockGitService{
					CreateCommitMocks: []CreateCommitMock{
						{OutHash: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
					},
					MoveBranchMocks: []MoveBranchMock{
						{OutError: errors.New("git error")},
					},
				}, nil
			},
			gitOpen: nil,
			args: []string{
				"-language=go",
				"-type=http-service",
				"-layout=vertical",
				"-module=github.com/octocat/service",
				"-docker=octocat",
			},
			inputs:           "",
			expectedExitCode: command.GitError,
		},
		{
			name: "Microrepo_GitAddRemoteFails",
			repo: &MockRepoService{
				DownloadTarArchiveMocks: []DownloadTarArchiveMock{
					{OutResponse: &github.Response{}},
				},
			},
			arch: &MockArchiveService{
				ExtractMocks: []ExtractMock{
					{OutError: nil},
				},
			},
			edit: &MockEditService{
				ReplaceInDirMocks: []ReplaceInDirMock{
					{OutError: nil},
				},
			},
			gitInit: func(string) (gitService, error) {
				return &MockGitService{
					CreateCommitMocks: []CreateCommitMock{
						{OutHash: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
					},
					MoveBranchMocks: []MoveBranchMock{
						{OutError: nil},
					},
					AddRemoteMocks: []AddRemoteMock{
						{OutError: errors.New("git error")},
					},
				}, nil
			},
			gitOpen: nil,
			args: []string{
				"-language=go",
				"-type=http-service",
				"-layout=vertical",
				"-module=github.com/octocat/service",
				"-docker=octocat",
			},
			inputs:           "",
			expectedExitCode: command.GitError,
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
			c.services.edit = tc.edit
			c.funcs.gitInit = tc.gitInit
			c.funcs.gitOpen = tc.gitOpen

			exitCode := c.run(tc.args)

			assert.Equal(t, tc.expectedExitCode, exitCode)
		})
	}
}
