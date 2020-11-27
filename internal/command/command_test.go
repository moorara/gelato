package command

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunPreflightChecks(t *testing.T) {
	tests := []struct {
		name                    string
		environment             map[string]string
		ctx                     context.Context
		checklist               PreflightChecklist
		expectedError           error
		expectWorkingDirectory  bool
		expectGoVersion         bool
		expectedGitRemoteDomain string
		expectedGitRemotePath   string
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
			name:        "AllChecks",
			environment: map[string]string{},
			ctx:         context.Background(),
			checklist: PreflightChecklist{
				Go:  true,
				Git: true,
			},
			expectedError:           nil,
			expectWorkingDirectory:  true,
			expectGoVersion:         true,
			expectedGitRemoteDomain: "github.com",
			expectedGitRemotePath:   "moorara/gelato",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for key, val := range tc.environment {
				err := os.Setenv(key, val)
				assert.NoError(t, err)
				defer os.Unsetenv(key)
			}

			info, err := RunPreflightChecks(tc.ctx, tc.checklist)

			if tc.expectedError != nil {
				assert.Zero(t, info)
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectWorkingDirectory, info.Context.WorkingDirectory != "")
				assert.Equal(t, tc.expectGoVersion, info.Go.Version != "")
				assert.Equal(t, tc.expectedGitRemoteDomain, info.Git.Remote.Domain)
				assert.Equal(t, tc.expectedGitRemotePath, info.Git.Remote.Path)
			}
		})
	}
}
