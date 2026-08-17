package main

import (
	"context"

	"dagger/goreleaser/internal/dagger"
)

type Goreleaser struct {
	Worktree *dagger.Directory
}

func New(
	// +defaultPath="."
	source *dagger.Directory,
) *Goreleaser {
	return &Goreleaser{Worktree: source}
}

func (m *Goreleaser) run(ctx context.Context, execArgs []string) (string, error) {
	return dag.Container().
		From("ghcr.io/goreleaser/goreleaser:latest").
		WithMountedDirectory("/data", m.Worktree).
		WithWorkdir("/data").
		WithExec(execArgs, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		Stdout(ctx)
}
