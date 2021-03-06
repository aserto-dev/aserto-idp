//go:build mage
// +build mage

package main

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/aserto-dev/mage-loot/common"
	"github.com/aserto-dev/mage-loot/deps"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func init() {
	// Set go version for docker builds
	os.Setenv("GO_VERSION", "1.17")
	// Set private repositories
	os.Setenv("GOPRIVATE", "github.com/aserto-dev")
}

// Generate generates all code.
func Generate() error {
	return common.GenerateWith([]string{
		filepath.Dir(deps.GoBinPath("mockgen")),
		filepath.Dir(deps.GoBinPath("wire")),
	})
}

// Build builds all binaries in ./cmd.
func Build() error {
	return common.BuildReleaser()
}

// Cleans the bin director
func Clean() error {
	return os.RemoveAll("dist")
}

// Release releases the project.
func Release() error {
	return common.Release()
}

// BuildAll builds all binaries in ./cmd for
// all configured operating systems and architectures.
func BuildAll() error {
	return common.BuildAllReleaser("--snapshot")
}

func Deps() {
	deps.GetAllDeps()
}

// Lint runs linting for the entire project.
func Lint() error {
	return common.Lint()
}

// Test runs all tests and generates a code coverage report.
func Test() error {
	return common.Test()
}

// All runs all targets in the appropriate order.
// The targets are run in the following order:
// deps, generate, lint, test, build, dockerImage
func All() error {
	mg.SerialDeps(Deps, Generate, Lint, Test, Build)
	return nil
}

func Run() error {
	return sh.RunV("./dist//aserto-idp_" + runtime.GOOS + "_" + runtime.GOARCH + "/aserto-idp")
}
