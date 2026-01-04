package main

import (
	"context"
	"fmt"
)

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

	modifiedFiles, err := m.GetModifiedFiles(ctx, false, "origin/main")
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
