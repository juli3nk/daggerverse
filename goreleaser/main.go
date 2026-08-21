package main

import (
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

func (m *Goreleaser) run(execArgs []string) *dagger.Container {
	return dag.Container().
		From("ghcr.io/goreleaser/goreleaser:latest").
		WithMountedDirectory("/data", m.Worktree).
		WithWorkdir("/data").
		WithExec(execArgs, dagger.ContainerWithExecOpts{UseEntrypoint: true})
}
