package main

import "context"

// ScanDependencies checks for outdated dependencies
func (m *Go) ScanDependencies(ctx context.Context) (string, error) {
	return dag.Container().
		From("golang:"+m.Version).
		WithMountedDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec([]string{"go", "list",
			"-u",
			"-m",
			"-json", "all",
		}).
		Stdout(ctx)
}
