package main

import (
	"context"

	"dagger/jsonfile/internal/dagger"
)

type Jsonfile struct{}

// +check
func (m *Jsonfile) Lint(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
	// +optional
	filedir []string,
) error {
	execArgs := []string{
		"jsonlint",
		"--diff",
	}

	if len(filedir) > 0 {
		execArgs = append(execArgs, filedir...)
	} else {
		execArgs = append(execArgs, ".")
	}

	_, err := dag.Container().
		From("pipelinecomponents/jsonlint:latest").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec(execArgs).
		Sync(ctx)

	return err
}
