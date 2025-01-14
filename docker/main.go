package main

import (
	"context"
  "strings"

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
  platform dagger.Platform,
	// +optional
	// +default="Dockerfile"
	dockerfile string,
	// +optional
	buildArgs []string,
	// +optional
	secrets []*dagger.Secret,
) *dagger.Container {
	opts := dagger.ContainerBuildOpts{
		Dockerfile: dockerfile,
	}
	if len(buildArgs) > 0 {
    var args []dagger.BuildArg

    for _, arg := range buildArgs {
      nv := strings.Split(arg, "=")

      args = append(args, dagger.BuildArg{Name: nv[0], Value: nv[1]})
    }
		opts.BuildArgs = args
	}
	if len(secrets) > 0 {
		opts.Secrets = secrets
	}

  ctr := dag.Container(dagger.ContainerOpts{Platform: platform}).
    Build(source, opts)

  return ctr
}
