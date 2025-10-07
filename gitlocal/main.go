package main

import (
	"context"
	"fmt"
	"strings"

	"dagger/gitlocal/internal/dagger"
)

type Gitlocal struct {
	Worktree *dagger.Directory
}

func New(source *dagger.Directory) *Gitlocal {
	return &Gitlocal{Worktree: source}
}

// Uncommitted returns if there are uncommitted files
func (m *Gitlocal) Uncommitted(ctx context.Context) (bool, error) {
	result, err := m.git().
		WithDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec([]string{"status", "--short"}, dagger.ContainerWithExecOpts{
			UseEntrypoint: true,
		}).
		Stdout(ctx)
	if err != nil {
		return false, err
	}

	return len(strings.TrimSpace(result)) > 0, nil
}

// GetModifiedFiles returns list of modified files
func (m *Gitlocal) GetModifiedFiles(
	ctx context.Context,
	// +optional
	// +default=false
	compareWithWorkingTree bool,
) ([]string, error) {
	var execArgs []string
	if compareWithWorkingTree {
		execArgs = []string{"diff", "--name-only", "HEAD"}
	} else {
		execArgs = []string{"diff-tree", "--no-commit-id", "--name-only", "-r", "HEAD"}
	}

	result, err := m.git().
		WithDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec(execArgs, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		Stdout(ctx)
	if err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(result), "\n")
	var filtered []string
	for _, f := range files {
		if f != "" {
			filtered = append(filtered, f)
		}
	}

	return filtered, nil
}

// GetLatestCommit returns the short SHA of HEAD
func (m *Gitlocal) GetLatestCommit(ctx context.Context) (string, error) {
	result, err := m.git().
		WithDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec([]string{"rev-parse", "--short", "HEAD"}, dagger.ContainerWithExecOpts{
			UseEntrypoint: true,
		}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

// GetLatestTag returns the most recent tag on HEAD (sorted by version)
func (m *Gitlocal) GetLatestTag(ctx context.Context) (string, error) {
	result, err := m.git().
		WithDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec([]string{
			"describe",
			"--tags",
			"--abbrev=0",
		}, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

// GetAllTags returns all tags pointing to HEAD
func (m *Gitlocal) GetAllTags(ctx context.Context) ([]string, error) {
	result, err := m.git().
		WithDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec([]string{
			"tag",
			"--points-at",
			"HEAD",
		}, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		Stdout(ctx)
	if err != nil {
		return nil, err
	}

	tags := strings.Split(strings.TrimSpace(result), "\n")
	var filtered []string
	for _, tag := range tags {
		if tag != "" {
			filtered = append(filtered, tag)
		}
	}

	return filtered, nil
}

// GetBranch returns the current branch name
func (m *Gitlocal) GetBranch(ctx context.Context) (string, error) {
	result, err := m.git().
		WithDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec([]string{
			"rev-parse",
			"--abbrev-ref",
			"HEAD",
		}, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

// GetCommitMessage returns the commit message of HEAD
func (m *Gitlocal) GetCommitMessage(ctx context.Context) (string, error) {
	result, err := m.git().
		WithDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec([]string{
			"log",
			"-1",
			"--pretty=%B",
		}, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

// GetCommitAuthor returns the author of HEAD
func (m *Gitlocal) GetCommitAuthor(ctx context.Context) (string, error) {
	result, err := m.git().
		WithDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec([]string{
			"log",
			"-1",
			"--pretty=%an <%ae>",
		}, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

// GetCommitDate returns the date of HEAD
func (m *Gitlocal) GetCommitDate(ctx context.Context) (string, error) {
	result, err := m.git().
		WithDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec([]string{
			"log",
			"-1",
			"--pretty=%cI",
		}, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

type GitInfo struct {
	Commit        string
	Branch        string
	Tag           string
	Dirty         bool
	Message       string
	Author        string
	Date          string
	ModifiedFiles []string
}

// Info returns all git information as a structured object
func (m *Gitlocal) Info(ctx context.Context) (*GitInfo, error) {
	commit, err := m.GetLatestCommit(ctx)
	if err != nil {
		return nil, err
	}

	branch, err := m.GetBranch(ctx)
	if err != nil {
		return nil, err
	}

	tag, _ := m.GetLatestTag(ctx) // Ignore error if no tag exists

	dirty, err := m.Uncommitted(ctx)
	if err != nil {
		return nil, err
	}

	message, err := m.GetCommitMessage(ctx)
	if err != nil {
		return nil, err
	}

	author, err := m.GetCommitAuthor(ctx)
	if err != nil {
		return nil, err
	}

	date, err := m.GetCommitDate(ctx)
	if err != nil {
		return nil, err
	}

	modifiedFiles, err := m.GetModifiedFiles(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get modified files: %w", err)
	}

	return &GitInfo{
		Commit:        commit,
		Branch:        branch,
		Tag:           tag,
		Dirty:         dirty,
		Message:       message,
		Author:        author,
		Date:          date,
		ModifiedFiles: modifiedFiles,
	}, nil
}

func (m *Gitlocal) git() *dagger.Container {
	return dag.Apko().Wolfi().
		WithPackages([]string{"git"}).
		Container().
		WithEntrypoint([]string{"git"})
}
