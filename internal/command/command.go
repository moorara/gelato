package command

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"golang.org/x/sync/errgroup"

	"github.com/moorara/gelato/pkg/shell"
)

const (
	// Success is the exit code when a command execution is successful.
	Success int = iota
	// SpecError is the exit code when reading the spec file fails.
	SpecError
	// FlagError is the exit code when an undefined or invalid flag is provided to a command.
	FlagError
	// PreflightError is the exit code when a preflight check fails.
	PreflightError
	// OSError is the exit code when an OS operation fails.
	OSError
	// GoError is the exit code when a go command fails.
	GoError
	// GitError is the exit code when a git command fails.
	GitError
	// GitHubError is the exit code when a GitHub operation fails.
	GitHubError
	// ChangelogError is the exit code when generating the changelog fails.
	ChangelogError
	// UnsupportedError is the exit code when a capability is not supported.
	UnsupportedError
	// MiscError is the exit code when a miscellaneous operation fails.
	MiscError
)

var (
	releaseRE    = regexp.MustCompile(`^v?([0-9]+)\.([0-9]+)\.([0-9]+)$`)
	prereleaseRE = regexp.MustCompile(`^(v?([0-9]+)\.([0-9]+)\.([0-9]+))-([0-9]+)-g([0-9a-f]+)$`)
)

// PreflightChecklist is a list of common preflight checks for commands.
type PreflightChecklist struct {
	Go bool
}

// PreflightInfo is a list of common preflight information for commands.
type PreflightInfo struct {
	WorkingDirectory string
	GoVersion        string
}

// RunPreflightChecks runs a list of preflight checks to ensure they are fulfilled.
// It returns a list of preflight information.
func RunPreflightChecks(ctx context.Context, checklist PreflightChecklist) (PreflightInfo, error) {
	var workingDirectory, goVersion string

	// RUN PREFLIGHT CHECKS

	group, ctx := errgroup.WithContext(ctx)

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

	if err := group.Wait(); err != nil {
		return PreflightInfo{}, err
	}

	return PreflightInfo{
		WorkingDirectory: workingDirectory,
		GoVersion:        goVersion,
	}, nil
}
