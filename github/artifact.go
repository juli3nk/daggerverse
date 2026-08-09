package main

import (
	"context"
	"fmt"

	"dagger/publisher/internal/dagger"
)

// CreateChecksums generates checksums for all files in a directory
func (m *Publisher) CreateChecksums(
	ctx context.Context,
	// Directory containing files
	artifacts *dagger.Directory,
	// Checksum algorithm (sha256, sha512, md5)
	// +optional
	// +default="sha256"
	algorithm string,
) (*dagger.File, error) {
	checksumCmd := fmt.Sprintf("%ssum", algorithm)

	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/artifacts", artifacts).
		WithWorkdir("/artifacts").
		WithExec([]string{"sh", "-c", fmt.Sprintf("%s * > checksums.txt", checksumCmd)}).
		File("/artifacts/checksums.txt"), nil
}

// CreateArchive creates a tar.gz archive from a directory
func (m *Publisher) CreateArchive(
	ctx context.Context,
	// Directory to archive
	directory *dagger.Directory,
	// Archive name (without extension)
	name string,
) *dagger.File {
	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/src", directory).
		WithWorkdir("/tmp").
		WithExec([]string{"tar", "-czf", fmt.Sprintf("%s.tar.gz", name), "-C", "/src", "."}).
		File(fmt.Sprintf("/tmp/%s.tar.gz", name))
}

// VerifyChecksums verifies checksums of files
func (m *Publisher) VerifyChecksums(
	ctx context.Context,
	// Directory containing files and checksums.txt
	artifacts *dagger.Directory,
	// Checksum algorithm
	// +optional
	// +default="sha256"
	algorithm string,
) (string, error) {
	checksumCmd := fmt.Sprintf("%ssum", algorithm)

	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/artifacts", artifacts).
		WithWorkdir("/artifacts").
		WithExec([]string{checksumCmd, "-c", "checksums.txt"}).
		Stdout(ctx)
}
