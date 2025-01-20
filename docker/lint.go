package main

import (
	"context"

	"dagger/docker/internal/dagger"
)

func (m *Docker) Lint(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	dockerfiles []string,
	// +optional
	format string,
	// +optional
	ignore []string,
	// +optional
	failureThreshold string,
) (string, error) {
	if len(dockerfiles) == 0 {
		dockerfiles = []string{"Dockerfile"}
	}

	execArgs := []string{
		"hadolint",
	}
	execArgs = append(execArgs, dockerfiles...)

	if len(format) > 0 {
		execArgs = append(execArgs, "--format", format)
	}
	if len(ignore) > 0 {
		for _, ign := range ignore {
			execArgs = append(execArgs, "--ignore", ign)
		}
	}
	if len(failureThreshold) > 0 {
		execArgs = append(execArgs, "--failure-threshold", failureThreshold)
	}

	return dag.Container().
		From("hadolint/hadolint:latest").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec(execArgs).
		Stdout(ctx)
}
