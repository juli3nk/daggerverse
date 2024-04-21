package main

import (
  "context"
)

type Json struct {}

func (m *Json) Lint(
  ctx context.Context,
  dir *Directory,
  source Optional[string],
) (string, error) {
  containerImage := "cytopia/jsonlint:stable"

  src := source.GetOr(".")

  return dag.Container().
    From(containerImage).
    WithMountedDirectory("/mnt", dir).
    WithWorkdir("/mnt").
    WithExec([]string{src}).
    Stdout(ctx)
}
