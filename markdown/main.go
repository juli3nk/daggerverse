package main

import (
	"context"

	"dagger/markdown/internal/dagger"
)

type Markdown struct{}

// Validates markdown files
// +check
func (m *Markdown) Lint(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
) error {
	_, err := dag.Container().
		From("tmknom/markdownlint:latest").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"markdownlint", "."}).
		Sync(ctx)

	return err
}
