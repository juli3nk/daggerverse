package main

import (
	"context"
	"fmt"
	"strings"

	"dagger/gitlocal/internal/dagger"
)

// LatestTag retourne le dernier tag atteignable depuis HEAD.
// Renvoie une chaîne vide si aucun tag n'existe (sans erreur).
func (m *Gitlocal) LatestTag(
	ctx context.Context,
	repo *dagger.Directory,
) (string, error) {
	stdout, stderr, exitCode, err := m.gitRaw(ctx, repo,
		[]string{"describe", "--tags", "--abbrev=0"},
	)
	if err != nil {
		return "", err
	}

	if exitCode != 0 {
		// Pas de tag : "No names found" ou "cannot describe"
		if strings.Contains(stderr, "No names found") ||
			strings.Contains(stderr, "cannot describe") {
			return "", nil
		}
		return "", fmt.Errorf("git describe --tags --abbrev=0 failed (exit %d): %s",
			exitCode, strings.TrimSpace(stderr))
	}

	return strings.TrimSpace(stdout), nil
}
