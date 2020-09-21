package main

import (
	"os"

	"github.com/mitchellh/cli"

	"github.com/moorara/gelato/internal/command/semver"
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
			ErrorColor:  cli.UiColorRed,
			WarnColor:   cli.UiColorYellow,
		},
	}

	c := cli.NewCLI("gelato", version.String())
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"semver": func() (cli.Command, error) {
			return semver.NewCommand(ui)
		},
	}

	code, err := c.Run()
	if err != nil {
		ui.Error(err.Error())
	}

	os.Exit(code)
}
