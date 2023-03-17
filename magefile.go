//go:build mage
// +build mage

package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var golangCILintVer = "v1.48.0" // https://github.com/golangci/golangci-lint/releases
var gosImportsVer = "v0.1.5"    // https://github.com/rinchsan/gosimports/releases/tag/v0.1.5

var errRunGoModTidy = errors.New("go.mod/sum not formatted, commit changes")

// Lint verifies code quality.
func Lint() error {
	if err := sh.RunV("go", "run", fmt.Sprintf("github.com/golangci/golangci-lint/cmd/golangci-lint@%s", golangCILintVer), "run"); err != nil {
		return err
	}

	if err := sh.RunV("go", "mod", "tidy"); err != nil {
		return err
	}

	if sh.Run("git", "diff", "--exit-code", "go.mod", "go.sum") != nil {
		return errRunGoModTidy
	}

	return nil
}

func build(goos string) error {
	if err := os.MkdirAll("build", 0755); err != nil {
		return err
	}

	suffix := ""
	env := map[string]string{}
	if goos != "" {
		suffix = "-" + goos
		env["GOOS"] = goos
	}

	return sh.RunWithV(env, "go", "build", "-o", "build/coraza-httpbin"+suffix, "cmd/coraza-httpbin/main.go")
}

// Build builds the project
func Build() error {
	return build("")
}

func BuildLinux() error {
	return build("linux")
}

func BuildDockerImage() {
	mg.Deps(BuildLinux)

	if err := sh.RunV("docker", "build", "-t", "ghcr.io/jcchavezs/coraza-httpbin", "."); err != nil {
		return
	}
}