package main

import "context"

func (m *Go) Lint(
	ctx context.Context,
) (string, error) {
	return dag.Container().
		From("golangci/golangci-lint:latest").
		WithMountedDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec([]string{"golangci-lint", "run", "-v"}).
		Stdout(ctx)
}
