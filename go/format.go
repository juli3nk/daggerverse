package main

import (
	"context"
	"fmt"
)

// Fmt formats Go code (check mode)
func (m *Go) Fmt(
	ctx context.Context,
	// +optional
	filedir []string,
) (string, error) {
	execArgs := []string{"gofmt", "-l"}
	if len(filedir) > 0 {
		execArgs = append(execArgs, filedir...)
	} else {
		execArgs = append(execArgs, ".")
	}

	return dag.Container().
		From(fmt.Sprintf("golang:%s", m.Version)).
		WithMountedDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec(execArgs).
		Stdout(ctx)
}
