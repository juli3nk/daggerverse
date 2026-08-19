package main

import (
	"context"
	"fmt"
)

// Verify checks if go.mod and go.sum are tidy
// +check
func (m *Go) ModCheck(ctx context.Context) error {
	// ── 1. Container de base avec cache module ──
	base := dag.Container().
		From(fmt.Sprintf("golang:%s", m.Version)).
		WithMountedDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume(fmt.Sprintf("go-mod-%s", m.Version)))

	// ── 2. Download + Verify ──
	_, err := base.
		WithExec([]string{"go", "mod", "download"}).
		WithExec([]string{"go", "mod", "verify"}).
		Sync(ctx)
	if err != nil {
		return fmt.Errorf("go mod download/verify failed: %w", err)
	}

	// ── 3. Capture les fichiers ORIGINAUX (avant tidy) ──
	//    On lit depuis m.Worktree (immutable) pour comparer après
	origMod, err := m.Worktree.File("go.mod").Contents(ctx)
	if err != nil {
		return fmt.Errorf("reading go.mod: %w", err)
	}

	// go.sum peut ne pas exister (ex: pas de deps externe)
	origSum := ""
	if f, err := m.Worktree.File("go.sum"); err == nil {
		origSum, _ = f.Contents(ctx)
	}

	// ── 4. go mod tidy dans le conteneur ──
	tidyCtr := base.WithExec([]string{"go", "mod", "tidy"})
	_, err = tidyCtr.Sync(ctx)
	if err != nil {
		return fmt.Errorf("go mod tidy failed: %w", err)
	}

	// ── 5. Capture les fichiers APRÈS tidy ──
	newMod, err := tidyCtr.File("/src/go.mod").Contents(ctx)
	if err != nil {
		return fmt.Errorf("reading tidied go.mod: %w", err)
	}

	newSum := ""
	if f, err := tidyCtr.File("/src/go.sum"); err == nil {
		newSum, _ = f.Contents(ctx)
	}

	// ── 6. Comparaison ──
	if origMod != newMod || origSum != newSum {
		return fmt.Errorf("go.mod and/or go.sum are not tidy. Run 'go mod tidy' and commit the changes")
	}

	return nil
}
