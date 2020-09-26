package build

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/mitchellh/cli"
	"golang.org/x/sync/errgroup"

	"github.com/moorara/gelato/internal/command"
	"github.com/moorara/gelato/internal/spec"
	"github.com/moorara/gelato/pkg/shell"
)

const (
	buildTimeout  = 5 * time.Minute
	buildSynopsis = `Builds artifacts`
	buildHelp     = `
	Use this command for building artifacts.

	Flags:  -cross-compile:

  Examples:  gelato build
  `
)

const (
	versionPath = "./version"
)

// cmd implements the cli.Command interface.
type cmd struct {
	ui      cli.Ui
	spec    spec.Spec
	outputs struct{}
}

// NewCommand creates a build command.
func NewCommand(ui cli.Ui, spec spec.Spec) (cli.Command, error) {
	return &cmd{
		ui:   ui,
		spec: spec,
	}, nil
}

// Synopsis returns a short one-line synopsis of the command.
func (c *cmd) Synopsis() string {
	return buildSynopsis
}

// Help returns a long help text including usage, description, and list of flags for the command.
func (c *cmd) Help() string {
	return buildHelp
}

// Run runs the actual command with the given command-line arguments.
func (c *cmd) Run(args []string) int {
	fs := flag.NewFlagSet("build", flag.ContinueOnError)
	fs.Usage = func() {
		c.ui.Output(c.Help())
	}

	if err := fs.Parse(args); err != nil {
		return command.FlagError
	}

	ctx, cancel := context.WithTimeout(context.Background(), buildTimeout)
	defer cancel()

	// RUN PREFLIGHT CHECKS

	checklist := command.PreflightChecklist{
		Go:  true,
		Git: true,
	}

	_, err := command.RunPreflightChecks(ctx, checklist)
	if err != nil {
		c.ui.Error(err.Error())
		return command.PreflightError
	}

	// GET THE CURRENT SEMANTIC VERSION

	semver, err := command.ResolveSemanticVersion(ctx)
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitError
	}

	// GET GIT & GO INFORMATION

	var gitSHA, gitBranch, goVersion, versionPkg string
	group, groupCtx := errgroup.WithContext(ctx)

	group.Go(func() (err error) {
		_, gitSHA, err = shell.Run(groupCtx, "git", "rev-parse", "HEAD")
		return err
	})

	group.Go(func() (err error) {
		_, gitBranch, err = shell.Run(groupCtx, "git", "rev-parse", "--abbrev-ref", "HEAD")
		return err
	})

	group.Go(func() (err error) {
		_, goVersion, err = shell.Run(groupCtx, "go", "version")
		return err
	})

	group.Go(func() (err error) {
		_, versionPkg, err = shell.Run(groupCtx, "go", "list", versionPath)
		return err
	})

	if err := group.Wait(); err != nil {
		c.ui.Error(err.Error())
		return command.MiscError
	}

	// CONSTRUCT LD FLAGS

	buildTime := time.Now().UTC().Format("2006-01-02 15:04:05 MST")
	buildTool := "Gelato"
	if c.spec.GelatoVersion != "" {
		buildTool += " " + c.spec.GelatoVersion
	}

	versionFlag := fmt.Sprintf("-X '%s.Version=%s'", versionPkg, semver)
	commitFlag := fmt.Sprintf("-X '%s.Commit=%s'", versionPkg, gitSHA[:7])
	branchFlag := fmt.Sprintf("-X '%s.Branch=%s'", versionPkg, gitBranch)
	goVersionFlag := fmt.Sprintf("-X '%s.GoVersion=%s'", versionPkg, goVersion)
	buildToolFlag := fmt.Sprintf("-X '%s.BuildTool=%s'", versionPkg, buildTool)
	buildTimeFlag := fmt.Sprintf("-X '%s.BuildTime=%s'", versionPkg, buildTime)
	ldFlags := fmt.Sprintf("%s %s %s %s %s %s", versionFlag, commitFlag, branchFlag, goVersionFlag, buildToolFlag, buildTimeFlag)

	//

	c.ui.Output(ldFlags)

	return command.Success
}
