package update

import (
	"os"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/command"
)

func TestNewCommand(t *testing.T) {
	ui := new(cli.MockUi)
	c, err := NewCommand(ui)

	assert.NotNil(t, c)
	assert.NoError(t, err)
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
		environment      map[string]string
		args             []string
		expectedExitCode int
	}{
		{
			name:             "UndefinedFlag",
			environment:      map[string]string{},
			args:             []string{"--undefined"},
			expectedExitCode: command.FlagError,
		},
		{
			name:             "PreflightCheckFails",
			environment:      map[string]string{},
			args:             []string{},
			expectedExitCode: command.PreflightError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for key, val := range tc.environment {
				err := os.Setenv(key, val)
				assert.NoError(t, err)
				defer os.Unsetenv(key)
			}

			ui := new(cli.MockUi)
			c := &cmd{
				ui: ui,
			}

			exitCode := c.Run(tc.args)

			assert.Equal(t, tc.expectedExitCode, exitCode)
		})
	}
}
