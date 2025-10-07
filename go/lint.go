package main

import "context"

// Lint runs golangci-lint
func (m *Go) Lint(
	ctx context.Context,
	// +optional
	filedir []string,
) (string, error) {
	execArgs := []string{"golangci-lint", "run", "-v"}
	if len(filedir) > 0 {
		execArgs = append(execArgs, filedir...)
	}

	return dag.Container().
		From("golangci/golangci-lint:latest").
		WithMountedDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec(execArgs).
		Stdout(ctx)
}
