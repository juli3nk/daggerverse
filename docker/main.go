package main

import (
	"context"

	"dagger/docker/internal/dagger"
)

type Docker struct{}

func (m *Docker) Lint(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	dockerfile []string,
	// +optional
	format string,
	// +optional
	ignore []string,
	// +optional
	failureThreshold string,
) (string, error) {
	containerImage := "hadolint/hadolint:latest"

	args := []string{
		"hadolint",
	}
	args = append(args, dockerfile...)

	if len(format) > 0 {
		args = append(args, "--format", format)
	}
	if len(ignore) > 0 {
		for _, ign := range ignore {
			args = append(args, "--ignore", ign)
		}
	}
	if len(failureThreshold) > 0 {
		args = append(args, "--failure-threshold", failureThreshold)
	}

	return dag.Container().
		From(containerImage).
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec(args).
		Stdout(ctx)
}

func (m *Docker) Build(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	// +default="."
	dir *dagger.Directory,
	// +optional
	// +default="Dockerfile"
	dockerfile string,
	// +optional
	buildArgs []string,
	// +optional
	secrets []*dagger.Secret,
) (*dagger.Container, error) {  
	opts := ContainerBuildOpts{
		Dockerfile: dockerfile,
	}
	if len(buildArgs) > 0 {
		opts.BuildArgs: buildArgs,
	}
	if len(secrets) > 0 {
		opts.Secrets: secrets,
	}

	return dag.Container().
	  Build(dir, opts).
	  Stdout(ctx)
}
