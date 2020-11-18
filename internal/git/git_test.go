package git

import (
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{
			name: "OK",
			path: ".",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g, err := New(tc.path)

			assert.NoError(t, err)
			assert.NotNil(t, g)
			assert.NotNil(t, g.repo)
		})
	}
}

func TestGit_GetRemoteInfo(t *testing.T) {
	repo, err := git.PlainOpen("../..")
	assert.NoError(t, err)

	tests := []struct {
		name           string
		expectedDomain string
		expectedPath   string
		expectedError  string
	}{
		{
			name:           "OK",
			expectedDomain: "github.com",
			expectedPath:   "moorara/gelato",
			expectedError:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := &Git{
				repo: repo,
			}

			domain, path, err := g.GetRemoteInfo()

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedDomain, domain)
				assert.Equal(t, tc.expectedPath, path)
			} else {
				assert.Empty(t, domain)
				assert.Empty(t, path)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
