package main

import (
	"context"
	"fmt"

	"dagger/go/internal/dagger"
)

// GenerateDocs generates API documentation
func (m *Go) GenerateDocs(ctx context.Context) *dagger.Directory {
	return dag.Container().
		From(fmt.Sprintf("golang:%s", m.Version)).
		WithMountedDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec([]string{"/bin/sh", "-c", "mkdir -p docs && for pkg in $(go list ./...); do echo \"# $pkg\" >> docs/api.txt; go doc -all \"$pkg\" >> docs/api.txt; done"}).
		Directory("/src/docs")
}
