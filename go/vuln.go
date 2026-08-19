package main

import (
	"context"
	"fmt"
)

// Vulncheck scanne le module Go et ses dépendances via l'outil officiel.
// Échoue (error non nil) si une vulnérabilité connue est détectée.
func (m *Go) Vulncheck(ctx context.Context) error {
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
