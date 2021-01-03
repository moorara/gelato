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
    -module      The Go module name for the new application
    -language    the programming language of the new application (values: go{{if .App.Language}}, default: {{.App.Language}}{{end}})
    -type        the type of the new application (values: cli|http-service|grpc-service{{if .App.Type}}, default: {{.App.Type}}){{end}})
    -layout      the layout of the new application (values: vertical|horizontal{{if .App.Layout}}, default: {{.App.Layout}}){{end}})

  Examples:
    gelato app
    gelato app -module=github.com/octocat/service -type=http-service -layout=vertical
  `
)

const (
	templateOwner = "moorara"
	templateRepo  = "gelato"
	defaultRef    = "main"
)

type (
	repoService interface {
		DownloadTarArchive(context.Context, string, io.Writer) (*github.Response, error)
	}

	archiveService interface {
		Extract(string, io.Reader, archive.Selector) error
	}

	editService interface {
		ReplaceInDir(string, []edit.ReplaceSpec) error
	}
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

	client := github.NewClient(token)
	repo := client.Repo(templateOwner, templateRepo)

	c.services.repo = repo
	c.services.arch = archive.NewTarArchive(log.Info)
	c.services.edit = edit.NewEditor(log.Trace)

	return c.run(args)
}

// run in an auxiliary method, so we can test the business logic with mock dependencies.
func (c *Command) run(args []string) int {
	flags := struct {
		module string
		docker string
	}{}

	fs := c.spec.App.FlagSet()
	fs.StringVar(&flags.module, "module", flags.module, "")
	fs.StringVar(&flags.docker, "docker", flags.docker, "")
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
		c.spec.App.Language, err = c.ui.Ask(fmt.Sprintf("Application Language (%s): ", langOptions))
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
		c.spec.App.Type, err = c.ui.Ask(fmt.Sprintf("Application Type (%s): ", typeOptions))
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
		c.spec.App.Layout, err = c.ui.Ask(fmt.Sprintf("Application Layout (%s): ", layoutOptions))
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
		flags.module, err = c.ui.Ask("Go module name: ")
		if err != nil {
			c.ui.Error(fmt.Sprintf("invalid module name: %s", err))
			return command.InputError
		}

		if flags.module == "" {
			c.ui.Error(fmt.Sprintf("unsupported module name: %s", flags.module))
			return command.UnsupportedError
		}
	}

	if flags.docker == "" {
		flags.docker, err = c.ui.Ask("Docker ID: ")
		if err != nil {
			c.ui.Error(fmt.Sprintf("invalid Docker ID: %s", err))
			return command.InputError
		}

		if flags.docker == "" {
			c.ui.Error(fmt.Sprintf("unsupported Docker ID: %s", flags.docker))
			return command.UnsupportedError
		}
	}

	// ==============================> DOWNLOAD REPO ARCHIVE <==============================

	buf := new(bytes.Buffer)

	ref := c.spec.Gelato.Revision
	if ref == "" {
		ref = defaultRef
	}

	c.ui.Output(fmt.Sprintf("Downloading templates revision %s ...", ref))

	_, err = c.services.repo.DownloadTarArchive(ctx, ref, buf)
	if err != nil {
		c.ui.Error(fmt.Sprintf("Failed to download repository archive: %s", err))
		return command.GitHubError
	}

	// ==============================> EXTRACT REPO ARCHIVE <==============================

	c.ui.Output(fmt.Sprintf("Extracting templates revision %s ...", ref))

	appName := filepath.Base(flags.module)

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

	err = c.services.arch.Extract(info.Context.WorkingDirectory, buf, func(path string) (string, bool) {
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

	// ==============================> EDIT <==============================

	c.ui.Output(fmt.Sprintf("Finishing %s ...", appName))

	specs := []edit.ReplaceSpec{
		// Edit module name and import paths
		{
			PathRE: regexp.MustCompile(`(go.mod|\.go|\.proto)$`),
			OldRE: regexp.MustCompile(fmt.Sprintf(`%s/%s`,
				c.spec.App.Layout,
				c.spec.App.Type,
			)),
			New: flags.module,
		},
		// Edit application name
		{
			PathRE: regexp.MustCompile(`(main\.go|\.gitignore|\.dockerignore|Makefile|Dockerfile|Dockerfile\.test|docker-compose\.yml)$`),
			OldRE:  regexp.MustCompile(c.spec.App.Type),
			New:    appName,
		},
		// Edit Docker image
		{
			PathRE: regexp.MustCompile(`Makefile$`),
			OldRE:  regexp.MustCompile(`dockerid`),
			New:    flags.docker,
		},
	}

	err = c.services.edit.ReplaceInDir(appName, specs)
	if err != nil {
		c.ui.Error(err.Error())
		return command.OSError
	}

	// ==============================> DONE <==============================

	return command.Success
}
