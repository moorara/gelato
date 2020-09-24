package semver

import (
	"context"
	"flag"
	"regexp"
	"time"

	"github.com/mitchellh/cli"
	"golang.org/x/sync/errgroup"

	"github.com/moorara/gelato/internal/command"
	"github.com/moorara/gelato/pkg/semver"
	"github.com/moorara/gelato/pkg/shell"
)

const (
	semverTimeout  = 2 * time.Second
	semverSynopsis = `Prints the current semantic version`
	semverHelp     = `
  Use this command for getting the current semantic version.

  Examples:  gelato semver
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

	// GET GIT INFORMATION

	var gitStatus, gitCommitCount, gitSHA string
	group, groupCtx := errgroup.WithContext(ctx)

	group.Go(func() (err error) {
		_, gitStatus, err = shell.Run(groupCtx, "git", "status", "--porcelain")
		return err
	})

	group.Go(func() (err error) {
		_, gitCommitCount, err = shell.Run(groupCtx, "git", "rev-list", "--count", "HEAD")
		return err
	})

	group.Go(func() (err error) {
		_, gitSHA, err = shell.Run(groupCtx, "git", "rev-parse", "HEAD")
		return err
	})

	if err := group.Wait(); err != nil {
		c.ui.Error(err.Error())
		return command.GitError
	}

	// RESOLVE THE CURRENT SEMANTIC VERSION

	code, gitDescribe, err := shell.Run(ctx, "git", "describe", "--tags", "HEAD")
	if err != nil && code != 128 { // 128 is returned when there is no git tag
		c.ui.Error(err.Error())
		return command.GitError
	}

	releaseRE := regexp.MustCompile(`^v?([0-9]+)\.([0-9]+)\.([0-9]+)$`)
	prereleaseRE := regexp.MustCompile(`^(v?([0-9]+)\.([0-9]+)\.([0-9]+))-([0-9]+)-g([0-9a-f]+)$`)

	if len(gitDescribe) == 0 {
		// No git tag and no previous semantic version -> using the default initial semantic version

		c.outputs.semver = semver.SemVer{
			Major: 0, Minor: 1, Patch: 0,
			Prerelease: []string{gitCommitCount},
		}

		if gitStatus == "" {
			c.outputs.semver.AddPrerelease(gitSHA[:7])
		} else {
			c.outputs.semver.AddPrerelease("dev")
		}
	} else if subs := releaseRE.FindStringSubmatch(gitDescribe); len(subs) == 4 {
		// The tag points to the HEAD commit
		// Example: v0.2.7 --> subs = []string{"v0.2.7", "0", "2", "7"}

		c.outputs.semver, _ = semver.Parse(subs[0])

		if gitStatus != "" {
			c.outputs.semver = c.outputs.semver.Next()
			c.outputs.semver.AddPrerelease("0", "dev")
		}
	} else if subs := prereleaseRE.FindStringSubmatch(gitDescribe); len(subs) == 7 {
		// The tag is the most recent tag reachable from the HEAD commit
		// Example: v0.2.7-10-gabcdeff --> subs = []string{"v0.2.7-10-gabcdeff", "v0.2.7", "0", "2", "7", "10", "abcdeff"}

		c.outputs.semver, _ = semver.Parse(subs[1])
		c.outputs.semver = c.outputs.semver.Next()
		c.outputs.semver.AddPrerelease(subs[5])

		if gitStatus == "" {
			c.outputs.semver.AddPrerelease(subs[6])
		} else {
			c.outputs.semver.AddPrerelease("dev")
		}
	}

	c.ui.Output(c.outputs.semver.String())

	return command.Success
}
