package gen

import (
	"errors"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/command"
)

func TestNewCommand(t *testing.T) {
	ui := new(cli.MockUi)
	c, err := NewCommand(ui)

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

	assert.NotNil(t, c.services.builder)
	assert.NotNil(t, c.services.mocker)
}

func TestCommand_run(t *testing.T) {
	tests := []struct {
		name             string
		builder          *MockCompilerService
		mocker           *MockCompilerService
		args             []string
		expectedExitCode int
	}{
		{
			name:             "UndefinedFlag",
			args:             []string{"--undefined"},
			expectedExitCode: command.FlagError,
		},
		{
			name: "BuilderCompileFails",
			builder: &MockCompilerService{
				CompileMocks: []CompileMock{
					{OutError: errors.New("error on compiling")},
				},
			},
			args:             []string{},
			expectedExitCode: command.GenerationError,
		},
		{
			name: "MockerCompileFails",
			builder: &MockCompilerService{
				CompileMocks: []CompileMock{
					{OutError: nil},
				},
			},
			mocker: &MockCompilerService{
				CompileMocks: []CompileMock{
					{OutError: errors.New("error on compiling")},
				},
			},
			args:             []string{},
			expectedExitCode: command.GenerationError,
		},
		{
			name: "Success",
			builder: &MockCompilerService{
				CompileMocks: []CompileMock{
					{OutError: nil},
				},
			},
			mocker: &MockCompilerService{
				CompileMocks: []CompileMock{
					{OutError: nil},
				},
			},
			args:             []string{},
			expectedExitCode: command.Success,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &Command{ui: cli.NewMockUi()}
			c.services.builder = tc.builder
			c.services.mocker = tc.mocker

			exitCode := c.run(tc.args)

			assert.Equal(t, tc.expectedExitCode, exitCode)
		})
	}
}
