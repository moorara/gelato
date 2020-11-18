package build

import (
	"errors"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/command"
	"github.com/moorara/gelato/internal/spec"
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

func TestNewCommand(t *testing.T) {
	ui := new(cli.MockUi)
	spec := spec.Spec{}
	c, err := NewCommand(ui, spec)

	assert.NoError(t, err)
	assert.NotNil(t, c)

	cmd, ok := c.(*cmd)
	assert.True(t, ok)

	assert.NotNil(t, cmd.services.git)
}

func TestCmd_Synopsis(t *testing.T) {
	c := &cmd{}
	synopsis := c.Synopsis()

	assert.NotEmpty(t, synopsis)
}

func TestCmd_Help(t *testing.T) {
	c := &cmd{}
	help := c.Help()

	assert.NotEmpty(t, help)
}

func TestCmd_Run(t *testing.T) {
	tests := []struct {
		name             string
		git              *MockGitService
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
			name: "VersionPackageMissing",
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutHash: "3a1960ec0cec18d2dca14d270d11c5bc4138abf6", OutBranch: "main"},
				},
			},
			args:             []string{},
			expectedExitCode: command.VersionPkgError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &cmd{ui: new(cli.MockUi)}
			c.services.git = tc.git

			exitCode := c.Run(tc.args)

			assert.Equal(t, tc.expectedExitCode, exitCode)
		})
	}
}
