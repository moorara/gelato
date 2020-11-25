package update

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/mitchellh/cli"
	"github.com/moorara/go-github"

	"github.com/moorara/gelato/internal/command"
)

const (
	updateTimeout  = time.Minute
	updateSynopsis = `Update Gelato`
	updateHelp     = `
  Use this command for updating gelato to the latest release.

  Usage:  gelato update

  Examples:
    gelato update
  `
)

const (
	updateOwner = "moorara"
	updateRepo  = "gelato"
)

type repoService interface {
	LatestRelease(context.Context) (*github.Release, *github.Response, error)
	DownloadReleaseAsset(context.Context, string, string, string) (*github.Response, error)
}

// Command is the cli.Command implementation for update command.
type Command struct {
	ui       cli.Ui
	services struct {
		repo repoService
	}
	outputs struct{}
}

// NewCommand creates an update command.
func NewCommand(ui cli.Ui) (*Command, error) {
	return &Command{
		ui: ui,
	}, nil
}

// Synopsis returns a short one-line synopsis of the command.
func (c *Command) Synopsis() string {
	return updateSynopsis
}

// Help returns a long help text including usage, description, and list of flags for the command.
func (c *Command) Help() string {
	return updateHelp
}

// Run runs the actual command with the given command-line arguments.
func (c *Command) Run(args []string) int {
	// If no access token is provided, we try without it!
	token := os.Getenv("GELATO_GITHUB_TOKEN")

	client := github.NewClient(token)
	repo := client.Repo(updateOwner, updateRepo)

	c.services.repo = repo

	return c.run(args)
}

// run in an auxiliary method, so we can test the business logic with mock dependencies.
func (c *Command) run(args []string) int {
	fs := flag.NewFlagSet("update", flag.ContinueOnError)
	fs.Usage = func() {
		c.ui.Output(c.Help())
	}

	if err := fs.Parse(args); err != nil {
		return command.FlagError
	}

	ctx, cancel := context.WithTimeout(context.Background(), updateTimeout)
	defer cancel()

	// ==============================> RUN PREFLIGHT CHECKS <==============================

	checklist := command.PreflightChecklist{}

	_, err := command.RunPreflightChecks(ctx, checklist)
	if err != nil {
		c.ui.Error(err.Error())
		return command.PreflightError
	}

	// ==============================> GET THE LATEST RELEASE <==============================

	c.ui.Output("Finding the latest release of Gelato ...")

	release, _, err := c.services.repo.LatestRelease(ctx)
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitHubError
	}

	// ==============================> DOWNLOAD THE LATEST BINARY <==============================

	c.ui.Output(fmt.Sprintf("Downloading Gelato %s ...", release.TagName))

	assetName := fmt.Sprintf("gelato-%s-%s", runtime.GOOS, runtime.GOARCH)

	binPath, err := exec.LookPath(os.Args[0])
	if err != nil {
		c.ui.Error(fmt.Sprintf("Cannot find the path for Gelato binary: %s", err))
		return command.OSError
	}

	_, err = c.services.repo.DownloadReleaseAsset(ctx, release.TagName, assetName, binPath)
	if err != nil {
		c.ui.Error(fmt.Sprintf("Failed to download and update Gelato binary: %s", err))
		return command.GitHubError
	}

	c.ui.Info(fmt.Sprintf("🍨 Gelato %s written to %s", release.Name, binPath))

	// ==============================> DONE <==============================

	return command.Success
}
