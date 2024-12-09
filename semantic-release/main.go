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
  var args []string

  containerImage := "ghcr.io/juli3nk/semantic-release:main"

  if len(repositoryUrl) > 0 {
    args = append(args, "--repository-url", repositoryUrl)
  }
  if dryRun {
    args = append(args, "--dry-run", "true")
  } else {
    args = append(args, "--dry-run", "false")
  }
  if ci {
    args = append(args, "--ci", "true")
  } else {
    args = append(args, "--ci", "false")
  }
  if debug {
    args = append(args, "--debug", "true")
  } else {
    args = append(args, "--debug", "false")
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
