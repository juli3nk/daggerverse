package main

import (
	"fmt"
	"runtime"
	"strings"

	"dagger/go/internal/dagger"
)

func (m *Go) Build(
	// The binary name to build
	name string,
	// Go packages
	packages []string,
	// +optional
	// +default="1"
	cgoEnabled string,
  // +optional
  // +default="false"
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
	if len(arch) == 0 {
		arch = runtime.GOARCH
	}

  binaryName := fmt.Sprintf("%s-%s-%s", name, os, arch)
	binaryPath := fmt.Sprintf("build/%s", binaryName)

	goBuildLdflags := ldflags

  if musl {
    goBuildLdflags = append(goBuildLdflags,
      "-linkmode",
      "external",
    )
  }

	goBuildLdflags = append(goBuildLdflags,
		"-extldflags",
		"-static",
		"-s",
		"-w",
	)

	// WithEnvVariable("GOMODCACHE", "/go/pkg/mod").

  ctr := dag.Container().
    From(fmt.Sprintf("golang:%s", m.Version))

  if musl {
    envCC := "musl-gcc"

    ctr = ctr.WithExec([]string{"apt-get", "update"}).WithExec([]string{
      "apt-get", "install", "--no-install-recommends", "--yes",
      "musl",
      "musl-dev",
      "musl-tools",
    })

    if arch == "arm64" {
      ctr = ctr.WithExec([]string{"/bin/sh", "-c", `curl -sfL https://musl.cc/aarch64-linux-musl-cross.tgz | tar -xzC /opt`})

      envCC = "/opt/aarch64-linux-musl-cross/bin/aarch64-linux-musl-gcc"
    }

    ctr = ctr.WithEnvVariable("CC", envCC)
  }

  args := []string{
		"go",
		"build",
		"-o", binaryPath,
		"-ldflags",
    strings.Join(goBuildLdflags, " "),
	}
	args = append(args, packages...)

	return ctr.
		WithMountedDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume(fmt.Sprintf("go-mod-%s", m.Version))).
		WithMountedCache("/go/build-cache", dag.CacheVolume(fmt.Sprintf("go-build-%s", m.Version))).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithEnvVariable("GOOS", os).
		WithEnvVariable("GOARCH", arch).
		WithEnvVariable("CGO_ENABLED", cgoEnabled).
		WithExec(args).
		File(binaryPath)
}
