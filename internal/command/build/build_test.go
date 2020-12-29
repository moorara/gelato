package build

import (
	"context"
	"errors"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/command"
	"github.com/moorara/gelato/internal/spec"
	"github.com/moorara/gelato/pkg/semver"
	"github.com/moorara/gelato/pkg/shell"
)

type (
	HEADMock struct {
		OutHash   string
		OutBranch string
		OutError  error
	}

	RemoteMock struct {
		InName    string
		OutDomain string
		OutPath   string
		OutError  error
	}

	MockGitService struct {
		HEADIndex int
		HEADMocks []HEADMock

		RemoteIndex int
		RemoteMocks []RemoteMock
	}
)

func (m *MockGitService) HEAD() (string, string, error) {
	i := m.HEADIndex
	m.HEADIndex++
	return m.HEADMocks[i].OutHash, m.HEADMocks[i].OutBranch, m.HEADMocks[i].OutError
}

func (m *MockGitService) Remote(name string) (string, string, error) {
	i := m.RemoteIndex
	m.RemoteIndex++
	m.RemoteMocks[i].InName = name
	return m.RemoteMocks[i].OutDomain, m.RemoteMocks[i].OutPath, m.RemoteMocks[i].OutError
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

func TestNewCommand(t *testing.T) {
	ui := new(cli.MockUi)
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
	c := &Command{ui: new(cli.MockUi)}
	c.Run([]string{"--undefined"})

	assert.NotNil(t, c.services.git)
	assert.NotNil(t, c.services.decorator)
	assert.NotNil(t, c.commands.semver)
}

func TestCommand_run(t *testing.T) {
	tests := []struct {
		name             string
		spec             spec.Spec
		git              *MockGitService
		decorator        *MockDecorateService
		goList           shell.RunnerFunc
		goBuild          shell.RunnerWithFunc
		semver           *MockSemverCommand
		args             []string
		expectedExitCode int
	}{
		{
			name: "UndefinedFlag",
			spec: spec.Spec{
				GelatoVersion: "v0.1.0",
				Build:         spec.Build{},
			},
			args:             []string{"--undefined"},
			expectedExitCode: command.FlagError,
		},
		{
			name: "GitHEADFails",
			spec: spec.Spec{
				GelatoVersion: "v0.1.0",
				Build:         spec.Build{},
			},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutError: errors.New("git error")},
				},
			},
			args:             []string{},
			expectedExitCode: command.GitError,
		},
		{
			name: "SemverRunFails",
			spec: spec.Spec{
				GelatoVersion: "v0.1.0",
				Build:         spec.Build{},
			},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutHash: "7813389d2b09cdf851665b7848daa212b27e4e82", OutBranch: "main"},
				},
			},
			goList: func(ctx context.Context, args ...string) (int, string, error) {
				return 1, "", errors.New("directory not found")
			},
			semver: &MockSemverCommand{
				RunMocks: []RunMock{
					{OutCode: command.GitError},
				},
			},
			args:             []string{},
			expectedExitCode: command.GitError,
		},
		{
			name: "Success_NoArtifact",
			spec: spec.Spec{
				GelatoVersion: "v0.1.0",
				Build:         spec.Build{},
			},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutHash: "7813389d2b09cdf851665b7848daa212b27e4e82", OutBranch: "main"},
				},
			},
			goList: func(ctx context.Context, args ...string) (int, string, error) {
				return 1, "github.com/octocat/Hello-World/version", nil
			},
			semver: &MockSemverCommand{
				RunMocks: []RunMock{
					{OutCode: command.Success},
				},
				SemVerMocks: []SemVerMock{
					{
						OutSemVer: semver.SemVer{},
					},
				},
			},
			args:             []string{},
			expectedExitCode: command.Success,
		},
		{
			name: "DecorateFails",
			spec: spec.Spec{
				GelatoVersion: "v0.1.0",
				Build: spec.Build{
					Decorate: true,
				},
			},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutHash: "7813389d2b09cdf851665b7848daa212b27e4e82", OutBranch: "main"},
				},
			},
			decorator: &MockDecorateService{
				DecorateMocks: []DecorateMock{
					{OutError: errors.New("error on decoration")},
				},
			},
			goList: func(ctx context.Context, args ...string) (int, string, error) {
				return 1, "github.com/octocat/Hello-World/version", nil
			},
			semver: &MockSemverCommand{
				RunMocks: []RunMock{
					{OutCode: command.Success},
				},
				SemVerMocks: []SemVerMock{
					{
						OutSemVer: semver.SemVer{},
					},
				},
			},
			args:             []string{},
			expectedExitCode: command.DecorationError,
		},
		{
			name: "Success_Decorated_NoArtifact",
			spec: spec.Spec{
				GelatoVersion: "v0.1.0",
				Build: spec.Build{
					Decorate: true,
				},
			},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutHash: "7813389d2b09cdf851665b7848daa212b27e4e82", OutBranch: "main"},
				},
			},
			decorator: &MockDecorateService{
				DecorateMocks: []DecorateMock{
					{OutError: nil},
				},
			},
			goList: func(ctx context.Context, args ...string) (int, string, error) {
				return 1, "github.com/octocat/Hello-World/version", nil
			},
			semver: &MockSemverCommand{
				RunMocks: []RunMock{
					{OutCode: command.Success},
				},
				SemVerMocks: []SemVerMock{
					{
						OutSemVer: semver.SemVer{},
					},
				},
			},
			args:             []string{},
			expectedExitCode: command.Success,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &Command{
				ui:   new(cli.MockUi),
				spec: tc.spec,
			}

			c.services.git = tc.git
			c.services.decorator = tc.decorator
			c.funcs.goList = tc.goList
			c.funcs.goBuild = tc.goBuild
			c.commands.semver = tc.semver

			exitCode := c.run(tc.args)

			assert.Equal(t, tc.expectedExitCode, exitCode)
		})
	}
}

