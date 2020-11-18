package command

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"golang.org/x/sync/errgroup"

	"github.com/moorara/gelato/pkg/semver"
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
	// VersionPkgError is the exit code when the version package is missing or invalid.
	VersionPkgError
	// UnsupportedError is the exit code when a capability is not supported.
	UnsupportedError
	// MiscError is the exit code when a miscellaneous operation fails.
	MiscError
)

var (
	releaseRE    = regexp.MustCompile(`^v?([0-9]+)\.([0-9]+)\.([0-9]+)$`)
	prereleaseRE = regexp.MustCompile(`^(v?([0-9]+)\.([0-9]+)\.([0-9]+))-([0-9]+)-g([0-9a-f]+)$`)
)

type (
	// PreflightChecklist is a list of preflight checks for commands.
	PreflightChecklist struct {
		Go  bool
		Git bool
	}

	// PreflightInfo is a list of preflight information for commands.
	PreflightInfo struct {
		WorkingDirectory string
		GoVersion        string
		GitVersion       string
	}
)

// RunPreflightChecks runs a list of preflight checks to ensure they are fulfilled.
// It returns a list of preflight information.
func RunPreflightChecks(ctx context.Context, checklist PreflightChecklist) (PreflightInfo, error) {
	var workingDirectory, goVersion, gitVersion string

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

	if checklist.Git {
		group.Go(func() (err error) {
			_, gitVersion, err = shell.Run(ctx, "git", "version")
			return err
		})
	}

	if err := group.Wait(); err != nil {
		return PreflightInfo{}, err
	}

	return PreflightInfo{
		WorkingDirectory: workingDirectory,
		GoVersion:        goVersion,
		GitVersion:       gitVersion,
	}, nil
}

// ResolveSemanticVersion returns the current semantic version.
func ResolveSemanticVersion(ctx context.Context) (semver.SemVer, error) {
	var sv semver.SemVer

	// GET GIT INFORMATION

	_, gitStatus, err := shell.Run(ctx, "git", "status", "--porcelain")
	if err != nil {
		return sv, err
	}

	_, gitSHA, err := shell.Run(ctx, "git", "rev-parse", "HEAD")
	if err != nil {
		return sv, err
	}

	_, gitCommitCount, err := shell.Run(ctx, "git", "rev-list", "--count", "HEAD")
	if err != nil {
		return sv, err
	}

	code, gitDescribe, err := shell.Run(ctx, "git", "describe", "--tags", "HEAD")
	if err != nil && code != 128 { // 128 is returned when there is no git tag
		return sv, err
	}

	// RESOLVE THE CURRENT SEMANTIC VERSION

	if len(gitDescribe) == 0 {
		// No git tag and no previous semantic version -> using the default initial semantic version

		sv = semver.SemVer{
			Major: 0, Minor: 1, Patch: 0,
			Prerelease: []string{gitCommitCount},
		}

		if gitStatus == "" {
			sv.AddPrerelease(gitSHA[:7])
		} else {
			sv.AddPrerelease("dev")
		}
	} else if subs := releaseRE.FindStringSubmatch(gitDescribe); len(subs) == 4 {
		// The tag points to the HEAD commit
		// Example: v0.2.7 --> subs = []string{"v0.2.7", "0", "2", "7"}

		sv, _ = semver.Parse(subs[0])

		if gitStatus != "" {
			sv = sv.Next()
			sv.AddPrerelease("0", "dev")
		}
	} else if subs := prereleaseRE.FindStringSubmatch(gitDescribe); len(subs) == 7 {
		// The tag is the most recent tag reachable from the HEAD commit
		// Example: v0.2.7-10-gabcdeff --> subs = []string{"v0.2.7-10-gabcdeff", "v0.2.7", "0", "2", "7", "10", "abcdeff"}

		sv, _ = semver.Parse(subs[1])
		sv = sv.Next()
		sv.AddPrerelease(subs[5])

		if gitStatus == "" {
			sv.AddPrerelease(subs[6])
		} else {
			sv.AddPrerelease("dev")
		}
	}

	return sv, nil
}
