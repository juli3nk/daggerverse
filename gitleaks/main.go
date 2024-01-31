package main

import (
  "context"
  "strconv"
)

type Gitleaks struct {}

const gitleaksImage = "zricethezav/gitleaks:latest"

func (m *Gitleaks) Detect(
  ctx context.Context,
  repo *Directory,
  exitCode Optional[int],
  source Optional[string],
) (string, error) {
  ec := exitCode.GetOr(0)
  src := source.GetOr(".")

  return dag.
    Container().
    From(gitleaksImage).
    WithMountedDirectory("/src", repo).
    WithWorkdir("/src").
    WithExec([]string{"detect", "--source", src, "--exit-code", strconv.Itoa(ec), "--verbose"}).
    Stdout(ctx)
}
