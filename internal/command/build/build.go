package build

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/mitchellh/cli"
	"golang.org/x/sync/errgroup"

	"github.com/moorara/gelato/internal/command"
	semvercmd "github.com/moorara/gelato/internal/command/semver"
	"github.com/moorara/gelato/internal/log"
	"github.com/moorara/gelato/internal/service/decorate"
	"github.com/moorara/gelato/internal/service/git"
	"github.com/moorara/gelato/internal/spec"
	"github.com/moorara/gelato/pkg/semver"
	"github.com/moorara/gelato/pkg/shell"
)

const (
	buildTimeout  = 5 * time.Minute
	buildSynopsis = `Build artifacts`
	buildHelp     = `
  Use this command for building artifacts.
  Currently, the build command only builds binaries for Go applications.

  By convention, It assumes the current directory is a main package if it contains a main.go file.
  It also assumes every directory inside cmd is a main package for a binary with the same name as the directory name.

  Decoration is an experimental feature to decorate the applications with horizontal layout.
  It wraps the controller, gateway, handler, and repository packages with a set of decorators.
  Decorators can be used for augmenting an application with observability, error reccovery, etc.

  Usage:  gelato build [flags]

  Flags:
    -cross-compile    build the binary for all platforms (default: {{.Build.CrossCompile}})
    -decorate         [EXPERIMENTAL] decorate the application before building

  Examples:
    gelato build
    gelato build -cross-compile
    gelato build -decorate
    gelato build -cross-compile -decorate
  `
)

const (
	decoratedDir = ".build"
	cmdDir       = "cmd"
	binPath      = "./bin/"
	versionPath  = "./version"
	timeFormat   = "2006-01-02 15:04:05 MST"
)

var (
	goVersionRE = regexp.MustCompile(`\d+\.\d+(\.\d+)?`)
)

type (
	gitService interface {
		HEAD() (string, string, error)
	}

	decorateService interface {
		Decorate(string) error
	}

	semverCommand interface {
		Run([]string) int
		SemVer() semver.SemVer
	}
)

// Artifact is a build artifacts.
type Artifact struct {
	Path  string
	Label string
}

// Command is the cli.Command implementation for build command.
type Command struct {
	sync.Mutex
	ui       cli.Ui
	spec     spec.Spec
	services struct {
		git       gitService
		decorator decorateService
	}
	funcs struct {
		goList  shell.RunnerFunc
		goBuild shell.RunnerWithFunc
	}
	commands struct {
		semver semverCommand
	}
	outputs struct {
		artifacts []Artifact
	}
}

// NewCommand creates a build command.
func NewCommand(ui cli.Ui, spec spec.Spec) (*Command, error) {
	return &Command{
		ui:   ui,
		spec: spec,
	}, nil
}

// Synopsis returns a short one-line synopsis of the command.
func (c *Command) Synopsis() string {
	return buildSynopsis
}

// Help returns a long help text including usage, description, and list of flags for the command.
func (c *Command) Help() string {
	var buf bytes.Buffer
	t := template.Must(template.New("help").Parse(buildHelp))
	_ = t.Execute(&buf, c.spec)
	return buf.String()
}

// Run runs the actual command with the given command-line arguments.
// This method is used as a proxy for creating dependencies and the actual command execution is delegated to the run method for testing purposes.
func (c *Command) Run(args []string) int {
	git, err := git.Open(".")
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitError
	}

	semver, _ := semvercmd.NewCommand(&cli.MockUi{})

	c.services.git = git
	c.services.decorator = decorate.New(decoratedDir, log.Info)
	c.funcs.goList = shell.Runner("go", "list", versionPath)
	c.funcs.goBuild = shell.RunnerWith("go", "build")
	c.commands.semver = semver

	return c.run(args)
}

