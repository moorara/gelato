package release

import (
	"context"
	"flag"
	"time"

	"github.com/mitchellh/cli"

	"github.com/moorara/gelato/internal/command"
	"github.com/moorara/gelato/internal/github"
	"github.com/moorara/gelato/internal/spec"
)

const (
	releaseTimeout  = 10 * time.Minute
	releaseSynopsis = `Creates a new release`
	releaseHelp     = `
	Use this command for creating a new release.
	The initial semantic version is always 0.1.0.

	Flags:

  Examples:  gelato release
  `
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
	return releaseHelp
}

// Run runs the actual command with the given command-line arguments.
func (c *cmd) Run(args []string) int {
	fs := flag.NewFlagSet("release", flag.ContinueOnError)
	fs.Usage = func() {
		c.ui.Output(c.Help())
	}

	if err := fs.Parse(args); err != nil {
		return command.FlagError
	}

	ctx, cancel := context.WithTimeout(context.Background(), releaseTimeout)
	defer cancel()

	// RUN PREFLIGHT CHECKS

	checklist := command.PreflightChecklist{
		Git:         true,
		GitHubToken: true,
	}

	info, err := command.RunPreflightChecks(ctx, checklist)
	if err != nil {
		c.ui.Error(err.Error())
		return command.PreflightError
	}

	// CREATE A GITHUB CLIENT

	_, err = github.New(info.GitHubToken, github.ScopeRepo)
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitHubError
	}

	return command.Success
}
