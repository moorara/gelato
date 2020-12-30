package spec

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	specFiles        = []string{"gelato.yml", "gelato.yaml", "gelato.json"}
	defaultPlatforms = []string{"linux-386", "linux-amd64", "linux-arm", "linux-arm64", "darwin-amd64", "windows-386", "windows-amd64"}
)

// Spec is the model for all specifications.
type Spec struct {
	GelatoVersion string `json:"-" yaml:"-"`

	Version string  `json:"version" yaml:"version"`
	App     App     `json:"app" yaml:"app"`
	Build   Build   `json:"build" yaml:"build"`
	Release Release `json:"release" yaml:"release"`
}

// FromFile reads and returns specifications from a file.
// If no spec file is found, an empty spec will be returned.
func FromFile() (Spec, error) {
	var spec Spec
	zero := Spec{}

	for _, specFile := range specFiles {
		file, err := os.Open(specFile)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return zero, err
		}
		defer file.Close()

		if ext := filepath.Ext(specFile); ext == ".yml" || ext == ".yaml" {
			err = yaml.NewDecoder(file).Decode(&spec)
		} else if ext == ".json" {
			err = json.NewDecoder(file).Decode(&spec)
		} else {
			return zero, errors.New("unknown spec file")
		}

		if err != nil {
			return zero, err
		}

		return spec, nil
	}

	return spec, nil
}

// WithDefaults returns a new object with default values.
func (s Spec) WithDefaults() Spec {
	if s.Version == "" {
		s.Version = "1.0"
	}

	s.App = s.App.WithDefaults()
	s.Build = s.Build.WithDefaults()
	s.Release = s.Release.WithDefaults()

	return s
}

// App has the specifications for an application.
type App struct {
	Language string `json:"language" yaml:"language"`
	Type     string `json:"type" yaml:"type"`
	Layout   string `json:"layout" yaml:"layout"`
}

const (
	// AppLanguageGo represents Go programming language.
	AppLanguageGo = "go"

	// AppTypeCLI represents a command-line application.
	AppTypeCLI = "cli"
	// AppTypeHTTPService represents an HTTP service.
	AppTypeHTTPService = "http-service"
	// AppTypeGRPCService represents a gRPC service.
	AppTypeGRPCService = "grpc-service"

	// AppLayoutVertical represents a vertical application layout.
	AppLayoutVertical = "vertical"
	// AppLayoutHorizontal represents a horizontal application layout (a.k.a. onion architecture).
	AppLayoutHorizontal = "horizontal"
)

// WithDefaults returns a new object with default values.
func (a App) WithDefaults() App {
	if a.Language == "" {
		a.Language = AppLanguageGo
	}

	if a.Type == "" {
		a.Type = AppTypeHTTPService
	}

	if a.Layout == "" {
		a.Layout = AppLayoutVertical
	}

	return a
}

// FlagSet returns a flag set for the app command arguments.
func (a *App) FlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("app", flag.ContinueOnError)
	fs.StringVar(&a.Language, "language", a.Language, "")
	fs.StringVar(&a.Type, "type", a.Type, "")
	fs.StringVar(&a.Layout, "layout", a.Layout, "")

	return fs
}

// Build has the specifications for the build command.
type Build struct {
	CrossCompile bool     `json:"crossCompile" yaml:"cross_compile"`
	Decorate     bool     `json:"decorate" yaml:"decorate"`
	Platforms    []string `json:"platforms" yaml:"platforms"`
}

// WithDefaults returns a new object with default values.
func (b Build) WithDefaults() Build {
	if len(b.Platforms) == 0 {
		b.Platforms = defaultPlatforms
	}

	return b
}

// FlagSet returns a flag set for the build command arguments.
func (b *Build) FlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("build", flag.ContinueOnError)
	fs.BoolVar(&b.CrossCompile, "cross-compile", b.CrossCompile, "")
	fs.BoolVar(&b.Decorate, "decorate", b.Decorate, "")

	return fs
}

// Release has the specifications for the release command.
type Release struct {
	Artifacts bool `json:"artifacts" yaml:"artifacts"`
}

// WithDefaults returns a new object with default values.
func (r Release) WithDefaults() Release {
	return r
}

// FlagSet returns a flag set for the release command arguments.
func (r *Release) FlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("release", flag.ContinueOnError)
	fs.BoolVar(&r.Artifacts, "artifacts", r.Artifacts, "")

	return fs
}
