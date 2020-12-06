package decorate

import (
	"context"
	"flag"
	"time"

	"github.com/mitchellh/cli"

	"github.com/moorara/gelato/internal/command"
	"github.com/moorara/gelato/internal/decorate"
	"github.com/moorara/gelato/internal/log"
	"github.com/moorara/gelato/internal/spec"
)

const (
	decorateTimeout  = time.Minute
	decorateSynopsis = `[WIP] Decorate an application`
	decorateHelp     = `
  Use this command for creating a decorated version of an application.
  Currently, this command can only decorate Go applications.

  A decorator decorates existing code blocks (packages, types, and functions) with extra capabilities.
  Decorators can be used for augmenting an applications with a whole range of different functionalities.
  They can be used for instrumenting an application, enabling observability, improving resiliency and reliability, hardening security, etc.

  Usage:  gelato decorate [flags]

  Flags:
    -trace:  show all logs (default: false)
    -debug:  show debug logs and above (default: false)
    -info:   show info logs and above (default: false)
    -warn:   show warn logs and above (default: false)
    -error:  only show error logs (default: true)

  Examples:
    gelato decorate
  `
)

type decoratorService interface {
	Decorate(log.Level, string) error
}

// Command is the cli.Command implementation for decorate command.
type Command struct {
	ui       cli.Ui
	spec     spec.App
	services struct {
		decorator decoratorService
	}
}

// NewCommand creates a decorate command.
func NewCommand(ui cli.Ui, spec spec.App) (*Command, error) {
	return &Command{
		ui:   ui,
		spec: spec,
	}, nil
}

// Synopsis returns a short one-line synopsis of the command.
func (c *Command) Synopsis() string {
	return decorateSynopsis
}

// Help returns a long help text including usage, description, and list of flags for the command.
func (c *Command) Help() string {
	return decorateHelp
}

// Run runs the actual command with the given command-line arguments.
// This method is used as a proxy for creating dependencies and the actual command execution is delegated to the run method for testing purposes.
func (c *Command) Run(args []string) int {
	decorator := decorate.New()
	c.services.decorator = decorator

	return c.run(args)
}

// run in an auxiliary method, so we can test the business logic with mock dependencies.
func (c *Command) run(args []string) int {
	flags := struct {
		trace, debug, info, warn, err bool
	}{}

	fs := flag.NewFlagSet("decorate", flag.ContinueOnError)
	fs.BoolVar(&flags.trace, "trace", false, "")
	fs.BoolVar(&flags.debug, "debug", false, "")
	fs.BoolVar(&flags.info, "info", false, "")
	fs.BoolVar(&flags.warn, "warn", false, "")
	fs.BoolVar(&flags.err, "error", true, "")
	fs.Usage = func() {
		c.ui.Output(c.Help())
	}

	if err := fs.Parse(args); err != nil {
		return command.FlagError
	}

	ctx, cancel := context.WithTimeout(context.Background(), decorateTimeout)
	defer cancel()

	// ==============================> RUN PREFLIGHT CHECKS <==============================

	checklist := command.PreflightChecklist{}

	info, err := command.RunPreflightChecks(ctx, checklist)
	if err != nil {
		c.ui.Error(err.Error())
		return command.PreflightError
	}

	// ==============================> DECORATE <==============================

	level := log.None
	switch {
	case flags.trace:
		level = log.Trace
	case flags.debug:
		level = log.Debug
	case flags.info:
		level = log.Info
	case flags.warn:
		level = log.Warn
	case flags.err:
		level = log.Error
	}

	err = c.services.decorator.Decorate(level, info.Context.WorkingDirectory)
	if err != nil {
		c.ui.Error(err.Error())
		return command.DecorationError
	}

	// ==============================> DONE <==============================

	return command.Success
}
