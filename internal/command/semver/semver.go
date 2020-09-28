package semver

import (
	"context"
	"flag"
	"time"

	"github.com/mitchellh/cli"

	"github.com/moorara/gelato/internal/command"
	"github.com/moorara/gelato/pkg/semver"
)

const (
	semverTimeout  = 2 * time.Second
	semverSynopsis = `Print the current semantic version`
	semverHelp     = `
  Use this command for getting the current semantic version.

  Usage:  gelato semver

  Examples:
    gelato semver
  `
)

// cmd implements the cli.Command interface.
type cmd struct {
	ui      cli.Ui
	outputs struct {
		semver semver.SemVer
	}
}

// NewCommand creates a semver command.
func NewCommand(ui cli.Ui) (cli.Command, error) {
	return &cmd{
		ui: ui,
	}, nil
}

// Synopsis returns a short one-line synopsis of the command.
func (c *cmd) Synopsis() string {
	return semverSynopsis
}

// Help returns a long help text including usage, description, and list of flags for the command.
func (c *cmd) Help() string {
	return semverHelp
}

// Run runs the actual command with the given command-line arguments.
func (c *cmd) Run(args []string) int {
	fs := flag.NewFlagSet("semver", flag.ContinueOnError)
	fs.Usage = func() {
		c.ui.Output(c.Help())
	}

	if err := fs.Parse(args); err != nil {
		return command.FlagError
	}

	ctx, cancel := context.WithTimeout(context.Background(), semverTimeout)
	defer cancel()

	// RUN PREFLIGHT CHECKS

	checklist := command.PreflightChecklist{
		Git: true,
	}

	_, err := command.RunPreflightChecks(ctx, checklist)
	if err != nil {
		c.ui.Error(err.Error())
		return command.PreflightError
	}

	// RESOLVE THE CURRENT SEMANTIC VERSION

	c.outputs.semver, err = command.ResolveSemanticVersion(ctx)
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitError
	}

	c.ui.Output(c.outputs.semver.String())

	return command.Success
}
