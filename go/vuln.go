package main

import (
	"context"
	"fmt"
)

// Vuln scans the Go module and its dependencies using the official tool.
// Fails (returns a non-nil error) if a known vulnerability is detected.
// +check
func (m *Go) Vuln(ctx context.Context) error {
	_, err := dag.Container().
		From(fmt.Sprintf("golang:%s", m.Version)).
		WithMountedDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume(fmt.Sprintf("go-mod-%s", m.Version))).
		WithExec([]string{"go", "install", "golang.org/x/vuln/cmd/govulncheck@latest"}).
		WithExec([]string{"govulncheck", "./..."}).
		Sync(ctx)

	return err
}
