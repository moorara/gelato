package decorator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOriginalPkgName(t *testing.T) {
	tests := []struct {
		name         string
		pkgName      string
		expectedName string
	}{
		{
			name:         "OK",
			pkgName:      "controller",
			expectedName: "_controller",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			name := getOriginalPkgName(tc.pkgName)

			assert.Equal(t, tc.expectedName, name)
		})
	}
}

func TestGetDecoratedPkgName(t *testing.T) {
	tests := []struct {
		name         string
		pkgName      string
		expectedName string
	}{
		{
			name:         "OK",
			pkgName:      "controller",
			expectedName: "_controller",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			name := getDecoratedPkgName(tc.pkgName)

			assert.Equal(t, tc.expectedName, name)
		})
	}
}

func TestIsMainPkg(t *testing.T) {
	tests := []struct {
		name           string
		pkgName        string
		expectedResult bool
	}{
		{
			name:           "MainPackage",
			pkgName:        "main",
			expectedResult: true,
		},
		{
			name:           "NotMainPackage",
			pkgName:        "controller",
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := isMainPkg(tc.pkgName)

			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestIsDecoratablePkg(t *testing.T) {
	tests := []struct {
		name           string
		importPath     string
		expectedResult bool
	}{
		{
			name:           "HandlerPackage",
			importPath:     "github.com/octocat/service/internal/handler",
			expectedResult: true,
		},
		{
			name:           "HandlerSubpackage",
			importPath:     "github.com/octocat/service/internal/handler/lookup",
			expectedResult: true,
		},
		{
			name:           "ControllerPackage",
			importPath:     "github.com/octocat/service/internal/controller",
			expectedResult: true,
		},
		{
			name:           "ControllerSubpackage",
			importPath:     "github.com/octocat/service/internal/controller/lookup",
			expectedResult: true,
		},
		{
			name:           "GatewayPackage",
			importPath:     "github.com/octocat/service/internal/gateway",
			expectedResult: true,
		},
		{
			name:           "GatewaySubpackage",
			importPath:     "github.com/octocat/service/internal/gateway/lookup",
			expectedResult: true,
		},
		{
			name:           "RepositoryPackage",
			importPath:     "github.com/octocat/service/internal/repository",
			expectedResult: true,
		},
		{
			name:           "RepositorySubpackage",
			importPath:     "github.com/octocat/service/internal/repository/lookup",
			expectedResult: true,
		},
		{
			name:           "NotDecoratablePackage",
			importPath:     "github.com/octocat/service/internal/mapper",
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := isDecoratablePkg(tc.importPath)

			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
