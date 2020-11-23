package build

import (
	"errors"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/command"
	"github.com/moorara/gelato/internal/spec"
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
	assert.NotNil(t, c.services.git)
	assert.NotNil(t, c.commands.semver)
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
	tests := []struct {
		name             string
		git              *MockGitService
		semver           *MockSemverCommand
		args             []string
		expectedExitCode int
	}{
		{
			name:             "UndefinedFlag",
			args:             []string{"--undefined"},
			expectedExitCode: command.FlagError,
		},
		{
			name: "GitHEADFails",
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
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutHash: "7813389d2b09cdf851665b7848daa212b27e4e82", OutBranch: "main"},
				},
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
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutHash: "7813389d2b09cdf851665b7848daa212b27e4e82", OutBranch: "main"},
				},
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
			c := &Command{ui: new(cli.MockUi)}
			c.services.git = tc.git
			c.commands.semver = tc.semver

			exitCode := c.Run(tc.args)

			assert.Equal(t, tc.expectedExitCode, exitCode)
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
