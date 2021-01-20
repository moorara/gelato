package gen

import (
	"context"
	"flag"
	"time"

	"github.com/mitchellh/cli"

	"github.com/moorara/gelato/internal/command"
	"github.com/moorara/gelato/internal/log"
	"github.com/moorara/gelato/internal/service/generate"
)

const (
	genTimeout  = time.Minute
	genSynopsis = `Generate test helpers`
	genHelp     = `
  Use this command for generating test helpers (mocks, factories, builders, etc.).

  Usage:  gelato gen

  Examples:
    gelato gen
  `
)

type generateService interface {
	Generate(string) error
}

// Command is the cli.Command implementation for gen command.
type Command struct {
	ui       cli.Ui
	services struct {
		generator generateService
	}
	outputs struct{}
}

// NewCommand creates a gen command.
func NewCommand(ui cli.Ui) (*Command, error) {
	return &Command{
		ui: ui,
	}, nil
}

// Synopsis returns a short one-line synopsis of the command.
func (c *Command) Synopsis() string {
	return genSynopsis
}

// Help returns a long help text including usage, description, and list of flags for the command.
func (c *Command) Help() string {
	return genHelp
}

// Run runs the actual command with the given command-line arguments.
// This method is used as a proxy for creating dependencies and the actual command execution is delegated to the run method for testing purposes.
func (c *Command) Run(args []string) int {
	c.services.generator = generate.New(log.Trace)

	return c.run(args)
}

// run in an auxiliary method, so we can test the business logic with mock dependencies.
func (c *Command) run(args []string) int {
	fs := flag.NewFlagSet("gen", flag.ContinueOnError)
	fs.Usage = func() {
		c.ui.Output(c.Help())
	}

	if err := fs.Parse(args); err != nil {
		return command.FlagError
	}

	ctx, cancel := context.WithTimeout(context.Background(), genTimeout)
	defer cancel()

	// ==============================> RUN PREFLIGHT CHECKS <==============================

	checklist := command.PreflightChecklist{}

	info, err := command.RunPreflightChecks(ctx, checklist)
	if err != nil {
		c.ui.Error(err.Error())
		return command.PreflightError
	}

	// ==============================> TODO: <==============================

	if err := c.services.generator.Generate(info.WorkingDirectory); err != nil {
		c.ui.Error(err.Error())
		return command.GenerationError
	}

	// ==============================> DONE <==============================

	return command.Success
}
