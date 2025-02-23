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

func (m *Gitlocal) GetModifiedFiles(
	ctx context.Context,
	compareWithWorkingTree bool,
) ([]string, error) {
	var execArgs []string
	if compareWithWorkingTree {
		execArgs = []string{
			"git",
			"diff",
			"--name-only",
			"HEAD",
		}
	} else {
		execArgs = []string{
			"git",
			"diff-tree",
			"--no-commit-id",
			"--name-only",
			"-r",
			"HEAD",
		}
	}

	result, err := m.container().
		WithDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec(execArgs).
		Stdout(ctx)
	if err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(result), "\n")

	return files, nil
}

func (m *Gitlocal) GetLatestCommit() (string, error) {
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

func (m *Gitlocal) GetLatestTag() (string, error) {
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

func (m *Gitlocal) container() *dagger.Container {
	return dag.Apko().Wolfi().
		WithPackages([]string{"git"}).
		Container()
}
