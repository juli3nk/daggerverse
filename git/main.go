package main

import (
	"context"
	"strings"

	"dagger/git/internal/dagger"
)

type Git struct {
	Worktree *dagger.Directory
}

func New(
	source *dagger.Directory,
) *Git {
	return &Git{Worktree: source}
}

// Returns if there is uncommited files
func (m *Git) Uncommited() (bool, error) {
	execArgs := []string{
		"git",
		"status",
		"--short",
	}

	result, err := m.container().
		WithDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec(execArgs).
		Stdout(context.TODO())
	if err != nil {
		return false, err
	}

	if len(result) == 0 {
		return false, nil
	}

	return true, nil
}

func (m *Git) GetLatestCommit() (string, error) {
	execArgs := []string{
		"git",
		"rev-parse",
		"--short",
		"HEAD",
	}

	result, err := m.container().
		WithDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec(execArgs).
		Stdout(context.TODO())
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

func (m *Git) GetLatestTag() (string, error) {
	execArgs := []string{
		"git",
		"tag",
		"--list",
		"--contains",
		"HEAD",
	}

	result, err := m.container().
		WithDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec(execArgs).
		Stdout(context.TODO())
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

func (m *Git) container() *dagger.Container {
	return dag.
		Wolfi().
		Container(dagger.WolfiContainerOpts{
			Packages: []string{"git"},
		})
}
