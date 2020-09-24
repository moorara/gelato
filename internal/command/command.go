package command

import (
	"context"
	"errors"
	"fmt"
	"os"

	"golang.org/x/sync/errgroup"

	"github.com/moorara/gelato/pkg/shell"
)

const (
	// Success is the exit code when a command execution is successful.
	Success int = iota

	// FlagError is the exit code when an undefined or invalid flag is provided to a command.
	FlagError

	// PreflightError is the exit code when a preflight check fails.
	PreflightError

	// GoError is the exit code when a go command fails.
	GoError

	// GitError is the exit code when a git command fails.
	GitError
)

type (
	// PreflightChecklist is a list of preflight checks for commands.
	PreflightChecklist struct {
		Go          bool
		Git         bool
		GitHubToken bool
	}

	// PreflightInfo is a list of preflight information for commands.
	PreflightInfo struct {
		WorkingDirectory string
		GoVersion        string
		GitVersion       string
		GitHubToken      string
	}
)

// RunPreflightChecks runs a list of preflight checks to ensure they are fulfilled.
// It returns a list of preflight information.
func RunPreflightChecks(ctx context.Context, checklist PreflightChecklist) (PreflightInfo, error) {
	var workingDirectory, goVersion, gitVersion, githubToken string
	group := new(errgroup.Group)

	// RUN PREFLIGHT CHECKS

	// Get the current working directory and add it to the context
	group.Go(func() (err error) {
		workingDirectory, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("error on getting the current working directory: %s", err)
		}
		return nil
	})

	if checklist.Go {
		group.Go(func() (err error) {
			_, goVersion, err = shell.Run(ctx, "go", "version")
			return err
		})
	}

	if checklist.Git {
		group.Go(func() (err error) {
			_, gitVersion, err = shell.Run(ctx, "git", "version")
			return err
		})
	}

	// Get the GitHub token and add it to the context
	if checklist.GitHubToken {
		group.Go(func() error {
			githubToken = os.Getenv("GELATO_GITHUB_TOKEN")
			if githubToken == "" {
				return errors.New("GELATO_GITHUB_TOKEN environment variable not set")
			}
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return PreflightInfo{}, err
	}

	return PreflightInfo{
		WorkingDirectory: workingDirectory,
		GoVersion:        goVersion,
		GitVersion:       gitVersion,
		GitHubToken:      githubToken,
	}, nil
}
