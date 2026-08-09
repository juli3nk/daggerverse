package main

import (
	"context"
	"fmt"

	"dagger/syft/internal/dagger"
)

type Syft struct{}

// SBOM generates Software Bill of Materials
func (m *Syft) SBOM(
	ctx context.Context,
	source *dagger.Directory,
	// Output format (spdx-json, cyclonedx-json, syft-json)
	// +optional
	// +default="spdx-json"
	format string,
) (*dagger.File, error) {
	return dag.Container().
		From("anchore/syft:latest").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{".", "-o", fmt.Sprintf("%s=/src/sbom.%s", format, format}), dagger.ContainerWithExecOpts{
			UseEntrypoint: true,
		}).
		File(fmt.Sprintf("/src/sbom.%s", format)), nil
}
