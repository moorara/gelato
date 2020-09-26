package command

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/moorara/gelato/pkg/semver"
	"github.com/stretchr/testify/assert"
)

func TestRunPreflightChecks(t *testing.T) {
	tests := []struct {
		name                   string
		environment            map[string]string
		ctx                    context.Context
		checklist              PreflightChecklist
		expectedError          error
		expectWorkingDirectory bool
		expectGoVersion        bool
		expectGitVersion       bool
		expectGitHubToken      bool
	}{
		{
			name:                   "NoCheck",
			environment:            map[string]string{},
			ctx:                    context.Background(),
			checklist:              PreflightChecklist{},
			expectedError:          nil,
			expectWorkingDirectory: true,
		},
		{
			name:        "GitHubTokenCheckFails",
			environment: map[string]string{},
			ctx:         context.Background(),
			checklist: PreflightChecklist{
				GitHubToken: true,
			},
			expectedError: errors.New("GELATO_GITHUB_TOKEN environment variable not set"),
		},
		{
			name: "AllChecks",
			environment: map[string]string{
				"GELATO_GITHUB_TOKEN": "github-token",
			},
			ctx: context.Background(),
			checklist: PreflightChecklist{
				Go:          true,
				Git:         true,
				GitHubToken: true,
			},
			expectedError:          nil,
			expectWorkingDirectory: true,
			expectGoVersion:        true,
			expectGitVersion:       true,
			expectGitHubToken:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for key, val := range tc.environment {
				err := os.Setenv(key, val)
				assert.NoError(t, err)
				defer os.Unsetenv(key)
			}

			preflightInfo, err := RunPreflightChecks(tc.ctx, tc.checklist)

			if tc.expectedError != nil {
				assert.Zero(t, preflightInfo)
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectWorkingDirectory, preflightInfo.WorkingDirectory != "")
				assert.Equal(t, tc.expectGoVersion, preflightInfo.GoVersion != "")
				assert.Equal(t, tc.expectGitVersion, preflightInfo.GitVersion != "")
				assert.Equal(t, tc.expectGitHubToken, preflightInfo.GitHubToken != "")
			}
		})
	}
}

func TestResolveSemanticVersion(t *testing.T) {
	tests := []struct {
		name           string
		ctx            context.Context
		expectedSemVer semver.SemVer
		expectedError  error
	}{
		{
			name:          "Success",
			ctx:           context.Background(),
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			semver, err := ResolveSemanticVersion(tc.ctx)

			if tc.expectedError == nil {
				assert.NotEmpty(t, semver)
				assert.NoError(t, err)
			} else {
				assert.Empty(t, semver)
				assert.Equal(t, tc.expectedError, err)
			}
		})
	}
}
