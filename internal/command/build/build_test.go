package build

import (
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/command"
	"github.com/moorara/gelato/internal/spec"
)

func TestNewCommand(t *testing.T) {
	ui := new(cli.MockUi)
	spec := spec.Spec{}
	c, err := NewCommand(ui, spec)

	assert.NoError(t, err)
	assert.NotNil(t, c)
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
		args             []string
		expectedExitCode int
	}{
		{
			name:             "UndefinedFlag",
			args:             []string{"--undefined"},
			expectedExitCode: command.FlagError,
		},
		{
			name:             "VersionPackageMissing",
			args:             []string{},
			expectedExitCode: command.VersionPkgError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &cmd{ui: new(cli.MockUi)}

			exitCode := c.Run(tc.args)

			assert.Equal(t, tc.expectedExitCode, exitCode)
		})
	}
}
