package main

import (
	"fmt"

	"dagger/builder/internal/dagger"

	cplatforms "github.com/containerd/platforms"
)

// CreateMinimalRootfs creates a minimal Alpine-based rootfs with CA certs and nobody user
func (m *Builder) CreateMinimalRootfs(alpineVersion string) *dagger.Directory {
	return dag.Container().
		From(fmt.Sprintf("alpine:%s", alpineVersion)).
		WithExec([]string{"apk", "--update", "add", "ca-certificates"}).
		WithExec([]string{"mkdir", "-p", "/output/etc/ssl/certs"}).
		WithExec([]string{"echo", "nobody:x:65534:"},
			dagger.ContainerWithExecOpts{RedirectStdout: "/output/etc/group"}).
		WithExec([]string{"echo", "nobody:x:65534:65534:nobody:/:"},
			dagger.ContainerWithExecOpts{RedirectStdout: "/output/etc/passwd"}).
		WithExec([]string{"cp", "/etc/ssl/certs/ca-certificates.crt", "/output/etc/ssl/certs/"}).
		Directory("/output")
}

// FileMapping represents a file to copy into the image
type FileMapping struct {
	// Destination path in the container
	Path string
	// Source file
	Source *dagger.File
}

type ImageBuildConfig struct {
	Type          string
	Directory     string
	BinaryName    string
	InjectVersion bool
	Packages      []string
	Entrypoint    []string
	Description   string
	Port          int
	// Ajouts génériques
	ExtraFiles []FileMapping
	BaseImage  string // Pour images non-Go
	User       string // Défaut: "nobody:nobody"
}

// Assemble binary + rootfs dans une image from scratch
func (m *Builder) Create(
	platform string,
	binary *dagger.File,
	binaryPath string,
	config ImageBuildConfig,
	scratch bool,
	rootfs *dagger.Directory,
	from string,
) (*dagger.Container, error) {
	// Parse and format the platform string
	var platformStr dagger.Platform
	if platform != "" {
		// Parse the platform string to validate it
		parsed, err := cplatforms.Parse(platform)
		if err != nil {
			return nil, fmt.Errorf("invalid platform %q: %w", platform, err)
		}
		platformStr = dagger.Platform(cplatforms.Format(parsed))
	} else {
		// Default to linux/amd64
		platformStr = dagger.Platform("linux/amd64")
	}

	image := dag.Container(dagger.ContainerOpts{Platform: platformStr})

	if scratch {
		image = image.WithRootfs(rootfs)
	} else {
		image = image.From(from)
	}

	image = image.WithFile(binaryPath, binary).
		WithEntrypoint(config.Entrypoint).
		WithUser("nobody:nobody")

	if config.Port > 0 {
		image = image.WithExposedPort(config.Port)
	}

	return image, nil
}
