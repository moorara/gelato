package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"

	"github.com/moorara/gelato/internal/command"
	"github.com/moorara/gelato/internal/command/app"
	"github.com/moorara/gelato/internal/command/build"
	"github.com/moorara/gelato/internal/command/release"
	"github.com/moorara/gelato/internal/command/semver"
	"github.com/moorara/gelato/internal/command/update"
	"github.com/moorara/gelato/internal/spec"
	"github.com/moorara/gelato/version"
)

func main() {
	ui := &cli.ConcurrentUi{
		Ui: &cli.ColoredUi{
			Ui: &cli.BasicUi{
				Reader:      os.Stdin,
				Writer:      os.Stdout,
				ErrorWriter: os.Stderr,
			},
			OutputColor: cli.UiColorNone,
			InfoColor:   cli.UiColorGreen,
			WarnColor:   cli.UiColorYellow,
			ErrorColor:  cli.UiColorRed,
		},
	}

	// Read the spec from file if any
	spec, err := spec.FromFile()
	if err != nil {
		ui.Error(fmt.Sprintf("Cannot read the spec file: %s", err))
		os.Exit(command.SpecError)
	}

	spec = spec.WithDefaults()
	spec.Gelato.Version = version.Version
	spec.Gelato.Revision = version.Commit

	c := cli.NewCLI("gelato", version.String())
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"app": func() (cli.Command, error) {
			return app.NewCommand(ui, spec)
		},
		"semver": func() (cli.Command, error) {
			return semver.NewCommand(ui)
		},
		"build": func() (cli.Command, error) {
			return build.NewCommand(ui, spec)
		},
		"release": func() (cli.Command, error) {
			return release.NewCommand(ui, spec)
		},
		"update": func() (cli.Command, error) {
			return update.NewCommand(ui)
		},
	}

	code, err := c.Run()
	if err != nil {
		ui.Error(err.Error())
	}

	os.Exit(code)
}
