package main

import (
	"context"

	"dagger/yaml/internal/dagger"
)

type Yaml struct{}

// +check
func (m *Yaml) Fmt(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
	// +optional
	filedir []string,
) error {
	var execArgs []string

	if len(filedir) > 0 {
		execArgs = append(execArgs, filedir...)
	} else {
		execArgs = append(execArgs, ".")
	}

	_, err := dag.Container().
		From("cytopia/yamlfmt:stable").
		WithMountedDirectory("/mnt", source).
		WithWorkdir("/mnt").
		WithExec(execArgs).
		Sync(ctx)

	return err
}

// +check
func (m *Yaml) Lint(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
	// +optional
	filedir []string,
) error {
	execArgs := []string{
		"yamllint",
		"--diff",
	}

	if len(filedir) > 0 {
		execArgs = append(execArgs, filedir...)
	} else {
		execArgs = append(execArgs, ".")
	}

	_, err := dag.Container().
		From("pipelinecomponents/yamllint:latest").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec(execArgs).
		Sync(ctx)

	return err
}