func TestCommand_buildAll(t *testing.T) {
	tests := []struct {
		name          string
		buildSpec     spec.Build
		goBuild       shell.RunnerWithFunc
		ctx           context.Context
		ldFlags       string
		mainPkg       string
		output        string
		expectedError string
	}{
		{
			name: "WithoutCrossCompile_BuildFails",
			buildSpec: spec.Build{
				CrossCompile: false,
			},
			goBuild: func(ctx context.Context, opts shell.RunOptions, args ...string) (int, string, error) {
				return 1, "", errors.New("directory not found")
			},
			ctx:           context.Background(),
			ldFlags:       `-X "github.com/octocat/Hello-World/version.Version=v0.1.0"`,
			mainPkg:       "./cmd/app",
			output:        "./bin/app",
			expectedError: "directory not found",
		},
		{
			name: "WithoutCrossCompile_BuildSucceeds",
			buildSpec: spec.Build{
				CrossCompile: false,
			},
			goBuild: func(ctx context.Context, opts shell.RunOptions, args ...string) (int, string, error) {
				return 0, "", nil
			},
			ctx:           context.Background(),
			ldFlags:       `-X "github.com/octocat/Hello-World/version.Version=v0.1.0"`,
			mainPkg:       "./cmd/app",
			output:        "./bin/app",
			expectedError: "",
		},
		{
			name: "WithCrossCompile_BuildFails",
			buildSpec: spec.Build{
				CrossCompile: true,
				Platforms:    []string{"linux-amd64", "darwin-amd64"},
			},
			goBuild: func(ctx context.Context, opts shell.RunOptions, args ...string) (int, string, error) {
				return 1, "", errors.New("directory not found")
			},
			ctx:           context.Background(),
			ldFlags:       `-X "github.com/octocat/Hello-World/version.Version=v0.1.0"`,
			mainPkg:       "./cmd/app",
			output:        "./bin/app",
			expectedError: "directory not found",
		},
		{
			name: "WithCrossCompile_BuildSucceeds",
			buildSpec: spec.Build{
				CrossCompile: true,
				Platforms:    []string{"linux-amd64", "darwin-amd64"},
			},
			goBuild: func(ctx context.Context, opts shell.RunOptions, args ...string) (int, string, error) {
				return 0, "", nil
			},
			ctx:           context.Background(),
			ldFlags:       `-X "github.com/octocat/Hello-World/version.Version=v0.1.0"`,
			mainPkg:       "./cmd/app",
			output:        "./bin/app",
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &Command{
				ui: new(cli.MockUi),
				spec: spec.Spec{
					Build: tc.buildSpec,
				},
			}

			c.funcs.goBuild = tc.goBuild

			err := c.buildAll(tc.ctx, tc.ldFlags, tc.mainPkg, tc.output)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestCommand_Artifacts(t *testing.T) {
	artifacts := []Artifact{
		{"bin/app", "linux"},
	}

	c := &Command{}
	c.outputs.artifacts = artifacts

	assert.Equal(t, artifacts, c.Artifacts())
}
