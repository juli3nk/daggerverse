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
	// The arch to build for
	// +optional
	arch string,
	// The os to build for
	// +optional
	// +default="linux"
	os string,
	// Enable CGO (requires a C toolchain)
	// +optional
	// +default=false
	cgo bool,
	// Use musl-gcc for static C linking (implies static binary, only valid with cgo and os=linux)
	// +optional
	// +default=false
	musl bool,
	// Omit the symbol table (reduces size, breaks pprof function names and go tool nm)
	// +optional
	// +default=false
	stripSymbols bool,
	// Omit DWARF debug info (reduces size, breaks delve/gdb source-level debugging)
	// +optional
	// +default=false
	stripDebug bool,
	// +optional
	ldflags []string,
) (*dagger.File, error) {
	if arch == "" {
		arch = runtime.GOARCH
	}

	binaryName := fmt.Sprintf("%s-%s-%s", name, os, arch)
	binaryPath := fmt.Sprintf("build/%s", binaryName)

	// Validation : musl sans CGO n'a pas de sens
	if musl && !cgo {
		return nil, fmt.Errorf("musl requires cgo to be enabled")
	}

	// Validation : musl est principalement une solution Linux
	if musl && os != "linux" {
		return nil, fmt.Errorf("musl static builds are only supported for linux, got %q", os)
	}

	cgoEnabled := "0"
	if cgo {
		cgoEnabled = "1"
	}

	// Construction des ldflags
	goBuildLdflags := make([]string, len(ldflags))
	copy(goBuildLdflags, ldflags)

	// Si musl, on force le linking statique via le linker externe
	if musl {
		goBuildLdflags = append(goBuildLdflags,
			"-linkmode", "external",
			"-extldflags", "-static",
		)
	}

	if stripSymbols {
		goBuildLdflags = append(goBuildLdflags, "-s")
	}
	if stripDebug {
		goBuildLdflags = append(goBuildLdflags, "-w")
	}

	ctr := m.buildEnv(os, arch, cgoEnabled, musl)

	args := []string{
		"go",
		"build",
		"-o", binaryPath,
		"-ldflags", strings.Join(goBuildLdflags, " "),
	}
	args = append(args, packages...)

	return ctr.WithExec(args).File(binaryPath), nil
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
	// +default=false
	cgo bool,
	// +optional
	// +default=false
	musl bool,
	// +optional
	// +default=false
	stripSymbols bool,
	// +optional
	// +default=false
	stripDebug bool,
	// +optional
	ldflags []string,
) (*dagger.Directory, error) {
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

		binary, err := m.Build(ctx, name, packages, arch, os, cgo, musl, stripSymbols, stripDebug, ldflags)
		if err != nil {
			return nil, err
		}
		binaryName := fmt.Sprintf("%s-%s-%s", name, os, arch)

		output = output.WithFile(fmt.Sprintf("bin/%s/%s", platform, binaryName), binary)
	}

	return output, nil
}

// baseContainer creates a base container with Go toolchain configured
func (m *Go) buildEnv(goos, goarch, cgoEnabled string, musl bool) *dagger.Container {
	ctr := dag.Container().From(fmt.Sprintf("golang:%s", m.Version))

	// Installation de Zig (une seule fois, cross-compiler universel)
	ctr = ctr.WithExec([]string{"sh", "-c", `
        curl -sfL https://ziglang.org/download/0.13.0/zig-linux-$(uname -m)-0.13.0.tar.xz |
        tar -xJ --strip-components=1 -C /usr/local
    `})

	target := fmt.Sprintf("%s-%s", goarch, goos)
	if goos == "darwin" {
		target = fmt.Sprintf("%s-macos", goarch) // zig syntax
	}

	cc := fmt.Sprintf("zig cc -target %s", target)
	if musl {
		// Zig intègre musl, pas besoin de l'installer
		cc = fmt.Sprintf("zig cc -target %s-musl", target)
	}

	return ctr.
		WithMountedDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume(fmt.Sprintf("go-mod-%s", m.Version))).
		WithMountedCache("/go/build-cache", dag.CacheVolume(fmt.Sprintf("go-build-%s", m.Version))).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithEnvVariable("GOOS", goos).
		WithEnvVariable("GOARCH", goarch).
		WithEnvVariable("CGO_ENABLED", cgoEnabled).
		WithEnvVariable("CC", cc)
}
