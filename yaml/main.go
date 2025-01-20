package main

import (
	"context"

	"dagger/yaml/internal/dagger"
)

type Yaml struct{}

func (m *Yaml) Fmt(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	filedir []string,
) (string, error) {
	var execArgs []string

	if len(filedir) > 0 {
		execArgs = append(execArgs, filedir...)
	} else {
		execArgs = append(execArgs, ".")
	}

	return dag.Container().
		From("cytopia/yamlfmt:stable").
		WithMountedDirectory("/mnt", source).
		WithWorkdir("/mnt").
		WithExec(execArgs).
		Stdout(ctx)
}

func (m *Yaml) Lint(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	filedir []string,
) (string, error) {
	execArgs := []string{
		"yamllint",
		"--diff",
	}

	if len(filedir) > 0 {
		execArgs = append(execArgs, filedir...)
	} else {
		execArgs = append(execArgs, ".")
	}

	return dag.Container().
		From("pipelinecomponents/yamllint:latest").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec(execArgs).
		Stdout(ctx)
}
