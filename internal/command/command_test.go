package command

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCheckAndCreateContext(t *testing.T) {
	tests := []struct {
		name          string
		environment   map[string]string
		checklist     PreflightChecklist
		timeout       time.Duration
		expectedError error
	}{
		{
			name:          "DefaultChecks",
			environment:   map[string]string{},
			checklist:     PreflightChecklist{},
			timeout:       2 * time.Second,
			expectedError: nil,
		},
		{
			name: "AllChecks",
			environment: map[string]string{
				"GELATO_GITHUB_TOKEN": "github-token",
			},
			checklist: PreflightChecklist{
				Go:          true,
				Git:         true,
				GitHubToken: true,
			},
			timeout:       2 * time.Second,
			expectedError: nil,
		},
		{
			name:        "GitHubTokenCheckFails",
			environment: map[string]string{},
			checklist: PreflightChecklist{
				Go:          true,
				Git:         true,
				GitHubToken: true,
			},
			timeout:       2 * time.Second,
			expectedError: errors.New("GELATO_GITHUB_TOKEN environment variable not set"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for key, val := range tc.environment {
				err := os.Setenv(key, val)
				assert.NoError(t, err)
				defer os.Unsetenv(key)
			}

			ctx, cancel, err := CheckAndCreateContext(tc.checklist, tc.timeout)
			if err == nil {
				defer cancel()
			}

			if tc.expectedError == nil {
				assert.NotNil(t, ctx)
				assert.NotNil(t, cancel)
				assert.NoError(t, err)
			} else {
				assert.Nil(t, ctx)
				assert.Nil(t, cancel)
				assert.Equal(t, tc.expectedError, err)
			}
		})
	}
}
