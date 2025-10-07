package main

import (
	"context"

	"dagger/go/internal/dagger"
)

// GenerateDocs generates API documentation
func (m *Go) GenerateDocs(ctx context.Context) *dagger.Directory {
	return dag.Container().
		From("golang:"+m.Version).
		WithMountedDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec([]string{"go", "doc", "-all"}).
		Directory("/src/docs")
}
