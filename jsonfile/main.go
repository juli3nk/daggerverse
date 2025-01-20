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
	filedir []string,
) (string, error) {
	execArgs := []string{
		"jsonlint",
		"--diff",
	}

	if len(filedir) > 0 {
		execArgs = append(execArgs, filedir...)
	} else {
		execArgs = append(execArgs, ".")
	}

	return dag.Container().
		From("pipelinecomponents/jsonlint:latest").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec(execArgs).
		Stdout(ctx)
}
