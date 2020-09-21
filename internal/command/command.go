package command

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

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

// ContextKey is the type for the keys added to a context.Context.
type ContextKey string

const (
	// WorkingDirectoryKey is the context key for the current working directory.
	WorkingDirectoryKey = ContextKey("WorkingDirectory")

	// GitHubTokenKey is the context key for the GitHub token.
	GitHubTokenKey = ContextKey("GitHubToken")
)

// PreflightChecklist is a list of preflight checks for commands.
type PreflightChecklist struct {
	Go          bool
	Git         bool
	GitHubToken bool
}

// CheckAndCreateContext does two things:
//   1. Runs a list of preflight checks to ensure they are fulfilled.
//   2. Creates a new context with a timeout and contextual information.
//
// The current working directory is always available on the context via workingDirectoryKey.
// If preflightChecklist.GitHubToken set to true, GitHub token will be available on the context via gitHubTokenKey.
func CheckAndCreateContext(checklist PreflightChecklist, timeout time.Duration) (context.Context, context.CancelFunc, error) {
	var workingDirectory, githubToken string
	group := new(errgroup.Group)

	// RUN PREFLIGHT CHECKS

	// Get the current working directory and add it to the context
	group.Go(func() error {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error on getting the current working directory: %s", err)
		}
		workingDirectory = wd
		return nil
	})

	if checklist.Go {
		group.Go(func() error {
			_, _, err := shell.Run(context.Background(), "go", "version")
			return err
		})
	}

	if checklist.Git {
		group.Go(func() error {
			_, _, err := shell.Run(context.Background(), "git", "version")
			return err
		})
	}

	// Get the GitHub token and add it to the context
	if checklist.GitHubToken {
		group.Go(func() error {
			val := os.Getenv("GELATO_GITHUB_TOKEN")
			if val == "" {
				return errors.New("GELATO_GITHUB_TOKEN environment variable not set")
			}
			githubToken = val
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return nil, nil, err
	}

	// CREATE CONTEXT

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)

	if workingDirectory != "" {
		ctx = context.WithValue(ctx, WorkingDirectoryKey, workingDirectory)
	}

	if githubToken != "" {
		ctx = context.WithValue(ctx, GitHubTokenKey, githubToken)
	}

	return ctx, cancel, nil
}
