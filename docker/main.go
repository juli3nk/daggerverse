package main

import (
  "context"
)

type Docker struct {}

func (m *Docker) Lint(
  ctx context.Context,
  dir *Directory,
  source Optional[string],
) (string, error) {
  containerImage := "hadolint/hadolint:latest"

  src := source.GetOr(".")

  return dag.Container().
    From(containerImage).
    WithMountedDirectory("/mnt", dir).
    WithWorkdir("/mnt").
    WithExec([]string{src}).
    Stdout(ctx)
}

func (m *Docker) Build(
  ctx context.Context,
  dir *Directory,
  dockerfile Optional[string],
) (string, error) {
  df := dockerfile.GetOr("Dockerfile")

  opts := ContainerBuildOpts{}

  return dag.Container().
    Build(dir).
    Stdout(ctx)
}

func (m *Docker) Push(
  ctx context.Context,
  dir *Directory,
  address string,
) (string, error) {
  df := address.GetOr("Dockerfile")

  opts := ContainerPublishOpts{}

  return dag.Container().
    WithRegistryAuth(address string, username string, secret *Secret).
    Publish(ctx, addr, opts).
    Stdout(ctx)
}
