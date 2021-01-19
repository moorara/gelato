package app

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/mitchellh/cli"
	"github.com/moorara/go-github"

	"github.com/moorara/gelato/internal/command"
	"github.com/moorara/gelato/internal/log"
	"github.com/moorara/gelato/internal/service/archive"
	"github.com/moorara/gelato/internal/service/edit"
	"github.com/moorara/gelato/internal/service/git"
	"github.com/moorara/gelato/internal/spec"
)

const (
	appTimeout  = 2 * time.Minute
	appSynopsis = `Create an application`
	appHelp     = `
  Use this command for creating a new application.
  Currently, the app command can only create Go applications.

  Usage:  gelato app [flags]

  Flags:
    -language    the programming language of the new application (values: go{{if .App.Language}}, default: {{.App.Language}}{{end}})
    -type        the type of the new application (values: cli|http-service|grpc-service{{if .App.Type}}, default: {{.App.Type}}){{end}})
    -layout      the layout of the new application (values: vertical|horizontal{{if .App.Layout}}, default: {{.App.Layout}}){{end}})
    -module      the Go module name for the new application
    -docker      the Docker ID for the Docker image of the new application
    -owners      a list of GitHub usernames, teams, or emails as code owners separated by space

  Examples:
    gelato app
    gelato app -type=http-service -layout=vertical -module=github.com/octocat/service -docker=octocat -owners=@octocat
  `
)

const (
	templateOwner = "moorara"
	templateRepo  = "gelato"
	makeSubmod    = "make"
)

type (
	repoService interface {
		DownloadTarArchive(context.Context, string, io.Writer) (*github.Response, error)
	}

	archiveService interface {
		Extract(string, io.Reader, archive.Selector) error
	}

	editService interface {
		Remove(...string) error
		Move(bool, ...edit.MoveSpec) error
		Append(bool, ...edit.AppendSpec) error
		ReplaceInDir(string, ...edit.ReplaceSpec) error
	}

	gitService interface {
		Path() (string, error)
		Remote(string) (string, string, error)
		Submodule(string) (git.Submodule, error)
		UpdateSubmodules() error
	}

	detectGitFunc func(string) (string, error)
	gitFunc       func(string) (gitService, error)
)

// Command is the cli.Command implementation for app command.
type Command struct {
	ui       cli.Ui
	spec     spec.Spec
	services struct {
		repo repoService
		arch archiveService
		edit editService
	}
	funcs struct {
		detectGit detectGitFunc
		gitInit   gitFunc
		gitOpen   gitFunc
	}
	outputs struct{}
}

// NewCommand creates an app command.
func NewCommand(ui cli.Ui, spec spec.Spec) (*Command, error) {
	return &Command{
		ui:   ui,
		spec: spec,
	}, nil
}

// Synopsis returns a short one-line synopsis of the command.
func (c *Command) Synopsis() string {
	return appSynopsis
}

// Help returns a long help text including usage, description, and list of flags for the command.
func (c *Command) Help() string {
	var buf bytes.Buffer
	t := template.Must(template.New("help").Parse(appHelp))
	_ = t.Execute(&buf, c.spec)
	return buf.String()
}

// Run runs the actual command with the given command-line arguments.
// This method is used as a proxy for creating dependencies and the actual command execution is delegated to the run method for testing purposes.
func (c *Command) Run(args []string) int {
	// If no access token is provided, we try without it!
	token := os.Getenv("GELATO_GITHUB_TOKEN")

	c.services.repo = github.NewClient(token).Repo(templateOwner, templateRepo)
	c.services.arch = archive.NewTarArchive(log.Info)
	c.services.edit = edit.NewEditor(log.Info)

	c.funcs.detectGit = git.DetectGit

	c.funcs.gitInit = func(path string) (gitService, error) {
		return git.Init(path)
	}

	c.funcs.gitOpen = func(path string) (gitService, error) {
		return git.Open(path)
	}

	return c.run(args)
}

