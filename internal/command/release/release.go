package release

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/mitchellh/cli"
	"github.com/moorara/go-github"

	changelog "github.com/moorara/changelog/generate"
	changelogSpec "github.com/moorara/changelog/spec"

	"github.com/moorara/gelato/internal/command"
	buildcmd "github.com/moorara/gelato/internal/command/build"
	semvercmd "github.com/moorara/gelato/internal/command/semver"
	"github.com/moorara/gelato/internal/service/git"
	"github.com/moorara/gelato/internal/spec"
	"github.com/moorara/gelato/pkg/semver"
)

const (
	releaseTimeout  = 10 * time.Minute
	releaseSynopsis = `Create a release`
	releaseHelp     = `
  Use this command for creating a new release.
  The initial semantic version is always 0.1.0.

  Currently, the release command only supports GitHub repositories.
  It also assumes the remote repository is named origin.

  Usage:  gelato release [flags]

  Flags:
    -patch        create a patch version release (default: true)
    -minor        create a minor version release (default: false)
    -major        create a major version release (default: false)
    -comment      add a description for the release
    -artifacts    build the artifacts and include them in the release (default: {{.Release.Artifacts}})

  Examples:
    gelato release
    gelato release -patch
    gelato release -minor
    gelato release -major
    gelato release -artifacts
    gelato release -comment="Fixing Bugs!"
    gelato release -minor -comment "New Features!"
    gelato release -major -comment "Breaking Changes!"
  `
)

const (
	remoteName = "origin"
)

var (
	h2Regex = regexp.MustCompile(`##[^\n]*\n`)
)

type (
	gitService interface {
		Remote(string) (string, string, error)
		HEAD() (string, string, error)
		IsClean() (bool, error)
		CreateCommit(string, ...string) (string, error)
		CreateTag(string, string, string) (string, error)
		Pull(context.Context) error
		Push(context.Context, string) error
		PushTag(context.Context, string, string) error
	}

	usersService interface {
		User(context.Context) (*github.User, *github.Response, error)
	}

	repoService interface {
		Get(context.Context) (*github.Repository, *github.Response, error)
		Permission(context.Context, string) (github.Permission, *github.Response, error)
		BranchProtection(context.Context, string, bool) (*github.Response, error)
		CreateRelease(context.Context, github.ReleaseParams) (*github.Release, *github.Response, error)
		UpdateRelease(context.Context, int, github.ReleaseParams) (*github.Release, *github.Response, error)
		UploadReleaseAsset(context.Context, int, string, string) (*github.ReleaseAsset, *github.Response, error)
	}

	changelogService interface {
		Generate(context.Context, changelogSpec.Spec) (string, error)
	}

	semverCommand interface {
		Run([]string) int
		SemVer() semver.SemVer
	}

	buildCommand interface {
		Run([]string) int
		Artifacts() []buildcmd.Artifact
	}
)

// Command is the cli.Command implementation for release command.
type Command struct {
	ui   cli.Ui
	spec spec.Spec
	data struct {
		owner         string
		repo          string
		changelogSpec changelogSpec.Spec
	}
	services struct {
		git       gitService
		users     usersService
		repo      repoService
		changelog changelogService
	}
	commands struct {
		semver semverCommand
		build  buildCommand
	}
	outputs struct{}
}

// NewCommand creates a release command.
func NewCommand(ui cli.Ui, spec spec.Spec) (*Command, error) {
	return &Command{
		ui:   ui,
		spec: spec,
	}, nil
}

// Synopsis returns a short one-line synopsis of the command.
func (c *Command) Synopsis() string {
	return releaseSynopsis
}

// Help returns a long help text including usage, description, and list of flags for the command.
func (c *Command) Help() string {
	var buf bytes.Buffer
	t := template.Must(template.New("help").Parse(releaseHelp))
	_ = t.Execute(&buf, c.spec)
	return buf.String()
}

