package main

import (
	"context"
	"encoding/json"
	"fmt"

	"dagger/markdown/internal/dagger"
)

type Markdown struct{}

// Validates markdown files
// +check
func (m *Markdown) Lint(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
	// If empty, lints the entire module.
	// Otherwise, expects a JSON array of paths: '["a.md","b.md"]'
	// +optional
	pathsJson string,
) error {
	execArgs := []string{"markdownlint"}

	if pathsJson != "" {
		var paths []string
		if err := json.Unmarshal([]byte(pathsJson), &paths); err != nil {
			return fmt.Errorf("invalid pathsJson (expected JSON array): %w", err)
		}
		if len(paths) > 0 {
			execArgs = append(execArgs, paths...)
		}
	} else {
		execArgs = append(execArgs, ".")
	}

	_, err := dag.Container().
		From("tmknom/markdownlint:latest").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec(execArgs).
		Sync(ctx)

	return err
}