// run in an auxiliary method, so we can test the business logic with mock dependencies.
func (c *Command) run(args []string) int {
	flags := struct {
		module string
		docker string
		owners string
	}{}

	fs := c.spec.App.FlagSet()
	fs.StringVar(&flags.module, "module", flags.module, "")
	fs.StringVar(&flags.docker, "docker", flags.docker, "")
	fs.StringVar(&flags.owners, "owners", flags.owners, "")
	fs.Usage = func() {
		c.ui.Output(c.Help())
	}

	if err := fs.Parse(args); err != nil {
		return command.FlagError
	}

	ctx, cancel := context.WithTimeout(context.Background(), appTimeout)
	defer cancel()

	// ==============================> RUN PREFLIGHT CHECKS <==============================

	checklist := command.PreflightChecklist{}

	info, err := command.RunPreflightChecks(ctx, checklist)
	if err != nil {
		c.ui.Error(err.Error())
		return command.PreflightError
	}

	// ==============================> GET INPUTS <==============================

	if c.spec.App.Language == "" {
		langOptions := strings.Join([]string{spec.AppLanguageGo}, ", ")
		c.spec.App.Language, err = c.ui.Ask(fmt.Sprintf("Application Language (%s):", langOptions))
		if err != nil {
			c.ui.Error(fmt.Sprintf("invalid application language: %s", err))
			return command.InputError
		}

		// Only Go applications are supported
		if c.spec.App.Language != spec.AppLanguageGo {
			c.ui.Error(fmt.Sprintf("unsupported application language: %s", c.spec.App.Language))
			return command.UnsupportedError
		}
	}

	if c.spec.App.Type == "" {
		typeOptions := strings.Join([]string{spec.AppTypeCLI, spec.AppTypeHTTPService, spec.AppTypeGRPCService}, ", ")
		c.spec.App.Type, err = c.ui.Ask(fmt.Sprintf("Application Type (%s):", typeOptions))
		if err != nil {
			c.ui.Error(fmt.Sprintf("invalid application type: %s", err))
			return command.InputError
		}

		// Only HTTP and gRPC services are supported
		if c.spec.App.Type != spec.AppTypeHTTPService && c.spec.App.Type != spec.AppTypeGRPCService {
			c.ui.Error(fmt.Sprintf("unsupported application type: %s", c.spec.App.Type))
			return command.UnsupportedError
		}
	}

	if c.spec.App.Layout == "" {
		layoutOptions := strings.Join([]string{spec.AppLayoutVertical, spec.AppLayoutHorizontal}, ", ")
		c.spec.App.Layout, err = c.ui.Ask(fmt.Sprintf("Application Layout (%s):", layoutOptions))
		if err != nil {
			c.ui.Error(fmt.Sprintf("invalid application layout: %s", err))
			return command.InputError
		}

		// Only vertical and horizontal layouts are supported
		if c.spec.App.Layout != spec.AppLayoutVertical && c.spec.App.Layout != spec.AppLayoutHorizontal {
			c.ui.Error(fmt.Sprintf("unsupported application layout: %s", c.spec.App.Layout))
			return command.UnsupportedError
		}
	}

	if flags.module == "" {
		flags.module, err = c.ui.Ask("Go module name:")
		if err != nil {
			c.ui.Error(fmt.Sprintf("invalid module name: %s", err))
			return command.InputError
		}

		// TODO: validate module name using a regular expression
		if flags.module == "" {
			c.ui.Error(fmt.Sprintf("unsupported module name: %s", flags.module))
			return command.UnsupportedError
		}
	}

	if flags.docker == "" {
		flags.docker, err = c.ui.Ask("Docker ID:")
		if err != nil {
			c.ui.Error(fmt.Sprintf("invalid Docker ID: %s", err))
			return command.InputError
		}

		if flags.docker == "" {
			c.ui.Error(fmt.Sprintf("unsupported Docker ID: %s", flags.docker))
			return command.UnsupportedError
		}
	}

	if flags.owners == "" {
		flags.owners, err = c.ui.Ask("GitHub code owners separated by space:")
		if err != nil {
			c.ui.Error(fmt.Sprintf("invalid GitHub code owners: %s", err))
			return command.InputError
		}

		if flags.owners == "" {
			c.ui.Error(fmt.Sprintf("unsupported GitHub code owners: %s", flags.owners))
			return command.UnsupportedError
		}
	}

	appName := filepath.Base(flags.module)
	appPath := filepath.Join(info.WorkingDirectory, appName)

	// ==============================> DOWNLOAD REPO ARCHIVE <==============================

	buf := new(bytes.Buffer)

	ref := c.spec.Gelato.Revision
	if ref == "" {
		ref = "main"
	}

	c.ui.Output(fmt.Sprintf("Downloading templates revision %s ...", ref))

	_, err = c.services.repo.DownloadTarArchive(ctx, ref, buf)
	if err != nil {
		c.ui.Error(fmt.Sprintf("Failed to download repository archive: %s", err))
		return command.GitHubError
	}

	// ==============================> EXTRACT REPO ARCHIVE <==============================

	c.ui.Output(fmt.Sprintf("Extracting templates revision %s ...", ref))

	targetPath := filepath.Join("templates",
		c.spec.App.Language,
		c.spec.App.Layout,
		c.spec.App.Type,
	)

	dirRegex, err := regexp.Compile(fmt.Sprintf("%s-%s-[0-9a-f]{7,40}/%s", templateOwner, templateRepo, targetPath))
	if err != nil {
		c.ui.Error(fmt.Sprintf("Cannot create regex for directory: %s", err))
		return command.MiscError
	}

	err = c.services.arch.Extract(info.WorkingDirectory, buf, func(path string) (string, bool) {
		if !strings.Contains(path, targetPath) {
			return "", false
		}

		path = dirRegex.ReplaceAllString(path, appName)
		return path, true
	})

	if err != nil {
		c.ui.Error(err.Error())
		return command.ExtractionError
	}

	// ==============================> OPEN GIT REPO <==============================

	var git gitService

	monorepo := false
	if _, err := c.funcs.detectGit(info.WorkingDirectory); err == nil {
		monorepo = true
	}

	if !monorepo {
		if git, err = c.funcs.gitInit(appPath); err != nil {
			c.ui.Error(fmt.Sprintf("Failed to init git repo: %s", err))
			return command.GitError
		}
	} else {
		if git, err = c.funcs.gitOpen(info.WorkingDirectory); err != nil {
			c.ui.Error(fmt.Sprintf("Failed to open git repo: %s", err))
			return command.GitError
		}
	}

	// ==============================> RESOLVE MAKE SUBMODULE <==============================

	repoPath, err := git.Path()
	if err != nil {
		c.ui.Error(fmt.Sprintf("Failed to get git repository path: %s", err))
		return command.GitError
	}

	submod, err := git.Submodule(makeSubmod)
	if err != nil {
		c.ui.Error(fmt.Sprintf("Failed to get make git submodule: %s", err))
		return command.GitError
	}

	// Resolve the absolute path for make git submodule
	makeAbsPath, err := filepath.Abs(filepath.Join(repoPath, submod.Path))
	if err != nil {
		c.ui.Error(fmt.Sprintf("Failed to resolve make submodule absolute path: %s", err))
		return command.MiscError
	}

	// Resolve the relative path for make git submodule
	makeRelPath, err := filepath.Rel(appPath, makeAbsPath)
	if err != nil {
		c.ui.Error(fmt.Sprintf("Failed to resolve make submodule relative path: %s", err))
		return command.MiscError
	}

	// Resolve the relative path of application
	relAppPath, err := filepath.Rel(repoPath, appPath)
	if err != nil {
		c.ui.Error(fmt.Sprintf("Failed to resolve application relative path: %s", err))
		return command.MiscError
	}

	// ==============================> EDIT FILES <==============================

	c.ui.Output(fmt.Sprintf("Finishing %s ...", appName))

	specs := []edit.ReplaceSpec{
		// Edit module name and import paths
		{
			PathRE: regexp.MustCompile(`(\.go|\.proto|go.mod|README.md)$`),
			OldRE: regexp.MustCompile(fmt.Sprintf(`%s/%s`,
				c.spec.App.Layout,
				c.spec.App.Type,
			)),
			New: flags.module,
		},
		// Edit application name
		{
			PathRE: regexp.MustCompile(`(main\.go|\.gitignore|\.dockerignore|Makefile|Dockerfile|Dockerfile\.test|docker-compose\.yml|\.md)$`),
			OldRE:  regexp.MustCompile(c.spec.App.Type),
			New:    appName,
		},
		// Edit Docker image
		{
			PathRE: regexp.MustCompile(`Makefile$`),
			OldRE:  regexp.MustCompile(`dockerid`),
			New:    flags.docker,
		},
		// Edit make submodule path
		{
			PathRE: regexp.MustCompile(`Makefile$`),
			OldRE:  regexp.MustCompile(`\.\./\.\./make`),
			New:    makeRelPath,
		},
		// Edit GitHub code owners
		{
			PathRE: regexp.MustCompile(`CODEOWNERS$`),
			OldRE:  regexp.MustCompile(`@octocat`),
			New:    flags.owners,
		},
	}

	if err := c.services.edit.ReplaceInDir(appPath, specs...); err != nil {
		c.ui.Error(err.Error())
		return command.OSError
	}

	// ==============================> PREPARE GIT REPO <==============================

	monorepoWorkflowPath := filepath.Join(appPath, ".github", "workflows", "monorepo.yml")

	if !monorepo {
		specs := []edit.ReplaceSpec{
			// Edit README
			{
				PathRE: regexp.MustCompile(`README.md$`),
				OldRE:  regexp.MustCompile(`REPO_URL`),
				New:    fmt.Sprintf("https://%s", flags.module),
			},
			// Edit README
			{
				PathRE: regexp.MustCompile(`README.md$`),
				OldRE:  regexp.MustCompile(`WORKFLOW_NAME`),
				New:    "Main",
			},
		}

		if err := c.services.edit.ReplaceInDir(appPath, specs...); err != nil {
			c.ui.Error(err.Error())
			return command.OSError
		}

		if err := c.services.edit.Remove(monorepoWorkflowPath); err != nil {
			c.ui.Error(fmt.Sprintf("Failed to remove: %s", err))
			return command.OSError
		}

		if err := git.UpdateSubmodules(); err != nil {
			c.ui.Error(fmt.Sprintf("Failed to add update git submodules: %s", err))
			return command.GitError
		}
	} else {
		repoDomain, repoFullName, err := git.Remote("origin")
		if err != nil {
			c.ui.Error(fmt.Sprintf("Failed to get git remote url: %s", err))
			return command.GitError
		}

		specs := []edit.ReplaceSpec{
			// Edit monorepo workflow
			{
				PathRE: regexp.MustCompile(`monorepo.yml$`),
				OldRE:  regexp.MustCompile(`APP_NAME`),
				New:    appName,
			},
			// Edit monorepo workflow
			{
				PathRE: regexp.MustCompile(`monorepo.yml$`),
				OldRE:  regexp.MustCompile(`RELATIVE_PATH`),
				New:    relAppPath,
			},
			// Edit README
			{
				PathRE: regexp.MustCompile(`README.md$`),
				OldRE:  regexp.MustCompile(`REPO_URL`),
				New:    fmt.Sprintf("https://%s/%s", repoDomain, repoFullName),
			},
			// Edit README
			{
				PathRE: regexp.MustCompile(`README.md$`),
				OldRE:  regexp.MustCompile(`WORKFLOW_NAME`),
				New:    appName,
			},
		}

		if err := c.services.edit.ReplaceInDir(appPath, specs...); err != nil {
			c.ui.Error(err.Error())
			return command.OSError
		}

		// Move workflow file
		moveWorkflow := edit.MoveSpec{
			Src:  monorepoWorkflowPath,
			Dest: filepath.Join(repoPath, ".github", "workflows", fmt.Sprintf("%s.yml", appName)),
		}

		if err := c.services.edit.Move(true, moveWorkflow); err != nil {
			c.ui.Error(fmt.Sprintf("Failed to move: %s", err))
			return command.OSError
		}

		// Add code owners
		appendCodeOwner := edit.AppendSpec{
			Path:    filepath.Join(repoPath, ".github", "CODEOWNERS"),
			Content: fmt.Sprintf("/%s/  %s", relAppPath, flags.owners),
		}

		if err := c.services.edit.Append(true, appendCodeOwner); err != nil {
			c.ui.Error(fmt.Sprintf("Failed to append: %s", err))
			return command.OSError
		}

		// Remove irrelevant directories and files
		githubDir := filepath.Join(appPath, ".github")
		gitmodFile := filepath.Join(appPath, ".gitmodules")
		if err := c.services.edit.Remove(githubDir, gitmodFile); err != nil {
			c.ui.Error(fmt.Sprintf("Failed to remove: %s", err))
			return command.OSError
		}
	}

	// ==============================> DONE <==============================

	c.ui.Info(fmt.Sprintf("%s is ready.", appName))

	if !monorepo {
		c.ui.Warn("Please create the first commit and rename your branch to main.")
		c.ui.Warn("Make sure you configure the remote repository (description, topics, settings, etc.).")
	} else {
		c.ui.Warn("Please create a new commit on a new branch.")
	}

	return command.Success
}
