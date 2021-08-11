// +build mage

package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/aserto-dev/mage-loot/common"
	"github.com/aserto-dev/mage-loot/deps"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pkg/errors"
)

func init() {
	// Set go version for docker builds
	os.Setenv("GO_VERSION", "1.16")
	// Set private repositories
	os.Setenv("GOPRIVATE", "github.com/aserto-dev")
}

// Generate generates all code.
func Generate() error {
	return common.Generate()
}

// Build builds all binaries in ./cmd.
func Build() error {
	flags, err := ldflags()
	if err != nil {
		return err
	}

	return common.Build(flags...)
}

// BuildAll builds all binaries in ./cmd for
// all configured operating systems and architectures.
func BuildAll() error {
	return common.BuildAll()
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

func ldflags() ([]string, error) {
	commit, err := common.Commit()
	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate git commit")
	}
	version, err := common.Version()
	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate version")
	}

	date := time.Now().UTC().Format(time.RFC3339)

	ldbase := "github.com/aserto-dev/aserto-idp/pkg/version"
	ldflags := fmt.Sprintf(`-X %s.ver=%s -X %s.commit=%s -X %s.date=%s`,
		ldbase, version, ldbase, commit, ldbase, date)

	return []string{"-ldflags", ldflags}, nil
}

func Run() error {
	return sh.RunV("./bin/" + runtime.GOOS + "-" + runtime.GOARCH + "/aserto-idp")
}
