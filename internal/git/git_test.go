package git

import (
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		expectedError string
	}{
		{
			name:          "Unknown",
			path:          "/unknown",
			expectedError: "repository does not exist",
		},
		{
			name:          "Success",
			path:          ".",
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g, err := New(tc.path)

			if tc.expectedError != "" {
				assert.Nil(t, g)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, g)
				assert.NotNil(t, g.repo)
			}
		})
	}
}

func TestParseRemoteURL(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		expectedDomain string
		expectedPath   string
		expectedError  string
	}{
		{
			name:          "Empty",
			url:           "",
			expectedError: "invalid git remote url: ",
		},
		{
			name:          "Invalid",
			url:           "octocat/Hello-World",
			expectedError: "invalid git remote url: octocat/Hello-World",
		},
		{
			name:           "SSH",
			url:            "git@github.com:octocat/Hello-World",
			expectedDomain: "github.com",
			expectedPath:   "octocat/Hello-World",
		},
		{
			name:           "SSH_git",
			url:            "git@github.com:octocat/Hello-World.git",
			expectedDomain: "github.com",
			expectedPath:   "octocat/Hello-World",
		},
		{
			name:           "HTTPS",
			url:            "https://github.com/octocat/Hello-World",
			expectedDomain: "github.com",
			expectedPath:   "octocat/Hello-World",
		},
		{
			name:           "HTTPS_git",
			url:            "https://github.com/octocat/Hello-World.git",
			expectedDomain: "github.com",
			expectedPath:   "octocat/Hello-World",
		},
		{
			name:          "SSHVariant",
			url:           "ssh://git@github.com/octocat/Hello-World",
			expectedError: "invalid git remote url: ssh://git@github.com/octocat/Hello-World",
		},
		{
			name:          "SSHVariant_git",
			url:           "ssh://git@github.com/octocat/Hello-World.git",
			expectedError: "invalid git remote url: ssh://git@github.com/octocat/Hello-World.git",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			domain, path, err := parseRemoteURL(tc.url)

			if tc.expectedError != "" {
				assert.Empty(t, domain)
				assert.Empty(t, path)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedDomain, domain)
				assert.Equal(t, tc.expectedPath, path)
			}
		})
	}
}

func TestGit_Remote(t *testing.T) {
	repo, err := git.PlainOpen("../..")
	assert.NoError(t, err)

	tests := []struct {
		name           string
		remoteName     string
		expectedDomain string
		expectedPath   string
		expectedError  string
	}{
		{
			name:          "Unknown",
			remoteName:    "unknown",
			expectedError: "remote not found",
		},
		{
			name:           "Success",
			remoteName:     "origin",
			expectedDomain: "github.com",
			expectedPath:   "moorara/gelato",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := &Git{repo: repo}

			domain, path, err := g.Remote(tc.remoteName)

			if tc.expectedError != "" {
				assert.Empty(t, domain)
				assert.Empty(t, path)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedDomain, domain)
				assert.Equal(t, tc.expectedPath, path)
			}
		})
	}
}

func TestGit_IsClean(t *testing.T) {
	repo, err := git.PlainOpen("../..")
	assert.NoError(t, err)

	g := &Git{repo: repo}

	_, err = g.IsClean()
	assert.NoError(t, err)
}

func TestGit_HEAD(t *testing.T) {
	repo, err := git.PlainOpen("../..")
	assert.NoError(t, err)

	g := &Git{repo: repo}

	hash, branch, err := g.HEAD()
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEmpty(t, branch)
}

func TestGit_Tag(t *testing.T) {
	// TODO: uncomment after releasing
	// repo, err := git.PlainOpen("../..")
	// assert.NoError(t, err)

	// g := &Git{repo: repo}

	// tag, err := g.Tag("v0.1.0")
	// assert.NoError(t, err)
	// assert.NotEmpty(t, tag)
}

func TestGit_Tags(t *testing.T) {
	// TODO: uncomment after releasing
	// repo, err := git.PlainOpen("../..")
	// assert.NoError(t, err)

	// g := &Git{repo: repo}

	// tags, err := g.Tags()
	// assert.NoError(t, err)
	// assert.NotEmpty(t, tags)
}

func TestGit_CreateTag(t *testing.T) {
	// CreateTag has side effects!
}

func TestGit_CommitsIn(t *testing.T) {
	repo, err := git.PlainOpen("../..")
	assert.NoError(t, err)

	g := &Git{repo: repo}

	commits, err := g.CommitsIn("HEAD")
	assert.NoError(t, err)
	assert.NotEmpty(t, commits)
}

func TestGit_CreateCommit(t *testing.T) {
	// CreateCommit has side effects!
}

func TestGit_Pull(t *testing.T) {
	// Pull has side effects!
}

func TestGit_Push(t *testing.T) {
	// Push has side effects!
}

func TestGit_PushTag(t *testing.T) {
	// PushTag has side effects!
}
