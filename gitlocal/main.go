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
	result, err := m.git(ctx, []string{"status", "--short"})
	if err != nil {
		return false, err
	}

	return len(strings.TrimSpace(result)) > 0, nil
}

// GetModifiedFiles returns list of modified files
func (m *Gitlocal) GetModifiedFiles(
	ctx context.Context,
	compareWithWorkingTree bool,
	// +optional
	// +default="HEAD"
	compare string,
) ([]string, error) {
	var execArgs []string
	if compareWithWorkingTree || strings.Contains(compare, "..") {
		// Range de commits (ex: "origin/main..HEAD")
		execArgs = []string{"diff", "--name-only", compare}
	} else {
		// Fichiers du commit spécifié (marche même pour le 1er commit)
		execArgs = []string{"show", "--name-only", "--pretty=format:", compare}
	}

	result, err := m.git(ctx, execArgs)
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
	result, err := m.git(ctx, []string{"rev-parse", "--short", "HEAD"})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

// GetLatestTag returns the most recent tag on HEAD (sorted by version)
func (m *Gitlocal) GetLatestTag(ctx context.Context) (string, error) {
	result, err := m.git(ctx, []string{"describe", "--tags", "--abbrev=0"})
	if err != nil {
		if strings.Contains(err.Error(), "No names found") {
			return "", nil
		}
		return "", err
	}

	return strings.TrimSpace(result), nil
}

// GetAllTags returns all tags pointing to HEAD
func (m *Gitlocal) GetAllTags(ctx context.Context) ([]string, error) {
	result, err := m.git(ctx, []string{"tag", "--points-at", "HEAD"})
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
	result, err := m.git(ctx, []string{"rev-parse", "--abbrev-ref", "HEAD"})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

// GetCommitMessage returns the commit message of HEAD
func (m *Gitlocal) GetCommitMessage(ctx context.Context) (string, error) {
	result, err := m.git(ctx, []string{"log", "-1", "--pretty=%B"})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

// GetCommitAuthor returns the author of HEAD
func (m *Gitlocal) GetCommitAuthor(ctx context.Context) (string, error) {
	result, err := m.git(ctx, []string{"log", "-1", "--pretty=%an <%ae>"})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

// GetCommitDate returns the date of HEAD
func (m *Gitlocal) GetCommitDate(ctx context.Context) (string, error) {
	result, err := m.git(ctx, []string{"log", "-1", "--pretty=%cI"})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

func (m *Gitlocal) git(ctx context.Context, execArgs []string) (string, error) {
	fullArgs := append([]string{"--no-pager"}, execArgs...)

	ctr := dag.Apko().Wolfi().
		WithPackages([]string{"git"}).
		Container().
		WithDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithEntrypoint([]string{"git"}).
		WithExec(fullArgs, dagger.ContainerWithExecOpts{
			UseEntrypoint: true,
			Expect:        dagger.ReturnTypeAny,
		})

	exitCode, err := ctr.ExitCode(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get exit code: %w", err)
	}

	if exitCode != 0 {
		stderr, _ := ctr.Stderr(ctx)
		return "", fmt.Errorf("git command failed (exit %d): %s",
			exitCode, strings.TrimSpace(stderr))
	}

	result, err := ctr.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to read git stdout: %w", err)
	}

	return result, nil
}
