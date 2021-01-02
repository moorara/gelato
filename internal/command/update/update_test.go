package update

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/moorara/go-github"
	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/command"
)

type (
	LatestReleaseMock struct {
		InContext   context.Context
		OutRelease  *github.Release
		OutResponse *github.Response
		OutError    error
	}

	DownloadReleaseAssetMock struct {
		InContext    context.Context
		InReleaseTag string
		InAssetName  string
		InWriter     io.Writer
		OutResponse  *github.Response
		OutError     error
	}

	MockRepoService struct {
		LatestReleaseIndex int
		LatestReleaseMocks []LatestReleaseMock

		DownloadReleaseAssetIndex int
		DownloadReleaseAssetMocks []DownloadReleaseAssetMock
	}
)

func (m *MockRepoService) LatestRelease(ctx context.Context) (*github.Release, *github.Response, error) {
	i := m.LatestReleaseIndex
	m.LatestReleaseIndex++
	m.LatestReleaseMocks[i].InContext = ctx
	return m.LatestReleaseMocks[i].OutRelease, m.LatestReleaseMocks[i].OutResponse, m.LatestReleaseMocks[i].OutError
}

func (m *MockRepoService) DownloadReleaseAsset(ctx context.Context, releaseTag, assetName string, writer io.Writer) (*github.Response, error) {
	i := m.DownloadReleaseAssetIndex
	m.DownloadReleaseAssetIndex++
	m.DownloadReleaseAssetMocks[i].InContext = ctx
	m.DownloadReleaseAssetMocks[i].InReleaseTag = releaseTag
	m.DownloadReleaseAssetMocks[i].InAssetName = assetName
	m.DownloadReleaseAssetMocks[i].InWriter = writer
	return m.DownloadReleaseAssetMocks[i].OutResponse, m.DownloadReleaseAssetMocks[i].OutError
}

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

	assert.NotNil(t, c.services.repo)
}

func TestCommand_run(t *testing.T) {
	tests := []struct {
		name             string
		repo             *MockRepoService
		args             []string
		expectedExitCode int
	}{
		{
			name:             "UndefinedFlag",
			repo:             &MockRepoService{},
			args:             []string{"--undefined"},
			expectedExitCode: command.FlagError,
		},
		{
			name: "LatestReleaseFails",
			repo: &MockRepoService{
				LatestReleaseMocks: []LatestReleaseMock{
					{OutError: errors.New("error on getting the latest GitHub release")},
				},
			},
			args:             []string{},
			expectedExitCode: command.GitHubError,
		},
		/* {
			name: "DownloadReleaseAssetFails",
			repo: &MockRepoService{
				LatestReleaseMocks: []LatestReleaseMock{
					{
						OutRelease: &github.Release{
							Name:    "1.0.0",
							TagName: "v1.0.0",
						},
						OutResponse: &github.Response{},
					},
				},
				DownloadReleaseAssetMocks: []DownloadReleaseAssetMock{
					{OutError: errors.New("error on downloading the release asset")},
				},
			},
			args:             []string{},
			expectedExitCode: command.GitHubError,
		},
		{
			name: "Success",
			repo: &MockRepoService{
				LatestReleaseMocks: []LatestReleaseMock{
					{
						OutRelease: &github.Release{
							Name:    "1.0.0",
							TagName: "v1.0.0",
						},
						OutResponse: &github.Response{},
					},
				},
				DownloadReleaseAssetMocks: []DownloadReleaseAssetMock{
					{
						OutResponse: &github.Response{},
					},
				},
			},
			args:             []string{},
			expectedExitCode: command.Success,
		}, */
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &Command{ui: new(cli.MockUi)}
			c.services.repo = tc.repo

			exitCode := c.run(tc.args)

			assert.Equal(t, tc.expectedExitCode, exitCode)
		})
	}
}
