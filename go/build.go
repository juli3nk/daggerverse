package main

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"dagger/go/internal/dagger"
)

// WithEnvVariable("GOMODCACHE", "/go/pkg/mod").

// Build builds a single binary for specific platform
func (m *Go) Build(
	ctx context.Context,
	// The binary name to build
	name string,
	// Go packages
	packages []string,
	// +optional
	// +default="1"
	cgoEnabled string,
	// +optional
	// +default=false
	musl bool,
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
	if arch == "" {
		arch = runtime.GOARCH
	}

	binaryName := fmt.Sprintf("%s-%s-%s", name, os, arch)
	binaryPath := fmt.Sprintf("build/%s", binaryName)

	goBuildLdflags := ldflags
	if musl {
		goBuildLdflags = append(goBuildLdflags, "-linkmode", "external")
	}
	goBuildLdflags = append(goBuildLdflags, "-extldflags", "-static", "-s", "-w")

	ctr := m.baseContainer(os, arch, cgoEnabled, musl)

	args := []string{
		"go",
		"build",
		"-o", binaryPath,
		"-ldflags",
		strings.Join(goBuildLdflags, " "),
	}
	args = append(args, packages...)

	return ctr.WithExec(args).File(binaryPath)
}

// BuildMulti builds binaries for multiple platforms
func (m *Go) BuildMulti(
	ctx context.Context,
	// The binary name to build
	name string,
	// Go packages
	packages []string,
	// Platforms to build for (format: os/arch)
	// +optional
	// +default=["linux/amd64","linux/arm64"]
	platforms []string,
	// +optional
	ldflags []string,
	// +optional
	// +default=false
	musl bool,
) *dagger.Directory {
	if len(platforms) == 0 {
		platforms = []string{"linux/amd64", "linux/arm64"}
	}

	output := dag.Directory()

	for _, platform := range platforms {
		parts := strings.Split(platform, "/")
		if len(parts) != 2 {
			continue
		}
		os, arch := parts[0], parts[1]

		cgoEnabled := "0"
		if musl {
			cgoEnabled = "1"
		}

		binary := m.Build(ctx, name, packages, cgoEnabled, musl, ldflags, arch, os)
		binaryName := fmt.Sprintf("%s-%s-%s", name, os, arch)

		output = output.WithFile(fmt.Sprintf("bin/%s/%s", platform, binaryName), binary)
	}

	return output
}

// baseContainer creates a base container with Go toolchain configured
func (m *Go) baseContainer(goos, goarch, cgoEnabled string, musl bool) *dagger.Container {
	ctr := dag.Container().From(fmt.Sprintf("golang:%s", m.Version))

	if musl {
		envCC := "musl-gcc"
		ctr = ctr.
			WithExec([]string{"apt-get", "update"}).
			WithExec([]string{
				"apt-get", "install", "--no-install-recommends", "--yes",
				"musl", "musl-dev", "musl-tools",
			})

		if goarch == "arm64" {
			ctr = ctr.WithExec([]string{
				"/bin/sh", "-c",
				`curl -sfL https://musl.cc/aarch64-linux-musl-cross.tgz | tar -xzC /opt`,
			})
			envCC = "/opt/aarch64-linux-musl-cross/bin/aarch64-linux-musl-gcc"
		}

		ctr = ctr.WithEnvVariable("CC", envCC)
	}

	return ctr.
		WithMountedDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume(fmt.Sprintf("go-mod-%s", m.Version))).
		WithMountedCache("/go/build-cache", dag.CacheVolume(fmt.Sprintf("go-build-%s", m.Version))).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithEnvVariable("GOOS", goos).
		WithEnvVariable("GOARCH", goarch).
		WithEnvVariable("CGO_ENABLED", cgoEnabled)
}
