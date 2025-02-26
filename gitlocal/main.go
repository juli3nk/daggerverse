package main

import (
	"context"
	"strings"

	"dagger/gitlocal/internal/dagger"
)

type Gitlocal struct {
	Worktree *dagger.Directory
}

func New(
	source *dagger.Directory,
) *Gitlocal {
	return &Gitlocal{Worktree: source}
}

// Returns if there is uncommited files
func (m *Gitlocal) Uncommited() (bool, error) {
	opts := dagger.ContainerWithExecOpts{
		UseEntrypoint: true,
	}
	execArgs := []string{
		"status",
		"--short",
	}

	result, err := m.git().
		WithDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec(execArgs, opts).
		Stdout(context.TODO())
	if err != nil {
		return false, err
	}

	if len(result) == 0 {
		return false, nil
	}

	return true, nil
}

func (m *Gitlocal) GetModifiedFiles(
	ctx context.Context,
	compareWithWorkingTree bool,
) ([]string, error) {
	opts := dagger.ContainerWithExecOpts{
		UseEntrypoint: true,
	}
	var execArgs []string
	if compareWithWorkingTree {
		execArgs = []string{
			"diff",
			"--name-only",
			"HEAD",
		}
	} else {
		execArgs = []string{
			"diff-tree",
			"--no-commit-id",
			"--name-only",
			"-r",
			"HEAD",
		}
	}

	result, err := m.git().
		WithDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec(execArgs, opts).
		Stdout(ctx)
	if err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(result), "\n")

	return files, nil
}

func (m *Gitlocal) GetLatestCommit() (string, error) {
	opts := dagger.ContainerWithExecOpts{
		UseEntrypoint: true,
	}
	execArgs := []string{
		"rev-parse",
		"--short",
		"HEAD",
	}

	result, err := m.git().
		WithDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec(execArgs, opts).
		Stdout(context.TODO())
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

func (m *Gitlocal) GetLatestTag() (string, error) {
	opts := dagger.ContainerWithExecOpts{
		UseEntrypoint: true,
	}
	execArgs := []string{
		"tag",
		"--list",
		"--contains",
		"HEAD",
	}

	result, err := m.git().
		WithDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec(execArgs, opts).
		Stdout(context.TODO())
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

func (m *Gitlocal) git() *dagger.Container {
	return dag.Apko().Wolfi().
		WithPackages([]string{"git"}).
		Container().
		WithEntrypoint([]string{"git"})
}
