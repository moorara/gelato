package app

import (
	"bufio"
	"io"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/command"
)

func TestNewCommand(t *testing.T) {
	ui := cli.NewMockUi()
	c, err := NewCommand(ui, "0.1.0")

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

	// assert.NotNil(t, c.services)
}

func TestCommand_run(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		inputs           string
		expectedExitCode int
	}{
		{
			name:             "UndefinedFlag",
			args:             []string{"--undefined"},
			expectedExitCode: command.FlagError,
		},
		{
			name:             "InvalidAppLang",
			args:             []string{},
			inputs:           "",
			expectedExitCode: command.InputError,
		},
		{
			name:             "EmptyAppLang",
			args:             []string{},
			inputs:           "\n",
			expectedExitCode: command.UnsupportedError,
		},
		{
			name:             "UnsupportedAppLang",
			args:             []string{},
			inputs:           "javascript\n",
			expectedExitCode: command.UnsupportedError,
		},
		{
			name:             "InvalidAppType",
			args:             []string{},
			inputs:           "go\n",
			expectedExitCode: command.InputError,
		},
		{
			name:             "EmptyAppType",
			args:             []string{},
			inputs:           "go\n\n",
			expectedExitCode: command.UnsupportedError,
		},
		{
			name:             "UnsupportedAppType",
			args:             []string{},
			inputs:           "go\ncli\n",
			expectedExitCode: command.UnsupportedError,
		},
		{
			name:             "InvalidAppType",
			args:             []string{},
			inputs:           "go\nhttp-service\n",
			expectedExitCode: command.InputError,
		},
		{
			name:             "EmptyAppType",
			args:             []string{},
			inputs:           "go\nhttp-service\n\n",
			expectedExitCode: command.UnsupportedError,
		},
		{
			name:             "UnsupportedAppType",
			args:             []string{},
			inputs:           "go\nhttp-service\ndiagonal\n",
			expectedExitCode: command.UnsupportedError,
		},
		{
			name:             "InvalidModuleName",
			args:             []string{},
			inputs:           "go\nhttp-service\nvertical\n",
			expectedExitCode: command.InputError,
		},
		{
			name:             "EmptyModuleName",
			args:             []string{},
			inputs:           "go\nhttp-service\nvertical\n\n",
			expectedExitCode: command.UnsupportedError,
		},
		{
			name:             "OK",
			args:             []string{},
			inputs:           "go\nhttp-service\nvertical\ngithub.com/octocat/service\n",
			expectedExitCode: command.Success,
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

			exitCode := c.run(tc.args)

			assert.Equal(t, tc.expectedExitCode, exitCode)
		})
	}
}
