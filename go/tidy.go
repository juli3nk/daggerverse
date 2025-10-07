package main

import (
	"context"
	"fmt"
)

// Verify checks if go.mod and go.sum are tidy
func (m *Go) Verify(ctx context.Context) (string, error) {
	return dag.Container().
		From(fmt.Sprintf("golang:%s", m.Version)).
		WithMountedDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume(fmt.Sprintf("go-mod-%s", m.Version))).
		WithExec([]string{"go", "mod", "verify"}).
		Stdout(ctx)
}
