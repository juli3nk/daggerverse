package main

import (
  "context"
)

type Yaml struct {}

func (m *Yaml) Fmt(
  ctx context.Context,
  dir *Directory,
  write Optional[int],
  source Optional[string],
) (string, error) {
  containerImage := "cytopia/yamlfmt:stable"

  w := write.GetOr(false)
  src := source.GetOr(".")

  cmdArgs := []string{src}

  if w {
    cmdArgs = append(cmdArgs, "-w")
  }

  return dag.Container().
    From(containerImage).
    WithMountedDirectory("/mnt", dir).
    WithWorkdir("/mnt").
    WithExec(cmdArgs).
    Stdout(ctx)
}

func (m *Yaml) Lint(
  ctx context.Context,
  dir *Directory,
  source Optional[string],
) (string, error) {
  containerImage := "cytopia/yamllint:stable"

  src := source.GetOr(".")

  return dag.Container().
    From(containerImage).
    WithMountedDirectory("/mnt", dir).
    WithWorkdir("/mnt").
    WithExec([]string{src}).
    Stdout(ctx)
}
