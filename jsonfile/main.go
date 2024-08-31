package main

import (
	"context"

	"dagger/jsonfile/internal/dagger"
)

type Jsonfile struct{}

func (m *Jsonfile) Lint(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	// +default="."
	filedir string,
) (string, error) {
	containerImage := "pipelinecomponents/jsonlint:latest"

	args := []string{
		"jsonlint",
		"--diff",
		filedir,
	}

	return dag.Container().
		From(containerImage).
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec(args).
		Stdout(ctx)
}
