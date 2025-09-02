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
	// +default=false
	dryRun bool,
	// +optional
	// +default=true
	ci bool,
	// +optional
	// +default=false
	debugMode bool,
) (string, error) {
	var execArgs []string

	if repositoryUrl != "" {
		execArgs = append(execArgs, "--repository-url", repositoryUrl)
	}

	if dryRun {
		execArgs = append(execArgs, "--dry-run")
	}

	if ci {
		execArgs = append(execArgs, "--ci")
	} else {
		execArgs = append(execArgs, "--no-ci")
	}

	if debugMode {
		execArgs = append(execArgs, "--debug")
	}

	secretRepoToken, err := repoTokenSecret.Plaintext(ctx)
	if err != nil {
		return "", err
	}

	return dag.Container().
		From("ghcr.io/juli3nk/semantic-release:main").
		WithMountedDirectory("/data", source).
		WithWorkdir("/data").
		WithEnvVariable("CI", "true").
		WithEnvVariable(repoTokenEnvVarName, secretRepoToken).
		WithExec(execArgs, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		Stdout(ctx)
}
