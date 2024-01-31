package main

import (
  "context"
)

type Shellcheck struct {}

const shellcheckImage = "koalaman/shellcheck:stable"

func (m *Shellcheck) Detect(
  ctx context.Context,
  repo *Directory,
  exitCode Optional[int],
  source Optional[string],
) (string, error) {
  //ec := exitCode.GetOr(0)
  src := source.GetOr(".")

  return dag.
    Container().
    From(shellcheckImage).
    WithMountedDirectory("/src", repo).
    WithWorkdir("/src").
    WithExec([]string{src}).
    Stdout(ctx)
}
