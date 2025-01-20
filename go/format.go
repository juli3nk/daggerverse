package main

import (
	"context"
)

func (m *Go) Fmt(
	ctx context.Context,
	// +optional
	filedir []string,
) (string, error) {
	execArgs := []string{
		"gofmt",
		"-l",
	}

	if len(filedir) > 0 {
		execArgs = append(execArgs, filedir...)
	}

	return dag.Container().
		From("golang:latest").
		WithMountedDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec(execArgs).
		Stdout(ctx)
}
