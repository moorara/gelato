package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromFile(t *testing.T) {
	tests := []struct {
		name          string
		specFiles     []string
		expectedSpec  Spec
		expectedError string
	}{
		{
			name:         "NoSpecFile",
			specFiles:    []string{"test/null"},
			expectedSpec: Spec{},
		},
		{
			name:          "UnknownFile",
			specFiles:     []string{"test/unknown"},
			expectedError: "unknown spec file",
		},
		{
			name:          "EmptyJSON",
			specFiles:     []string{"test/empty.json"},
			expectedError: "EOF",
		},
		{
			name:          "EmptyYAML",
			specFiles:     []string{"test/empty.yaml"},
			expectedError: "EOF",
		},
		{
			name:          "InvalidJSON",
			specFiles:     []string{"test/invalid.json"},
			expectedError: "invalid character",
		},
		{
			name:          "InvalidYAML",
			specFiles:     []string{"test/invalid.yaml"},
			expectedError: "cannot unmarshal",
		},
		{
			name:      "ValidJSON",
			specFiles: []string{"test/valid.json"},
			expectedSpec: Spec{
				Version: "1.0",
				App: App{
					Language: AppLanguageGo,
					Type:     AppTypeGRPCService,
					Layout:   AppLayoutHorizontal,
				},
				Build: Build{
					CrossCompile: true,
					Decorate:     true,
					Platforms:    []string{"linux-386", "linux-amd64", "linux-arm", "linux-arm64", "darwin-amd64", "windows-386", "windows-amd64"},
				},
				Release: Release{
					Artifacts: true,
				},
			},
		},
		{
			name:      "ValidYAML",
			specFiles: []string{"test/valid.yaml"},
			expectedSpec: Spec{
				Version: "1.0",
				App: App{
					Language: AppLanguageGo,
					Type:     AppTypeGRPCService,
					Layout:   AppLayoutHorizontal,
				},
				Build: Build{
					CrossCompile: true,
					Decorate:     true,
					Platforms:    []string{"linux-386", "linux-amd64", "linux-arm", "linux-arm64", "darwin-amd64", "windows-386", "windows-amd64"},
				},
				Release: Release{
					Artifacts: true,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			specFiles = tc.specFiles
			spec, err := FromFile()

			if tc.expectedError != "" {
				assert.Contains(t, err.Error(), tc.expectedError)
				assert.Equal(t, Spec{}, spec)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedSpec, spec)
			}
		})
	}
}

func TestSpecWithDefaults(t *testing.T) {
	tests := []struct {
		name         string
		spec         Spec
		expectedSpec Spec
	}{
		{
			"DefaultsRequired",
			Spec{},
			Spec{
				GelatoVersion: "",
				Version:       "1.0",
				App: App{
					Language: AppLanguageGo,
					Type:     "",
					Layout:   "",
				},
				Build: Build{
					CrossCompile: false,
					Decorate:     false,
					Platforms:    defaultPlatforms,
				},
				Release: Release{
					Artifacts: false,
				},
			},
		},
		{
			"DefaultsNotRequired",
			Spec{
				GelatoVersion: "0.1.0",
				Version:       "2.0",
				App: App{
					Language: AppLanguageGo,
					Type:     AppTypeGRPCService,
					Layout:   AppLayoutHorizontal,
				},
				Build: Build{
					CrossCompile: true,
					Decorate:     true,
					Platforms:    []string{"linux-amd64", "darwin-amd64", "windows-amd64"},
				},
				Release: Release{
					Artifacts: true,
				},
			},
			Spec{
				GelatoVersion: "0.1.0",
				Version:       "2.0",
				App: App{
					Language: AppLanguageGo,
					Type:     AppTypeGRPCService,
					Layout:   AppLayoutHorizontal,
				},
				Build: Build{
					CrossCompile: true,
					Decorate:     true,
					Platforms:    []string{"linux-amd64", "darwin-amd64", "windows-amd64"},
				},
				Release: Release{
					Artifacts: true,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedSpec, tc.spec.WithDefaults())
		})
	}
}

func TestAppWithDefaults(t *testing.T) {
	tests := []struct {
		name        string
		app         App
		expectedApp App
	}{
		{
			"DefaultsRequired",
			App{},
			App{
				Language: AppLanguageGo,
				Type:     "",
				Layout:   "",
			},
		},
		{
			"DefaultsNotRequired",
			App{
				Language: AppLanguageGo,
				Type:     AppTypeGRPCService,
				Layout:   AppLayoutHorizontal,
			},
			App{
				Language: AppLanguageGo,
				Type:     AppTypeGRPCService,
				Layout:   AppLayoutHorizontal,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedApp, tc.app.WithDefaults())
		})
	}
}

func TestBuildWithDefaults(t *testing.T) {
	tests := []struct {
		name          string
		build         Build
		expectedBuild Build
	}{
		{
			"DefaultsRequired",
			Build{},
			Build{
				CrossCompile: false,
				Decorate:     false,
				Platforms:    defaultPlatforms,
			},
		},
		{
			"DefaultsNotRequired",
			Build{
				CrossCompile: true,
				Decorate:     true,
				Platforms:    []string{"linux-amd64", "darwin-amd64", "windows-amd64"},
			},
			Build{
				CrossCompile: true,
				Decorate:     true,
				Platforms:    []string{"linux-amd64", "darwin-amd64", "windows-amd64"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedBuild, tc.build.WithDefaults())
		})
	}
}

func TestBuildFlagSet(t *testing.T) {
	tests := []struct {
		build Build
	}{
		{
			build: Build{},
		},
		{
			build: Build{
				CrossCompile: true,
				Decorate:     true,
				Platforms:    []string{"linux-amd64", "darwin-amd64", "windows-amd64"},
			},
		},
	}

	for _, tc := range tests {
		fs := tc.build.FlagSet()

		assert.NotNil(t, fs)
	}
}

func TestReleaseWithDefaults(t *testing.T) {
	tests := []struct {
		name            string
		release         Release
		expectedRelease Release
	}{
		{
			"DefaultsRequired",
			Release{},
			Release{
				Artifacts: false,
			},
		},
		{
			"DefaultsNotRequired",
			Release{
				Artifacts: true,
			},
			Release{
				Artifacts: true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedRelease, tc.release.WithDefaults())
		})
	}
}

func TestReleaseFlagSet(t *testing.T) {
	tests := []struct {
		release Release
	}{
		{
			release: Release{},
		},
		{
			release: Release{
				Artifacts: true,
			},
		},
	}

	for _, tc := range tests {
		fs := tc.release.FlagSet()

		assert.NotNil(t, fs)
	}
}