// run in an auxiliary method, so we can test the business logic with mock dependencies.
func (c *Command) run(args []string) int {
	fs := c.spec.Build.FlagSet()
	fs.Usage = func() {
		c.ui.Output(c.Help())
	}

	if err := fs.Parse(args); err != nil {
		return command.FlagError
	}

	ctx, cancel := context.WithTimeout(context.Background(), buildTimeout)
	defer cancel()

	// ==============================> RUN PREFLIGHT CHECKS <==============================

	checklist := command.PreflightChecklist{
		Go:  true,
		Git: true,
	}

	info, err := command.RunPreflightChecks(ctx, checklist)
	if err != nil {
		c.ui.Error(err.Error())
		return command.PreflightError
	}

	// ==============================> GET GIT & GO INFORMATION <==============================

	gitSHA, gitBranch, err := c.services.git.HEAD()
	if err != nil {
		c.ui.Error(err.Error())
		return command.GitError
	}

	_, versionPkg, err := c.funcs.goList(ctx)
	if err != nil {
		c.ui.Warn(err.Error())
	}

	// ==============================> GET THE SEMANTIC VERSION <==============================

	// Run semver command
	code := c.commands.semver.Run(nil)
	if code != command.Success {
		return code
	}

	semver := c.commands.semver.SemVer()

	// ==============================> CONSTRUCT LD FLAGS <==============================

	var ldFlags string

	// Construct the LD flags only if the version package exist
	if versionPkg != "" {
		goVersion := goVersionRE.FindString(info.Go.Version)
		buildTime := time.Now().UTC().Format(timeFormat)
		buildTool := "Gelato"

		if c.spec.Gelato.Version != "" {
			buildTool += " " + c.spec.Gelato.Version
		}

		ldFlags = strings.Join([]string{
			fmt.Sprintf(`-X "%s.Version=%s"`, versionPkg, semver),
			fmt.Sprintf(`-X "%s.Commit=%s"`, versionPkg, gitSHA[:7]),
			fmt.Sprintf(`-X "%s.Branch=%s"`, versionPkg, gitBranch),
			fmt.Sprintf(`-X "%s.GoVersion=%s"`, versionPkg, goVersion),
			fmt.Sprintf(`-X "%s.BuildTool=%s"`, versionPkg, buildTool),
			fmt.Sprintf(`-X "%s.BuildTime=%s"`, versionPkg, buildTime),
		}, " ")
	}

	// ==============================> DECORATE <==============================

	if c.spec.Build.Decorate {
		if err := c.services.decorator.Decorate(info.WorkingDirectory); err != nil {
			c.ui.Error(err.Error())
			return command.DecorationError
		}
	}

	// ==============================> BUILD BINARIES <==============================

	var cmdPath string
	if c.spec.Build.Decorate {
		cmdPath = fmt.Sprintf("./%s/%s/", decoratedDir, cmdDir)
	} else {
		cmdPath = fmt.Sprintf("./%s/", cmdDir)
	}

	// By convention, we assume every directory inside cmd is a main package for a binary with the same name as the directory name.
	if _, err = os.Stat(cmdPath); err == nil {
		files, err := ioutil.ReadDir(cmdPath)
		if err != nil {
			c.ui.Error(err.Error())
			return command.OSError
		}

		for _, file := range files {
			if file.IsDir() {
				mainPkg := cmdPath + file.Name()
				output := binPath + file.Name()

				if err := c.buildAll(ctx, ldFlags, mainPkg, output); err != nil {
					c.ui.Error(err.Error())
					return command.GoError
				}
			}
		}
	}

	// We also assume the current directory is a main package if it contains a main.go file.
	if _, err = os.Stat("./main.go"); err == nil {
		mainPkg := "."
		output := binPath + filepath.Base(info.WorkingDirectory)

		if err := c.buildAll(ctx, ldFlags, mainPkg, output); err != nil {
			c.ui.Error(err.Error())
			return command.GoError
		}
	}

	if len(c.outputs.artifacts) == 0 {
		c.ui.Warn("No main package found.")
		c.ui.Warn("Run gelato build -help for more information.")
	}

	// ==============================> DONE <==============================

	return command.Success
}

func (c *Command) buildAll(ctx context.Context, ldFlags, mainPkg, output string) error {
	if !c.spec.Build.CrossCompile {
		return c.build(ctx, "", "", ldFlags, mainPkg, output)
	}

	// Cross-compiling
	group, groupCtx := errgroup.WithContext(ctx)
	for _, platform := range c.spec.Build.Platforms {
		output := output + "-" + platform
		vals := strings.Split(platform, "-")

		group.Go(func() error {
			return c.build(groupCtx, vals[0], vals[1], ldFlags, mainPkg, output)
		})
	}

	return group.Wait()
}

func (c *Command) build(ctx context.Context, os, arch, ldFlags, mainPkg, output string) error {
	opts := shell.RunOptions{
		Environment: map[string]string{
			"GOOS":   os,
			"GOARCH": arch,
		},
	}

	args := []string{}
	if ldFlags != "" {
		args = append(args, "-ldflags", ldFlags)
	}
	if output != "" {
		args = append(args, "-o", output)
	}
	args = append(args, mainPkg)

	_, _, err := c.funcs.goBuild(ctx, opts, args...)
	if err != nil {
		return err
	}

	c.Mutex.Lock()
	c.outputs.artifacts = append(c.outputs.artifacts, Artifact{
		Path: output,
	})
	c.Mutex.Unlock()

	c.ui.Output("🍨 " + output)

	return nil
}

// Artifacts returns the build artifacts after the command is run.
func (c *Command) Artifacts() []Artifact {
	return c.outputs.artifacts
}
