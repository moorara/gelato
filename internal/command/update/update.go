package update

import (
	"flag"
	"time"

	"github.com/mitchellh/cli"

	"github.com/moorara/gelato/internal/command"
)

const (
	updateTimeout  = time.Minute
	updateSynopsis = `Updates Gelato`
	updateHelp     = `
  Use this command for updating gelato to the latest release.

  Examples:  gelato update
  `
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

	// RUN PREFLIGHT CHECKS AND CREATE A CONTEXT

	checklist := command.PreflightChecklist{
		GitHubToken: true,
	}

	_, cancel, err := command.CheckAndCreateContext(checklist, updateTimeout)
	if err != nil {
		c.ui.Error(err.Error())
		return command.PreflightError
	}
	defer cancel()

	// CREATE A GITHUB CLIENT

	// GET THE LATEST RELEASE FROM GITHUB

	// DOWNLOAD THE LATEST BINARY FROM GITHUB

	// WRITE THE LATEST BINARY TO DISK

	return command.Success
}
