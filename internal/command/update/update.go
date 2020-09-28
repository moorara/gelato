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

	"github.com/moorara/gelato/internal/command"
	"github.com/moorara/gelato/internal/github"
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

// cmd implements the cli.Command interface.
type cmd struct {
	ui      cli.Ui
	outputs struct{}
}

// NewCommand creates an update command.
func NewCommand(ui cli.Ui) (cli.Command, error) {
	return &cmd{
		ui: ui,
	}, nil
}

// Synopsis returns a short one-line synopsis of the command.
func (c *cmd) Synopsis() string {
	return updateSynopsis
}

// Help returns a long help text including usage, description, and list of flags for the command.
func (c *cmd) Help() string {
	return updateHelp
}

// Run runs the actual command with the given command-line arguments.
func (c *cmd) Run(args []string) int {
	fs := flag.NewFlagSet("update", flag.ContinueOnError)
	fs.Usage = func() {
		c.ui.Output(c.Help())
	}

	if err := fs.Parse(args); err != nil {
		return command.FlagError
	}

	ctx, cancel := context.WithTimeout(context.Background(), updateTimeout)
	defer cancel()

	// RUN PREFLIGHT CHECKS

	checklist := command.PreflightChecklist{
		GitHubToken: true,
	}

	info, err := command.RunPreflightChecks(ctx, checklist)
	if err != nil {
		c.ui.Error(err.Error())
		return command.PreflightError
	}

	// CREATE A GITHUB CLIENT

	github, err := github.New(info.GitHubToken)
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitHubError
	}

	// GET THE LATEST RELEASE FROM GITHUB

	c.ui.Output("⬇ Finding the latest release of Gelato ...")

	release, err := github.GetLatestRelease(ctx, updateOwner, updateRepo)
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitHubError
	}

	// DOWNLOAD THE LATEST BINARY FROM GITHUB AND WRITE IT TO DISK

	c.ui.Output(fmt.Sprintf("⬇ Downloading Gelato %s ...", release.TagName))

	var downloadURL string
	assetName := fmt.Sprintf("gelato-%s-%s", runtime.GOOS, runtime.GOARCH)
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			downloadURL = asset.DownloadURL
		}
	}

	if downloadURL == "" {
		c.ui.Error(fmt.Sprintf("Cannot find download URL for %s", assetName))
		return command.GitHubError
	}

	binPath, err := exec.LookPath(os.Args[0])
	if err != nil {
		c.ui.Error(fmt.Sprintf("Cannot find binary path for Gelato: %s", err))
		return command.OSError
	}

	err = github.DownloadReleaseAsset(ctx, downloadURL, binPath)
	if err != nil {
		c.ui.Error(fmt.Sprintf("Failed to download and update Gelato binary: %s", err))
		return command.GitHubError
	}

	c.ui.Info(fmt.Sprintf("🍨 Gelato %s written to %s", release.Name, binPath))

	return command.Success
}
