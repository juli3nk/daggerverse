package main

import (
	"context"
	"fmt"
	"strings"

	"dagger/gitlocal/internal/dagger"
)

type Gitlocal struct{}

func (m *Gitlocal) gitRaw(
	ctx context.Context,
	worktree *dagger.Directory,
	execArgs []string,
) (stdout, stderr string, exitCode int, err error) {
	fullArgs := append([]string{"--no-pager"}, execArgs...)

	ctr := dag.Apko().Wolfi().
		WithPackages([]string{"git"}).
		Container().
		WithDirectory("/src", worktree).
		WithWorkdir("/src").
		WithEntrypoint([]string{"git"}).
		WithExec(fullArgs, dagger.ContainerWithExecOpts{
			UseEntrypoint: true,
			Expect:        dagger.ReturnTypeAny,
		})

	exitCode, err = ctr.ExitCode(ctx)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to get exit code: %w", err)
	}

	stdout, err = ctr.Stdout(ctx)
	if err != nil {
		return "", "", exitCode, fmt.Errorf("failed to read git stdout: %w", err)
	}

	stderr, err = ctr.Stderr(ctx)
	if err != nil {
		return "", "", exitCode, fmt.Errorf("failed to read git stderr: %w", err)
	}

	return stdout, stderr, exitCode, nil
}

// Version “strict” pour les commandes où *tout* exitCode != 0 est une erreur.
func (m *Gitlocal) git(
	ctx context.Context,
	worktree *dagger.Directory,
	execArgs []string,
) (string, error) {
	stdout, stderr, exitCode, err := m.gitRaw(ctx, worktree, execArgs)
	if err != nil {
		return "", err
	}

	if exitCode != 0 {
		return "", fmt.Errorf("git command failed (exit %d): %s",
			exitCode, strings.TrimSpace(stderr))
	}

	return stdout, nil
}
