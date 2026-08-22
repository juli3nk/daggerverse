package main

import (
	"context"

	"dagger/semantic-release/internal/dagger"
)

type SemanticRelease struct{}

func (m *SemanticRelease) Run(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
	repoTokenEnv string,
	repoToken *dagger.Secret,
	// +optional
	repoUrl string,
	// +optional
	// +default=false
	dryRun bool,
	// +optional
	// +default=true
	ci bool,
	// +optional
	// +default=false
	enableDebug bool,
) (string, error) {
	var execArgs []string

	if repoUrl != "" {
		execArgs = append(execArgs, "--repository-url", repoUrl)
	}

	if dryRun {
		execArgs = append(execArgs, "--dry-run")
	}

	if ci {
		execArgs = append(execArgs, "--ci")
	} else {
		execArgs = append(execArgs, "--no-ci")
	}

	if enableDebug {
		execArgs = append(execArgs, "--debug")
	}

	return dag.Container().
		From("ghcr.io/juli3nk/semantic-release:main").
		WithMountedDirectory("/data", source).
		WithWorkdir("/data").
		WithEnvVariable("CI", "true").
		WithSecretVariable(repoTokenEnv, repoToken).
		WithExec(execArgs, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		Stdout(ctx)
}
