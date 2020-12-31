package release

import (
	"errors"
	"os"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"

	changelogSpec "github.com/moorara/changelog/spec"

	"github.com/moorara/gelato/internal/command"
	buildcmd "github.com/moorara/gelato/internal/command/build"
	"github.com/moorara/gelato/internal/spec"
	"github.com/moorara/gelato/pkg/semver"
	"github.com/moorara/go-github"
)

func TestNewCommand(t *testing.T) {
	ui := cli.NewMockUi()
	spec := spec.Spec{}
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
	tests := []struct {
		name             string
		environment      map[string]string
		expectedExitCode int
	}{
		{
			name: "NoToken",
			environment: map[string]string{
				"GELATO_GITHUB_TOKEN": "",
			},
			expectedExitCode: command.GitHubError,
		},
		{
			name: "OK",
			environment: map[string]string{
				"GELATO_GITHUB_TOKEN": "github-access-token",
			},
			expectedExitCode: command.FlagError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for key, val := range tc.environment {
				err := os.Setenv(key, val)
				assert.NoError(t, err)
				defer os.Unsetenv(key)
			}

			c := &Command{ui: cli.NewMockUi()}

			exitCode := c.Run([]string{"--undefined"})

			assert.Equal(t, tc.expectedExitCode, exitCode)

			if tc.expectedExitCode == command.Success {
				assert.Equal(t, "moorara", c.data.owner)
				assert.Equal(t, "gelato", c.data.repo)
				assert.NotEmpty(t, c.data.changelogSpec)
				assert.NotNil(t, c.services.git)
				assert.NotNil(t, c.services.users)
				assert.NotNil(t, c.services.repo)
				assert.NotNil(t, c.services.changelog)
				assert.NotNil(t, c.commands.semver)
				assert.NotNil(t, c.commands.build)
			}
		})
	}
}

