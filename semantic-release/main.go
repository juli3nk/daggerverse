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
  // +default=false
  ci bool,
  // +optional
  // +default=false
  debug bool,
) (string, error) {
  var args []string

  containerImage := "ghcr.io/juli3nk/semantic-release:main"

  if len(repositoryUrl) > 0 {
    args = append(args, "--repository-url", repositoryUrl)
  }
  if dryRun {
    args = append(args, "--dry-run")
  }
  if ci {
    args = append(args, "--ci")
  }
  if debug {
    args = append(args, "--debug")
  }

  secretRepoToken, err := repoTokenSecret.Plaintext(ctx)
  if err != nil {
    return "", nil
  }

	return dag.Container().
    From(containerImage).
    WithMountedDirectory("/data", source).
		WithWorkdir("/data").
    WithEnvVariable(repoTokenEnvVarName, secretRepoToken).
		WithExec(args, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
    Stdout(ctx)
}
