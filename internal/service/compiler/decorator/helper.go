package decorator

import "strings"

const (
	mainPkg       = "main"
	handlerPkg    = "handler"
	controllerPkg = "controller"
	gatewayPkg    = "gateway"
	repositoryPkg = "repository"
)

// getOriginalPkgName returns a package names that can be used for referencing the original package.
func getOriginalPkgName(name string) string {
	return "_" + name
}

// getDecoratedPkgName returns a package name that can be used for referencing the decorated package.
func getDecoratedPkgName(name string) string {
	return "_" + name
}

// isMainPkg determines
func isMainPkg(name string) bool {
	return name == mainPkg
}

func isDecoratablePkg(importPath string) bool {
	isHandlerPkg := strings.HasSuffix(importPath, "/"+handlerPkg) || strings.Contains(importPath, "/"+handlerPkg+"/")
	isControllerPkg := strings.HasSuffix(importPath, "/"+controllerPkg) || strings.Contains(importPath, "/"+controllerPkg+"/")
	isGatewayPkg := strings.HasSuffix(importPath, "/"+gatewayPkg) || strings.Contains(importPath, "/"+gatewayPkg+"/")
	isRepositoryPkg := strings.HasSuffix(importPath, "/"+repositoryPkg) || strings.Contains(importPath, "/"+repositoryPkg+"/")
	return isHandlerPkg || isControllerPkg || isGatewayPkg || isRepositoryPkg
}
