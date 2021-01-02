[![Go Doc][godoc-image]][godoc-url]
[![Build Status][workflow-image]][workflow-url]
[![Go Report Card][goreport-image]][goreport-url]
[![Test Coverage][coverage-image]][coverage-url]
[![Maintainability][maintainability-image]][maintainability-url]

# Gelato

Gelato is an opinionated tool and framework for Go applications.
My vision for Gelato is an application framework for generating, building, and releasing Go applications and services
with built-in secure defaults, observability, and best practices.

## Why?

_TBD_

## Quick Start

### Install

```
brew install moorara/brew/gelato
```

For other platforms, you can download the binary from the [latest release](https://github.com/moorara/gelato/releases/latest).

#### Dependencies

  - [go](https://golang.org)

### Examples

```bash
# Show available commands
gelato -help

# Self-update Gelato
gelato update

# Show the current semantic version
gelato semver

# Build a Go application
gelato build

# Release a repository
gelato release
```

### Spec File

You can check in a file in your repository for configuring how Gelato commands are executed.
JSON format is also supported.

<details>
  <summary>gelato.yaml</summary>

```yaml
version: "1.0"

build:
  cross_compile: true
  platforms:
    - linux-386
    - linux-amd64
    - linux-arm
    - linux-arm64
    - darwin-amd64
    - windows-386
    - windows-amd64

release:
  artifacts: true
```
</details>

## Versioning

Gelato uses Semantic Versioning 2.0.0 as described [here](https://semver.org).
It supports injecting build metadata into your binaries by including a `version` package in your repository.

<details>
  <summary>version.go</summary>

```go
var (
  Version   string
  Commit    string
  Branch    string
  GoVersion string
  BuildTool string
  BuildTime string
)
```
</details>

## Commands

### `update`

`gelato update` updates Gelato to its latest version.
It downloads the latest release for your system from GitHub and replaces the local binary.

### `semver`

`gelato semver` resolves and prints the current semantic version.
This command can be used to get the current semantic version for building artifacts such as Docker image.

### `build`

`gelato build` compiles your binary and injects the build metadata into the `version` package (if any).

`gelato build -cross-compile` builds the binaries for all supported platforms.

`gelato build -decorate` decorates an application with a set of decorators.
Decoration is an experimental feature to decorate the applications with **horizontal layout**.
It wraps the `controller`, `gateway`, `handler`, and `repository` packages with a set of decorators.
Decorators can be used for augmenting an application with *observability*, *error reccovery*, etc.

### `release`

`gelato release` can be used for releasing a **GitHub** repository.
You can use `-patch`, `-minor`, or `-major` flags to release different semantic versions.
You can also use `-comment` flag to include a description for your release.

`GELATO_GITHUB_TOKEN` environment variable should be set to a [personal access token](https://github.com/settings/tokens) with `repo` scope.
The user who is generating the token should also have `Admin` permission to repositories.

The initial release is always `0.1.0`.


[godoc-url]: https://pkg.go.dev/github.com/moorara/gelato
[godoc-image]: https://pkg.go.dev/badge/github.com/moorara/gelato
[workflow-url]: https://github.com/moorara/gelato/actions
[workflow-image]: https://github.com/moorara/gelato/workflows/Main/badge.svg
[goreport-url]: https://goreportcard.com/report/github.com/moorara/gelato
[goreport-image]: https://goreportcard.com/badge/github.com/moorara/gelato
[coverage-url]: https://codeclimate.com/github/moorara/gelato/test_coverage
[coverage-image]: https://api.codeclimate.com/v1/badges/a2ea750b7fe5a629654c/test_coverage
[maintainability-url]: https://codeclimate.com/github/moorara/gelato/maintainability
[maintainability-image]: https://api.codeclimate.com/v1/badges/a2ea750b7fe5a629654c/maintainability
