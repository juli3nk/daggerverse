package main

import (
	"context"

	"dagger/semantic-release/internal/dagger"
)

type SemanticRelease struct{}

func (m *SemanticRelease) Run(
	ctx context.Context,
	source *dagger.Directory,
	repoTokenEnvVarName string,
	repoTokenSecret *dagger.Secret,
	// +optional
	repositoryUrl string,
	// +optional
	// +default=true
	dryRun bool,
	// +optional
	// +default=true
	ci bool,
	// +optional
	// +default=false
	debug bool,
) (string, error) {
	var execArgs []string

	if len(repositoryUrl) > 0 {
		execArgs = append(execArgs, "--repository-url", repositoryUrl)
	}
	if dryRun {
		execArgs = append(execArgs, "--dry-run", "true")
	} else {
		execArgs = append(execArgs, "--dry-run", "false")
	}
	if ci {
		execArgs = append(execArgs, "--ci", "true")
	} else {
		execArgs = append(execArgs, "--ci", "false")
	}
	if debug {
		execArgs = append(execArgs, "--debug", "true")
	} else {
		execArgs = append(execArgs, "--debug", "false")
	}

	secretRepoToken, err := repoTokenSecret.Plaintext(ctx)
	if err != nil {
		return "", nil
	}

	return dag.Container().
		From("ghcr.io/juli3nk/semantic-release:main").
		WithMountedDirectory("/data", source).
		WithWorkdir("/data").
		WithEnvVariable(repoTokenEnvVarName, secretRepoToken).
		WithExec(execArgs, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		Stdout(ctx)
}
