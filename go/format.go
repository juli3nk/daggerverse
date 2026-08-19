package main

import (
	"context"
	"encoding/json"
	"fmt"
)

// Fmt formats Go code (check mode)
// +check
func (m *Go) Fmt(
	ctx context.Context,
	// If empty, formats the entire module.
	// Otherwise, expects a JSON array of paths: '["a.go","b.go"]'
	// +optional
	pathsJson string,
) error {
	execArgs := []string{"gofmt", "-l"}

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
		From(fmt.Sprintf("golang:%s", m.Version)).
		WithMountedDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec(execArgs).
		Sync(ctx)

	return err
}
