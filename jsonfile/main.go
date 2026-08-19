package main

import (
	"context"
	"encoding/json"
	"fmt"

	"dagger/jsonfile/internal/dagger"
)

type Jsonfile struct{}

// +check
func (m *Jsonfile) Lint(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
	// If empty, lints the entire module.
	// Otherwise, expects a JSON array of paths: '["a.json","b.json"]'
	// +optional
	pathsJson string,
) error {
	execArgs := []string{"jsonlint", "--diff"}

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
		From("pipelinecomponents/jsonlint:latest").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec(execArgs).
		Sync(ctx)

	return err
}