func TestCommand_run(t *testing.T) {
	user := github.User{
		Login: "octocat",
	}

	repo := github.Repository{
		Name:          "Hello-World",
		FullName:      "octocat/Hello-World",
		DefaultBranch: "main",
	}

	version := semver.SemVer{
		Major: 0, Minor: 1, Patch: 0,
		Prerelease: []string{"2", "605a46c"},
	}

	draftRelease := github.Release{
		Name:       "0.1.0",
		TagName:    "v0.1.0",
		Target:     "main",
		Draft:      true,
		Prerelease: false,
	}

	artifacts := []buildcmd.Artifact{
		{
			Path:  "bin/app",
			Label: "linux",
		},
	}

	asset := github.ReleaseAsset{
		Name:  "bin/app",
		Label: "linux",
	}

	release := github.Release{
		Name:       "0.1.0",
		TagName:    "v0.1.0",
		Target:     "main",
		Draft:      false,
		Prerelease: false,
	}

	tests := []struct {
		name             string
		spec             spec.Spec
		git              *MockGitService
		users            *MockUsersService
		repo             *MockRepoService
		changelog        *MockChangelogService
		semver           *MockSemverCommand
		build            *MockBuildCommand
		args             []string
		expectedExitCode int
	}{
		{
			name:             "UndefinedFlag",
			args:             []string{"--undefined"},
			expectedExitCode: command.FlagError,
		},
		{
			name: "RepoGetFails",
			spec: spec.Spec{},
			repo: &MockRepoService{
				GetMocks: []GetMock{
					{OutError: errors.New("github error")},
				},
			},
			args:             []string{},
			expectedExitCode: command.GitHubError,
		},
		{
			name: "GitHEADFails",
			spec: spec.Spec{},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutError: errors.New("git error")},
				},
			},
			repo: &MockRepoService{
				GetMocks: []GetMock{
					{OutRepository: &repo, OutResponse: &github.Response{}},
				},
			},
			args:             []string{},
			expectedExitCode: command.GitError,
		},
		{
			name: "NotOnDefaultBranch",
			spec: spec.Spec{},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutBranch: "feature-branch"},
				},
			},
			repo: &MockRepoService{
				GetMocks: []GetMock{
					{OutRepository: &repo, OutResponse: &github.Response{}},
				},
			},
			args:             []string{},
			expectedExitCode: command.GitError,
		},
		{
			name: "GitIsCleanFails",
			spec: spec.Spec{},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutBranch: "main"},
				},
				IsCleanMocks: []IsCleanMock{
					{OutError: errors.New("git error")},
				},
			},
			repo: &MockRepoService{
				GetMocks: []GetMock{
					{OutRepository: &repo, OutResponse: &github.Response{}},
				},
			},
			args:             []string{},
			expectedExitCode: command.GitError,
		},
		{
			name: "RepoNotClean",
			spec: spec.Spec{},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutBranch: "main"},
				},
				IsCleanMocks: []IsCleanMock{
					{OutBool: false},
				},
			},
			repo: &MockRepoService{
				GetMocks: []GetMock{
					{OutRepository: &repo, OutResponse: &github.Response{}},
				},
			},
			args:             []string{},
			expectedExitCode: command.GitError,
		},
		{
			name: "UsersUserFails",
			spec: spec.Spec{},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutBranch: "main"},
				},
				IsCleanMocks: []IsCleanMock{
					{OutBool: true},
				},
			},
			users: &MockUsersService{
				UserMocks: []UserMock{
					{OutError: errors.New("github error")},
				},
			},
			repo: &MockRepoService{
				GetMocks: []GetMock{
					{OutRepository: &repo, OutResponse: &github.Response{}},
				},
			},
			args:             []string{},
			expectedExitCode: command.GitHubError,
		},
		{
			name: "RepoPermissionFails",
			spec: spec.Spec{},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutBranch: "main"},
				},
				IsCleanMocks: []IsCleanMock{
					{OutBool: true},
				},
			},
			users: &MockUsersService{
				UserMocks: []UserMock{
					{OutUser: &user, OutResponse: &github.Response{}},
				},
			},
			repo: &MockRepoService{
				GetMocks: []GetMock{
					{OutRepository: &repo, OutResponse: &github.Response{}},
				},
				PermissionMocks: []PermissionMock{
					{OutError: errors.New("github error")},
				},
			},
			args:             []string{},
			expectedExitCode: command.GitHubError,
		},
		{
			name: "InvalidUserPermission",
			spec: spec.Spec{},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutBranch: "main"},
				},
				IsCleanMocks: []IsCleanMock{
					{OutBool: true},
				},
			},
			users: &MockUsersService{
				UserMocks: []UserMock{
					{OutUser: &user, OutResponse: &github.Response{}},
				},
			},
			repo: &MockRepoService{
				GetMocks: []GetMock{
					{OutRepository: &repo, OutResponse: &github.Response{}},
				},
				PermissionMocks: []PermissionMock{
					{OutPermission: github.PermissionWrite, OutResponse: &github.Response{}},
				},
			},
			args:             []string{},
			expectedExitCode: command.GitHubError,
		},
		{
			name: "GitPullFails",
			spec: spec.Spec{},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutBranch: "main"},
				},
				IsCleanMocks: []IsCleanMock{
					{OutBool: true},
				},
				PullMocks: []PullMock{
					{OutError: errors.New("git error")},
				},
			},
			users: &MockUsersService{
				UserMocks: []UserMock{
					{OutUser: &user, OutResponse: &github.Response{}},
				},
			},
			repo: &MockRepoService{
				GetMocks: []GetMock{
					{OutRepository: &repo, OutResponse: &github.Response{}},
				},
				PermissionMocks: []PermissionMock{
					{OutPermission: github.PermissionAdmin, OutResponse: &github.Response{}},
				},
			},
			args:             []string{},
			expectedExitCode: command.GitError,
		},
		{
			name: "SemverRunFails",
			spec: spec.Spec{},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutBranch: "main"},
				},
				IsCleanMocks: []IsCleanMock{
					{OutBool: true},
				},
				PullMocks: []PullMock{
					{OutError: nil},
				},
			},
			users: &MockUsersService{
				UserMocks: []UserMock{
					{OutUser: &user, OutResponse: &github.Response{}},
				},
			},
			repo: &MockRepoService{
				GetMocks: []GetMock{
					{OutRepository: &repo, OutResponse: &github.Response{}},
				},
				PermissionMocks: []PermissionMock{
					{OutPermission: github.PermissionAdmin, OutResponse: &github.Response{}},
				},
			},
			semver: &MockSemverCommand{
				RunMocks: []SemverRunMock{
					{OutCode: command.GitError},
				},
			},
			args:             []string{},
			expectedExitCode: command.GitError,
		},
		{
			name: "CreateReleaseFails",
			spec: spec.Spec{},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutBranch: "main"},
				},
				IsCleanMocks: []IsCleanMock{
					{OutBool: true},
				},
				PullMocks: []PullMock{
					{OutError: nil},
				},
			},
			users: &MockUsersService{
				UserMocks: []UserMock{
					{OutUser: &user, OutResponse: &github.Response{}},
				},
			},
			repo: &MockRepoService{
				GetMocks: []GetMock{
					{OutRepository: &repo, OutResponse: &github.Response{}},
				},
				PermissionMocks: []PermissionMock{
					{OutPermission: github.PermissionAdmin, OutResponse: &github.Response{}},
				},
				CreateReleaseMocks: []CreateReleaseMock{
					{OutError: errors.New("github error")},
				},
			},
			semver: &MockSemverCommand{
				RunMocks: []SemverRunMock{
					{OutCode: command.Success},
				},
				SemVerMocks: []SemVerMock{
					{OutSemVer: version},
				},
			},
			args:             []string{},
			expectedExitCode: command.GitHubError,
		},
		{
			name: "ChangelogGenerateFails",
			spec: spec.Spec{},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutBranch: "main"},
				},
				IsCleanMocks: []IsCleanMock{
					{OutBool: true},
				},
				PullMocks: []PullMock{
					{OutError: nil},
				},
			},
			users: &MockUsersService{
				UserMocks: []UserMock{
					{OutUser: &user, OutResponse: &github.Response{}},
				},
			},
			repo: &MockRepoService{
				GetMocks: []GetMock{
					{OutRepository: &repo, OutResponse: &github.Response{}},
				},
				PermissionMocks: []PermissionMock{
					{OutPermission: github.PermissionAdmin, OutResponse: &github.Response{}},
				},
				CreateReleaseMocks: []CreateReleaseMock{
					{OutRelease: &draftRelease, OutResponse: &github.Response{}},
				},
			},
			changelog: &MockChangelogService{
				GenerateMocks: []GenerateMock{
					{OutError: errors.New("changelog error")},
				},
			},
			semver: &MockSemverCommand{
				RunMocks: []SemverRunMock{
					{OutCode: command.Success},
				},
				SemVerMocks: []SemVerMock{
					{OutSemVer: version},
				},
			},
			args:             []string{},
			expectedExitCode: command.ChangelogError,
		},
		{
			name: "CreateCommitFails",
			spec: spec.Spec{},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutBranch: "main"},
				},
				IsCleanMocks: []IsCleanMock{
					{OutBool: true},
				},
				PullMocks: []PullMock{
					{OutError: nil},
				},
				CreateCommitMocks: []CreateCommitMock{
					{OutError: errors.New("git error")},
				},
			},
			users: &MockUsersService{
				UserMocks: []UserMock{
					{OutUser: &user, OutResponse: &github.Response{}},
				},
			},
			repo: &MockRepoService{
				GetMocks: []GetMock{
					{OutRepository: &repo, OutResponse: &github.Response{}},
				},
				PermissionMocks: []PermissionMock{
					{OutPermission: github.PermissionAdmin, OutResponse: &github.Response{}},
				},
				CreateReleaseMocks: []CreateReleaseMock{
					{OutRelease: &draftRelease, OutResponse: &github.Response{}},
				},
			},
			changelog: &MockChangelogService{
				GenerateMocks: []GenerateMock{
					{OutContent: "changelog content"},
				},
			},
			semver: &MockSemverCommand{
				RunMocks: []SemverRunMock{
					{OutCode: command.Success},
				},
				SemVerMocks: []SemVerMock{
					{OutSemVer: version},
				},
			},
			args:             []string{},
			expectedExitCode: command.GitError,
		},
		{
			name: "CreateTagFails",
			spec: spec.Spec{},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutBranch: "main"},
				},
				IsCleanMocks: []IsCleanMock{
					{OutBool: true},
				},
				PullMocks: []PullMock{
					{OutError: nil},
				},
				CreateCommitMocks: []CreateCommitMock{
					{OutHash: "6e8c7d217faab1d88905d4c75b4e7995a42c81d5"},
				},
				CreateTagMocks: []CreateTagMock{
					{OutError: errors.New("git error")},
				},
			},
			users: &MockUsersService{
				UserMocks: []UserMock{
					{OutUser: &user, OutResponse: &github.Response{}},
				},
			},
			repo: &MockRepoService{
				GetMocks: []GetMock{
					{OutRepository: &repo, OutResponse: &github.Response{}},
				},
				PermissionMocks: []PermissionMock{
					{OutPermission: github.PermissionAdmin, OutResponse: &github.Response{}},
				},
				CreateReleaseMocks: []CreateReleaseMock{
					{OutRelease: &draftRelease, OutResponse: &github.Response{}},
				},
			},
			changelog: &MockChangelogService{
				GenerateMocks: []GenerateMock{
					{OutContent: "changelog content"},
				},
			},
			semver: &MockSemverCommand{
				RunMocks: []SemverRunMock{
					{OutCode: command.Success},
				},
				SemVerMocks: []SemVerMock{
					{OutSemVer: version},
				},
			},
			args:             []string{},
			expectedExitCode: command.GitError,
		},
		{
			name: "BuildRunFails",
			spec: spec.Spec{
				Release: spec.Release{
					Artifacts: true,
				},
			},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutBranch: "main"},
				},
				IsCleanMocks: []IsCleanMock{
					{OutBool: true},
				},
				PullMocks: []PullMock{
					{OutError: nil},
				},
				CreateCommitMocks: []CreateCommitMock{
					{OutHash: "6e8c7d217faab1d88905d4c75b4e7995a42c81d5"},
				},
				CreateTagMocks: []CreateTagMock{
					{OutHash: "a3580a0f64b08ba6085d530c828c40b8aa082c1e"},
				},
			},
			users: &MockUsersService{
				UserMocks: []UserMock{
					{OutUser: &user, OutResponse: &github.Response{}},
				},
			},
			repo: &MockRepoService{
				GetMocks: []GetMock{
					{OutRepository: &repo, OutResponse: &github.Response{}},
				},
				PermissionMocks: []PermissionMock{
					{OutPermission: github.PermissionAdmin, OutResponse: &github.Response{}},
				},
				CreateReleaseMocks: []CreateReleaseMock{
					{OutRelease: &draftRelease, OutResponse: &github.Response{}},
				},
			},
			changelog: &MockChangelogService{
				GenerateMocks: []GenerateMock{
					{OutContent: "changelog content"},
				},
			},
			semver: &MockSemverCommand{
				RunMocks: []SemverRunMock{
					{OutCode: command.Success},
				},
				SemVerMocks: []SemVerMock{
					{OutSemVer: version},
				},
			},
			build: &MockBuildCommand{
				RunMocks: []BuildRunMock{
					{OutCode: command.GoError},
				},
			},
			args:             []string{},
			expectedExitCode: command.GoError,
		},
		{
			name: "UploadReleaseAssetFails",
			spec: spec.Spec{
				Release: spec.Release{
					Artifacts: true,
				},
			},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutBranch: "main"},
				},
				IsCleanMocks: []IsCleanMock{
					{OutBool: true},
				},
				PullMocks: []PullMock{
					{OutError: nil},
				},
				CreateCommitMocks: []CreateCommitMock{
					{OutHash: "6e8c7d217faab1d88905d4c75b4e7995a42c81d5"},
				},
				CreateTagMocks: []CreateTagMock{
					{OutHash: "a3580a0f64b08ba6085d530c828c40b8aa082c1e"},
				},
			},
			users: &MockUsersService{
				UserMocks: []UserMock{
					{OutUser: &user, OutResponse: &github.Response{}},
				},
			},
			repo: &MockRepoService{
				GetMocks: []GetMock{
					{OutRepository: &repo, OutResponse: &github.Response{}},
				},
				PermissionMocks: []PermissionMock{
					{OutPermission: github.PermissionAdmin, OutResponse: &github.Response{}},
				},
				CreateReleaseMocks: []CreateReleaseMock{
					{OutRelease: &draftRelease, OutResponse: &github.Response{}},
				},
				UploadReleaseAssetMocks: []UploadReleaseAssetMock{
					{OutError: errors.New("github error")},
				},
			},
			changelog: &MockChangelogService{
				GenerateMocks: []GenerateMock{
					{OutContent: "changelog content"},
				},
			},
			semver: &MockSemverCommand{
				RunMocks: []SemverRunMock{
					{OutCode: command.Success},
				},
				SemVerMocks: []SemVerMock{
					{OutSemVer: version},
				},
			},
			build: &MockBuildCommand{
				RunMocks: []BuildRunMock{
					{OutCode: command.Success},
				},
				ArtifactsMocks: []ArtifactsMock{
					{OutArtifacts: artifacts},
				},
			},
			args:             []string{},
			expectedExitCode: command.GitHubError,
		},
		{
			name: "EnableBranchProtectionFails",
			spec: spec.Spec{
				Release: spec.Release{
					Artifacts: true,
				},
			},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutBranch: "main"},
				},
				IsCleanMocks: []IsCleanMock{
					{OutBool: true},
				},
				PullMocks: []PullMock{
					{OutError: nil},
				},
				CreateCommitMocks: []CreateCommitMock{
					{OutHash: "6e8c7d217faab1d88905d4c75b4e7995a42c81d5"},
				},
				CreateTagMocks: []CreateTagMock{
					{OutHash: "a3580a0f64b08ba6085d530c828c40b8aa082c1e"},
				},
			},
			users: &MockUsersService{
				UserMocks: []UserMock{
					{OutUser: &user, OutResponse: &github.Response{}},
				},
			},
			repo: &MockRepoService{
				GetMocks: []GetMock{
					{OutRepository: &repo, OutResponse: &github.Response{}},
				},
				PermissionMocks: []PermissionMock{
					{OutPermission: github.PermissionAdmin, OutResponse: &github.Response{}},
				},
				CreateReleaseMocks: []CreateReleaseMock{
					{OutRelease: &draftRelease, OutResponse: &github.Response{}},
				},
				UploadReleaseAssetMocks: []UploadReleaseAssetMock{
					{OutReleaseAsset: &asset, OutResponse: &github.Response{}},
				},
				BranchProtectionMocks: []BranchProtectionMock{
					{OutError: errors.New("github error")},
				},
			},
			changelog: &MockChangelogService{
				GenerateMocks: []GenerateMock{
					{OutContent: "changelog content"},
				},
			},
			semver: &MockSemverCommand{
				RunMocks: []SemverRunMock{
					{OutCode: command.Success},
				},
				SemVerMocks: []SemVerMock{
					{OutSemVer: version},
				},
			},
			build: &MockBuildCommand{
				RunMocks: []BuildRunMock{
					{OutCode: command.Success},
				},
				ArtifactsMocks: []ArtifactsMock{
					{OutArtifacts: artifacts},
				},
			},
			args:             []string{},
			expectedExitCode: command.GitHubError,
		},
		{
			name: "PushFails",
			spec: spec.Spec{
				Release: spec.Release{
					Artifacts: true,
				},
			},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutBranch: "main"},
				},
				IsCleanMocks: []IsCleanMock{
					{OutBool: true},
				},
				PullMocks: []PullMock{
					{OutError: nil},
				},
				CreateCommitMocks: []CreateCommitMock{
					{OutHash: "6e8c7d217faab1d88905d4c75b4e7995a42c81d5"},
				},
				CreateTagMocks: []CreateTagMock{
					{OutHash: "a3580a0f64b08ba6085d530c828c40b8aa082c1e"},
				},
				PushMocks: []PushMock{
					{OutError: errors.New("git error")},
				},
			},
			users: &MockUsersService{
				UserMocks: []UserMock{
					{OutUser: &user, OutResponse: &github.Response{}},
				},
			},
			repo: &MockRepoService{
				GetMocks: []GetMock{
					{OutRepository: &repo, OutResponse: &github.Response{}},
				},
				PermissionMocks: []PermissionMock{
					{OutPermission: github.PermissionAdmin, OutResponse: &github.Response{}},
				},
				CreateReleaseMocks: []CreateReleaseMock{
					{OutRelease: &draftRelease, OutResponse: &github.Response{}},
				},
				UploadReleaseAssetMocks: []UploadReleaseAssetMock{
					{OutReleaseAsset: &asset, OutResponse: &github.Response{}},
				},
				BranchProtectionMocks: []BranchProtectionMock{
					{OutResponse: &github.Response{}},
					{OutResponse: &github.Response{}},
				},
			},
			changelog: &MockChangelogService{
				GenerateMocks: []GenerateMock{
					{OutContent: "changelog content"},
				},
			},
			semver: &MockSemverCommand{
				RunMocks: []SemverRunMock{
					{OutCode: command.Success},
				},
				SemVerMocks: []SemVerMock{
					{OutSemVer: version},
				},
			},
			build: &MockBuildCommand{
				RunMocks: []BuildRunMock{
					{OutCode: command.Success},
				},
				ArtifactsMocks: []ArtifactsMock{
					{OutArtifacts: artifacts},
				},
			},
			args:             []string{},
			expectedExitCode: command.GitError,
		},
		{
			name: "PushTagFails",
			spec: spec.Spec{
				Release: spec.Release{
					Artifacts: true,
				},
			},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutBranch: "main"},
				},
				IsCleanMocks: []IsCleanMock{
					{OutBool: true},
				},
				PullMocks: []PullMock{
					{OutError: nil},
				},
				CreateCommitMocks: []CreateCommitMock{
					{OutHash: "6e8c7d217faab1d88905d4c75b4e7995a42c81d5"},
				},
				CreateTagMocks: []CreateTagMock{
					{OutHash: "a3580a0f64b08ba6085d530c828c40b8aa082c1e"},
				},
				PushMocks: []PushMock{
					{OutError: nil},
				},
				PushTagMocks: []PushTagMock{
					{OutError: errors.New("git error")},
				},
			},
			users: &MockUsersService{
				UserMocks: []UserMock{
					{OutUser: &user, OutResponse: &github.Response{}},
				},
			},
			repo: &MockRepoService{
				GetMocks: []GetMock{
					{OutRepository: &repo, OutResponse: &github.Response{}},
				},
				PermissionMocks: []PermissionMock{
					{OutPermission: github.PermissionAdmin, OutResponse: &github.Response{}},
				},
				CreateReleaseMocks: []CreateReleaseMock{
					{OutRelease: &draftRelease, OutResponse: &github.Response{}},
				},
				UploadReleaseAssetMocks: []UploadReleaseAssetMock{
					{OutReleaseAsset: &asset, OutResponse: &github.Response{}},
				},
				BranchProtectionMocks: []BranchProtectionMock{
					{OutResponse: &github.Response{}},
					{OutResponse: &github.Response{}},
				},
			},
			changelog: &MockChangelogService{
				GenerateMocks: []GenerateMock{
					{OutContent: "changelog content"},
				},
			},
			semver: &MockSemverCommand{
				RunMocks: []SemverRunMock{
					{OutCode: command.Success},
				},
				SemVerMocks: []SemVerMock{
					{OutSemVer: version},
				},
			},
			build: &MockBuildCommand{
				RunMocks: []BuildRunMock{
					{OutCode: command.Success},
				},
				ArtifactsMocks: []ArtifactsMock{
					{OutArtifacts: artifacts},
				},
			},
			args:             []string{},
			expectedExitCode: command.GitError,
		},
		{
			name: "UpdateReleaseFails",
			spec: spec.Spec{
				Release: spec.Release{
					Artifacts: true,
				},
			},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutBranch: "main"},
				},
				IsCleanMocks: []IsCleanMock{
					{OutBool: true},
				},
				PullMocks: []PullMock{
					{OutError: nil},
				},
				CreateCommitMocks: []CreateCommitMock{
					{OutHash: "6e8c7d217faab1d88905d4c75b4e7995a42c81d5"},
				},
				CreateTagMocks: []CreateTagMock{
					{OutHash: "a3580a0f64b08ba6085d530c828c40b8aa082c1e"},
				},
				PushMocks: []PushMock{
					{OutError: nil},
				},
				PushTagMocks: []PushTagMock{
					{OutError: nil},
				},
			},
			users: &MockUsersService{
				UserMocks: []UserMock{
					{OutUser: &user, OutResponse: &github.Response{}},
				},
			},
			repo: &MockRepoService{
				GetMocks: []GetMock{
					{OutRepository: &repo, OutResponse: &github.Response{}},
				},
				PermissionMocks: []PermissionMock{
					{OutPermission: github.PermissionAdmin, OutResponse: &github.Response{}},
				},
				CreateReleaseMocks: []CreateReleaseMock{
					{OutRelease: &draftRelease, OutResponse: &github.Response{}},
				},
				UploadReleaseAssetMocks: []UploadReleaseAssetMock{
					{OutReleaseAsset: &asset, OutResponse: &github.Response{}},
				},
				BranchProtectionMocks: []BranchProtectionMock{
					{OutResponse: &github.Response{}},
					{OutResponse: &github.Response{}},
				},
				UpdateReleaseMocks: []UpdateReleaseMock{
					{OutError: errors.New("github error")},
				},
			},
			changelog: &MockChangelogService{
				GenerateMocks: []GenerateMock{
					{OutContent: "changelog content"},
				},
			},
			semver: &MockSemverCommand{
				RunMocks: []SemverRunMock{
					{OutCode: command.Success},
				},
				SemVerMocks: []SemVerMock{
					{OutSemVer: version},
				},
			},
			build: &MockBuildCommand{
				RunMocks: []BuildRunMock{
					{OutCode: command.Success},
				},
				ArtifactsMocks: []ArtifactsMock{
					{OutArtifacts: artifacts},
				},
			},
			args:             []string{},
			expectedExitCode: command.GitHubError,
		},
		{
			name: "Success_PatchRelease",
			spec: spec.Spec{
				Release: spec.Release{
					Artifacts: true,
				},
			},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutBranch: "main"},
				},
				IsCleanMocks: []IsCleanMock{
					{OutBool: true},
				},
				PullMocks: []PullMock{
					{OutError: nil},
				},
				CreateCommitMocks: []CreateCommitMock{
					{OutHash: "6e8c7d217faab1d88905d4c75b4e7995a42c81d5"},
				},
				CreateTagMocks: []CreateTagMock{
					{OutHash: "a3580a0f64b08ba6085d530c828c40b8aa082c1e"},
				},
				PushMocks: []PushMock{
					{OutError: nil},
				},
				PushTagMocks: []PushTagMock{
					{OutError: nil},
				},
			},
			users: &MockUsersService{
				UserMocks: []UserMock{
					{OutUser: &user, OutResponse: &github.Response{}},
				},
			},
			repo: &MockRepoService{
				GetMocks: []GetMock{
					{OutRepository: &repo, OutResponse: &github.Response{}},
				},
				PermissionMocks: []PermissionMock{
					{OutPermission: github.PermissionAdmin, OutResponse: &github.Response{}},
				},
				CreateReleaseMocks: []CreateReleaseMock{
					{OutRelease: &draftRelease, OutResponse: &github.Response{}},
				},
				UploadReleaseAssetMocks: []UploadReleaseAssetMock{
					{OutReleaseAsset: &asset, OutResponse: &github.Response{}},
				},
				BranchProtectionMocks: []BranchProtectionMock{
					{OutResponse: &github.Response{}},
					{OutResponse: &github.Response{}},
				},
				UpdateReleaseMocks: []UpdateReleaseMock{
					{OutRelease: &release, OutResponse: &github.Response{}},
				},
			},
			changelog: &MockChangelogService{
				GenerateMocks: []GenerateMock{
					{OutContent: "changelog content"},
				},
			},
			semver: &MockSemverCommand{
				RunMocks: []SemverRunMock{
					{OutCode: command.Success},
				},
				SemVerMocks: []SemVerMock{
					{OutSemVer: version},
				},
			},
			build: &MockBuildCommand{
				RunMocks: []BuildRunMock{
					{OutCode: command.Success},
				},
				ArtifactsMocks: []ArtifactsMock{
					{OutArtifacts: artifacts},
				},
			},
			args:             []string{"-comment", "Release description"},
			expectedExitCode: command.Success,
		},
		{
			name: "Success_MinorRelease",
			spec: spec.Spec{
				Release: spec.Release{
					Artifacts: true,
				},
			},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutBranch: "main"},
				},
				IsCleanMocks: []IsCleanMock{
					{OutBool: true},
				},
				PullMocks: []PullMock{
					{OutError: nil},
				},
				CreateCommitMocks: []CreateCommitMock{
					{OutHash: "6e8c7d217faab1d88905d4c75b4e7995a42c81d5"},
				},
				CreateTagMocks: []CreateTagMock{
					{OutHash: "a3580a0f64b08ba6085d530c828c40b8aa082c1e"},
				},
				PushMocks: []PushMock{
					{OutError: nil},
				},
				PushTagMocks: []PushTagMock{
					{OutError: nil},
				},
			},
			users: &MockUsersService{
				UserMocks: []UserMock{
					{OutUser: &user, OutResponse: &github.Response{}},
				},
			},
			repo: &MockRepoService{
				GetMocks: []GetMock{
					{OutRepository: &repo, OutResponse: &github.Response{}},
				},
				PermissionMocks: []PermissionMock{
					{OutPermission: github.PermissionAdmin, OutResponse: &github.Response{}},
				},
				CreateReleaseMocks: []CreateReleaseMock{
					{OutRelease: &draftRelease, OutResponse: &github.Response{}},
				},
				UploadReleaseAssetMocks: []UploadReleaseAssetMock{
					{OutReleaseAsset: &asset, OutResponse: &github.Response{}},
				},
				BranchProtectionMocks: []BranchProtectionMock{
					{OutResponse: &github.Response{}},
					{OutResponse: &github.Response{}},
				},
				UpdateReleaseMocks: []UpdateReleaseMock{
					{OutRelease: &release, OutResponse: &github.Response{}},
				},
			},
			changelog: &MockChangelogService{
				GenerateMocks: []GenerateMock{
					{OutContent: "changelog content"},
				},
			},
			semver: &MockSemverCommand{
				RunMocks: []SemverRunMock{
					{OutCode: command.Success},
				},
				SemVerMocks: []SemVerMock{
					{OutSemVer: version},
				},
			},
			build: &MockBuildCommand{
				RunMocks: []BuildRunMock{
					{OutCode: command.Success},
				},
				ArtifactsMocks: []ArtifactsMock{
					{OutArtifacts: artifacts},
				},
			},
			args:             []string{"-minor", "-comment", "Release description"},
			expectedExitCode: command.Success,
		},
		{
			name: "Success_MajorRelease",
			spec: spec.Spec{
				Release: spec.Release{
					Artifacts: true,
				},
			},
			git: &MockGitService{
				HEADMocks: []HEADMock{
					{OutBranch: "main"},
				},
				IsCleanMocks: []IsCleanMock{
					{OutBool: true},
				},
				PullMocks: []PullMock{
					{OutError: nil},
				},
				CreateCommitMocks: []CreateCommitMock{
					{OutHash: "6e8c7d217faab1d88905d4c75b4e7995a42c81d5"},
				},
				CreateTagMocks: []CreateTagMock{
					{OutHash: "a3580a0f64b08ba6085d530c828c40b8aa082c1e"},
				},
				PushMocks: []PushMock{
					{OutError: nil},
				},
				PushTagMocks: []PushTagMock{
					{OutError: nil},
				},
			},
			users: &MockUsersService{
				UserMocks: []UserMock{
					{OutUser: &user, OutResponse: &github.Response{}},
				},
			},
			repo: &MockRepoService{
				GetMocks: []GetMock{
					{OutRepository: &repo, OutResponse: &github.Response{}},
				},
				PermissionMocks: []PermissionMock{
					{OutPermission: github.PermissionAdmin, OutResponse: &github.Response{}},
				},
				CreateReleaseMocks: []CreateReleaseMock{
					{OutRelease: &draftRelease, OutResponse: &github.Response{}},
				},
				UploadReleaseAssetMocks: []UploadReleaseAssetMock{
					{OutReleaseAsset: &asset, OutResponse: &github.Response{}},
				},
				BranchProtectionMocks: []BranchProtectionMock{
					{OutResponse: &github.Response{}},
					{OutResponse: &github.Response{}},
				},
				UpdateReleaseMocks: []UpdateReleaseMock{
					{OutRelease: &release, OutResponse: &github.Response{}},
				},
			},
			changelog: &MockChangelogService{
				GenerateMocks: []GenerateMock{
					{OutContent: "changelog content"},
				},
			},
			semver: &MockSemverCommand{
				RunMocks: []SemverRunMock{
					{OutCode: command.Success},
				},
				SemVerMocks: []SemVerMock{
					{OutSemVer: version},
				},
			},
			build: &MockBuildCommand{
				RunMocks: []BuildRunMock{
					{OutCode: command.Success},
				},
				ArtifactsMocks: []ArtifactsMock{
					{OutArtifacts: artifacts},
				},
			},
			args:             []string{"-major", "-comment", "Release description"},
			expectedExitCode: command.Success,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &Command{
				ui:   cli.NewMockUi(),
				spec: tc.spec,
			}

			c.data.owner = "octocat"
			c.data.repo = "Hello-World"
			c.data.changelogSpec = changelogSpec.Spec{
				General: changelogSpec.General{
					File: "CHANGELOG.md",
				},
			}

			c.services.git = tc.git
			c.services.users = tc.users
			c.services.repo = tc.repo
			c.services.changelog = tc.changelog
			c.commands.semver = tc.semver
			c.commands.build = tc.build

			exitCode := c.run(tc.args)

			assert.Equal(t, tc.expectedExitCode, exitCode)
		})
	}
}
