package decorate

import (
	"errors"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/command"
	"github.com/moorara/gelato/internal/log"
	"github.com/moorara/gelato/internal/spec"
)

type (
	DecorateMock struct {
		InLevel  log.Level
		InPath   string
		OutError error
	}

	MockDecoratorService struct {
		DecorateIndex int
		DecorateMocks []DecorateMock
	}
)

func (m *MockDecoratorService) Decorate(level log.Level, path string) error {
	i := m.DecorateIndex
	m.DecorateIndex++
	m.DecorateMocks[i].InLevel = level
	m.DecorateMocks[i].InPath = path
	return m.DecorateMocks[i].OutError
}

func TestNewCommand(t *testing.T) {
	ui := new(cli.MockUi)
	spec := spec.App{}
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

	assert.NotNil(t, c.services.decorator)
}

func TestCommand_run(t *testing.T) {
	tests := []struct {
		name             string
		spec             spec.App
		decorator        *MockDecoratorService
		args             []string
		expectedExitCode int
	}{
		{
			name:             "UndefinedFlag",
			spec:             spec.App{},
			args:             []string{"--undefined"},
			expectedExitCode: command.FlagError,
		},
		{
			name: "DecorateFails",
			spec: spec.App{},
			decorator: &MockDecoratorService{
				DecorateMocks: []DecorateMock{
					{OutError: errors.New("decoration error")},
				},
			},
			args:             []string{},
			expectedExitCode: command.DecorationError,
		},
		{
			name: "Success_Trace",
			spec: spec.App{},
			decorator: &MockDecoratorService{
				DecorateMocks: []DecorateMock{
					{OutError: nil},
				},
			},
			args:             []string{"-trace"},
			expectedExitCode: command.Success,
		},
		{
			name: "Success_Debug",
			spec: spec.App{},
			decorator: &MockDecoratorService{
				DecorateMocks: []DecorateMock{
					{OutError: nil},
				},
			},
			args:             []string{"-debug"},
			expectedExitCode: command.Success,
		},
		{
			name: "Success_Info",
			spec: spec.App{},
			decorator: &MockDecoratorService{
				DecorateMocks: []DecorateMock{
					{OutError: nil},
				},
			},
			args:             []string{"-info"},
			expectedExitCode: command.Success,
		},
		{
			name: "Success_Warn",
			spec: spec.App{},
			decorator: &MockDecoratorService{
				DecorateMocks: []DecorateMock{
					{OutError: nil},
				},
			},
			args:             []string{"-warn"},
			expectedExitCode: command.Success,
		},
		{
			name: "Success_Error",
			spec: spec.App{},
			decorator: &MockDecoratorService{
				DecorateMocks: []DecorateMock{
					{OutError: nil},
				},
			},
			args:             []string{"-error"},
			expectedExitCode: command.Success,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &Command{
				ui:   new(cli.MockUi),
				spec: tc.spec,
			}

			c.services.decorator = tc.decorator

			exitCode := c.run(tc.args)

			assert.Equal(t, tc.expectedExitCode, exitCode)
		})
	}
}
