package app

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/cli"

	"github.com/moorara/gelato/internal/command"
	"github.com/moorara/gelato/internal/spec"
)

const (
	appTimeout  = 2 * time.Minute
	appSynopsis = `Create an application`
	appHelp     = `
  Use this command for creating a new application.
  Currently, the app command can only create Go applications.

  Usage:  gelato app

  Examples:
    gelato app
  `
)

// Command is the cli.Command implementation for app command.
type Command struct {
	ui       cli.Ui
	version  string
	services struct{}
	outputs  struct{}
}

// NewCommand creates an app command.
func NewCommand(ui cli.Ui, version string) (*Command, error) {
	return &Command{
		ui:      ui,
		version: version,
	}, nil
}

// Synopsis returns a short one-line synopsis of the command.
func (c *Command) Synopsis() string {
	return appSynopsis
}

// Help returns a long help text including usage, description, and list of flags for the command.
func (c *Command) Help() string {
	return appHelp
}

// Run runs the actual command with the given command-line arguments.
// This method is used as a proxy for creating dependencies and the actual command execution is delegated to the run method for testing purposes.
func (c *Command) Run(args []string) int {
	return c.run(args)
}

// run in an auxiliary method, so we can test the business logic with mock dependencies.
func (c *Command) run(args []string) int {
	fs := flag.NewFlagSet("app", flag.ContinueOnError)
	fs.Usage = func() {
		c.ui.Output(c.Help())
	}

	if err := fs.Parse(args); err != nil {
		return command.FlagError
	}

	ctx, cancel := context.WithTimeout(context.Background(), appTimeout)
	defer cancel()

	// ==============================> RUN PREFLIGHT CHECKS <==============================

	checklist := command.PreflightChecklist{}

	_, err := command.RunPreflightChecks(ctx, checklist)
	if err != nil {
		c.ui.Error(err.Error())
		return command.PreflightError
	}

	// ==============================> GET INPUTS <==============================

	langOptions := strings.Join([]string{spec.AppLanguageGo}, ", ")
	appLang, err := c.ui.Ask(fmt.Sprintf("Application Language (%s): ", langOptions))
	if err != nil {
		c.ui.Error(fmt.Sprintf("invalid application language: %s", err))
		return command.InputError
	}

	// Only Go applications are supported
	if appLang != spec.AppLanguageGo {
		c.ui.Error(fmt.Sprintf("unsupported application language: %s", appLang))
		return command.UnsupportedError
	}

	typeOptions := strings.Join([]string{spec.AppTypeCLI, spec.AppTypeHTTPService, spec.AppTypeGRPCService}, ", ")
	appType, err := c.ui.Ask(fmt.Sprintf("Application Type (%s): ", typeOptions))
	if err != nil {
		c.ui.Error(fmt.Sprintf("invalid application type: %s", err))
		return command.InputError
	}

	// Only HTTP and gRPC services are supported
	if appType != spec.AppTypeHTTPService && appType != spec.AppTypeGRPCService {
		c.ui.Error(fmt.Sprintf("unsupported application type: %s", appType))
		return command.UnsupportedError
	}

	layoutOptions := strings.Join([]string{spec.AppLayoutVertical, spec.AppLayoutHorizontal}, ", ")
	appLayout, err := c.ui.Ask(fmt.Sprintf("Application Layout (%s): ", layoutOptions))
	if err != nil {
		c.ui.Error(fmt.Sprintf("invalid application layout: %s", err))
		return command.InputError
	}

	// Only vertical and horizontal layouts are supported
	if appLayout != spec.AppLayoutVertical && appLayout != spec.AppLayoutHorizontal {
		c.ui.Error(fmt.Sprintf("unsupported application layout: %s", appLayout))
		return command.UnsupportedError
	}

	modName, err := c.ui.Ask("Go module name: ")
	if err != nil {
		c.ui.Error(fmt.Sprintf("invalid module name: %s", err))
		return command.InputError
	}

	if modName == "" {
		c.ui.Error(fmt.Sprintf("unsupported module name: %s", modName))
		return command.UnsupportedError
	}

	// ==============================> DONE <==============================

	return command.Success
}
