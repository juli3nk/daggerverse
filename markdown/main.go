package main

import (
	"context"
	"dagger/markdown/internal/dagger"
)

type Markdown struct{}

// ValidateMarkdown validates markdown files
func (m *Markdown) Lint(
	ctx context.Context,
	source *dagger.Directory,
) (string, error) {
	return dag.Container().
		From("tmknom/markdownlint:latest").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"markdownlint", "."}).
		Stdout(ctx)
}