// Run runs the actual command with the given command-line arguments.
// This method is used as a proxy for creating dependencies and the actual command execution is delegated to the run method for testing purposes.
func (c *Command) Run(args []string) int {
	git, err := git.Open(".")
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitError
	}

	// TODO: should we check for remote names other than origin?
	domain, path, err := git.Remote(remoteName)
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitError
	}

	if domain != "github.com" {
		c.ui.Error(fmt.Sprintf("unsupported Git platform: %s", domain))
		return command.GitHubError
	}

	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		c.ui.Error("unexpected GitHub repository: cannot parse owner and repo")
		return command.GitHubError
	}
	ownerName, repoName := parts[0], parts[1]

	token := os.Getenv("GELATO_GITHUB_TOKEN")
	if token == "" {
		c.ui.Error("GELATO_GITHUB_TOKEN environment variable not set")
		return command.GitHubError
	}

	client := github.NewClient(token)
	repo := client.Repo(ownerName, repoName)

	cs, err := changelogSpec.Default().FromFile()
	if err != nil {
		c.ui.Error(err.Error())
		return command.ChangelogError
	}

	cs = cs.WithRepo(domain, path)
	cs.Repo.AccessToken = token
	chlogLogger := newLogger(c.ui)

	changelog, err := changelog.New(cs, chlogLogger)
	if err != nil {
		c.ui.Error(err.Error())
		return command.ChangelogError
	}

	semver, _ := semvercmd.NewCommand(&cli.MockUi{})
	build, _ := buildcmd.NewCommand(c.ui, c.spec)

	c.data.owner = ownerName
	c.data.repo = repoName
	c.data.changelogSpec = cs

	c.services.git = git
	c.services.users = client.Users
	c.services.repo = repo
	c.services.changelog = changelog
	c.commands.semver = semver
	c.commands.build = build

	return c.run(args)
}

