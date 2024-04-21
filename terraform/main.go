package main

import (
  "context"
)

type Terraform struct {}

const containerImage = "hashicorp/terraform:latest"

func (m *Terraform) Fmt(
  ctx context.Context,
  dir *Directory,
  write Optional[int],
  source Optional[string],
) (string, error) {
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

func (m *Terraform) Lint(
  ctx context.Context,
  dir *Directory,
  source Optional[string],
) (string, error) {
  src := source.GetOr(".")

  return dag.Container().
    From(containerImage).
    WithMountedDirectory("/mnt", dir).
    WithWorkdir("/mnt").
    WithExec([]string{src}).
    Stdout(ctx)
}
