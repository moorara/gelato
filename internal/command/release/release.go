package release

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/mitchellh/cli"

	"github.com/moorara/gelato/internal/command"
	"github.com/moorara/gelato/internal/github"
	"github.com/moorara/gelato/internal/spec"
	"github.com/moorara/gelato/pkg/shell"
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

// cmd implements the cli.Command interface.
type cmd struct {
	ui      cli.Ui
	spec    spec.Spec
	outputs struct{}
}

// NewCommand creates a release command.
func NewCommand(ui cli.Ui, spec spec.Spec) (cli.Command, error) {
	return &cmd{
		ui:   ui,
		spec: spec,
	}, nil
}

// Synopsis returns a short one-line synopsis of the command.
func (c *cmd) Synopsis() string {
	return releaseSynopsis
}

// Help returns a long help text including usage, description, and list of flags for the command.
func (c *cmd) Help() string {
	var buf bytes.Buffer
	t := template.Must(template.New("help").Parse(releaseHelp))
	_ = t.Execute(&buf, c.spec)
	return buf.String()
}

// Run runs the actual command with the given command-line arguments.
func (c *cmd) Run(args []string) int {
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

	// RUN PREFLIGHT CHECKS

	c.ui.Output("Running preflight checks ...")

	checklist := command.PreflightChecklist{
		Git:         true,
		GitHubToken: true,
	}

	info, err := command.RunPreflightChecks(ctx, checklist)
	if err != nil {
		c.ui.Error(err.Error())
		return command.PreflightError
	}

	// CHECK THE GIT REPO

	var domain, owner, repo string

	_, gitRemoteURL, err := shell.Run(ctx, "git", "remote", "get-url", "--push", "origin")
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitError
	}

	if subs := sshRE.FindStringSubmatch(gitRemoteURL); len(subs) == 4 || len(subs) == 5 {
		// Git remote url is using SSH protocol
		// Example: git@github.com:moorara/cherry.git --> subs = []string{"git@github.com:moorara/cherry.git", "github.com", "moorara", "cherry", ".git"}
		domain, owner, repo = subs[1], subs[2], subs[3]
	} else if subs := httpsRE.FindStringSubmatch(gitRemoteURL); len(subs) == 4 || len(subs) == 5 {
		// Git remote url is using HTTPS protocol
		// Example: https://github.com/moorara/cherry.git --> subs = []string{"https://github.com/moorara/cherry.git", "github.com", "moorara", "cherry", ".git"}
		domain, owner, repo = subs[1], subs[2], subs[3]
	} else {
		c.ui.Error(fmt.Sprintf("Invalid git remote url: %s", gitRemoteURL))
		return command.GitError
	}

	if strings.ToLower(domain) != "github.com" {
		c.ui.Error(fmt.Sprintf("Unsupported remote repository: %s", domain))
		return command.UnsupportedError
	}

	_, gitBranch, err := shell.Run(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitError
	}

	if gitBranch != "master" {
		c.ui.Error("A repository can be released only from the master branch.")
		return command.GitError
	}

	_, gitStatus, err := shell.Run(ctx, "git", "status", "--porcelain")
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitError
	}

	if gitStatus != "" {
		c.ui.Error("Working directory is not clean and has uncommitted changes.")
		return command.GitError
	}

	// CREATE A GITHUB CLIENT AND CHECK PERMISSION

	c.ui.Output("Checking GitHub permission ...")

	gh, err := github.New(info.GitHubToken)
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitHubError
	}

	user, err := gh.GetUser(ctx)
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitHubError
	}

	perm, err := gh.GetRepoPermission(ctx, owner, repo, user.Login)
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitHubError
	}

	if perm != github.PermissionAdmin {
		c.ui.Error(fmt.Sprintf("The GitHub token does not have admin permission for releasing %s/%s/%s", domain, owner, repo))
		return command.GitHubError
	}

	// UPDATE MASTER BRANCH

	c.ui.Output("Pulling the latest changes on the master branch ...")

	_, _, err = shell.Run(ctx, "git", "pull")
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitError
	}

	// RESOLVE THE RELEASE SEMANTIC VERSION

	// CREATE A DRAFT RELEASE

	// GENERATE CHANGELOG

	// CREATE RELEASE COMMIT & TAG

	// BUILDING AND UPLOADING ARTIFACTS

	// TEMPORARILY ENABLE PUSH TO MASTER

	// PUSH RELEASE COMMIT & TAG

	// PUBLISH THE DRAFT RELEASE

	return command.Success
}
