package main

import (
	"fmt"
	"runtime"
	"strings"

	"dagger/go/internal/dagger"
)

// Returns a binary
func (m *Go) Build(
	// The binary name to build
	name string,
	// Go packages
	packages []string,
	// +optional
	ldflags []string,
	// The arch to build for
	// +optional
	arch string,
	// The os to build for
	// +optional
	// +default="linux"
	os string,
) *dagger.File {
	if len(arch) == 0 {
		arch = runtime.GOARCH
	}

	goBuildLdflags := ldflags
	goBuildLdflags = append(goBuildLdflags,
		"-extldflags",
		"-static",
		"-s",
		"-w",
	)

	binaryName := fmt.Sprintf("%s-%s-%s", name, os, arch)
	binaryPath := fmt.Sprintf("build/%s", binaryName)

	args := []string{
		"go",
		"build",
		"-o", binaryPath,
		"-ldflags", strings.Join(goBuildLdflags, " "),
	}
	args = append(args, packages...)

	// WithEnvVariable("GOMODCACHE", "/go/pkg/mod").

	return dag.Container().
		From(fmt.Sprintf("golang:%s", m.Version)).
		WithMountedDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume(fmt.Sprintf("go-mod-%s", m.Version))).
		WithMountedCache("/go/build-cache", dag.CacheVolume(fmt.Sprintf("go-build-%s", m.Version))).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithEnvVariable("GOOS", os).
		WithEnvVariable("GOARCH", arch).
		WithExec(args).
		File(binaryPath)
}
