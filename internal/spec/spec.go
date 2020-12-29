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
	Language AppLanguage `json:"language" yaml:"language"`
	Type     AppType     `json:"type" yaml:"type"`
	Layout   AppLayout   `json:"layout" yaml:"layout"`
}

// AppLanguage specifies the programming language of an application.
type AppLanguage string

const (
	// AppLanguageGo represents Go programming language.
	AppLanguageGo AppLanguage = "go"
)

// AppType specifies the type of an application.
type AppType string

const (
	// AppTypeGeneric represents a generic application.
	AppTypeGeneric AppType = "generic"
	// AppTypeCLI represents a command-line application.
	AppTypeCLI AppType = "cli"
	// AppTypeService represents a backend service.
	AppTypeService AppType = "service"
)

// AppLayout specifies the layout of an application code.
type AppLayout string

const (
	// AppLayoutCustom represents a custom application.
	AppLayoutCustom AppLayout = "custom"
	// AppLayoutVertical represents a vertical application layout.
	AppLayoutVertical AppLayout = "vertical"
	// AppLayoutHorizontal represents a horizontal application layout (a.k.a. onion architecture).
	AppLayoutHorizontal AppLayout = "horizontal"
)

// WithDefaults returns a new object with default values.
func (a App) WithDefaults() App {
	if a.Language == "" {
		a.Language = AppLanguageGo
	}

	if a.Type == "" {
		a.Type = AppTypeGeneric
	}

	if a.Layout == "" {
		a.Layout = AppLayoutCustom
	}

	return a
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
