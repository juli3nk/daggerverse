package main

import (
	"context"
	"encoding/json"
	"fmt"

	"dagger/yaml/internal/dagger"
)

type Yaml struct{}

// +check
func (m *Yaml) Fmt(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
	// If empty, formats the entire module.
	// Otherwise, expects a JSON array of paths: '["a.go","b.go"]'
	// +optional
	pathsJson string,
) error {
	var execArgs []string

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
		From("cytopia/yamlfmt:stable").
		WithMountedDirectory("/mnt", source).
		WithWorkdir("/mnt").
		WithExec(execArgs).
		Sync(ctx)

	return err
}

// +check
func (m *Yaml) Lint(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
	// If empty, lints the entire module.
	// Otherwise, expects a JSON array of paths: '["a.go","b.go"]'
	// +optional
	pathsJson string,
) error {
	execArgs := []string{"yamllint", "--diff"}

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
		From("pipelinecomponents/yamllint:latest").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec(execArgs).
		Sync(ctx)

	return err
}
