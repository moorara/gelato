package release

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/mitchellh/cli"
	"github.com/moorara/go-github"

	"github.com/moorara/gelato/internal/command"
	"github.com/moorara/gelato/internal/git"
	"github.com/moorara/gelato/internal/spec"
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
    -patch:      create a patch version release (default: true)
    -minor:      create a minor version release (default: false)
    -major:      create a major version release (default: false)
    -comment:    add a description for the release
    -artifacts:  build the artifacts and include them in the release (default: {{.Release.Artifacts}})

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

var (
	sshRE   = regexp.MustCompile(`^git@([A-Za-z][0-9A-Za-z-]+[0-9A-Za-z]\.[A-Za-z]{2,}):([A-Za-z][0-9A-Za-z-]+[0-9A-Za-z])/([A-Za-z][0-9A-Za-z-]+[0-9A-Za-z])(.git)?$`)
	httpsRE = regexp.MustCompile(`^https://([A-Za-z][0-9A-Za-z-]+[0-9A-Za-z]\.[A-Za-z]{2,})/([A-Za-z][0-9A-Za-z-]+[0-9A-Za-z])/([A-Za-z][0-9A-Za-z-]+[0-9A-Za-z])(.git)?$`)
)

type (
	gitService interface {
		Remote(string) (string, string, error)
		HEAD() (string, string, error)
		IsClean() (bool, error)
		Pull(ctx context.Context) error
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
)

// Command is the cli.Command implementation for release command.
type Command struct {
	ui       cli.Ui
	spec     spec.Spec
	owner    string
	repo     string
	services struct {
		git   gitService
		users usersService
		repo  repoService
	}
	outputs struct{}
}

// NewCommand creates a release command.
func NewCommand(ui cli.Ui, spec spec.Spec) (*Command, error) {
	g, err := git.New(".")
	if err != nil {
		return nil, err
	}

	// TODO: should we check for other remote names too?
	domain, path, err := g.Remote("origin")
	if err != nil {
		return nil, err
	}

	if domain != "github.com" {
		return nil, fmt.Errorf("unsupported Git platform: %s", domain)
	}

	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		return nil, errors.New("unexpected GitHub repository: cannot parse owner and repo")
	}

	token := os.Getenv("GELATO_GITHUB_TOKEN")
	if token == "" {
		return nil, errors.New("GELATO_GITHUB_TOKEN environment variable not set")
	}

	client := github.NewClient(token)
	repo := client.Repo(parts[0], parts[1])

	c := &Command{
		ui:    ui,
		spec:  spec,
		owner: parts[0],
		repo:  parts[1],
	}

	c.services.git = g
	c.services.users = client.Users
	c.services.repo = repo

	return c, nil
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
func (c *Command) Run(args []string) int {
	params := struct {
		patch, minor, major bool
		comment             string
	}{}

	fs := c.spec.Release.FlagSet()
	fs.BoolVar(&params.patch, "patch", true, "")
	fs.BoolVar(&params.minor, "minor", false, "")
	fs.BoolVar(&params.major, "major", false, "")
	fs.StringVar(&params.comment, "comment", "", "")
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

	// ==============================> CHECK GIT REPO <==============================

	_, gitBranch, err := c.services.git.HEAD()
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitError
	}

	// TODO: find out default branch
	if gitBranch != "main" {
		c.ui.Error("A repository can be released only from the master branch.")
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
		c.ui.Error(fmt.Sprintf("The GitHub token does not have admin permission for releasing %s/%s", c.owner, c.repo))
		return command.GitHubError
	}

	// ==============================> UPDATE DEFAULT BRANCH <==============================

	c.ui.Output("Pulling the latest changes on the master branch ...")

	err = c.services.git.Pull(ctx)
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitError
	}

	// RESOLVE THE RELEASE SEMANTIC VERSION

	// CREATE A DRAFT RELEASE

	// GENERATE CHANGELOG

	// CREATE RELEASE COMMIT & TAG

	// BUILDING AND UPLOADING ARTIFACTS

	// TEMPORARILY ENABLE PUSH TO DEFAULT BRANCH (MASTER)

	// PUSH RELEASE COMMIT & TAG

	// PUBLISH THE DRAFT RELEASE

	return command.Success
}
