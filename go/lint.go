package main

import (
	"context"
	"encoding/json"
	"fmt"
)

// Lint runs golangci-lint
// +check
func (m *Go) Lint(
	ctx context.Context,
	// If empty, lints the entire module.
	// Otherwise, expects a JSON array of paths: '["a.go","b.go"]'
	// +optional
	pathsJson string,
) error {
	execArgs := []string{"golangci-lint", "run", "-v"}

	if pathsJson != "" {
		var paths []string
		if err := json.Unmarshal([]byte(pathsJson), &paths); err != nil {
			return fmt.Errorf("invalid pathsJson (expected JSON array): %w", err)
		}
		if len(paths) > 0 {
			execArgs = append(execArgs, paths...)
		}
	}

	_, err := dag.Container().
		From("golangci/golangci-lint:latest").
		WithMountedDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec(execArgs).
		Sync(ctx)

	return err
}