// run in an auxiliary method, so we can test the business logic with mock dependencies.
func (c *Command) run(args []string) int {
	flags := struct {
		patch, minor, major bool
		comment             string
	}{}

	fs := c.spec.Release.FlagSet()
	fs.BoolVar(&flags.patch, "patch", true, "")
	fs.BoolVar(&flags.minor, "minor", false, "")
	fs.BoolVar(&flags.major, "major", false, "")
	fs.StringVar(&flags.comment, "comment", "", "")
	fs.Usage = func() {
		c.ui.Output(c.Help())
	}

	if err := fs.Parse(args); err != nil {
		return command.FlagError
	}

	ctx, cancel := context.WithTimeout(context.Background(), releaseTimeout)
	defer cancel()

	// ==============================> RUN PREFLIGHT CHECKS <==============================

	c.ui.Output("Running preflight checks ...")

	checklist := command.PreflightChecklist{}

	_, err := command.RunPreflightChecks(ctx, checklist)
	if err != nil {
		c.ui.Error(err.Error())
		return command.PreflightError
	}

	// ==============================> FIND OUT DEFAULT BRANCH <==============================

	repo, _, err := c.services.repo.Get(ctx)
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitHubError
	}

	// ==============================> CHECK GIT REPO <==============================

	_, gitBranch, err := c.services.git.HEAD()
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitError
	}

	if gitBranch != repo.DefaultBranch {
		c.ui.Error("The repository can only be released from the default branch.")
		c.ui.Error(fmt.Sprintf("  Git Branch:      %s", gitBranch))
		c.ui.Error(fmt.Sprintf("  Default Branch:  %s", repo.DefaultBranch))
		return command.GitError
	}

	isClean, err := c.services.git.IsClean()
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitError
	}

	if !isClean {
		c.ui.Error("Working directory is not clean and has uncommitted changes.")
		return command.GitError
	}

	// ==============================> CHECK GITHUB PERMISSION <==============================

	c.ui.Output("Checking GitHub permission ...")

	user, _, err := c.services.users.User(ctx)
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitHubError
	}

	perm, _, err := c.services.repo.Permission(ctx, user.Login)
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitHubError
	}

	if perm != github.PermissionAdmin {
		c.ui.Error(fmt.Sprintf("The GitHub token does not have admin permission for releasing %s/%s", c.data.owner, c.data.repo))
		return command.GitHubError
	}

	// ==============================> UPDATE DEFAULT BRANCH <==============================

	c.ui.Info(fmt.Sprintf("Pulling the latest changes on the %s branch ...", gitBranch))

	err = c.services.git.Pull(ctx)
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitError
	}

	// ==============================> RESOLVE SEMANTIC VERSION <==============================

	// Run semver command
	code := c.commands.semver.Run(nil)
	if code != command.Success {
		return code
	}

	var version semver.SemVer

	switch {
	case flags.major:
		version = c.commands.semver.SemVer().ReleaseMajor()
	case flags.minor:
		version = c.commands.semver.SemVer().ReleaseMinor()
	case flags.patch:
		fallthrough
	default:
		version = c.commands.semver.SemVer().ReleasePatch()
	}

	tagName := "v" + version.String()

	// ==============================> CREATE A DRAFT RELEASE <==============================

	c.ui.Info(fmt.Sprintf("Creating the draft release %s ...", version))

	params := github.ReleaseParams{
		Name:       version.String(),
		TagName:    tagName,
		Target:     gitBranch,
		Draft:      true,
		Prerelease: false,
	}

	release, _, err := c.services.repo.CreateRelease(ctx, params)
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitHubError
	}

	// ==============================> GENERATE CHANGELOG <==============================

	c.ui.Info(fmt.Sprintf("Creating/Updating the changelog (%s) ...", c.data.changelogSpec.General.File))

	c.data.changelogSpec.Tags.Future = tagName

	changelog, err := c.services.changelog.Generate(ctx, c.data.changelogSpec)
	if err != nil {
		c.ui.Error(err.Error())
		return command.ChangelogError
	}

	// Remove the H2 title
	changelog = h2Regex.ReplaceAllString(changelog, "")
	changelog = strings.TrimLeft(changelog, "\n")

	// ==============================> CREATE RELEASE COMMIT & TAG <==============================

	c.ui.Info(fmt.Sprintf("Creating the release commit and tag %s ...", version))

	message := fmt.Sprintf("Release %s", version)

	commit, err := c.services.git.CreateCommit(message, c.data.changelogSpec.General.File)
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitError
	}

	_, err = c.services.git.CreateTag(commit, tagName, message)
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitError
	}

	// ==============================> BUILD AND UPLOAD ARTIFACTS  <==============================

	if c.spec.Release.Artifacts {
		c.ui.Output("Building artifacts ...")

		// Run build command
		code = c.commands.build.Run(nil)
		if code != command.Success {
			return code
		}

		c.ui.Info(fmt.Sprintf("Uploading artifacts to release %s ...", release.Name))

		group, groupCtx := errgroup.WithContext(ctx)

		for _, artifact := range c.commands.build.Artifacts() {
			artifact := artifact // https://golang.org/doc/faq#closures_and_goroutines
			group.Go(func() error {
				_, _, err := c.services.repo.UploadReleaseAsset(groupCtx, release.ID, artifact.Path, artifact.Label)
				return err
			})
		}

		if err := group.Wait(); err != nil {
			c.ui.Error(err.Error())
			return command.GitHubError
		}
	}

	// ==============================> ENABLE PUSH TO DEFAULT BRANCH <==============================

	c.ui.Warn(fmt.Sprintf("Temporarily enabling push to %s branch ...", gitBranch))

	_, err = c.services.repo.BranchProtection(ctx, gitBranch, false)
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitHubError
	}

	// Make sure we re-enable the branch protection
	defer func() {
		c.ui.Warn(fmt.Sprintf("🔒 Re-disabling push to %s branch ...", gitBranch))
		_, err := c.services.repo.BranchProtection(ctx, gitBranch, true)
		if err != nil {
			c.ui.Error(err.Error())
			os.Exit(command.GitHubError)
		}
	}()

	// ==============================> PUSH RELEASE COMMIT & TAG <==============================

	c.ui.Info(fmt.Sprintf("Pushing release commit %s ...", version))

	err = c.services.git.Push(ctx, remoteName)
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitError
	}

	c.ui.Info(fmt.Sprintf("Pushing release tag %s ...", tagName))

	err = c.services.git.PushTag(ctx, remoteName, tagName)
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitError
	}

	// ==============================> PUBLISH THE RELEASE <==============================

	c.ui.Info(fmt.Sprintf("Publishing release %s ...", release.Name))

	if flags.comment != "" {
		changelog = fmt.Sprintf("%s\n\n%s", flags.comment, changelog)
	}

	params = github.ReleaseParams{
		Name:       release.Name,
		TagName:    release.TagName,
		Target:     release.Target,
		Draft:      false,
		Prerelease: false,
		Body:       changelog,
	}

	release, _, err = c.services.repo.UpdateRelease(ctx, release.ID, params)
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitHubError
	}

	// ==============================> DONE <==============================

	return command.Success
}
