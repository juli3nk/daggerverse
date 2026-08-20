package main

import (
	"context"
	"fmt"
	"strings"

	"dagger/gitlocal/internal/dagger"
)

// LatestTag returns the latest reachable tag from HEAD.
// Returns an empty string if no tags exist (without error).
func (m *Gitlocal) LatestTag(
	ctx context.Context,
	// +defaultPath="."
	repo *dagger.Directory,
) (string, error) {
	stdout, stderr, exitCode, err := gitRaw(ctx, repo,
		[]string{"describe", "--tags", "--abbrev=0"},
	)
	if err != nil {
		return "", err
	}

	if exitCode != 0 {
		// No tag: "No names found" or "cannot describe"
		if strings.Contains(stderr, "No names found") ||
			strings.Contains(stderr, "cannot describe") {
			return "", nil
		}
		return "", fmt.Errorf("git describe --tags --abbrev=0 failed (exit %d): %s",
			exitCode, strings.TrimSpace(stderr))
	}

	return strings.TrimSpace(stdout), nil
}
