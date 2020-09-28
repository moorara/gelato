package release

import (
	"bytes"
	"context"
	"text/template"
	"time"

	"github.com/mitchellh/cli"

	"github.com/moorara/gelato/internal/command"
	"github.com/moorara/gelato/internal/github"
	"github.com/moorara/gelato/internal/spec"
)

const (
	releaseTimeout  = 10 * time.Minute
	releaseSynopsis = `Create a release`
	releaseHelp     = `
  Use this command for creating a new release.
  The initial semantic version is always 0.1.0.

  Currently, the release command only supports GitHub repositories.
  It also assumes the remote repository is named origin.

  Usage:  gelato release [flags]

  Flags:
    -patch:      create a patch version release (default: true)
    -minor:      create a minor version release (default: false)
    -major:      create a major version release (default: false)
    -comment:    add a description for the release
    -artifacts:  build the artifacts and include them in the release (default: {{.Release.Artifacts}})

  Examples:
    gelato release
    gelato release -patch
    gelato release -minor
    gelato release -major
    gelato release -artifacts
    gelato release -comment="Fixing Bugs!"
    gelato release -minor -comment "New Features!"
    gelato release -major -comment "Breaking Changes!"
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
	var buf bytes.Buffer
	t := template.Must(template.New("help").Parse(releaseHelp))
	_ = t.Execute(&buf, c.spec)
	return buf.String()
}

// Run runs the actual command with the given command-line arguments.
func (c *cmd) Run(args []string) int {
	params := struct {
		patch, minor, major bool
		comment             string
	}{}

	fs := c.spec.Release.FlagSet()
	fs.BoolVar(&params.patch, "patch", true, "")
	fs.BoolVar(&params.minor, "minor", false, "")
	fs.BoolVar(&params.major, "major", false, "")
	fs.StringVar(&params.comment, "comment", "", "")
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

	// TODO:

	return command.Success
}
