package shell

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name             string
		ctx              context.Context
		command          string
		args             []string
		expectedExitCode int
		expectedOutput   string
		expectedError    error
	}{
		{
			name:             "NotFound",
			ctx:              context.Background(),
			command:          "unknown",
			args:             []string{},
			expectedExitCode: -1,
			expectedOutput:   "",
			expectedError:    errors.New("error on running unknown: exec: \"unknown\": executable file not found in $PATH: "),
		},
		{
			name:             "Error",
			ctx:              context.Background(),
			command:          "cat",
			args:             []string{"null"},
			expectedExitCode: 1,
			expectedOutput:   "",
			expectedError:    errors.New("error on running cat null: exit status 1: cat: null: No such file or directory"),
		},
		{
			name:             "Success",
			ctx:              context.Background(),
			command:          "echo",
			args:             []string{"foo", "bar"},
			expectedExitCode: 0,
			expectedOutput:   "foo bar",
			expectedError:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			code, out, err := Run(tc.ctx, tc.command, tc.args...)

			assert.Equal(t, tc.expectedExitCode, code)
			assert.Equal(t, tc.expectedOutput, out)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestRunWith(t *testing.T) {
	tests := []struct {
		name             string
		ctx              context.Context
		opts             RunOptions
		command          string
		args             []string
		expectedExitCode int
		expectedOutput   string
		expectedError    error
	}{
		{
			name:             "NotFound",
			ctx:              context.Background(),
			opts:             RunOptions{},
			command:          "unknown",
			args:             []string{},
			expectedExitCode: -1,
			expectedOutput:   "",
			expectedError:    errors.New("error on running unknown: exec: \"unknown\": executable file not found in $PATH: "),
		},
		{
			name:             "Error",
			ctx:              context.Background(),
			opts:             RunOptions{},
			command:          "cat",
			args:             []string{"null"},
			expectedExitCode: 1,
			expectedOutput:   "",
			expectedError:    errors.New("error on running cat null: exit status 1: cat: null: No such file or directory"),
		},
		{
			name: "Success",
			ctx:  context.Background(),
			opts: RunOptions{
				Environment: map[string]string{
					"PLACEHOLDER": "foo bar",
				},
			},
			command:          "printenv",
			args:             []string{"PLACEHOLDER"},
			expectedExitCode: 0,
			expectedOutput:   "foo bar",
			expectedError:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			code, out, err := RunWith(tc.ctx, tc.opts, tc.command, tc.args...)

			assert.Equal(t, tc.expectedExitCode, code)
			assert.Equal(t, tc.expectedOutput, out)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestRunner(t *testing.T) {
	tests := []struct {
		name             string
		ctx              context.Context
		command          string
		args             []string
		runArgs          []string
		expectedExitCode int
		expectedOutput   string
		expectedError    error
	}{
		{
			name:             "NotFound",
			ctx:              context.Background(),
			command:          "unknown",
			args:             []string{},
			runArgs:          []string{},
			expectedExitCode: -1,
			expectedOutput:   "",
			expectedError:    errors.New("error on running unknown: exec: \"unknown\": executable file not found in $PATH: "),
		},
		{
			name:             "Error",
			ctx:              context.Background(),
			command:          "cat",
			args:             []string{"null"},
			runArgs:          []string{},
			expectedExitCode: 1,
			expectedOutput:   "",
			expectedError:    errors.New("error on running cat null: exit status 1: cat: null: No such file or directory"),
		},
		{
			name:             "Success",
			ctx:              context.Background(),
			command:          "echo",
			args:             []string{"foo", "bar"},
			runArgs:          []string{"baz"},
			expectedExitCode: 0,
			expectedOutput:   "foo bar baz",
			expectedError:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			run := Runner(tc.ctx, tc.command, tc.args...)
			code, out, err := run(tc.runArgs...)

			assert.Equal(t, tc.expectedExitCode, code)
			assert.Equal(t, tc.expectedOutput, out)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestRunnerWith(t *testing.T) {
	tests := []struct {
		name             string
		ctx              context.Context
		opts             RunOptions
		command          string
		args             []string
		runArgs          []string
		expectedExitCode int
		expectedOutput   string
		expectedError    error
	}{
		{
			name:             "NotFound",
			ctx:              context.Background(),
			opts:             RunOptions{},
			command:          "unknown",
			args:             []string{},
			runArgs:          []string{},
			expectedExitCode: -1,
			expectedOutput:   "",
			expectedError:    errors.New("error on running unknown: exec: \"unknown\": executable file not found in $PATH: "),
		},
		{
			name:             "Error",
			ctx:              context.Background(),
			opts:             RunOptions{},
			command:          "cat",
			args:             []string{"null"},
			runArgs:          []string{},
			expectedExitCode: 1,
			expectedOutput:   "",
			expectedError:    errors.New("error on running cat null: exit status 1: cat: null: No such file or directory"),
		},
		{
			name: "Success",
			ctx:  context.Background(),
			opts: RunOptions{
				Environment: map[string]string{
					"PLACEHOLDER": "foo bar baz",
				},
			},
			command:          "printenv",
			args:             []string{},
			runArgs:          []string{"PLACEHOLDER"},
			expectedExitCode: 0,
			expectedOutput:   "foo bar baz",
			expectedError:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			run := RunnerWith(tc.ctx, tc.opts, tc.command, tc.args...)
			code, out, err := run(tc.runArgs...)

			assert.Equal(t, tc.expectedExitCode, code)
			assert.Equal(t, tc.expectedOutput, out)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestRunnerFunc_WithArgs(t *testing.T) {
	tests := []struct {
		name             string
		runnerFunc       RunnerFunc
		args             []string
		runArgs          []string
		expectedExitCode int
		expectedOutput   string
		expectedError    error
	}{
		{
			name:             "NotFound",
			runnerFunc:       Runner(context.Background(), "unknown"),
			args:             []string{},
			runArgs:          []string{},
			expectedExitCode: -1,
			expectedOutput:   "",
			expectedError:    errors.New("error on running unknown: exec: \"unknown\": executable file not found in $PATH: "),
		},
		{
			name:             "Error",
			runnerFunc:       Runner(context.Background(), "cat"),
			args:             []string{"null"},
			runArgs:          []string{},
			expectedExitCode: 1,
			expectedOutput:   "",
			expectedError:    errors.New("error on running cat null: exit status 1: cat: null: No such file or directory"),
		},
		{
			name:             "Success",
			runnerFunc:       Runner(context.Background(), "echo"),
			args:             []string{"foo", "bar"},
			runArgs:          []string{"baz"},
			expectedExitCode: 0,
			expectedOutput:   "foo bar baz",
			expectedError:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			run := tc.runnerFunc.WithArgs(tc.args...)
			code, out, err := run(tc.runArgs...)

			assert.Equal(t, tc.expectedExitCode, code)
			assert.Equal(t, tc.expectedOutput, out)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
